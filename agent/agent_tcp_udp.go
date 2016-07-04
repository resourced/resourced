package agent

import (
	"crypto/tls"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/narqo/go-dogstatsd-parser"
	"github.com/rcrowley/go-metrics"

	resourced_config "github.com/resourced/resourced/config"
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

func (a *Agent) HandleGraphite(dataInBytes []byte) {
	for _, singleMetric := range strings.Split(string(dataInBytes), "\n") {
		if singleMetric == "" {
			continue
		}

		dataInChunks := strings.Split(singleMetric, " ")

		logFields := logrus.Fields{
			"Metric": singleMetric,
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

					a.saveRawKeyValueMetricToResultDB(key, value)
				}
			} else {
				logFields["Error"] = err
				logrus.WithFields(logFields).Error("Failed to parse Graphite metric")
			}
		}
	}
}

func (a *Agent) HandleStatsD(dataInBytes []byte) {
	for _, singleMetric := range strings.Split(string(dataInBytes), "\n") {
		if singleMetric == "" {
			continue
		}

		logFields := logrus.Fields{
			"Metric": singleMetric,
		}

		statsdMetric, err := dogstatsd.Parse(singleMetric)
		if err != nil {
			logFields["Error"] = err
			logrus.WithFields(logFields).Error("Failed to parse StatsD metric")
			continue
		}

		// Don't do anything if there are no value to store.
		if statsdMetric.Name == "" || statsdMetric.Value == nil {
			continue
		}

		// Update log information
		logFields["Key"] = statsdMetric.Name
		if statsdMetric.Type != dogstatsd.Gauge {
			logFields["Value"] = statsdMetric.Value.(int64)
		}

		if statsdMetric.Type == dogstatsd.Counter {
			c := a.StatsDMetrics.GetOrRegister(statsdMetric.Name, metrics.NewCounter())
			if c != nil {
				c.(metrics.Counter).Inc(statsdMetric.Value.(int64))
				logrus.WithFields(logFields).Info("Increment StatsD counter")
			}

		} else if statsdMetric.Type == dogstatsd.Gauge {
			g := a.StatsDMetrics.GetOrRegister(statsdMetric.Name, metrics.NewGaugeFloat64())
			if g != nil {
				logFields["Value"] = statsdMetric.Value.(float64)
				g.(metrics.GaugeFloat64).Update(statsdMetric.Value.(float64))
				logrus.WithFields(logFields).Info("Update StatsD gauge")
			}

		} else if statsdMetric.Type == dogstatsd.Histogram {
			h := a.StatsDMetrics.GetOrRegister(statsdMetric.Name, metrics.NewHistogram(metrics.NewUniformSample(a.GeneralConfig.MetricReceiver.HistogramReservoirSize)))
			if h != nil {
				h.(metrics.Histogram).Update(statsdMetric.Value.(int64))
				logrus.WithFields(logFields).Info("Update StatsD historgram")
			}

		} else if statsdMetric.Type == dogstatsd.Meter {
			m := a.StatsDMetrics.GetOrRegister(statsdMetric.Name, metrics.NewMeter())
			if m != nil {
				m.(metrics.Meter).Mark(statsdMetric.Value.(int64))
				logrus.WithFields(logFields).Info("Mark StatsD meter")
			}

		} else if statsdMetric.Type == dogstatsd.Timer {
			t := a.StatsDMetrics.GetOrRegister(statsdMetric.Name, metrics.NewTimer())
			if t != nil {
				// t.Time(func() {})
				t.(metrics.Timer).Update(statsdMetric.Value.(time.Duration))
				logrus.WithFields(logFields).Info("Update StatsD timer")
			}
		}
	}
}

func (a *Agent) HandleLog(dataInBytes []byte) {
	a.TCPLogDB.Append("Loglines", string(dataInBytes))
}
