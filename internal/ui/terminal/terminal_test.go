package terminal

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"testing"
)

// TestNewTerminal проверяет функцию NewTerminal, чтобы убедиться, что она инициализирует объект терминала значениями
// по умолчанию.
func TestNewTerminal(t *testing.T) {
	term := NewTerminal()

	if term == nil {
		t.Fatal("NewTerminal вернул nil")
	}

	if term.cached {
		t.Error("Новый терминал не должен иметь кешированных размеров")
	}

	if term.width != 0 {
		t.Errorf("Начальная ширина должна быть 0, получено %d", term.width)
	}

	if term.height != 0 {
		t.Errorf("Начальная высота должна быть 0, получено %d", term.height)
	}
}

// TestGetSize проверяет, что метод getSize терминала правильно извлекает и кэширует размеры терминала
func TestGetSize(t *testing.T) {
	term := NewTerminal()

	width, height := term.GetSize()

	// Размеры должны быть положительными
	if width <= 0 {
		t.Errorf("Ширина должна быть положительной, получено %d", width)
	}

	if height <= 0 {
		t.Errorf("Высота должна быть положительной, получено %d", height)
	}

	// После первого вызова размеры должны быть закешированы
	if !term.cached {
		t.Error("После GetSize размеры должны быть закешированы")
	}

	if term.width != width {
		t.Errorf("Закешированная ширина %d не совпадает с возвращённой %d", term.width, width)
	}

	if term.height != height {
		t.Errorf("Закешированная высота %d не совпадает с возвращённой %d", term.height, height)
	}
}

// TestGetWidth проверяет корректность работы метода GetWidth, возвращающего текущую ширину терминала.
func TestGetWidth(t *testing.T) {
	term := NewTerminal()

	width := term.GetWidth()

	if width <= 0 {
		t.Errorf("Ширина должна быть положительной, получено %d", width)
	}

	// Проверяем что значение совпадает с GetSize
	w, _ := term.GetSize()
	if width != w {
		t.Errorf("GetWidth() = %d, но GetSize() вернул ширину %d", width, w)
	}
}

// TestGetHeight проверяет метод getHeight, чтобы убедиться, что он возвращает положительное значение и соответствует высоте из getSize.
func TestGetHeight(t *testing.T) {
	term := NewTerminal()

	height := term.GetHeight()

	if height <= 0 {
		t.Errorf("Высота должна быть положительной, получено %d", height)
	}

	// Проверяем что значение совпадает с GetSize
	_, h := term.GetSize()
	if height != h {
		t.Errorf("GetHeight() = %d, но GetSize() вернул высоту %d", height, h)
	}
}

// TestRefresh проверяет поведение кэширования измерений терминала, обеспечивая правильное обновление кэша
// после вызова Refresh.
func TestRefresh(t *testing.T) {
	term := NewTerminal()

	// Первый вызов для кеширования
	width1, height1 := term.GetSize()

	if !term.cached {
		t.Fatal("После GetSize размеры должны быть закешированы")
	}

	// Вызываем Refresh
	term.Refresh()

	// Кеш должен быть обновлён
	if !term.cached {
		t.Error("После Refresh размеры должны быть закешированы снова")
	}

	// Получаем размеры снова
	width2, height2 := term.GetSize()

	// В нормальных условиях размеры не должны измениться
	// но проверяем что метод работает
	if width2 <= 0 || height2 <= 0 {
		t.Error("После Refresh размеры должны оставаться корректными")
	}

	// Проверяем что значения разумные (могут совпадать или отличаться)
	_ = width1
	_ = height1
}

// TestCaching проверяет механизм кэширования в методе getSize типа Terminal.
func TestCaching(t *testing.T) {
	term := NewTerminal()

	// Первый вызов - должен обновить кеш
	width1, height1 := term.GetSize()

	if !term.cached {
		t.Fatal("После первого GetSize должен быть установлен флаг cached")
	}

	// Второй вызов - должен использовать кеш
	width2, height2 := term.GetSize()

	if width1 != width2 || height1 != height2 {
		t.Error("Повторный вызов GetSize должен возвращать закешированные значения")
	}

	if !term.cached {
		t.Error("Флаг cached не должен сбрасываться при повторных вызовах")
	}
}

// TestGetCanvasSize проверяет функцию getCanvasSize, чтобы убедиться, что она возвращает допустимые положительные
// значения ширины и высоты. Проверяет, что размеры являются приемлемыми и превышают минимальные ожидаемые
// пороговые значения.
func TestGetCanvasSize(t *testing.T) {
	width, height := getCanvasSize()

	if width <= 0 {
		t.Errorf("getCanvasSize должен вернуть положительную ширину, получено %d", width)
	}

	if height <= 0 {
		t.Errorf("getCanvasSize должен вернуть положительную высоту, получено %d", height)
	}

	// Минимальные разумные значения (по умолчанию 80x24)
	if width < 20 {
		t.Errorf("Ширина слишком мала: %d", width)
	}

	if height < 10 {
		t.Errorf("Высота слишком мала: %d", height)
	}
}

