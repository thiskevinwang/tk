package cmd

import (
	"os"
	"path/filepath"

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
		Use:   "cdk-go",
		Short: `cdk-go is a toolkit for various daily tasks`,
		Long:  `cdk-go is a toolkit for various daily tasks`,
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yml)")

	rootCmd.Flags().BoolVarP(&aFlag, "apple", "a", false, "a flag")
	rootCmd.Flags().BoolVarP(&bFlag, "banana", "b", false, "a flag")
	rootCmd.Flags().BoolVarP(&cFlag, "cherry", "c", false, "a flag")
	rootCmd.Flags().BoolVarP(&dFlag, "durian", "d", false, "a flag")
}

func initConfig() {

	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(aurora.Red(err))
	}

	configDir := filepath.Join(home, ".mc")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	viper.AutomaticEnv()
	viper.ReadInConfig()

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
}
