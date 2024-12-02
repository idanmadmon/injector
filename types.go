package injector

import "time"

type SenderConfig struct {
	WorkersAmount int
	Timeout       time.Duration
	FileSize      int64
	FileAmount    int
}

type ReceiverConfig struct {
	InnerConfig interface{}
}
