package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "schedule",
	Short: "schedule",
	Long:  "schedule",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	pflags := rootCmd.PersistentFlags()
	pflags.StringVar(&cfgFile, "config", "/etc/trafic.yaml", "configuration file")

	pflags.String("log-tag", "", "a tag that is prepended to all log lines")
	viper.BindPFlag("log.tag", pflags.Lookup("log-tag"))

	pflags.String("flows-dir", "", "folder with flow configuration files")
	viper.BindPFlag("flows.dir", pflags.Lookup("flows-dir"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(filepath.Join("/etc", "/trafic", ".trafic"))
		viper.SetConfigName("trafic.yaml")
	}

	viper.SetEnvPrefix("TRAFIC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