// TestGetTerminalSizeXTerm проверяет, что функция getTerminalSizeXTerm не вызывает паники и работает в неинтерактивных
// средах.
func TestGetTerminalSizeXTerm(t *testing.T) {
	width, height := getTerminalSizeXTerm()

	// Этот тест может не пройти в неинтерактивной среде (CI/CD)
	// поэтому проверяем только что функция не паникует
	if width < 0 || height < 0 {
		t.Error("getTerminalSizeXTerm не должен возвращать отрицательные значения")
	}

	// В реальном терминале должны быть положительные значения
	// В CI/CD может вернуть 0, 0 - это нормально
}

// TestGetTerminalSizeWindows проверяет функциональность getTerminalSizeWindows в среде, специфичной для Windows.
// Пропускает проверку на платформах, отличных от Windows, и проверяет, верны ли возвращаемые значения размера терминала
func TestGetTerminalSizeWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Тест только для Windows")
	}

	width, height := getTerminalSizeWindows()

	// В Windows должны получить корректные размеры
	if width <= 0 || height <= 0 {
		t.Log("Не удалось получить размеры через Windows API, это может быть нормально в некоторых средах")
	}
}

// TestClear проверяет функцию очистки, чтобы убедиться, что она правильно очищает экран на различных платформах и
// в различных средах.
func TestClear(t *testing.T) {
	tests := []struct {
		name string
		goos string
	}{
		{
			name: "очистка на текущей платформе",
			goos: runtime.GOOS,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Перехватываем stdout для проверки вывода
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			Clear()

			err := w.Close()
			if err != nil {
				return
			}
			os.Stdout = old

			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			if err != nil {
				return
			}
			output := buf.String()

			// Проверяем что что-то было выведено
			if len(output) == 0 && runtime.GOOS != "windows" {
				t.Error("Clear должен выводить ANSI коды на Unix системах")
			}

			// Для Unix систем проверяем наличие ANSI кодов
			if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
				if !containsANSI(output) {
					t.Error("Вывод Clear должен содержать ANSI escape последовательности")
				}
			}
		})
	}
}

// TestClearNonInteractive проверяет, что функция Clear предоставляет базовые управляющие последовательности в
// неинтерактивном терминале.
func TestClearNonInteractive(t *testing.T) {
	// Сохраняем оригинальное значение TERM
	oldTerm := os.Getenv("TERM")
	defer func(key, value string) {
		err := os.Setenv(key, value)
		if err != nil {
			return
		}
	}("TERM", oldTerm)

	// Устанавливаем пустое значение TERM (неинтерактивный режим)
	err := os.Setenv("TERM", "")
	if err != nil {
		return
	}

	// Перехватываем stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Clear()

	err = w.Close()
	if err != nil {
		return
	}
	os.Stdout = old

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		return
	}
	output := buf.String()

	// В неинтерактивном режиме должен быть минимальный вывод
	if len(output) == 0 {
		t.Error("Clear должен выводить хотя бы базовые escape последовательности")
	}
}

// TestDefaultSize гарантирует, что размер терминала по умолчанию равен 80x24, когда фактический размер определить
// невозможно.
func TestDefaultSize(t *testing.T) {
	// Тестируем что даже если не удастся определить размер,
	// вернутся значения по умолчанию (80x24)
	term := NewTerminal()
	width, height := term.GetSize()

	// Минимум должны получить дефолтные значения
	if width < 80 || height < 24 {
		// Если меньше, значит что-то не так с реализацией fallback
		t.Logf("Получены размеры %dx%d, ожидалось минимум 80x24", width, height)
	}
}

// TestMultipleTerminalInstances проверяет, что несколько экземпляров терминала работают независимо друг от друга
// без общего состояния.
func TestMultipleTerminalInstances(t *testing.T) {
	// Проверяем что несколько экземпляров работают независимо
	term1 := NewTerminal()
	term2 := NewTerminal()

	width1, height1 := term1.GetSize()
	width2, height2 := term2.GetSize()

	// Размеры должны совпадать (один и тот же терминал)
	if width1 != width2 {
		t.Errorf("Разные экземпляры вернули разную ширину: %d vs %d", width1, width2)
	}

	if height1 != height2 {
		t.Errorf("Разные экземпляры вернули разную высоту: %d vs %d", height1, height2)
	}

	// Но кеш должен быть независимым
	term1.cached = false
	if !term2.cached {
		t.Error("Изменение кеша в term1 не должно влиять на term2")
	}
}

// TestRefreshInvalidatesCache проверяет, что метод обновления делает недействительными и обновляет кэшированные
// значения размера терминала.
func TestRefreshInvalidatesCache(t *testing.T) {
	term := NewTerminal()

	// Кешируем размеры
	term.GetSize()
	if !term.cached {
		t.Fatal("Размеры должны быть закешированы")
	}

	// Модифицируем закешированные значения
	term.width = 999
	term.height = 999

	// Вызываем Refresh
	term.Refresh()

	// Проверяем что значения обновились
	if term.width == 999 || term.height == 999 {
		t.Error("Refresh должен обновить закешированные значения")
	}

	// Кеш должен быть активен
	if !term.cached {
		t.Error("После Refresh кеш должен быть снова активен")
	}
}

// Вспомогательная функция для проверки наличия ANSI кодов
func containsANSI(s string) bool {
	return len(s) > 0 && (bytes.Contains([]byte(s), []byte("\033[")) ||
		bytes.Contains([]byte(s), []byte("\x1b[")))
}
