package renderer

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput перехватывает вывод в stdout для тестирования
// Возвращает строку со всем выведенным содержимым
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	err := w.Close()
	if err != nil {
		return ""
	}
	os.Stdout = old

	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		return ""
	}
	return buf.String()
}

// TestNewDrawTask проверяет создание нового объекта DrawTask
// Тест убеждается что все поля инициализированы корректно
func TestNewDrawTask(t *testing.T) {
	task := NewDrawTask()

	if task == nil {
		t.Fatal("NewDrawTask вернул nil")
	}

	// Проверяем что поля имеют нулевые значения
	if task.Width != 0 {
		t.Errorf("Начальная ширина должна быть 0, получено %d", task.Width)
	}

	if task.Height != 0 {
		t.Errorf("Начальная высота должна быть 0, получено %d", task.Height)
	}

	if task.Position.X != 0 || task.Position.Y != 0 {
		t.Errorf("Начальная позиция должна быть (0, 0), получено (%d, %d)",
			task.Position.X, task.Position.Y)
	}

	if task.Content != nil {
		t.Error("Начальный контент должен быть nil")
	}
}

// TestSetContent проверяет установку содержимого задачи отрисовки
// Тест проверяет что метод возвращает сам объект для цепочки вызовов (fluent API)
func TestSetContent(t *testing.T) {
	tests := []struct {
		name    string
		content []string
	}{
		{
			name:    "пустой массив",
			content: []string{},
		},
		{
			name:    "одна строка",
			content: []string{"hello"},
		},
		{
			name:    "несколько строк",
			content: []string{"line1", "line2", "line3"},
		},
		{
			name:    "строки с unicode",
			content: []string{"привет", "мир", "🎉"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := NewDrawTask()
			result := task.SetContent(tt.content)

			// Проверяем что метод возвращает сам объект
			if result != task {
				t.Error("SetContent должен возвращать сам объект для цепочки вызовов")
			}

			// Проверяем что содержимое установлено
			if len(task.Content) != len(tt.content) {
				t.Errorf("Ожидалось %d строк контента, получено %d",
					len(tt.content), len(task.Content))
			}

			// Проверяем каждую строку
			for i, line := range tt.content {
				if task.Content[i] != line {
					t.Errorf("Строка %d: ожидалось %q, получено %q",
						i, line, task.Content[i])
				}
			}
		})
	}
}

// TestSetSize проверяет установку размеров задачи отрисовки
// Если размеры меньше 1, должен автоматически вызваться SetAutoSize
func TestSetSize(t *testing.T) {
	tests := []struct {
		name          string
		width, height int
		content       []string
		expectAuto    bool
	}{
		{
			name:       "корректные размеры",
			width:      10,
			height:     5,
			content:    []string{"test"},
			expectAuto: false,
		},
		{
			name:       "нулевая ширина - автоматический размер",
			width:      0,
			height:     5,
			content:    []string{"test", "hello"},
			expectAuto: true,
		},
		{
			name:       "нулевая высота - автоматический размер",
			width:      10,
			height:     0,
			content:    []string{"test", "hello", "world"},
			expectAuto: true,
		},
		{
			name:       "оба нулевые - автоматический размер",
			width:      0,
			height:     0,
			content:    []string{"test"},
			expectAuto: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := NewDrawTask().SetContent(tt.content)
			result := task.SetSize(tt.width, tt.height)

			// Проверяем возврат self
			if result != task {
				t.Error("SetSize должен возвращать сам объект")
			}

			if tt.expectAuto {
				// При автоматическом размере проверяем что размеры установлены
				if task.Width == 0 && len(tt.content) > 0 && len(tt.content[0]) > 0 {
					t.Error("SetAutoSize должен установить ширину")
				}
				if task.Height == 0 && len(tt.content) > 0 {
					t.Error("SetAutoSize должен установить высоту")
				}
			} else {
				// При явном размере проверяем точность
				if task.Width != tt.width {
					t.Errorf("Ширина: ожидалось %d, получено %d", tt.width, task.Width)
				}
				if task.Height != tt.height {
					t.Errorf("Высота: ожидалось %d, получено %d", tt.height, task.Height)
				}
			}
		})
	}
}

// TestSetAutoSize проверяет автоматическое определение размеров
// Ширина определяется по самой длинной строке, высота - по количеству строк
func TestSetAutoSize(t *testing.T) {
	tests := []struct {
		name           string
		content        []string
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "одна короткая строка",
			content:        []string{"hi"},
			expectedWidth:  2,
			expectedHeight: 1,
		},
		{
			name:           "несколько строк разной длины",
			content:        []string{"short", "medium line", "x"},
			expectedWidth:  11, // "medium line"
			expectedHeight: 3,
		},
		{
			name:           "unicode символы",
			content:        []string{"привет", "мир"},
			expectedWidth:  6, // "привет"
			expectedHeight: 2,
		},
		{
			name:           "пустой массив",
			content:        []string{},
			expectedWidth:  0,
			expectedHeight: 0,
		},
		{
			name:           "строки с эмодзи",
			content:        []string{"hello 🎉", "test"},
			expectedWidth:  8, // "hello 🎉" - 8 рун
			expectedHeight: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := NewDrawTask().SetContent(tt.content)
			result := task.SetAutoSize()

			// Проверяем возврат self
			if result != task {
				t.Error("SetAutoSize должен возвращать сам объект")
			}

			// Проверяем ширину
			if task.Width != tt.expectedWidth {
				t.Errorf("Ширина: ожидалось %d, получено %d", tt.expectedWidth, task.Width)
			}

			// Проверяем высоту
			if task.Height != tt.expectedHeight {
				t.Errorf("Высота: ожидалось %d, получено %d", tt.expectedHeight, task.Height)
			}
		})
	}
}

