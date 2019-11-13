package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/oidc-mock/internal/entrypoint"
	"go.aporeto.io/oidc-mock/internal/versions"
)

func main() {

	if err := oidcMockCmd.Execute(); err != nil {
		panic(err)
	}
}

var config = &entrypoint.Configuration{}

// oidcMockCmd represents the base command when called without any subcommands
var oidcMockCmd = &cobra.Command{
	Use:   "oidcmock [options]",
	Short: "oidcmock command line interface",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		config.LogLevel = viper.GetString("log-level")
		config.LogFormat = viper.GetString("log-format")
		config.ServerIP = viper.GetString("server-ip")
		config.ServerPort = viper.GetString("server-port")
		config.DevelopmentMode = viper.GetBool("dev")
		time.Local = time.UTC
		return nil
	},
	RunE: func(ccmd *cobra.Command, args []string) (err error) {
		if viper.GetBool("version") {
			fmt.Println(versions.GetVersions())
			return nil
		}

		entrypoint.StartServer(config)
		return nil
	},
}

func init() {

	cobra.OnInitialize(initConfig)

	// Defaults
	oidcMockCmd.Flags().StringVar(&config.LogLevel, "log-level", "info", "Set the log-level between info, debug, trace")
	oidcMockCmd.Flags().StringVar(&config.LogFormat, "log-format", "human", "Set the log-format between console, json")

	oidcMockCmd.Flags().StringVar(&config.ServerIP, "server-ip", "192.168.100.1", "Set the default server ip")
	oidcMockCmd.Flags().StringVar(&config.ServerPort, "server-port", ":6999", "Set the default server port")

	oidcMockCmd.Flags().BoolVar(&config.DevelopmentMode, "dev", false, "Enable development mode")

	oidcMockCmd.Flags().StringVar(&config.PrivateKeyPath, "private-key", ".data/oidc.rsa", "Set the default private key")
	oidcMockCmd.Flags().StringVar(&config.PublicKeyPath, "public-key", ".data/oidc.rsa.pub", "Set the default public key")

	oidcMockCmd.Flags().BoolP("version", "v", false, "Show version.")

	if err := viper.BindPFlags(oidcMockCmd.Flags()); err != nil {
		panic(err)
	}
}

func initConfig() {

	viper.SetEnvPrefix("OIDCMOCK")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}
