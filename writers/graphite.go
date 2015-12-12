package writers

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/marpaia/graphite-golang"
	"github.com/nytlabs/gojsonexplode"
)

func init() {
	Register("NewGraphite", NewGraphite)
}

// NewGraphite is NewGraphite constructor.
func NewGraphite() IWriter {
	return &Graphite{}
}

// Graphite is a writer that serialize readers data to New Relic Insights.
type Graphite struct {
	Base
	Host string
	Port int
	grph *graphite.Graphite
}

// ToJson returns flattened data in JSON
func (g *Graphite) ToJson() ([]byte, error) {
	if g.Data == nil {
		return nil, errors.New("Data field is nil.")
	}

	dataInJson, err := json.Marshal(g.Data)
	if err != nil {
		return nil, err
	}

	return gojsonexplode.Explodejson(dataInJson, ".")
}

// Run sends data to remote graphite server
func (g *Graphite) Run() error {
	flattenData := make(map[string]interface{})

	flattenDataJson, err := g.ToJson()
	if err != nil {
		return err
	}

	err = json.Unmarshal(flattenDataJson, &flattenData)
	if err != nil {
		return err
	}

	if g.grph == nil {
		grph, err := graphite.NewGraphite(g.Host, g.Port)
		if err != nil {
			return err
		}
		g.grph = grph
	}

	metrics := make([]graphite.Metric, len(flattenData))

	index := 0
	for key, _ := range flattenData {
		// TODO: Why... is value a string?
		metrics[index] = graphite.NewMetric(key, "1", time.Now().Unix())
		index = index + 1
	}

	return g.grph.SendMetrics(metrics)
}
