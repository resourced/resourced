package agent

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	gocache "github.com/patrickmn/go-cache"
	"github.com/rcrowley/go-metrics"

	"github.com/resourced/resourced/libmap"
)

func (a *Agent) NewMetricsRegistryForSelf() metrics.Registry {
	r := metrics.NewRegistry()
	metrics.RegisterDebugGCStats(r)
	metrics.RegisterRuntimeMemStats(r)

	go metrics.CaptureDebugGCStats(r, time.Second*60)
	go metrics.CaptureRuntimeMemStats(r, time.Second*60)

	return r
}

func (a *Agent) buildResultDBPayloadFromKeyValueMetric(key string, value interface{}) map[string]interface{} {
	var subkey string

	chunks := strings.Split(key, ".")
	prefix := chunks[0]
	dataPath := "/r/" + prefix

	existingRecord, existingRecordExists := a.ResultDB.Get(dataPath)

	nestedMap := libmap.NewTSafeNestedMapInterface(nil)

	if existingRecordExists {
		existingData, existingDataExists := existingRecord.(map[string]interface{})["Data"]
		if existingDataExists {
			// NOTE: this is not great, but totally solve the problem:
			// panic: interface conversion: interface {} is *interface {}, not map[string]interface {}
			asJSON, err := json.Marshal(existingData)
			if err == nil {
				var existingDataAsMapInterface map[string]interface{}

				if err := json.Unmarshal(asJSON, &existingDataAsMapInterface); err == nil {
					nestedMap = libmap.NewTSafeNestedMapInterface(existingDataAsMapInterface)
				}
			}
		}
	}

	subkey = strings.Replace(key, prefix+".", "", 1)
	nestedMap.Set(subkey, value)

	// Build record and envelopes the data
	record := make(map[string]interface{})
	record["UnixNano"] = time.Now().UnixNano()
	record["Path"] = dataPath
	record["Data"] = nestedMap.All()

	host, err := a.hostData()
	if err == nil {
		record["Host"] = host
	}

	return record
}

func (a *Agent) saveRawKeyValueMetricToResultDB(key string, value interface{}) {
	recordToBeSaved := a.buildResultDBPayloadFromKeyValueMetric(key, value)

	chunks := strings.Split(key, ".")
	prefix := chunks[0]
	dataPath := "/r/" + prefix

	a.ResultDB.Set(dataPath, recordToBeSaved, gocache.DefaultExpiration)
}

// Gather all aggregated StatsD metrics and store them in-memory storage.
func (a *Agent) flushStatsDMetricsToResultDBOnce() {
	percentiles := []float64{0.5, 0.75, 0.95, 0.99, 0.999}
	durationUnit := time.Nanosecond
	du := float64(durationUnit)

	a.StatsDMetrics.Each(func(name string, i interface{}) {
		switch metric := i.(type) {
		case metrics.Counter:
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.count", name), metric.Count())
		case metrics.Gauge:
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.value", name), metric.Value())
		case metrics.GaugeFloat64:
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.value", name), metric.Value())
		case metrics.Histogram:
			h := metric.Snapshot()
			ps := h.Percentiles(percentiles)

			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.count", name), h.Count())
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.min", name), h.Min())
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.max", name), h.Max())
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.mean", name), h.Mean())
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.std-dev", name), h.StdDev())

			for psIdx, psKey := range percentiles {
				key := strings.Replace(strconv.FormatFloat(psKey*100.0, 'f', -1, 64), ".", "", 1)
				a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.%s-percentile", name, key), ps[psIdx])
			}
		case metrics.Meter:
			m := metric.Snapshot()

			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.count", name), m.Count())
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.mean", name), m.RateMean())
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.1-minute", name), m.Rate1())
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.5-minute", name), m.Rate5())
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.15-minute", name), m.Rate15())

		case metrics.Timer:
			t := metric.Snapshot()
			ps := t.Percentiles(percentiles)

			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.count", name), t.Count())
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.min", name), t.Min()/int64(du))
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.max", name), t.Max()/int64(du))
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.mean", name), t.Mean()/du)
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.std-dev", name), t.StdDev()/du)

			for psIdx, psKey := range percentiles {
				key := strings.Replace(strconv.FormatFloat(psKey*100.0, 'f', -1, 64), ".", "", 1)
				a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.%s-percentile", name, key), ps[psIdx]/du)
			}

			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.mean-rate", name), t.RateMean())
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.1-minute", name), t.Rate1())
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.5-minute", name), t.Rate5())
			a.saveRawKeyValueMetricToResultDB(fmt.Sprintf("%s.15-minute", name), t.Rate15())
		}
	})
}

func (a *Agent) FlushStatsDMetricsToResultDB(statsInterval time.Duration) {
	for _ = range time.Tick(statsInterval) {
		a.flushStatsDMetricsToResultDBOnce()
	}
}
