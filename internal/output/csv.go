package output

import (
	"encoding/csv"
	"os"
)

func csvFormat() {
	file, err := os.Create("output.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Запись заголовка
	writer.Write([]string{"Имя", "Возраст", "Город"})

	if err := writer.Error(); err != nil {
		panic(err)
	}
}
