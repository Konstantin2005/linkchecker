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

		fmt.Println("json")
	case "csv":

		fmt.Println("csv")
	case "text":
		f := NewFormatter()

		f.PrintResult(config)
		f.PrintSummary(Sum)
		f.PrintError(Sum.ProblemLinks)

	}

}
