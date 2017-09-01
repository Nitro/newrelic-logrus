package newrelic_logrus

import (
	"testing"

	"github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_EndToEnd(t *testing.T) {
	Convey("Logging the right level sends stuff to New Relic", t, func() {
		config := newrelic.NewConfig("Logrus Hook Tester", "1234567890123456789012345678901234567890")
		config.Enabled = false // Don't actually talk to New Relic, or validate license
		app, err := newrelic.NewApplication(config)

		So(err, ShouldBeNil)

		hook := NewNewRelicLogrusHook(
			app,
			[]logrus.Level{logrus.WarnLevel, logrus.ErrorLevel},
		)

		logrus.AddHook(hook)

		So(func() {
				logrus.WithFields(logrus.Fields{"key": "value"}).Warn("This is a test message")
			}, ShouldNotPanic,
		)
		So(hook.didFire, ShouldBeTrue)
	})
}
