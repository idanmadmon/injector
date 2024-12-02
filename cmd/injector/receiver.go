package injector

import (
	"fmt"

	"github.com/spf13/cobra"
)

var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "create a http receiver",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("this is reciever\n")
	},
}

func init() {
	rootCmd.AddCommand(receiverCmd)
}
