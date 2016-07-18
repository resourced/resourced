// Package agent runs readers, writers, and HTTP server.
package agent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	gocache "github.com/patrickmn/go-cache"
	"github.com/rcrowley/go-metrics"
	"github.com/satori/go.uuid"

	"github.com/cskr/pubsub"
	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/executors"
	"github.com/resourced/resourced/host"
	"github.com/resourced/resourced/libmap"
	"github.com/resourced/resourced/libstring"
	"github.com/resourced/resourced/libtime"
	"github.com/resourced/resourced/loggers"
	"github.com/resourced/resourced/readers"
	"github.com/resourced/resourced/writers"
)

// New is the constructor for Agent struct.
func New() (*Agent, error) {
	agent := &Agent{}

	agent.ID = uuid.NewV4().String()

	err := agent.setConfigs()
	if err != nil {
		return nil, err
	}

	err = agent.setTags()
	if err != nil {
		return nil, err
	}

	err = agent.setAccessTokens()
	if err != nil {
		return nil, err
	}

	agent.ResultDB = gocache.New(time.Duration(agent.GeneralConfig.TTL)*time.Second, 10*time.Second)
	agent.ExecutorCounterDB = libmap.NewTSafeMapCounter(nil)

	agent.LiveLogPubSub = pubsub.New(int(agent.GeneralConfig.LogReceiver.ChannelCapacity))
	agent.LiveLogSubscribers = make(map[string]chan interface{})

	agent.StatsDMetrics = metrics.NewRegistry()

	return agent, err
}

// Agent struct carries most of the functionality of ResourceD.
// It collects information through readers and serve them up as HTTP+JSON.
type Agent struct {
	ID                 string
	Tags               map[string]string
	AccessTokens       []string
	Configs            *resourced_config.Configs
	GeneralConfig      resourced_config.GeneralConfig
	DbPath             string
	StatsDMetrics      metrics.Registry
	ResultDB           *gocache.Cache
	ExecutorCounterDB  *libmap.TSafeMapCounter
	LiveLogPubSub      *pubsub.PubSub
	LiveLogSubscribers map[string]chan interface{}
}

// Run executes a reader/writer/executor/log config.
func (a *Agent) Run(config resourced_config.Config) (output []byte, err error) {
	if config.GoStruct != "" && config.Kind == "reader" {
		output, err = a.runGoStructReader(config)
	} else if config.GoStruct != "" && config.Kind == "writer" {
		output, err = a.runGoStructWriter(config)
	} else if config.GoStruct != "" && config.Kind == "executor" {
		output, err = a.runGoStructExecutor(config)
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":              err.Error(),
			"config.GoStruct":    config.GoStruct,
			"config.Path":        config.Path,
			"config.Interval":    config.Interval,
			"config.Kind":        config.Kind,
			"config.ReaderPaths": fmt.Sprintf("%s", config.ReaderPaths),
		}).Error("Failed to execute runGoStructReader/runGoStructWriter/runGoStructExecutor")
	}

	err = a.saveRun(config, output, err)

	return output, err
}

// initGoStructReader initialize and return IReader.
func (a *Agent) initGoStructReader(config resourced_config.Config) (readers.IReader, error) {
	return readers.NewGoStructByConfig(config)
}

// initGoStructWriter initialize and return IWriter.
func (a *Agent) initGoStructWriter(config resourced_config.Config) (writers.IWriter, error) {
	writer, err := writers.NewGoStructByConfig(config)
	if err != nil {
		return nil, err
	}

	// Set configs data.
	writer.SetConfigs(a.Configs)

	// Get readers data.
	readersData := make(map[string][]byte)

	for _, readerPath := range config.ReaderPaths {
		readerJsonBytes, err := a.GetRunByPath(config.PathWithKindPrefix("r", readerPath))
		if err == nil {
			readersData[readerPath] = readerJsonBytes
		}
	}

	writer.SetReadersDataInBytes(readersData)

	return writer, err
}

