package cmd

import (
	"encoding/json"
	"fmt"
	Command "four-key/command"
	"four-key/models"
	"four-key/settings"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path"
)

var setCommand = &cobra.Command{
	Use:   "set",
	Short: "set config",
	Long:  "set config",
	Run:   onSet,
}

func init() {
	rootCmd.AddCommand(setCommand)
	setCommand.Flags().StringP("output", "o", "", "Set output source of 4Key metrics results")
}

func onSet(cmd *cobra.Command, args []string) {
	output, err := cmd.Flags().GetString("output")
	commander := Command.ACommander()

	if output == "" || err != nil {
		fmt.Println(commander.Fatal("output parameter error please check and re run"))
	}

	document := models.Document{}
	existFile, err := ioutil.ReadFile(path.Join(commander.GetFourKeyPath(), settings.EnvironmentFileName))

	if err != nil {
		s, err := settings.Get()
		if s == nil {
			fmt.Println(commander.Fatal("The file does not exist!"))
			return
		}

		existFile, err = ioutil.ReadFile(path.Join(commander.GetFourKeyPath(), settings.EnvironmentFileName))
		if err != nil {
			fmt.Println(commander.Fatal("The file does not exist!"))
			return
		}
	}

	err = json.Unmarshal([]byte(string(existFile)), &document)
	if err != nil {
		fmt.Println(commander.Fatal("json parse error"))
		return
	}

	if output != "" {
		document.Output = output
	}

	err = writeFile(document)

	if err != nil {
		fmt.Println(commander.Fatal("write error. err -> ", err))
		return
	}
}
