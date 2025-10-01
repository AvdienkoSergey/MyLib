package terminal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput перехватывает stdout для тестирования
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

// TestNewCursorManager проверяет, что новый экземпляр CursorManager правильно инициализирован значениями по умолчанию.
func TestNewCursorManager(t *testing.T) {
	cm := NewCursorManager()

	if cm == nil {
		t.Fatal("NewCursorManager возвращает nil")
	}

	if cm.currentRow != 1 {
		t.Errorf("Ожидаем что currentRow будет 1, получено %d", cm.currentRow)
	}

	if cm.currentCol != 1 {
		t.Errorf("Ожидаем что currentCol будет 1, получено %d", cm.currentCol)
	}
}

// TestMoveTo проверяет метод moveTo, обеспечивая правильное перемещение курсора и обновляя его состояние.
func TestMoveTo(t *testing.T) {
	tests := []struct {
		name string
		col  int
		row  int
		want string
	}{
		{
			name: "Ожидаем перемещения в 1,1",
			col:  1,
			row:  1,
			want: "\033[1;1H",
		},
		{
			name: "Ожидаем перемещения в 10,5",
			col:  10,
			row:  5,
			want: "\033[5;10H",
		},
		{
			name: "Ожидаем перемещения в 100,50",
			col:  100,
			row:  50,
			want: "\033[50;100H",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewCursorManager()
			output := captureOutput(func() {
				cm.MoveTo(tt.col, tt.row)
			})

			if output != tt.want {
				t.Errorf("MoveTo(%d, %d) вернул = %q, хотели %q", tt.col, tt.row, output, tt.want)
			}

			if cm.currentRow != tt.row {
				t.Errorf("currentRow = %d, хотели %d", cm.currentRow, tt.row)
			}

			if cm.currentCol != tt.col {
				t.Errorf("currentCol = %d, хотели %d", cm.currentCol, tt.col)
			}
		})
	}
}

// TestClearLine проверяет, правильно ли функция ClearLine перемещает курсор в указанную строку и очищает строку.
func TestClearLine(t *testing.T) {
	cm := NewCursorManager()

	output := captureOutput(func() {
		cm.ClearLine(5)
	})

	expectedMoveTo := "\033[5;1H"
	expectedClear := "\033[K"
	expected := expectedMoveTo + expectedClear

	if output != expected {
		t.Errorf("ClearLine(5) возвращает = %q, хотели %q", output, expected)
	}
}

// TestClearArea проверяет, что функция ClearArea правильно очищает указанную область на выходе терминала.
// Проверяет выходные данные, захватывая стандартный вывод и сравнивая его с ожидаемыми escape-последовательностями
// и пробелами ANSI.
func TestClearArea(t *testing.T) {
	tests := []struct {
		name   string
		col    int
		row    int
		length int
	}{
		{
			name:   "очищает 5 символов начиная 10,5",
			col:    10,
			row:    5,
			length: 5,
		},
		{
			name:   "очищает 20 символов начиная 1,1",
			col:    1,
			row:    1,
			length: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewCursorManager()

			output := captureOutput(func() {
				cm.ClearArea(tt.col, tt.row, tt.length)
			})

			expectedMoveTo := fmt.Sprintf("\033[%d;%dH", tt.row, tt.col)
			expectedSpaces := strings.Repeat(" ", tt.length)
			expected := expectedMoveTo + expectedSpaces

			if output != expected {
				t.Errorf("ClearArea() возвращает = %q, хотели %q", output, expected)
			}
		})
	}
}

// TestClearRect проверяет метод clearRect, чтобы убедиться, что он правильно очищает прямоугольную область терминала
// с использованием пробелов
func TestClearRect(t *testing.T) {
	cm := NewCursorManager()

	output := captureOutput(func() {
		cm.ClearRect(5, 2, 10, 4)
	})

	// Проверяем, что очищаются строки 2, 3, 4
	// Каждая строка: MoveTo + 6 пробелов (10-5+1)
	if !strings.Contains(output, "\033[2;5H") {
		t.Error("ClearRect должен переместиться на строку 2")
	}
	if !strings.Contains(output, "\033[3;5H") {
		t.Error("ClearRect должен переместиться на строку 3")
	}
	if !strings.Contains(output, "\033[4;5H") {
		t.Error("ClearRect должен переместиться на строку 4")
	}

	// Проверяем количество пробелов (должно быть 6 на каждую строку)
	spaces := strings.Count(output, "      ")
	if spaces != 3 {
		t.Errorf("Ожидалось 3 группы по 6 пробелов, получено %d", spaces)
	}
}

