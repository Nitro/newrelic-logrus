package newrelic_logrus

import (
	"bytes"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewFormatter(t *testing.T) {
	Convey("NewFormatter() returns a properly configured Formatter", t, func() {
		fields := []string{"beowulf", "heorot"}
		jsonF := &logrus.JSONFormatter{}
		formatter := NewFormatter(fields, jsonF)

		So(formatter.fields, ShouldResemble,
			map[string]struct{}{
				"beowulf": struct{}{},
				"heorot":  struct{}{},
			},
		)
		So(formatter.formatter, ShouldEqual, jsonF)
	})
}

func Test_Format(t *testing.T) {
	Convey("Format()", t, func() {
		fields := []string{"beowulf", "heorot"}
		jsonF := &logrus.JSONFormatter{}
		formatter := NewFormatter(fields, jsonF)
		capture := &bytes.Buffer{}

		// Set up logrus for testing
		logrus.SetFormatter(formatter)
		logrus.SetOutput(capture)

		Reset(func() {
			logrus.SetFormatter(&logrus.TextFormatter{})
			logrus.SetOutput(os.Stdout)
		})

		Convey("filters out the New Relic transaction", func() {
			logrus.WithFields(logrus.Fields{"txn": "junk"}).Warn("Intentional warning")
			So(capture.String(), ShouldContainSubstring, "Intentional warning")
			So(capture.String(), ShouldNotContainSubstring, "txn")
		})

		Convey("filters out specified fields", func() {
			logrus.WithFields(logrus.Fields{"beowulf": "Geat", "heorot": "Hart"}).Warn("Intentional warning")
			So(capture.String(), ShouldContainSubstring, "Intentional warning")
			So(capture.String(), ShouldContainSubstring, `"beowulf":"[FILTERED]"`)
			So(capture.String(), ShouldContainSubstring, `"heorot":"[FILTERED]"`)
		})
	})
}
