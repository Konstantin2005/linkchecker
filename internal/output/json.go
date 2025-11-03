package output

import (
	"encoding/json"
	"linkchecker/config"
	"linkchecker/pkg"
	"os"
)

type FormatterJson interface {
	PrintResultJson(result config.Config)
	PrintSummaryJson(summary *config.Summary)
	PrintErrorJson(str map[string]config.CheckResult)
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

func (j *JSONFormatter) PrintErrorJson(str map[string]config.CheckResult) {
	data := map[string]any{
		"problem_links": []any{
			map[string]any{
				"url": str,
			},
		},
	}
	json.NewEncoder(os.Stdout).Encode(data)
}

func NewJSONFormatter(file *os.File) *JSONFormatter {
	return &JSONFormatter{enc: json.NewEncoder(file)}
}
