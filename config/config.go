package config

import (
	"time"
)

type Config struct {
	URL           string
	MaxDepth      int
	Timeout       time.Duration
	Workers       int
	OutputFormat  string
	Verbose       bool
	SkipSSLVerify bool
}
type CheckResult struct {
	URL          string
	StatusCode   int
	Error        string
	Workers      int
	Depth        int
	Referrer     string
	ResponseTime time.Duration
}

type Summary struct {
	TotalLinks   int
	CheckedLinks int
	Successful   int
	Errors       int
	ErrorByType  map[int]int
	Duration     time.Duration
}
