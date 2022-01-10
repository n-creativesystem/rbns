package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"

	sig "os/signal"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	envPrefix = "NS"
)

var (
	configFile     string
	SignalReceived = errors.New("signal received")
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rbns",
		Short: "Role based N Security(RBAC)",
		Long:  "Role based N Security(RBAC)",
	}
	cmd.AddCommand(newRunCmd())
	return cmd
}

func Execute() error {
	cmd := newRootCmd()
	return cmd.Execute()
}

func initialize(cmd *cobra.Command) {
	v := viper.GetViper()
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		home, _ := os.UserHomeDir()
		v.AddConfigPath(".")
		v.AddConfigPath(path.Join("/etc", ".rbns"))
		v.AddConfigPath(path.Join(home, ".rbns"))
		v.SetConfigName("config")
	}
	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// config file does not found in search path
		default:
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	bindFlags(cmd, v)
}

func signal(ctx context.Context) error {
	signals := []os.Signal{
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGSTOP,
	}
	osNotify := make(chan os.Signal, 1)
	sig.Notify(osNotify, signals...)
	select {
	case <-ctx.Done():
		sig.Reset()
		return nil
	case _ = <-osNotify:
		return SignalReceived
	}
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		name := strings.ToUpper(strcase.ToSnake(f.Name))
		v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, name))
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
