package server

import "time"

// Config holds only what we actually need
type Config struct {
	Port        string
	StopTimeout time.Duration
}

// Default returns sensible defaults
func Default() Config {
	return Config{
		Port:        "8080",
		StopTimeout: 10 * time.Second,
	}
}
