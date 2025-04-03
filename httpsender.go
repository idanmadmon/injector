package injector

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/docker/go-units"
)

const (
	timesBuffer = 20
)

type HttpSender struct {
	conf    SenderConfig
	url     string
	counter uint32
	file    string
	fileC   chan struct{}
	wg      sync.WaitGroup
}

func NewHttpSender(conf SenderConfig, url string) *HttpSender {
	return &HttpSender{
		conf: conf,
		url:  url,
	}
}

func (h *HttpSender) Start(ctx context.Context) {
	h.file = strings.Repeat("a", h.conf.FileSize)

	if h.conf.Timeout != 0 {
		var cancelCtx context.CancelFunc
		ctx, cancelCtx = context.WithTimeout(ctx, h.conf.Timeout+time.Millisecond*900)
		defer cancelCtx()
	}

	h.fileC = make(chan struct{}, timesBuffer*h.conf.FileAmount)
	defer close(h.fileC)

	for i := 0; i < h.conf.WorkersAmount; i++ {
		go h.sendWorker()
	}

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Gracefully shutdown - waiting for all files to finish")
			h.wg.Wait()
			fmt.Println("Sent total %d files of %s", atomic.LoadUint32(&h.counter), units.HumanSize(float64(h.conf.FileSize)))
			return
		case <-time.After(time.Second):
			for i := 0; i < h.conf.FileAmount; i++ {
				h.fileC <- struct{}{}
				h.wg.Add(1)
			}
			fmt.Printf("Sending %d files, channel status: %d / %d\n", h.conf.FileAmount, len(h.fileC), cap(h.fileC))
		}
	}
}

func (h *HttpSender) sendWorker() {
	for range h.fileC {
		h.sendHttpRequest()
	}
}

func (h *HttpSender) sendHttpRequest() {
	req, err := http.NewRequest("POST", h.url, strings.NewReader(h.file))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Got non 200 response: failed to read reaponse body, error", err)
		} else {
			fmt.Println("Got non 200 response:", string(body))
		}

		resp.Body.Close()
		return
	}

	resp.Body.Close()
	atomic.AddUint32(&h.counter, 1)
	h.wg.Done()
}
