package newrelic_logrus

import (
	"errors"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type NewRelicLogrusHook struct {
	Application *newrelic.Application
	LogLevels   []logrus.Level
	didFire     bool
}

func NewNewRelicLogrusHook(app *newrelic.Application, levels []logrus.Level) *NewRelicLogrusHook {
	return &NewRelicLogrusHook{
		Application: app,
		LogLevels:   levels,
	}
}

func (n *NewRelicLogrusHook) Levels() []logrus.Level {
	return n.LogLevels
}

// withTransaction either retrieves the current transaction from the Entry.Data
// or it starts a new one. The transaction is then passed into the func it was
// passed.
func (n *NewRelicLogrusHook) withTransaction(entry *logrus.Entry, fn func(txn *newrelic.Transaction) error) error {
	var txn *newrelic.Transaction

	if entry.Data["txn"] != nil {
		txn = entry.Data["txn"].(*newrelic.Transaction)
	} else {
		// Hacky. We don't know what transaction we're in so we
		// just start a new one specific to error reporting.
		txn = n.Application.StartTransaction("errorTxn")
		defer txn.End()
	}

	return fn(txn)
}

func (n *NewRelicLogrusHook) Fire(entry *logrus.Entry) error {
	n.didFire = true // for testing only... there's no way to get data out of NR agent

	return n.withTransaction(entry, func(txn *newrelic.Transaction) error {
		for k, v := range entry.Data {
			if k != "txn" {
				txn.AddAttribute(k, v)
			}
		}

		txn.NoticeError(errors.New(entry.Message))
		return nil
	})
}
