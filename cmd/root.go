package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/logrusorgru/aurora"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	aFlag bool
	bFlag bool
	cFlag bool
	dFlag bool

	Version = "v0.0.0"

	// rootCmd represents the base command when called without any subcommand
	rootCmd = &cobra.Command{
		Use:   "tk",
		Short: `tk is a toolkit for various daily tasks`,
		Long:  `tk is a toolkit for various daily tasks`,
		PreRun: func(rootCmd *cobra.Command, args []string) {

		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.Help()
				os.Exit(0)
			}
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(aurora.Red(err))
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yml)")

	rootCmd.Flags().BoolVarP(&aFlag, "apple", "a", false, "a flag")
}

func initConfig() {
	var logLevel = viper.GetString("log-level")

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(aurora.Red(err))
	}

	configDir := path.Join(home, ".tk")

	// create /Users/<user>/.tk if it doesn't exist
	// this allows viper.SafeWriteConfig() to work
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	// https://github.com/spf13/viper#reading-config-files
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	viper.AutomaticEnv()

	// this writes `{}` to /Users/<user>/.tk/config.yaml
	if err = viper.SafeWriteConfig(); err != nil {
		log.Error("SafeWriteConfig: ", err.Error())
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			// Should create on
			log.Warn(fmt.Sprintf("Config file not found"))

			if err = viper.WriteConfig(); err != nil {
				log.Error("WriteConfig: ", err.Error())
			}

		} else {
			// Config file was found but another error was produced
			log.Warn(aurora.Yellow("Config file was found but another error was produced"))
		}
	}

	if repos := viper.GetStringSlice("repos"); repos == nil {
		viper.Set("repos", []string{})
		viper.WriteConfig()
	}

	log.Info("ConfigFileUsed: ", aurora.Blue(viper.ConfigFileUsed()))

}
