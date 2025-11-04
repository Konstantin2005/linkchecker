package output

import (
	"encoding/json"
	"linkchecker/config"
	"linkchecker/pkg"
	"os"
	"time"
)

type FormatterJson interface {
	PrintResultJson(result config.Config)
	PrintSummaryJson(summary *config.Summary)
	PrintErrorJson(str map[string]config.CheckResult, duration time.Duration)
}

// JSONFormatter записывает всё строками JSON.
type JSONFormatter struct {
	enc *json.Encoder
}

func (j *JSONFormatter) PrintResultJson(r config.Config) {
	j.enc.Encode(map[string]any{
		"url": r.URL,
		"config": map[string]any{
			"depth":   r.MaxDepth,
			"workers": r.Workers,
			"timeout": r.Timeout,
		},
	})
}

func (j *JSONFormatter) PrintSummaryJson(s *config.Summary) {
	j.enc.Encode(map[string]any{
		"statistics": map[string]any{
			"total_links":   s.TotalLinks,
			"checked_links": s.CheckedLinks,
			"successful":    s.Successful,
			"errors":        pkg.SumError(s.ErrorByType),
			"error_types":   s.ErrorByType,
		},
	})
}

func (j *JSONFormatter) PrintErrorJson(problems map[string]config.CheckResult, dur time.Duration) {

	// объект «ссылка + результат»
	type linkWithResult struct {
		URL string `json:"url"`
		config.CheckResult
	}

	// корневой объект
	out := struct {
		ProblemLinks []linkWithResult `json:"problem_links"`
		Duration     string           `json:"duration"`
	}{
		Duration: dur.String(),
	}

	// заполняем слайс
	out.ProblemLinks = make([]linkWithResult, 0, len(problems))
	for url, res := range problems {
		out.ProblemLinks = append(out.ProblemLinks, linkWithResult{
			URL:         url,
			CheckResult: res,
		})
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	j.enc.Encode(out)
}

func NewJSONFormatter(file *os.File) *JSONFormatter {
	return &JSONFormatter{enc: json.NewEncoder(file)}
}
