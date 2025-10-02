package main

import (
	"Guess/internal/ui/terminal"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

// parseSGRMouse парсит SGR формат мыши: ESC[<btn;x;yM
func parseSGRMouse(data []byte) (button, x, y int, isPress bool) {
	// Убираем ESC[< и конечный символ
	s := string(data[3 : len(data)-1])
	parts := strings.Split(s, ";")
	if len(parts) != 3 {
		return 0, 0, 0, false
	}

	button, _ = strconv.Atoi(parts[0])
	x, _ = strconv.Atoi(parts[1])
	y, _ = strconv.Atoi(parts[2])

	// M = press, m = release
	isPress = data[len(data)-1] == 'M'

	return
}

func cleanup() {
	terminal.DisableRawMode()
	fmt.Print("\033[?1000l\033[?1006l") // Отключаем mouse tracking
	fmt.Print("\033[2J\033[H")          // Очищаем экран
}

func main() {
	// Перехват сигналов для корректного завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Включаем raw mode
	if err := terminal.EnableRawMode(); err != nil {
		fmt.Printf("Ошибка включения raw mode: %v\n", err)
		return
	}
	defer cleanup()

	// Включаем mouse tracking
	fmt.Print("\033[?1000h\033[?1006h")

	// Очищаем экран
	fmt.Print("\033[2J\033[H")
	fmt.Print("Кликните мышью. Нажмите 'q' для выхода.\r\n")

	// Горутина для обработки Ctrl+C
	go func() {
		<-sigChan
		cleanup()
		fmt.Print("Выход по Ctrl+C\n")
		os.Exit(0)
	}()

	// Основной цикл чтения событий
	for {
		data, err := terminal.ReadInput()
		if err != nil {
			fmt.Printf("\r\nОшибка чтения: %v\r\n", err)
			break
		}

		// Проверяем длину данных
		if len(data) == 0 {
			continue
		}

		// Выход по 'q' или Ctrl+C
		if data[0] == 3 || data[0] == 'q' {
			fmt.Print("\r\nВыход...\r\n")
			break
		}

		// Обработка SGR формата мыши: ESC[<B;X;YM
		if len(data) >= 9 && data[0] == 27 && data[1] == '[' && data[2] == '<' {
			button, x, y, isPress := parseSGRMouse(data)

			// Обрабатываем только нажатие левой кнопки
			if isPress && button == 0 {
				fmt.Printf("\r\nКЛИК: x=%d, y=%d\r\n", x, y)

				// Вызываем обработчик кликабельных областей
				clicked := terminal.HandleClick(x, y)
				if !clicked {
					fmt.Printf("Клик мимо всех областей\r\n")
				}
			}
		}
	}
}
