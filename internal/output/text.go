package output

import (
	"fmt"
	"linkchecker/config"
	"linkchecker/pkg"
	"net/http"
	"sort"
)

type Formatter interface {
	PrintResult(result config.Config)
	PrintSummary(summary *config.Summary)
	PrintError(ProblemLinks map[string]config.CheckResult)
}
type textFormatter struct{}

func (t textFormatter) PrintResult(r config.Config) {
	fmt.Println("=== Link Checker Results ===")
	fmt.Printf("Website: %s\n", r.URL)
	fmt.Printf("Depth: %d, Workers: %d, Timeout: %s\n\n", r.MaxDepth, r.Workers, r.Timeout)
}

func (t textFormatter) PrintSummary(r *config.Summary) {

	fmt.Println("📊 Statistics:")
	fmt.Printf("Total links found: %d\n", r.CheckedLinks)
	fmt.Printf("Successfully checked: %d\n", r.CheckedLinks)
	fmt.Printf("Successful (200): %d\n", r.ErrorByType[200])
	fmt.Printf("Errors: %d\n\n", pkg.SumError(r.ErrorByType))

	fmt.Printf("⏱️  Duration: %.1fs\n\n", r.Duration.Seconds())

}

func (t textFormatter) PrintError(problemLinks map[string]config.CheckResult) {
	if len(problemLinks) == 0 {
		fmt.Println("No problematic links found.")
		return
	}

	// Группируем по статус-коду
	groups := make(map[int][]struct {
		url      string
		referrer string
		depth    int
	})
	for linkURL, res := range problemLinks {
		if res.StatusCode >= 400 { // Только ошибки (4xx и 5xx)
			group := groups[res.StatusCode]
			referrerStr := ""
			if res.Referrer != nil {
				referrerStr = res.Referrer.String()
			}
			group = append(group, struct {
				url      string
				referrer string
				depth    int
			}{linkURL, referrerStr, res.Depth})
			groups[res.StatusCode] = group
		}
	}

	if len(groups) == 0 {
		fmt.Println("No problematic links with error status codes.")
		return
	}

	// Сортируем статус-коды по возрастанию
	var statusCodes []int
	for code := range groups {
		statusCodes = append(statusCodes, code)
	}
	sort.Ints(statusCodes)

	fmt.Println("❌ Problematic Links:")

	for _, code := range statusCodes {
		// Получаем стандартное описание статуса (например, "Not Found" для 404)
		description := http.StatusText(code)
		if description == "" {
			description = "Unknown Error" // На случай нестандартных кодов
		}
		fmt.Printf("%d %s:\n", code, description)

		links := groups[code]
		// Сортируем ссылки по URL для предсказуемости (опционально)
		sort.Slice(links, func(i, j int) bool {
			return links[i].url < links[j].url
		})

		for _, link := range links {
			fmt.Printf("  %s (from: %s, depth: %d)\n", link.url, link.referrer, link.depth)
		}
		fmt.Println() // Пустая строка между группами для читаемости
	}

}

func NewFormatter() Formatter {
	return &textFormatter{}
}
