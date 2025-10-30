package pkg

import (
	"fmt"
	"time"
)

func Loading(done chan bool) {

	// Запускаем спиннер в отдельной горутине
	go func() {
		spinner := []string{"|", "/", "-", "\\"}
		i := 0
		for {
			select {
			case <-done:
				fmt.Print("\r") // Очистка строки
				return
			default:
				fmt.Printf("\r%s идет обход...", spinner[i])
				i = (i + 1) % len(spinner)
				time.Sleep(100 * time.Millisecond)
			}
		}
		fmt.Print("\r")
	}()
}
