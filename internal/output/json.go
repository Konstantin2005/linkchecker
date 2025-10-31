package output

import (
	"encoding/json"
	"linkchecker/config"
	"os"
)

type FormatterJson interface {
	PrintResultJson(result config.Config)
	PrintSummaryJson(summary *config.Summary)
	PrintErrorJson(url, error string)
}

// JSONFormatter записывает всё строками JSON.
type JSONFormatter struct {
	enc *json.Encoder
}

func (j *JSONFormatter) PrintResultJson(r config.Config) {
	j.enc.Encode(map[string]any{
		"type": r,
	})
}

func (j *JSONFormatter) PrintSummaryJson(s *config.Summary) {
	j.enc.Encode(map[string]any{
		"total_links":   s.TotalLinks,
		"checked_links": s.CheckedLinks,
		"successful":    s.Successful,
		"errors":        s.Errors,
		"error_types":   s.ErrorByType,
	})
}

func (j *JSONFormatter) PrintErrorJson(url, err string) {
	j.enc.Encode(map[string]any{
		"type": "error",
		"url":  url,
		"err":  err,
	})

}

func NewJSONFormatter(file *os.File) *JSONFormatter {
	return &JSONFormatter{enc: json.NewEncoder(file)}
}