// TestWriteAt проверяет метод WriteAt, что текст правильно написан в определенных местах терминала
func TestWriteAt(t *testing.T) {
	tests := []struct {
		name string
		col  int
		row  int
		text string
	}{
		{
			name: "запись hello в позиции 10,5",
			col:  10,
			row:  5,
			text: "hello",
		},
		{
			name: "запись пустой строки",
			col:  1,
			row:  1,
			text: "",
		},
		{
			name: "запись с unicode",
			col:  5,
			row:  3,
			text: "привет 🎉",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewCursorManager()

			output := captureOutput(func() {
				cm.WriteAt(tt.col, tt.row, tt.text)
			})

			expectedMoveTo := fmt.Sprintf("\033[%d;%dH", tt.row, tt.col)
			expected := expectedMoveTo + tt.text

			if output != expected {
				t.Errorf("WriteAt() вывод = %q, ожидалось %q", output, expected)
			}
		})
	}
}

// TestUpdateArea проверяет поведение UpdateArea, сравнивая его выходные данные с ожидаемыми перемещениями курсора и
// текстом
func TestUpdateArea(t *testing.T) {
	cm := NewCursorManager()

	output := captureOutput(func() {
		cm.UpdateArea(10, 5, 8, "new")
	})

	// Должно быть: MoveTo + 8 пробелов + MoveTo + "new"
	expectedMoveTo := "\033[5;10H"
	expectedClear := strings.Repeat(" ", 8)
	expectedText := "new"
	expected := expectedMoveTo + expectedClear + expectedMoveTo + expectedText

	if output != expected {
		t.Errorf("UpdateArea() выход = %q, хотели %q", output, expected)
	}
}

// TestHideCursor проверяет функциональность метода HideCursor, чтобы убедиться, что он выдает правильную
// escape-последовательность ANSI.
func TestHideCursor(t *testing.T) {
	cm := NewCursorManager()

	output := captureOutput(func() {
		cm.HideCursor()
	})

	expected := "\033[?25l"
	if output != expected {
		t.Errorf("HideCursor() выход = %q, хотели %q", output, expected)
	}
}

// TestShowCursor проверяет метод ShowCursor, чтобы убедиться, что он выдает ожидаемую escape-последовательность ANSI
// для отображения курсора.
func TestShowCursor(t *testing.T) {
	cm := NewCursorManager()

	output := captureOutput(func() {
		cm.ShowCursor()
	})

	expected := "\033[?25h"
	if output != expected {
		t.Errorf("ShowCursor() выход = %q, хотели %q", output, expected)
	}
}

// TestCursorManagerState проверяет начальное состояние CursorManager и тестирует обновление состояния после вызова
// moveTo.
func TestCursorManagerState(t *testing.T) {
	cm := NewCursorManager()

	// Начальное состояние
	if cm.currentRow != 1 || cm.currentCol != 1 {
		t.Error("Начальная позиция должна быть 1,1")
	}

	// После MoveTo состояние должно обновиться
	captureOutput(func() {
		cm.MoveTo(10, 20)
	})

	if cm.currentRow != 20 {
		t.Errorf("После MoveTo(10, 20) currentRow = %d, ожидалось 20", cm.currentRow)
	}
	if cm.currentCol != 10 {
		t.Errorf("После MoveTo(10, 20) currentCol = %d, ожидалось 10", cm.currentCol)
	}
}

// TestClearRectBoundaries проверяет поведение clearRect, тестируя его способность очищать прямоугольные области консоли.
func TestClearRectBoundaries(t *testing.T) {
	tests := []struct {
		name                               string
		startCol, startRow, endCol, endRow int
		expectedRows                       int
		expectedWidth                      int
	}{
		{
			name:          "одна строка",
			startCol:      1,
			startRow:      1,
			endCol:        5,
			endRow:        1,
			expectedRows:  1,
			expectedWidth: 5,
		},
		{
			name:          "3x3 область",
			startCol:      1,
			startRow:      1,
			endCol:        3,
			endRow:        3,
			expectedRows:  3,
			expectedWidth: 3,
		},
		{
			name:          "1x10 область",
			startCol:      5,
			startRow:      10,
			endCol:        5,
			endRow:        19,
			expectedRows:  10,
			expectedWidth: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := NewCursorManager()

			output := captureOutput(func() {
				cm.ClearRect(tt.startCol, tt.startRow, tt.endCol, tt.endRow)
			})

			// Проверяем количество строк по количеству escape-последовательностей MoveTo
			moveCount := strings.Count(output, "\033[")
			if moveCount != tt.expectedRows {
				t.Errorf("Ожидалось очистить %d строк, получено %d", tt.expectedRows, moveCount)
			}

			// Проверяем ширину по длине пробелов
			spaces := strings.Repeat(" ", tt.expectedWidth)
			spaceCount := strings.Count(output, spaces)
			if spaceCount != tt.expectedRows {
				t.Errorf("Ожидалось %d групп по %d пробелов, получено %d", tt.expectedRows, tt.expectedWidth, spaceCount)
			}
		})
	}
}
