package cmd

import (
	"context"
	"injector"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	receiverAddress   string
	receiverReadSpeed int64
)

var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "create a http receiver",
	Run: func(cmd *cobra.Command, args []string) {
		r := injector.NewHttpReceiver(receiverAddress, receiverReadSpeed)

		ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancelCtx()

		go r.Start()

		<-ctx.Done()
		r.Stop()
	},
}

func init() {
	rootCmd.AddCommand(receiverCmd)

	receiverCmd.Flags().StringVarP(&receiverAddress, "address", "a", "0.0.0.0:1238", "address to listen on")
	receiverCmd.Flags().Int64VarP(&receiverReadSpeed, "speed", "s", 0, "bytes to read per second - default no rate limit")
}
