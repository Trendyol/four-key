package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "four-key",
	Short: "four-key Metrics Command",
	Long:  "four-key Metrics Command",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
