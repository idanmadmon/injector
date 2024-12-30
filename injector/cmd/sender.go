package cmd

import (
	"context"
	"fmt"
	"injector"
	"os"
	"os/signal"
	"syscall"
	"time"

	units "github.com/docker/go-units"
	"github.com/spf13/cobra"
)

var (
	senderWorkersAmount int
	senderTimeout       string
	senderFileAmount    int
	senderFileSize      string
	senderUrl           string
)

var senderCmd = &cobra.Command{
	Use:   "sender",
	Short: "create a http sender",
	Run: func(cmd *cobra.Command, args []string) {
		size, err := units.FromHumanSize(senderFileSize)
		if err != nil {
			fmt.Printf("file size given isn't valid: %s, error: %v\n", senderFileSize, err)
			os.Exit(1)
		}

		t, err := time.ParseDuration(senderTimeout)
		if senderTimeout != "" && err != nil {
			fmt.Printf("timeout given isn't valid: %s, error: %v\n", senderTimeout, err)
			os.Exit(1)
		}

		conf := injector.SenderConfig{
			WorkersAmount: senderWorkersAmount,
			Timeout:       t,
			FileSize:      int(size),
			FileAmount:    senderFileAmount,
		}

		ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancelCtx()

		sender := injector.NewHttpSender(conf, senderUrl)
		sender.Start(ctx)
	},
}

func init() {
	senderCmd.Flags().IntVarP(&senderWorkersAmount, "workers", "w", 30, "workers amount")
	senderCmd.Flags().StringVarP(&senderTimeout, "timeout", "t", "", "timeout - string human times - empty means no timeout")
	senderCmd.Flags().StringVarP(&senderFileSize, "size", "s", "", "filesize - string human size")
	senderCmd.Flags().IntVarP(&senderFileAmount, "amount", "a", 10, "fileamount")
	senderCmd.Flags().StringVarP(&senderUrl, "url", "u", "", "url to send to")
	rootCmd.AddCommand(senderCmd)
}
