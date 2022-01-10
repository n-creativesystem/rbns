package main

import (
	"github.com/n-creativesystem/rbns/cmd"
	"github.com/n-creativesystem/rbns/logger"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logger.Fatal(err, "application main")
	}
}
