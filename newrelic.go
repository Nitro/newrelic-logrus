package newrelic_logrus

import (
	"errors"

	"github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
)

type NewRelicLogrusHook struct {
	Application newrelic.Application
	LogLevels   []logrus.Level
	didFire     bool
}

func NewNewRelicLogrusHook(app newrelic.Application, levels []logrus.Level) *NewRelicLogrusHook {
	return &NewRelicLogrusHook{
		Application: app,
		LogLevels:   levels,
	}
}

func (n *NewRelicLogrusHook) Levels() []logrus.Level {
	return n.LogLevels
}

func (n *NewRelicLogrusHook) Fire(entry *logrus.Entry) error {
	// Hacky. We don't know what transaction we're in so we
	// just start a new one specific to error reporting.
	txn := n.Application.StartTransaction("errorTxn", nil, nil)
	for k, v := range entry.Data {
		txn.AddAttribute(k, v)
	}
	txn.NoticeError(errors.New(entry.Message))
	txn.End()

	n.didFire = true // for testing only... there's no way to get data out of NR agent

	return nil
}
