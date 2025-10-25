package output

import (
	"encoding/json"
	"linkchecker/config"
	"os"
)

// JSONFormatter записывает всё строками JSON.
type JSONFormatter struct {
	enc *json.Encoder
}

func NewJSONFormatter(file *os.File) *JSONFormatter {
	return &JSONFormatter{enc: json.NewEncoder(file)}
}

func (j *JSONFormatter) PrintSummary(s *config.Summary) {
	j.enc.Encode(map[string]any{
		"total_links":   s.TotalLinks,
		"checked_links": s.CheckedLinks,
		"successful":    s.Successful,
		"errors":        s.Errors,
		"error_types":   s.ErrorByType,
	})
}

func (j *JSONFormatter) PrintResult(r *config.CheckResult) {
	j.enc.Encode(map[string]any{
		"type":   "result",
		"url":    r.URL,
		"from":   r.Referrer,
		"status": r.StatusCode,
		"depth":  r.Depth,
	})
}

func (j *JSONFormatter) PrintError(url, err string) {
	j.enc.Encode(map[string]any{
		"type": "error",
		"url":  url,
		"err":  err,
	})

}
