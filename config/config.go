package config

import (
	"time"
)

type Problem struct {
	url     string
	fromURL string
}

type Config struct {
	URL           string        `json:"URL"`
	MaxDepth      int           `json:"MaxDepth"`
	Timeout       time.Duration `json:"Timeout"`
	Workers       int           `json:"Workers"`
	OutputFormat  string        `json:"OutputFormat"`
	Verbose       bool          `json:"Verbose"`
	SkipSSLVerify bool          `json:"SkipSSLVerify"`
}
type Summary struct {
	TotalLinks   int           `json:"TotalLinks"`
	CheckedLinks int           `json:"CheckedLinks"`
	Successful   int           `json:"Successful"`
	Errors       int           `json:"Errors"`
	ErrorByType  map[int]int   `json:"ErrorByType"`
	Duration     time.Duration `json:"Duration"`
}

type CheckResult struct {
	URL          string        `json:"URL"`
	StatusCode   int           `json:"StatusCode"`
	Error        string        `json:"Error"`
	Workers      int           `json:"Workers"`
	Depth        int           `json:"Depth"`
	Referrer     string        `json:"Referrer"`
	ResponseTime time.Duration `json:"ResponseTime"`
}
