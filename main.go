package main

import (
	"time"

	"github.com/n-creativesystem/rbns/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	if err := cmd.Execute(); err != nil {
		logrus.Fatalln(err)
	}
}
