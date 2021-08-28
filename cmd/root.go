package cmd

import (
	"context"
	"fmt"
	"os"
	"syscall"

	sig "os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Use:   "rbns",
		Short: "Role based N Security(RBAC)",
		Long:  "Role based N Security(RBAC)",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func initialize(before func(), after func()) func() {
	if before == nil {
		before = func() {}
	}
	if after == nil {
		after = func() {}
	}
	return func() {
		before()
		if configFile != "" {
			viper.SetConfigFile(configFile)
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			viper.AddConfigPath("/etc")
			viper.AddConfigPath(home)
			viper.SetConfigName("rbns")
		}

		viper.SetEnvPrefix("rbns")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			switch err.(type) {
			case viper.ConfigFileNotFoundError:
				// config file does not found in search path
			default:
				_, _ = fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
		after()
	}
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
	case s := <-osNotify:
		return fmt.Errorf("signal received: %v", s)
	}
}
