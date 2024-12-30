package injector

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	streamduplicator "github.com/idanmadmon/stream-duplicator"
)

type StreamDuplicatorReceiverSender struct {
	receiverAddress string
	senderUrls      []string
	maxOffsetDiff   int
	totalFiles      int64
	ctx             context.Context
	cancel          context.CancelFunc
	server          *http.Server
}

func NewHttpStreamDuplicatorReceiverSender(receiverAddress string, senderUrls []string, maxOffsetDiff int) *StreamDuplicatorReceiverSender {
	return &StreamDuplicatorReceiverSender{
		receiverAddress: receiverAddress,
		senderUrls:      senderUrls,
		maxOffsetDiff:   maxOffsetDiff,
	}
}

func (sd *StreamDuplicatorReceiverSender) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	sd.ctx = ctx
	sd.cancel = cancel

	mux := http.NewServeMux()
	mux.HandleFunc("/", sd.handlePost)

	sd.server = &http.Server{
		Addr:    sd.receiverAddress,
		Handler: mux,
	}

	fmt.Printf("Server is running on %s\n", sd.receiverAddress)
	go sd.statsPrinter()
	return sd.server.ListenAndServe()
}

func (sd *StreamDuplicatorReceiverSender) handlePost(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	duplicator, readers := streamduplicator.NewStreamDuplicatorWithMaxOffsetDiffWithReadersAmount(req.Body, sd.maxOffsetDiff, len(sd.senderUrls))

	for i, url := range sd.senderUrls {
		go sd.sendHttpRequest(url, readers[i])
	}

	duplicator.WaitForReaders()

	req.Body.Close()
	w.WriteHeader(http.StatusOK)
}

func (sd *StreamDuplicatorReceiverSender) sendHttpRequest(url string, reader *streamduplicator.Reader) {
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Printf("Error creating request to %s : %v\n", url, err)
		return
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request to %s : %v\n", url, err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Got non 200 response from %s : failed to read reaponse body, error %v\n", url, err)
		} else {
			fmt.Printf("Got non 200 response from %s : %v\n", url, string(body))
		}

		resp.Body.Close()
		return
	}

	resp.Body.Close()
	atomic.AddInt64(&sd.totalFiles, 1)
}

func (sd *StreamDuplicatorReceiverSender) statsPrinter() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var lastFiles int64

	for {
		select {
		case <-sd.ctx.Done():
			return
		case <-ticker.C:
			currentFiles := atomic.LoadInt64(&sd.totalFiles)
			filesThisSecond := currentFiles - lastFiles
			lastFiles = currentFiles
			fmt.Printf("Files sent this second: %d\n", filesThisSecond)
		}
	}
}

func (sd *StreamDuplicatorReceiverSender) Stop() {
	if sd.server == nil || sd.ctx == nil || sd.cancel == nil {
		fmt.Printf("Server isn't running...")
	}

	sd.cancel()
	if err := sd.server.Close(); err != nil {
		fmt.Printf("Error stopping server: %v\n", err)
	}

	sd.ctx = nil
	sd.cancel = nil
	sd.server = nil

	fmt.Printf("Server stopped. Total files sent: %d\n", atomic.LoadInt64(&sd.totalFiles))
	atomic.StoreInt64(&sd.totalFiles, 0)
}