// initResourcedMasterWriter initialize ResourceD Master specific IWriter.
func (a *Agent) initResourcedMasterWriter(config resourced_config.Config) (writers.IWriter, error) {
	var apiPath string

	if config.GoStruct == "ResourcedMasterHost" {
		apiPath = "/api/hosts"
	}

	urlFromConfigInterface, ok := config.GoStructFields["Url"]
	if !ok || urlFromConfigInterface == nil { // Check if Url is not defined in config
		config.GoStructFields["Url"] = a.GeneralConfig.ResourcedMaster.URL + apiPath

	} else { // Check if Url does not contain apiPath
		urlFromConfig := urlFromConfigInterface.(string)
		if !strings.HasSuffix(urlFromConfig, apiPath) {
			config.GoStructFields["Url"] = a.GeneralConfig.ResourcedMaster.URL + apiPath
		}
	}

	// Check if username is not defined
	// If so, set GeneralConfig.ResourcedMaster.AccessToken as default
	usernameFromConfigInterface, ok := config.GoStructFields["Username"]
	if !ok || usernameFromConfigInterface == nil {
		config.GoStructFields["Username"] = a.GeneralConfig.ResourcedMaster.AccessToken

	}

	return a.initGoStructWriter(config)
}

// initGoStructExecutor initialize and return IExecutor.
func (a *Agent) initGoStructExecutor(config resourced_config.Config) (executors.IExecutor, error) {
	executor, err := executors.NewGoStructByConfig(config)
	if err != nil {
		return nil, err
	}

	goodItems := libmap.AllNonExpiredCache(a.ResultDB)
	goodItemsInBytes := make(map[string][]byte)

	for key, item := range goodItems {
		itemInJson, err := json.Marshal(item.Object)
		if err == nil {
			goodItemsInBytes[key] = itemInJson
		}
	}

	executor.SetReadersDataInBytes(goodItemsInBytes)
	executor.SetCounterDB(a.ExecutorCounterDB)
	executor.SetTags(a.Tags)

	// Check if ResourcedMasterURL is not defined
	// If so, set GeneralConfig.ResourcedMaster.URL as default
	if config.ResourcedMasterURL == "" {
		executor.SetResourcedMasterURL(a.GeneralConfig.ResourcedMaster.URL)
	}

	// Check if ResourcedMasterAccessToken is not defined
	// If so, set GeneralConfig.ResourcedMaster.AccessToken as default
	if config.ResourcedMasterAccessToken == "" {
		executor.SetResourcedMasterAccessToken(a.GeneralConfig.ResourcedMaster.AccessToken)
	}

	host, err := a.hostData()
	if err != nil {
		return nil, err
	}

	executor.SetHostData(host)

	return executor, nil
}

// runGoStruct executes Run() fom IReader/IWriter/IExecutor and returns the output.
// Note that IWriter and IExecutor also implement IReader.
func (a *Agent) runGoStruct(readerOrWriterOrExecutor readers.IReader) ([]byte, error) {
	err := readerOrWriterOrExecutor.Run()
	if err != nil {
		errData := make(map[string]string)
		errData["Error"] = err.Error()
		return json.Marshal(errData)
	}

	return readerOrWriterOrExecutor.ToJson()
}

// runGoStructReader executes IReader and returns the output.
func (a *Agent) runGoStructReader(config resourced_config.Config) ([]byte, error) {
	// Initialize IReader
	reader, err := a.initGoStructReader(config)
	if err != nil {
		return nil, err
	}

	return a.runGoStruct(reader)
}

// runGoStructWriter executes IWriter and returns error if exists.
func (a *Agent) runGoStructWriter(config resourced_config.Config) ([]byte, error) {
	var writer writers.IWriter
	var err error

	// Initialize IWriter
	if strings.HasPrefix(config.GoStruct, "ResourcedMaster") {
		writer, err = a.initResourcedMasterWriter(config)
		if err != nil {
			return nil, err
		}

	} else {
		writer, err = a.initGoStructWriter(config)
		if err != nil {
			return nil, err
		}
	}

	err = writer.GenerateData()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":              err.Error(),
			"config.GoStruct":    config.GoStruct,
			"config.Path":        config.Path,
			"config.Interval":    config.Interval,
			"config.Kind":        config.Kind,
			"config.ReaderPaths": fmt.Sprintf("%s", config.ReaderPaths),
		}).Error("Failed to execute writer.GenerateData()")

		return nil, err
	}

	return a.runGoStruct(writer)
}

