// Package config implements the configuration for the application.
// It is initialized by the init() function and can be accessed by GetConfig()
package config

import (
	"flag"
	"io"
	"time"
)

var instance *config

const (
	ServiceName = "goph-keeper"
)

type config struct {
	// Server is the server configuration
	ServerIP  string
	PollTimer time.Duration
	DumpTimer time.Duration
}

// NewConfig creates a new configuration struct
func init() {
	// Initialize the config struct
	instance = &config{}

	// Set the default values
	flag.StringVar(&instance.ServerIP, "addr", "localhost:8080", "Server IP address default: localhost:8080")
	flag.DurationVar(&instance.PollTimer, "poll", 5*time.Second, "Poll timer default: 5s")
	flag.DurationVar(&instance.DumpTimer, "dump", 10*time.Second, "Dump timer default: 10s")

	// Parse the flags and ignore the rest
	flag.CommandLine.SetOutput(io.Discard)
}

// GetConfig returns the configuration initialized by init func
func GetConfig() *config {
	return instance
}
