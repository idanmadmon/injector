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
	streamDuplicatorReceiverAddress string
	streamDuplicatorSenderUrls      []string
	streamDuplicatorMaxOffsetDiff   int
)

var streamDuplicatorCmd = &cobra.Command{
	Use:   "stream-duplicator",
	Short: "create a stream-duplicator http receiver sender",
	Run: func(cmd *cobra.Command, args []string) {
		sd := injector.NewHttpStreamDuplicatorReceiverSender(streamDuplicatorReceiverAddress, streamDuplicatorSenderUrls, streamDuplicatorMaxOffsetDiff)

		ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancelCtx()

		go sd.Start()

		<-ctx.Done()
		sd.Stop()
	},
}

func init() {
	rootCmd.AddCommand(streamDuplicatorCmd)

	streamDuplicatorCmd.Flags().StringVarP(&streamDuplicatorReceiverAddress, "address", "a", "0.0.0.0:1236", "address to listen on")
	streamDuplicatorCmd.Flags().StringArrayVarP(&streamDuplicatorSenderUrls, "dest-url", "u", []string{}, "urls to send to, can enter multiple times")
	streamDuplicatorCmd.Flags().IntVarP(&streamDuplicatorMaxOffsetDiff, "maxoffset", "m", 0, "max offset between readers - default no limit")
}
