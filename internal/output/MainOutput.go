package output

import (
	"fmt"
	"linkchecker/config"
	"os"
)

func MainFormate(config config.Config, Sum *config.Summary) {

	str := config.OutputFormat

	switch str {
	case "json":
		j := NewJSONFormatter(os.Stdout)
		j.PrintResultJson(config)
		j.PrintSummaryJson(Sum)
		j.PrintErrorJson(Sum.ProblemLinks)

		f := NewFormatter()

		f.PrintResult(config)
		f.PrintSummary(Sum)
		f.PrintError(Sum.ProblemLinks)

		//c := NewFormatterCsv()
		//
		//c.PrintResultCsv(config)
		//c.PrintSummaryCsv(Sum)

	case "csv":

		fmt.Println("csv")
	case "text":

	}

}
