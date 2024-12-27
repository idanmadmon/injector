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
	address   string
	readSpeed int64
)

var receiverCmd = &cobra.Command{
	Use:   "receiver",
	Short: "create a http receiver",
	Run: func(cmd *cobra.Command, args []string) {
		r := injector.NewHttpReceiver(readSpeed)

		ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancelCtx()

		go r.Start(address)

		<-ctx.Done()
		r.Stop()
	},
}

func init() {
	rootCmd.AddCommand(receiverCmd)

	receiverCmd.Flags().StringVarP(&address, "address", "a", "0.0.0.0:1238", "address to listen on")
	receiverCmd.Flags().Int64VarP(&readSpeed, "speed", "s", 0, "bytes to read per second - default no rate limit")
}