// runGoStructExecutor executes IExecutor and returns the output.
func (a *Agent) runGoStructExecutor(config resourced_config.Config) ([]byte, error) {
	var executor executors.IExecutor
	var err error

	// Initialize IExecutor
	executor, err = a.initGoStructExecutor(config)
	if err != nil {
		return nil, err
	}

	return a.runGoStruct(executor)
}

// hostData builds host related information.
func (a *Agent) hostData() (*host.Host, error) {
	h, err := host.NewHostByHostname()
	if err != nil {
		return nil, err
	}

	h.Tags = a.Tags

	return h, nil
}

// saveRun gathers basic, host, and reader/witer information and save them into local storage.
func (a *Agent) saveRun(config resourced_config.Config, output []byte, err error) error {
	// Do not perform save if config.Path is empty.
	if config.Path == "" {
		return nil
	}

	// Do not perform save if both output and error are empty.
	if output == nil && err == nil {
		return nil
	}

	record := config.CommonJsonData()

	host, err := a.hostData()
	if err != nil {
		return err
	}
	record["Host"] = host

	if err == nil {
		runData := new(interface{})
		err = json.Unmarshal(output, &runData)
		if err != nil {
			return err
		}
		record["Data"] = runData

	} else {
		errMap := make(map[string]string)
		errMap["Error"] = err.Error()
		record["Data"] = errMap
	}

	a.ResultDB.Set(config.PathWithPrefix(), record, gocache.DefaultExpiration)

	return err
}

// GetRunByPath returns JSON data stored in local storage given path string.
func (a *Agent) GetRunByPath(path string) ([]byte, error) {
	valueInterface, found := a.ResultDB.Get(path)
	if found {
		return json.Marshal(valueInterface)
	}
	return nil, nil
}

// RunForever executes Run() in an infinite loop with a sleep of config.Interval.
func (a *Agent) RunForever(config resourced_config.Config) {
	go func(config resourced_config.Config) {
		waitTime, err := time.ParseDuration(config.Interval)
		if err != nil {
			waitTime, _ = time.ParseDuration("60s")
		}

		for range time.Tick(waitTime) {
			a.Run(config)
		}

		for {
			a.Run(config)
		}
	}(config)
}

