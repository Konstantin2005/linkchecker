package config

import (
	"net/url"
	"time"
)

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
	TotalLinks   int                    `json:"TotalLinks"`
	CheckedLinks int                    `json:"CheckedLinks"`
	Successful   int                    `json:"Successful"`
	Errors       int                    `json:"Errors"`
	Duration     time.Duration          `json:"Duration"`
	ErrorByType  map[int]int            `json:"ErrorByType"`
	ProblemLinks map[string]CheckResult `json:"ProblemLinks"`
}

type CheckResult struct {
	URL          string        `json:"URL"`
	StatusCode   int           `json:"StatusCode"`
	Error        error         `json:"Error"`
	Referrer     *url.URL      `json:"Referrer"`
	Depth        int           `json:"Depth"`
	ResponseTime time.Duration `json:"ResponseTime"`
}
