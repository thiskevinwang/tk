package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   `version`,
	Short: `Print the version number`,
	Long:  `Print the version number`,
	Run: func(versionCmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
