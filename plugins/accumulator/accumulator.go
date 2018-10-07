// TODO:
// - Rename accumulator to metrics.
// - Move to list folder.
// - Rename Find to FetchOne

package accumulator

// Tag for metric.
type Tag struct {
	Name  string
	Value string
}

// Metric is a collection with many Tags and Values.
type Metric struct {
	Key    string
	Tags   []Tag
	Values interface{}
}

// Value is a collection for specific metric.
type Value struct {
	Key   string
	Value interface{}
}

// Items is a collection of metrics
type Items []Metric

var items *Items

// Load is a singleton method to return same object.
func Load() *Items {
	if items == nil {
		items = &Items{}
	}
	return items
}

// Reset the metric accumulator.
func (l *Items) Reset() {
	*l = (*l)[:0]
}

// Count all metrics in accumulator.
func (l *Items) Count() int {
	return len(*l)
}

// Add is aggregator for metric in accumulator.
func (l *Items) Add(m Metric) {
	if !items.Unique(m) {
		*l = append(*l, m)
	} else {
		items.Accumulator(m)
	}
}

// Find and return specific metric.
func (l *Items) Find(key string, tagName string, tagValue string) (interface{}) {
	for itemIndex := 0; itemIndex < len(*l); itemIndex++ {
		if (*l)[itemIndex].Key == key {
			for _, metricTag := range (*l)[itemIndex].Tags {
				if metricTag.Name == tagName && metricTag.Value == tagValue {
					return (*l)[itemIndex].Values
				}
			}
		}
	}
	return -1
}

// Unique is a check to verify the metric key is one in the accumulator.
func (l *Items) Unique(m Metric) bool {
	for _, i := range *l {
		if i.Key == m.Key && TagsEquals(i.Tags, m.Tags) {
			return true
		}
	}
	return false
}

// Accumulator sum values when we have the same key.
func (l *Items) Accumulator(m Metric) {
	for itemIndex := 0; itemIndex < len(*l); itemIndex++ {
		if (*l)[itemIndex].Key == m.Key && TagsEquals((*l)[itemIndex].Tags, m.Tags) == true {
			for itemValueIndex, itemValue := range (*l)[itemIndex].Values.([]Value) {
				for _, metricValue := range m.Values.([]Value) {
					if itemValue.Key == metricValue.Key {
						sumValue := metricValue.Value.(uint)
						oldValue := (*l)[itemIndex].Values.([]Value)[itemValueIndex].Value.(uint)
						newValue := oldValue + sumValue

						(*l)[itemIndex].Values.([]Value)[itemValueIndex].Value = newValue

						break
					}
				}
			}
		}
	}
}

// TagsEquals verify two Tags are equals.
func TagsEquals(a, b []Tag) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