// TestSetAutoSizePanic проверяет что SetAutoSize паникует при nil контенте
// Это защита от неправильного использования API
func TestSetAutoSizePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("SetAutoSize должен паниковать при nil Content")
		}
	}()

	task := NewDrawTask()
	task.SetAutoSize() // Должно спаниковать
}

// TestSetWidth проверяет установку только ширины
func TestSetWidth(t *testing.T) {
	task := NewDrawTask()
	result := task.SetWidth(42)

	if result != task {
		t.Error("SetWidth должен возвращать сам объект")
	}

	if task.Width != 42 {
		t.Errorf("Ширина: ожидалось 42, получено %d", task.Width)
	}
}

// TestSetHeight проверяет установку только высоты
func TestSetHeight(t *testing.T) {
	task := NewDrawTask()
	result := task.SetHeight(24)

	if result != task {
		t.Error("SetHeight должен возвращать сам объект")
	}

	if task.Height != 24 {
		t.Errorf("Высота: ожидалось 24, получено %d", task.Height)
	}
}

// TestSetPosition проверяет установку позиции отрисовки
func TestSetPosition(t *testing.T) {
	tests := []struct {
		name string
		x, y int
	}{
		{
			name: "начало координат",
			x:    0,
			y:    0,
		},
		{
			name: "положительные координаты",
			x:    10,
			y:    20,
		},
		{
			name: "отрицательные координаты",
			x:    -5,
			y:    -10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := NewDrawTask()
			result := task.SetPosition(tt.x, tt.y)

			if result != task {
				t.Error("SetPosition должен возвращать сам объект")
			}

			if task.Position.X != tt.x {
				t.Errorf("X: ожидалось %d, получено %d", tt.x, task.Position.X)
			}

			if task.Position.Y != tt.y {
				t.Errorf("Y: ожидалось %d, получено %d", tt.y, task.Position.Y)
			}
		})
	}
}

// TestSetColorSchema проверяет установку цветовой схемы
func TestSetColorSchema(t *testing.T) {
	tests := []struct {
		name   string
		fg, bg string
	}{
		{
			name: "основные цвета",
			fg:   "red",
			bg:   "white",
		},
		{
			name: "пустые значения",
			fg:   "",
			bg:   "",
		},
		{
			name: "только foreground",
			fg:   "blue",
			bg:   "",
		},
		{
			name: "только background",
			fg:   "",
			bg:   "yellow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := NewDrawTask()
			result := task.SetColorSchema(tt.fg, tt.bg)

			if result != task {
				t.Error("SetColorSchema должен возвращать сам объект")
			}

			if task.ColorSchema.FG != tt.fg {
				t.Errorf("FG: ожидалось %q, получено %q", tt.fg, task.ColorSchema.FG)
			}

			if task.ColorSchema.BG != tt.bg {
				t.Errorf("BG: ожидалось %q, получено %q", tt.bg, task.ColorSchema.BG)
			}
		})
	}
}

// TestFluentAPI проверяет что все методы можно вызывать цепочкой
// Это паттерн Builder для удобного конфигурирования
func TestFluentAPI(t *testing.T) {
	task := NewDrawTask().
		SetContent([]string{"line1", "line2"}).
		SetPosition(10, 20).
		SetWidth(50).
		SetHeight(10).
		SetColorSchema("green", "black")

	if task.Position.X != 10 || task.Position.Y != 20 {
		t.Error("SetPosition не сработал в цепочке")
	}

	if task.Width != 50 {
		t.Error("SetWidth не сработал в цепочке")
	}

	if task.Height != 10 {
		t.Error("SetHeight не сработал в цепочке")
	}

	if len(task.Content) != 2 {
		t.Error("SetContent не сработал в цепочке")
	}

	if task.ColorSchema.FG != "green" || task.ColorSchema.BG != "black" {
		t.Error("SetColorSchema не сработал в цепочке")
	}
}

