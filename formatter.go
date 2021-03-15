package newrelic_logrus

import (
	"github.com/sirupsen/logrus"
)

// A Formatter implements the logrus.Formatter interface and allows
// filtering out certain fields.
type Formatter struct {
	formatter logrus.Formatter
	fields    map[string]struct{}
}

// NewFormatter creates a properly configure formatter
func NewFormatter(fieldList []string, formatter logrus.Formatter) *Formatter {
	fields := make(map[string]struct{})

	for _, v := range fieldList {
		fields[v] = struct{}{}
	}
	return &Formatter{
		formatter: formatter,
		fields:    fields,
	}
}

// Format takes a Logrus entry and then formats it into the correct output. It
// does that by manipulating the fields and then calling the underlying
// formatter it was constructed with.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Remove the New Relic transaction
	delete(entry.Data, "txn")

	// Fiter out any other fields we want to masque
	for k, _ := range entry.Data {
		if _, ok := f.fields[k]; ok {
			entry.Data[k] = "[FILTERED]"
		}
	}

	data, err := f.formatter.Format(entry)
	return data, err
}
