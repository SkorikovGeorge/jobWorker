package server

import "time"

type ServerConfig struct {
	Address         string
	Port            int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	IdleTimeout     time.Duration
}

var cfg = ServerConfig{
	Address:         "0.0.0.0",
	Port:            8080,
	ReadTimeout:     time.Second * 5,
	WriteTimeout:    time.Second * 5,
	ShutdownTimeout: time.Second * 30,
	IdleTimeout:     time.Second * 60,
}
