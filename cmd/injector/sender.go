package injector

import (
	"fmt"
	"os"

	units "github.com/docker/go-units"
	"github.com/spf13/cobra"
)

var (
	fileSize string
)

var senderCmd = &cobra.Command{
	Use:   "sender",
	Short: "create a http sender",
	Run: func(cmd *cobra.Command, args []string) {
		size, err := units.FromHumanSize(fileSize)
		if err != nil {
			fmt.Printf("file size given isn't valid: %s, error: %v\n", fileSize, err)
			os.Exit(1)
		}

		fmt.Printf("this is sender, fileSize: %s == %d\n", fileSize, size)
	},
}

func init() {
	senderCmd.Flags().StringVarP(&fileSize, "size", "s", "", "filesize - human size")
	rootCmd.AddCommand(senderCmd)
}
