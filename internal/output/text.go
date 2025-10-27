package output

import (
	"fmt"
	"linkchecker/config"
)

type Formatter interface {
	PrintResult(result config.Config)
	PrintSummary(summary *config.Summary)
	PrintError(url, error string)
}
type textFormatter struct{}

func (t textFormatter) PrintResult(r config.Config) {
	fmt.Println("=== Link Checker Results ===")
	fmt.Printf("Website: %s\n", r.URL)
	fmt.Printf("Depth: %d, Workers: %d, Timeout: %s\n\n", r.MaxDepth, r.Workers, r.Timeout)
}

func (t textFormatter) PrintSummary(r *config.Summary) {

	fmt.Println("📊 Statistics:")
	fmt.Printf("Total links found: %d\n", r.TotalLinks)
	fmt.Printf("Successfully checked: %d\n", r.CheckedLinks)
	fmt.Printf("Successful (200): %d\n", r.ErrorByType[200])
	fmt.Printf("Errors: %d\n\n", len(r.ErrorByType))

	fmt.Printf("⏱️  Duration: %.1fs\n\n", r.Duration.Seconds())

}

func (t textFormatter) PrintError(url, error string) {

}
func NewFormatter() Formatter {
	return &textFormatter{}
}
