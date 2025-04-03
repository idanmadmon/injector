package injector

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	ratereader "github.com/idanmadmon/rate-limited-reader"
)

type Receiver struct {
	conf       ReceiverConfig
	address    string
	totalBytes int64
	totalFiles int64
	ctx        context.Context
	cancel     context.CancelFunc
	server     *http.Server
}

func NewHttpReceiver(conf ReceiverConfig, address string) *Receiver {
	return &Receiver{
		conf:    conf,
		address: address,
	}
}

func (r *Receiver) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	r.ctx = ctx
	r.cancel = cancel

	mux := http.NewServeMux()
	mux.HandleFunc("/", r.handlePost)

	r.server = &http.Server{
		Addr:    r.address,
		Handler: mux,
	}

	fmt.Printf("Server is running on %s\n", r.address)

	if r.conf.ReadSpeed > 0 {
		fmt.Printf("Receiving with rate limit of %d bytes per second\n", r.conf.ReadSpeed)
	} else {
		fmt.Println("Receiving with no rate limit")
	}

	go r.statsPrinter()
	return r.server.ListenAndServe()
}

func (r *Receiver) handlePost(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	var bytesRead int64
	var err error

	reader := req.Body
	if r.conf.ReadSpeed > 0 {
		reader = ratereader.NewRateLimitedReadCloser(req.Body, r.conf.ReadSpeed)
	}

	// bytesRead, err = io.Copy(io.Discard, reader) - no, for showing bytes read over time
	buffer := make([]byte, 32*1024)
	var n int
	var totalBytesRead int64
	for err == nil {
		n, err = reader.Read(buffer)
		bytesRead = int64(n)
		totalBytesRead += bytesRead
		atomic.AddInt64(&r.totalBytes, bytesRead)
		r.verboseLog(fmt.Sprintf("Got %d bytes", bytesRead))
	}

	r.verboseLog(fmt.Sprintf("Finish reading file with %d bytes \n", totalBytesRead))

	if err != nil && err != io.EOF {
		fmt.Printf("got error: %v\n", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	req.Body.Close()
	atomic.AddInt64(&r.totalFiles, 1)
	w.WriteHeader(http.StatusOK)
}

func (r *Receiver) readWithLimit(body io.Reader) (int64, error) {
	buffer := make([]byte, r.conf.ReadSpeed)
	var totalBytes int64
	finished := false

	for !finished {
		select {
		case <-r.ctx.Done():
			return totalBytes, nil
		default:
			start := time.Now()

			n, err := body.Read(buffer)
			if n > 0 {
				totalBytes += int64(n)

				elapsed := time.Since(start)
				sleepDuration := time.Second*time.Duration(n)/time.Duration(r.conf.ReadSpeed) - elapsed

				if sleepDuration > 0 {
					time.Sleep(sleepDuration)
				}
			}

			if err == io.EOF {
				finished = true
			}
			if err != nil {
				return totalBytes, err
			}
		}
	}

	return totalBytes, nil
}

func (r *Receiver) statsPrinter() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	var lastBytes, lastFiles int64

	for {
		select {
		case <-r.ctx.Done():
			return
		case <-ticker.C:
			currentBytes := atomic.LoadInt64(&r.totalBytes)
			currentFiles := atomic.LoadInt64(&r.totalFiles)

			bytesThisSecond := currentBytes - lastBytes
			filesThisSecond := currentFiles - lastFiles

			lastBytes = currentBytes
			lastFiles = currentFiles

			fmt.Printf("Bytes this second: %d, Files this second: %d\n", bytesThisSecond, filesThisSecond)
		}
	}
}

func (r *Receiver) Stop() {
	if r.server == nil || r.ctx == nil || r.cancel == nil {
		fmt.Printf("Server isn't running...")
	}

	r.cancel()
	if err := r.server.Close(); err != nil {
		fmt.Printf("Error stopping server: %v\n", err)
	}

	r.ctx = nil
	r.cancel = nil
	r.server = nil

	fmt.Printf("Server stopped. Total bytes read: %d, Total files received: %d\n", atomic.LoadInt64(&r.totalBytes), atomic.LoadInt64(&r.totalFiles))

	atomic.StoreInt64(&r.totalBytes, 0)
	atomic.StoreInt64(&r.totalFiles, 0)
}

func (r *Receiver) verboseLog(msg string) {
	if r.conf.Verbose {
		fmt.Println(msg)
	}
}
