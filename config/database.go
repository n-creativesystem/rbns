package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
)

type DataBase struct {
	Dialector string `envconfig:"DIALECTOR" default:"postgres"`
	MasterDSN string `envconfig:"MASTER_DSN"`
	SlaveDSN  string `envconfig:"SLAVE_DSN"`
}

func NewDataBase() DataBase {
	var conf DataBase
	envconfig.MustProcess("", &conf)
	if conf.SlaveDSN == "" {
		conf.SlaveDSN = conf.MasterDSN
	}
	conf.MasterDSN = os.ExpandEnv(conf.MasterDSN)
	conf.SlaveDSN = os.ExpandEnv(conf.SlaveDSN)
	return conf
}