// TestColorToANSI проверяет конвертацию названий цветов в ANSI коды
func TestColorToANSI(t *testing.T) {
	tests := []struct {
		name     string
		color    string
		expected string
	}{
		{name: "черный", color: "black", expected: "30"},
		{name: "красный", color: "red", expected: "31"},
		{name: "зеленый", color: "green", expected: "32"},
		{name: "желтый", color: "yellow", expected: "33"},
		{name: "синий", color: "blue", expected: "34"},
		{name: "пурпурный", color: "magenta", expected: "35"},
		{name: "голубой", color: "cyan", expected: "36"},
		{name: "белый", color: "white", expected: "37"},
		{name: "неизвестный цвет", color: "unknown", expected: "37"}, // Дефолт
		{name: "пустая строка", color: "", expected: "37"},           // Дефолт
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := colorToANSI(tt.color)
			if result != tt.expected {
				t.Errorf("colorToANSI(%q) = %q, ожидалось %q",
					tt.color, result, tt.expected)
			}
		})
	}
}

// TestColorToBgANSI проверяет конвертацию названий цветов фона в ANSI коды
func TestColorToBgANSI(t *testing.T) {
	tests := []struct {
		name     string
		color    string
		expected string
	}{
		{name: "черный фон", color: "black", expected: "40"},
		{name: "красный фон", color: "red", expected: "41"},
		{name: "зеленый фон", color: "green", expected: "42"},
		{name: "желтый фон", color: "yellow", expected: "43"},
		{name: "синий фон", color: "blue", expected: "44"},
		{name: "пурпурный фон", color: "magenta", expected: "45"},
		{name: "голубой фон", color: "cyan", expected: "46"},
		{name: "белый фон", color: "white", expected: "47"},
		{name: "неизвестный цвет фона", color: "unknown", expected: "40"}, // Дефолт
		{name: "пустая строка фона", color: "", expected: "40"},           // Дефолт
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := colorToBgANSI(tt.color)
			if result != tt.expected {
				t.Errorf("colorToBgANSI(%q) = %q, ожидалось %q",
					tt.color, result, tt.expected)
			}
		})
	}
}

// TestApplyColors проверяет применение цветовой схемы через ANSI коды
func TestApplyColors(t *testing.T) {
	tests := []struct {
		name   string
		schema ColorSchema
		checks []string // Подстроки которые должны быть в выводе
	}{
		{
			name:   "только foreground",
			schema: ColorSchema{FG: "red", BG: ""},
			checks: []string{"\033[31m"},
		},
		{
			name:   "только background",
			schema: ColorSchema{FG: "", BG: "blue"},
			checks: []string{"\033[44m"},
		},
		{
			name:   "оба цвета",
			schema: ColorSchema{FG: "green", BG: "yellow"},
			checks: []string{"\033[32m", "\033[43m"},
		},
		{
			name:   "без цветов",
			schema: ColorSchema{FG: "", BG: ""},
			checks: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				applyColors(tt.schema)
			})

			for _, check := range tt.checks {
				if !strings.Contains(output, check) {
					t.Errorf("Вывод не содержит ожидаемый код %q. Вывод: %q",
						check, output)
				}
			}

			if len(tt.checks) == 0 && output != "" {
				t.Errorf("При пустой схеме не должно быть вывода, получено: %q", output)
			}
		})
	}
}

// TestDraw проверяет основную функцию отрисовки
// Проверяем что вывод содержит правильные ANSI коды и координаты
func TestDraw(t *testing.T) {
	tests := []struct {
		name   string
		task   *DrawTask
		checks []string // Подстроки которые должны присутствовать
	}{
		{
			name: "простая отрисовка",
			task: NewDrawTask().
				SetContent([]string{"test"}).
				SetPosition(5, 10),
			checks: []string{
				"\033[?25l",   // HideCursor
				"\033[10;5H",  // MoveTo(5, 10)
				"X: 5, Y: 10", // Контент
				"\033[0m",     // Сброс стилей
			},
		},
		{
			name: "с цветами",
			task: NewDrawTask().
				SetContent([]string{"colored"}).
				SetPosition(1, 1).
				SetColorSchema("red", "white"),
			checks: []string{
				"\033[31m", // Красный текст
				"\033[47m", // Белый фон
				"\033[0m",  // Сброс
			},
		},
		{
			name: "несколько строк",
			task: NewDrawTask().
				SetContent([]string{"line1", "line2", "line3"}).
				SetPosition(10, 5),
			checks: []string{
				"X: 10, Y: 5", // Первая строка
				"X: 10, Y: 6", // Вторая строка
				"X: 10, Y: 7", // Третья строка
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				tt.task.Draw()
			})

			for _, check := range tt.checks {
				if !strings.Contains(output, check) {
					t.Errorf("Вывод не содержит ожидаемую подстроку %q.\nПолный вывод: %q",
						check, output)
				}
			}
		})
	}
}

// TestDrawWithEmptyContent проверяет отрисовку с пустым контентом
// Не должно быть паники, должны применяться только общие настройки
func TestDrawWithEmptyContent(t *testing.T) {
	task := NewDrawTask().
		SetContent([]string{}).
		SetPosition(1, 1)

	output := captureOutput(func() {
		task.Draw()
	})

	// Проверяем базовые команды
	if !strings.Contains(output, "\033[?25l") {
		t.Error("Должен быть вызван HideCursor")
	}

	if !strings.Contains(output, "\033[0m") {
		t.Error("Должен быть сброс стилей")
	}
}
