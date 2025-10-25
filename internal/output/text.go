package output

import (
	"fmt"
	"linkchecker/config"
)

type Formatter interface {
	PrintResult(result *config.CheckResult)
	PrintSummary(summary *config.Summary)
	PrintError(url, error string)
}
type textFormatter struct{}

func (t textFormatter) PrintSummary(r *config.Summary) {

	fmt.Println("📊 Statistics:")
	fmt.Printf("Total links found: %d\n", r.TotalLinks)
	fmt.Printf("Successfully checked: %d\n", r.CheckedLinks)
	fmt.Printf("Successful (200): %d\n", r.ErrorByType[200])
	fmt.Printf("Errors: %d\n\n", r.Errors)

	fmt.Printf("⏱️  Duration: %.1fs\n\n", r.Duration.Seconds())

}

func (t textFormatter) PrintResult(r *config.CheckResult) {
	fmt.Println("=== Link Checker Results ===")
	fmt.Printf("Website: %s\n", r.URL)
	fmt.Printf("Depth: %d, Workers: %d, Timeout: %s\n\n", r.Depth, r.Workers, r.ResponseTime)
}

func (t textFormatter) PrintError(url, error string) {

}
func NewFormatter() Formatter {
	return &textFormatter{}
}

/*
func PrintError(url, error string) {
	if len(r.Problems) > 0 {
		fmt.Println("❌ Problematic Links:")

		// группировка по коду ответа
		/*group := map[int][]Problem{}
		for _, p := range r.Problems {
			group[p.Status] = append(group[p.Status], p)
		}

		for _, status := range []int{404, 500} { // порядок вывода
			if list, ok := group[status]; ok {
				switch status {
				case 404:
					fmt.Println("404 Not Found:")
				case 500:
					fmt.Println("500 Internal Server Error:")
				default:
					fmt.Printf("%d:\n", status)
				}
				for _, p := range list {
					fmt.Printf("  %s (from: %s, depth: %d)\n", p.URL, p.From, p.Depth)
				}
				fmt.Println()
			}
		}
	}

	// Дополнительные «content»-строки
	for _, p := range r.Problems {
		fmt.Printf("[URL \"%s\" has the following content: \"%s\"]\n", p.URL, p.Message)
		fmt.Printf("[URL \"%s,\" has the following content: \"%s\"]\n", p.From, p.Message)
	}
}
*/
