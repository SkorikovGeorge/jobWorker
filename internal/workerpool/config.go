package workerpool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

var Workers = NewWorkingPool()

type WorkersConfig struct {
	WorkersCount int
	RedisAddress string
	JobTimeout   time.Duration
	RetryDelay   time.Duration
	MaxRetries   int
	PauseDelay   time.Duration
	EmptyDelay   time.Duration
}

var cfg = WorkersConfig{
	WorkersCount: 10,
	RedisAddress: "127.0.0.1:6379",
	JobTimeout:   time.Second * 3,
	RetryDelay:   time.Second * 1,
	PauseDelay:   time.Second * 1,
	EmptyDelay:   time.Second * 1,
	MaxRetries:   2,
}

type WorkerPool struct {
	Config      *WorkersConfig
	ShutdownCtx context.Context
	CancelCtx   context.CancelFunc
	Paused      atomic.Bool
	wg          sync.WaitGroup
}

func NewWorkingPool() *WorkerPool {
	shutdownCtx, cancel := context.WithCancel(context.Background())
	wp := WorkerPool{
		Config:      &cfg,
		ShutdownCtx: shutdownCtx,
		CancelCtx:   cancel,
	}
	wp.SetupWorkers()
	return &wp
}
