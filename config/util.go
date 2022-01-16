package config

import (
	"github.com/n-creativesystem/rbns/ncsfw/logger"
	"github.com/spf13/pflag"
)

type flags struct {
	*pflag.FlagSet
}

func (f *flags) MustString(name string) string {
	value, err := f.GetString(name)
	if err != nil {
		logger.Panic(err, "")
	}
	return value
}

func (f *flags) MustInt(name string) int {
	value, err := f.GetInt(name)
	if err != nil {
		logger.Panic(err, "")
	}
	return value
}

func (f *flags) MustStringArray(name string) []string {
	value, err := f.GetStringArray(name)
	if err != nil {
		logger.Panic(err, "")
	}
	return value
}
