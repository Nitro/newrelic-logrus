package newrelic_logrus

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	newrelic "github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_EndToEnd(t *testing.T) {
	Convey("Testing End-to-End", t, func() {
		config := newrelic.NewConfig("Logrus Hook Tester", "1234567890123456789012345678901234567890")
		config.Enabled = false // Don't actually talk to New Relic, or validate license
		app, err := newrelic.NewApplication(config)

		So(err, ShouldBeNil)

		hook := NewNewRelicLogrusHook(
			app,
			[]logrus.Level{logrus.WarnLevel, logrus.ErrorLevel},
		)

		logrus.AddHook(hook)

		Convey("logging the right level sends stuff to New Relic", func() {
			log.SetOutput(ioutil.Discard)
			Reset(func() { log.SetOutput(os.Stdout) })

			logrus.WithFields(logrus.Fields{"key": "value"}).Warn("This is a test message")
			So(hook.didFire, ShouldBeTrue)
		})

		Convey("uses an existing transaction if it is passed", func() {
			capture := &bytes.Buffer{}
			log.SetOutput(capture)

			Reset(func() { log.SetOutput(os.Stdout) })

			txn := app.StartTransaction("testing", nil, nil)
			logrus.WithFields(logrus.Fields{"txn": txn, "key": "value"}).Warn("This is a test message")
			txn.End()

			So(hook.didFire, ShouldBeTrue)
			// Just make sure the point in the output string is our transaction
			So(capture.String(), ShouldContainSubstring, fmt.Sprintf("%v", txn))
		})
	})
}
