package config

import (
	"net/url"
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
	TotalLinks   int                    `json:"TotalLinks"`
	CheckedLinks int                    `json:"CheckedLinks"`
	Successful   int                    `json:"Successful"`
	Errors       int                    `json:"Errors"`
	Duration     time.Duration          `json:"Duration"`
	ErrorByType  map[int]int            `json:"ErrorByType"`
	ProblemLinks map[string]CheckResult `json:"ProblemLinks"`
}

type CheckResult struct {
	StatusCode   int           `json:"StatusCode"`
	Error        error         `json:"Error"`
	Workers      int           `json:"Workers"`
	Depth        int           `json:"Depth"`
	Referrer     *url.URL      `json:"Referrer"`
	ResponseTime time.Duration `json:"ResponseTime"`
}
