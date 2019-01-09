package main

import (
	"time"
)

// Options represents the command line options of the operator
type Options struct {
	ConfigFile     string
	KubeconfigFile string

	ResyncPeriod  time.Duration
	NodeThreshold int
}
