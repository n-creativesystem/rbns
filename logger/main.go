package logger

import (
	"strings"

	"github.com/n-creativesystem/rbns/utilsconv"
	"github.com/sirupsen/logrus"
)

const (
	TimestampFormat = "2006/01/02 - 15:04:05"
)

var (
	format = strings.ToLower(utilsconv.DefaultGetEnv("LOG_FORMAT_TYPE", "JSON"))
)

func New() *logrus.Logger {
	log := logrus.New()
	SetFormatter(log)
	return log
}

func SetFormatter(log *logrus.Logger) {
	var formatter logrus.Formatter = GetFormatter()
	log.SetFormatter(formatter)
}

func GetFormat() string {
	return format
}

func GetFormatter() logrus.Formatter {
	switch GetFormat() {
	case "text":
		return &logrus.TextFormatter{
			TimestampFormat: TimestampFormat,
		}
	default:
		return &logrus.JSONFormatter{
			TimestampFormat: TimestampFormat,
		}
	}
}
