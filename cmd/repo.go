package cmd

import (
	"fmt"
	"net/url"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(repoCmd)
}

var repoCmd = &cobra.Command{
	Use:   `repo`,
	Short: `Open github repository`,
	Long:  `Open github repository`,
	Run: func(versionCmd *cobra.Command, args []string) {
		github_url := "https://github.com"
		// items := []string{"thiskevinwang/thekevinwang.com", "thiskevinwang/blocstats", "thiskevinwang/tk", "thiskevinwang/cdk-go", "hashicorp/waypoint"}
		repos := viper.GetStringSlice("repos")
		index := -1
		var result string
		var err error

		for index < 0 {

			fmt.Println("No repositories found in:", aurora.Blue(viper.ConfigFileUsed()))
			prompt := promptui.SelectWithAdd{
				Label:    "Which repo do you want to open?",
				Items:    repos,
				AddLabel: "Add a repository...",
			}

			index, result, err = prompt.Run()

			if index == -1 {
				repos = append(repos, result)
				viper.Set("repos", repos)
				viper.WriteConfig()
			}
		}

		if err != nil {
			log.Fatal(fmt.Printf("Prompt failed %v\n", aurora.Red(err.Error())))
			return
		}

		fmt.Println(fmt.Sprintf("Opening: %s", aurora.Green(result)))

		baseUrl, err := url.Parse(github_url)
		if err != nil {
			log.Fatal(aurora.Red("Malformed URL: " + err.Error()))
			return
		}
		baseUrl.Path = result
		time.Sleep(500 * time.Millisecond)
		openbrowser(baseUrl.String())
	},
}
