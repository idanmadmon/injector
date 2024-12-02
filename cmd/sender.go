package main

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
	workersAmount int
	timeout       string
	fileAmount    int
	fileSize      string
	url           string
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

		t, err := time.ParseDuration(timeout)
		if timeout != "" && err != nil {
			fmt.Printf("timeout given isn't valid: %s, error: %v\n", fileSize, err)
			os.Exit(1)
		}

		conf := injector.SenderConfig{
			WorkersAmount: workersAmount,
			Timeout:       t,
			FileSize:      size,
			FileAmount:    fileAmount,
		}

		ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancelCtx()

		sender := injector.NewHttpSender(conf, url)
		sender.Start(ctx)
	},
}

func init() {
	senderCmd.Flags().IntVarP(&workersAmount, "workers", "w", 30, "workers amount")
	senderCmd.Flags().StringVarP(&timeout, "timeout", "t", "", "timeout - string human times - empty means no timeout")
	senderCmd.Flags().StringVarP(&fileSize, "size", "s", "", "filesize - string human size")
	senderCmd.Flags().IntVarP(&fileAmount, "amount", "a", 10, "fileamount")
	senderCmd.Flags().StringVarP(&url, "url", "u", "", "url to send to")
	rootCmd.AddCommand(senderCmd)
}
