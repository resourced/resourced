package agent

import (
	"crypto/tls"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/narqo/go-dogstatsd-parser"
	gocache "github.com/patrickmn/go-cache"
	"github.com/rcrowley/go-metrics"

	resourced_config "github.com/resourced/resourced/config"
	"github.com/resourced/resourced/libstring"
)

func (a *Agent) NewTCPServer(config resourced_config.ITCPServer, name string) (net.Listener, error) {
	if config.GetAddr() != "" {
		logFields := logrus.Fields{
			"LogReceiver.Addr": config.GetAddr(),
			"LogLevel":         a.GeneralConfig.LogLevel,
		}

		if config.GetCertFile() != "" && config.GetKeyFile() != "" {
			logFields["LogReceiver.CertFile"] = config.GetCertFile()
			logFields["LogReceiver.KeyFile"] = config.GetKeyFile()

			cert, err := tls.LoadX509KeyPair(config.GetCertFile(), config.GetKeyFile())
			if err != nil {
				logrus.WithFields(logFields).Fatal(err)
				return nil, err
			}

			logrus.WithFields(logFields).Info("Running " + name + "+SSL server")

			tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

			return tls.Listen("tcp", config.GetAddr(), tlsConfig)

		} else {
			logrus.WithFields(logFields).Info("Running " + name + " server")

			return net.Listen("tcp", config.GetAddr())
		}
	}

	return nil, nil
}

func (a *Agent) NewUDPServer(config resourced_config.ITCPServer, name string) (*net.UDPConn, error) {
	if config.GetAddr() != "" {
		logFields := logrus.Fields{
			"LogReceiver.Addr": config.GetAddr(),
			"LogLevel":         a.GeneralConfig.LogLevel,
		}

		udpAddr, err := net.ResolveUDPAddr("udp", config.GetAddr())
		if err != nil {
			logrus.WithFields(logFields).Error("Failed to run " + name + " server")
			return nil, err
		}

		logrus.WithFields(logFields).Info("Running " + name + " server")
		return net.ListenUDP("udp4", udpAddr)
	}

	return nil, nil
}

func (a *Agent) saveRawMetricToResultDB(key string, value interface{}) {
	var subkey string

	chunks := strings.Split(key, ".")
	prefix := chunks[0]
	dataPath := "/r/" + prefix
	data := make(map[string]interface{})

	hostnameIndex := libstring.FindHostnameChunkInMetricKey(key)
	if hostnameIndex == -1 {
		if len(chunks) > 1 {
			subkey = strings.Replace(key, prefix+".", "", 1)
		}

		data[subkey] = value

	} else {
		hostname := chunks[hostnameIndex]
		hostnameData := make(map[string]interface{})

		subkey = strings.Replace(key, strings.Join(chunks[0:hostnameIndex], ".")+".", "", 1)
		hostnameData[subkey] = value
		data[hostname] = hostnameData
	}

	// Create record envelope for data
	record := make(map[string]interface{})
	record["UnixNano"] = time.Now().UnixNano()
	record["Path"] = dataPath
	record["Data"] = data

	host, err := a.hostData()
	if err == nil {
		record["Host"] = host
	}

	a.ResultDB.Set(dataPath, record, gocache.DefaultExpiration)
}

func (a *Agent) HandleGraphite(dataInBytes []byte) {
	for _, data := range strings.Split(string(dataInBytes), "\n") {
		dataInChunks := strings.Split(data, " ")

		logFields := logrus.Fields{
			"Metric": string(dataInBytes),
		}

		if len(dataInChunks) >= 2 {
			key := dataInChunks[0]
			logFields["Key"] = key

			value, err := strconv.ParseFloat(dataInChunks[1], 64)
			if err == nil {
				// Loop through blacklist and set key-value if everything is good
				doSetValue := true

				for _, blacklistRegex := range a.GeneralConfig.MetricReceiver.BlacklistCompiled {
					if blacklistRegex.MatchString(key) {
						doSetValue = false
						break
					}
				}

				if doSetValue {
					logFields["Value"] = value
					logrus.WithFields(logFields).Info("Storing Graphite metric in memory")

					a.saveRawMetricToResultDB(key, value)
				}
			} else {
				logrus.WithFields(logFields).Info("Failed to parse Graphite metric")
			}
		}
	}
}

func (a *Agent) HandleStatsD(dataInBytes []byte) {
	statsdMetric, err := dogstatsd.Parse(string(dataInBytes))
	if err != nil {
		return
	}

	// Don't do anything if there are no value to store.
	if statsdMetric.Value == nil {
		return
	}

	if statsdMetric.Type == dogstatsd.Counter {
		c := metrics.NewCounter()
		a.StatsDMetrics.Register(statsdMetric.Name, c)
		c.Inc(statsdMetric.Value.(int64))

	} else if statsdMetric.Type == dogstatsd.Gauge {
		g := metrics.NewGauge()
		a.StatsDMetrics.Register(statsdMetric.Name, g)
		g.Update(statsdMetric.Value.(int64))

	} else if statsdMetric.Type == dogstatsd.Histogram {
		s := metrics.NewUniformSample(a.GeneralConfig.MetricReceiver.HistogramReservoirSize)
		h := metrics.NewHistogram(s)
		a.StatsDMetrics.Register(statsdMetric.Name, h)
		h.Update(statsdMetric.Value.(int64))

	} else if statsdMetric.Type == dogstatsd.Meter {
		m := metrics.NewMeter()
		a.StatsDMetrics.Register(statsdMetric.Name, m)
		m.Mark(statsdMetric.Value.(int64))

	} else if statsdMetric.Type == dogstatsd.Timer {
		// TODO(didip): Not sure what to do with Time() method here.
		// t := metrics.NewTimer()
		// a.StatsDMetrics.Register(statsdMetric.Name, t)
		// t.Time(func() {})
		// t.Update(statsdMetric.Value.(int64))
	}
}

func (a *Agent) HandleLog(dataInBytes []byte) {
	a.TCPLogDB.Append("Loglines", string(dataInBytes))
}
