package injector

import "time"

type SenderConfig struct {
	WorkersAmount int
	Timeout       time.Duration
	FileSize      int
	FileAmount    int
}

type ReceiverConfig struct {
	ReadSpeed    int64
	MaxBandwidth int64
	Verbose      bool
}
