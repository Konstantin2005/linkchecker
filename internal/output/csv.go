package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"linkchecker/config"
	"strconv"
	"sync"
)

type FormatterCsv interface {
	PrintResultCsv(result config.Config)
	PrintSummaryCsv(summary *config.Summary)
	PrintErrorCsv(url, error string)
}

type CSVFormatter struct {
	w             io.Writer  // куда писать
	mu            sync.Mutex // потокобезопасность
	errHeaderOnce sync.Once  // чтобы заголовок ошибок вывести один раз
}

func (f *CSVFormatter) PrintResultCsv(cfg config.Config) {
	f.mu.Lock()
	defer f.mu.Unlock()

	writer := csv.NewWriter(f.w)
	defer writer.Flush()

	_ = writer.Write([]string{"section", "key", "value"})
	_ = writer.Write([]string{"config", "website", cfg.URL})
	_ = writer.Write([]string{"config", "maxDepth", strconv.Itoa(cfg.MaxDepth)})
	_ = writer.Write([]string{"config", "workers", strconv.Itoa(cfg.Workers)})
	_ = writer.Write([]string{"config", "timeout", cfg.Timeout.String()})
	_ = writer.Write([]string{"config", "verbose", strconv.FormatBool(cfg.Verbose)})
	_ = writer.Write([]string{"config", "skipSSLVerify", strconv.FormatBool(cfg.SkipSSLVerify)})

}

func (f *CSVFormatter) PrintSummaryCsv(s *config.Summary) {
	f.mu.Lock()
	defer f.mu.Unlock()

	writer := csv.NewWriter(f.w)
	defer writer.Flush()

	_ = writer.Write([]string{"section", "metric", "value"})
	_ = writer.Write([]string{"summary", "totalLinks", strconv.Itoa(s.TotalLinks)})
	_ = writer.Write([]string{"summary", "checkedLinks", strconv.Itoa(s.CheckedLinks)})
	_ = writer.Write([]string{"summary", "successful200", strconv.Itoa(s.Successful)})
	_ = writer.Write([]string{"summary", "errors", strconv.Itoa(s.Errors)})
	_ = writer.Write([]string{"summary", "duration", s.Duration.String()})

	// детализируем ошибки по статус-коду
	for code, cnt := range s.ErrorByType {
		_ = writer.Write([]string{"summary", fmt.Sprintf("http%d", code), strconv.Itoa(cnt)})
	}

	// выводим сами проблемные ссылки
	if len(s.ProblemLinks) > 0 {
		_ = writer.Write([]string{}) // пустая строка-разделитель
		_ = writer.Write([]string{"section", "url", "statusCode", "referrer", "depth", "responseTime", "error"})
		for _, res := range s.ProblemLinks {
			ref := ""
			if res.Referrer != nil {
				ref = res.Referrer.String()
			}
			errStr := ""
			if res.Error != nil {
				errStr = res.Error.Error()
			}
			_ = writer.Write([]string{
				"problemLink",
				res.URL,
				strconv.Itoa(res.StatusCode),
				ref,
				strconv.Itoa(res.Depth),
				res.ResponseTime.String(),
				errStr,
			})
		}
	}

}

func (f *CSVFormatter) PrintErrorCsv(url, errMsg string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	writer := csv.NewWriter(f.w)
	defer writer.Flush()

	// Заголовок пишем только при первом вызове
	f.errHeaderOnce.Do(func() {
		_ = writer.Write([]string{"section", "url", "error"})
	})
	_ = writer.Write([]string{"error", url, errMsg})

}

func NewFormatterCsv() CSVFormatter {
	return CSVFormatter{}
}
