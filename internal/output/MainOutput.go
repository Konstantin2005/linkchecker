package output

import (
	"fmt"
	"linkchecker/config"
)

func MainFormate(config config.Config, Sum *config.Summary) {

	str := config.OutputFormat
	f := NewFormatter()

	switch str {
	case "json":
		fmt.Println("json")
	case "csv":
		fmt.Println("csv")
	case "text":
		fmt.Println("text")

		f.PrintResult(config)
		f.PrintSummary(Sum)

	}

}
