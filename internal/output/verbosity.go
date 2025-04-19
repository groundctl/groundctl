package output

import "github.com/sirupsen/logrus"

func SetVerbosity(verbosity int) {
	logrus.SetLevel(min(logrus.TraceLevel, logrus.InfoLevel+logrus.Level(verbosity)))
}
