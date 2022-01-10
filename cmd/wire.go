// +build wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/n-creativesystem/rbns/registry"
	"github.com/spf13/pflag"
)

func initializeRun(flags *pflag.FlagSet) (*servers, error) {
	wire.Build(registry.SuperSet, newServer)
	return &servers{}, nil
}
