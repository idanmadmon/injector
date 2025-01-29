package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "1.0.0"
var rootCmd = &cobra.Command{
	Use:     "injector",
	Version: version,
	Short:   "injector - a tool for sender and receiver tests",
	Long:    `injector - a tool for sender and receiver tests`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
