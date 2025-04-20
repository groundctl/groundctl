package output

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type Format int

const (
	Normal Format = iota
	Logrus
	JSON
)

var FormatIds = map[Format][]string{
	Normal: {"normal"},
	Logrus: {"logrus"},
	JSON:   {"json"},
}

var formatFormatters = map[Format]func() logrus.Formatter{
	Normal: newNormalFormatter,
	Logrus: newLogrusFormatter,
	JSON:   newJSONFormatter,
}

var ErrInvalidFormat = errors.New("invalid format")

func SetFormat(format Format) error {
	newFormatter, ok := formatFormatters[format]
	if !ok {
		return ErrInvalidFormat
	}
	logrus.SetFormatter(newFormatter())
	return nil
}

func newNormalFormatter() logrus.Formatter {
	return &NormalFormatter{}
}

func newLogrusFormatter() logrus.Formatter {
	return &logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		DisableColors:          false,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	}
}

func newJSONFormatter() logrus.Formatter {
	return &logrus.JSONFormatter{}
}

type NormalFormatter struct{}

func (f *NormalFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message + "\n"), nil
}