func (a *Agent) RunLoggerForever(config resourced_config.Config) {
	logger, err := loggers.NewGoStructByConfig(config)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"Error":           err.Error(),
			"config.GoStruct": config.GoStruct,
			"config.Path":     config.Path,
			"config.Interval": config.Interval,
			"config.Kind":     config.Kind,
		}).Error("Failed to instantiate logger struct")
	}

	// Collect log lines from various sources:
	//   * live:// is the TCP listener.
	//   * regular file
	if strings.HasPrefix(logger.GetSource(), "live://") {
		// Subscribe every target to consume log line from LiveLogPubSub
		for _, target := range logger.GetTargets() {
			subscriberKey := logger.(loggers.ILoggerChannel).PubSubKey(target.Endpoint)
			subscriberChannel := a.LiveLogPubSub.Sub(subscriberKey)

			a.LiveLogSubscribers[subscriberKey] = subscriberChannel

			// Collect log lines by subscribing to pubsub channel
			go func(subscriberKey string, subscriberChannel chan interface{}) {
				defer close(subscriberChannel)
				logger.(loggers.ILoggerChannel).RunBlockingChannel(subscriberKey, subscriberChannel)
			}(subscriberKey, subscriberChannel)
		}

	} else {
		// Collect log lines by watching the file via inotify
		go func() {
			logger.(loggers.ILoggerFile).RunBlockingFile(logger.GetSource())
		}()
	}

	// flushTime: Interval in between flushing to multiple targets.
	flushTime := libtime.ParseDurationWithDefault(config.Interval, "60s")
	for range time.Tick(flushTime) {
		var loglines []string

		// Fetch log lines from logger's buffer if source is file
		if logger.GetSource() != "live://" {
			loglines = logger.GetAndResetLoglines(logger.GetSource())
		}

		// Send log lines to various targets.
		for _, target := range logger.GetTargets() {
			go func(logger loggers.ILogger, target resourced_config.LogTargetConfig) {

				// Fetch log lines from logger's buffer if source is live://
				// This block is here because the PubSub key is a composite of logger.Source and target.Endpoint.
				if logger.GetSource() == "live://" {
					loglines = logger.GetAndResetLoglines(logger.(loggers.ILoggerChannel).PubSubKey(target.Endpoint))
				}

				go func(loglines []string) {
					outputJSON, err := json.Marshal(loglines)
					if err != nil {
						logrus.WithFields(logrus.Fields{"Error": err.Error()}).Error("Failed to marshal log lines to JSON for /logs/* HTTP endpoint")

						// Check if we have to prune in-memory log lines.
						if int64(logger.GetLoglinesLength(logger.GetSource())) > logger.GetBufferSize() {
							logger.ResetLoglines(logger.GetSource())
						}
					}

					a.saveRun(config, outputJSON, err)
				}(loglines)

				if strings.HasPrefix(target.Endpoint, "http://RESOURCED_MASTER_URL") {
					// Target is ResourceD Master
					go func(loglines []string) {
						loglines = logger.ProcessOutgoingLoglines(loglines, config.DenyList)

						masterURLPath := strings.Replace(target.Endpoint, "http://RESOURCED_MASTER_URL", "", 1)

						hostData, err := a.hostData()
						if err != nil {
							logger.LogErrorAndResetLoglinesIfNeeded(logger.GetSource(), err, "Failed to get host data for sending log lines to ResourceD Master")
							return
						}

						err = logger.SendLogToMaster(a.GeneralConfig.ResourcedMaster.AccessToken, a.GeneralConfig.ResourcedMaster.URL, masterURLPath, hostData, loglines, logger.GetSource())
						if err != nil {
							logger.LogErrorAndResetLoglinesIfNeeded(logger.GetSource(), err, "Failed to send log lines to ResourceD Master")
						}
					}(loglines)

				} else if strings.HasPrefix(target.Endpoint, "resourced+tcp://") {
					// Target is another ResourceD Agent
					go func(loglines []string) {
						loglines = logger.ProcessOutgoingLoglines(loglines, config.DenyList)

						anotherAgentEndpoint := strings.Replace(target.Endpoint, "resourced+tcp://", "", 1)

						err = logger.SendLogToAgent(anotherAgentEndpoint, 3, loglines, logger.GetSource())
						if err != nil {
							logger.LogErrorAndResetLoglinesIfNeeded(logger.GetSource(), err, "Failed to forward log lines to another agent")
						}
					}(loglines)

				} else if strings.HasPrefix(target.Endpoint, "file://") {
					// Target is local file
					go func(loglines []string) {
						targetFile := libstring.ExpandTildeAndEnv(strings.Replace(target.Endpoint, "file://", "", 1))

						loglines = logger.ProcessOutgoingLoglines(loglines, config.DenyList)

						err = logger.WriteToFile(targetFile, loglines)
						if err != nil {
							logger.LogErrorAndResetLoglinesIfNeeded(logger.GetSource(), err, "Failed to write log lines to file")
						}
					}(loglines)

				}
			}(logger.(loggers.ILogger), target)
		}
	}
}

// RunAllForever runs everything in an infinite loop.
func (a *Agent) RunAllForever() {
	for _, config := range a.Configs.Readers {
		a.RunForever(config)
	}

	for _, config := range a.Configs.Writers {
		a.RunForever(config)
	}

	for _, config := range a.Configs.Executors {
		a.RunForever(config)
	}

	for _, config := range a.Configs.Loggers {
		a.RunLoggerForever(config)
	}
}
