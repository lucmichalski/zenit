package newrelic

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

func Normalize(items *metrics.Items) map[string]map[string]interface{} {
	events := make(map[string]map[string]interface{})

	for _, m := range *items {
		switch v := m.Values.(type) {
		case int, uint, uint64, float64:
			if _, ok := events[m.Key]; !ok {
				events[m.Key] = make(map[string]interface{})
				events[m.Key]["host"] = config.File.General.Hostname
				events[m.Key]["eventType"] = common.ToCamel(m.Key)
			}

			for t := range m.Tags {
				if m.Tags[t].Name == "name" {
					events[m.Key][m.Tags[t].Value] = v
				} else {
					events[m.Key][m.Tags[t].Name] = m.Tags[t].Value
				}
			}
		case []metrics.Value:
			if _, ok := events[m.Key]; !ok {
				events[m.Key] = make(map[string]interface{})
				events[m.Key]["host"] = config.File.General.Hostname
				events[m.Key]["eventType"] = common.ToCamel(m.Key)
			}

			for y := range v {
				events[m.Key][v[y].Key] = v[y].Value
			}
		}
	}

	return events
}