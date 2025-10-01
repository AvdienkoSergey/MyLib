package terminal

import (
	"fmt"
	"strings"
)

// CursorManager управляет позиционированием курсора в терминале
type CursorManager struct {
	currentRow int
	currentCol int
}

// NewCursorManager создает новый менеджер курсора
func NewCursorManager() *CursorManager {
	return &CursorManager{
		currentRow: 1,
		currentCol: 1,
	}
}

// ===============================
// ОСНОВНЫЕ ОПЕРАЦИИ С КУРСОРОМ
// ===============================

// MoveTo перемещает курсор на указанную позицию (строка, столбец)
func (cm *CursorManager) MoveTo(col, row int) {
	fmt.Printf("\033[%d;%dH", row, col) // ANSI: row первый, col второй
	cm.currentRow = row
	cm.currentCol = col
}

// ===============================
// ОПЕРАЦИИ ОЧИСТКИ
// ===============================

// ClearLine очищает указанную строку
func (cm *CursorManager) ClearLine(row int) {
	cm.MoveTo(1, row)
	fmt.Print("\033[K") // Очистить от курсора до конца строки
}

// ClearArea очищает область от указанной позиции на заданную длину
func (cm *CursorManager) ClearArea(col, row, length int) {
	cm.MoveTo(col, row)
	fmt.Print(strings.Repeat(" ", length))
}

// ClearRect очищает прямоугольную область
func (cm *CursorManager) ClearRect(startCol, startRow, endCol, endRow int) {
	for row := startRow; row <= endRow; row++ {
		cm.MoveTo(startCol, row)
		width := endCol - startCol + 1
		fmt.Print(strings.Repeat(" ", width))
	}
}

// ===============================
// ОПЕРАЦИИ ВЫВОДА
// ===============================

// WriteAt записывает текст в указанной позиции
func (cm *CursorManager) WriteAt(col, row int, text string) {
	cm.MoveTo(col, row)
	fmt.Print(text)
}

// UpdateArea заменяет содержимое области новым текстом
func (cm *CursorManager) UpdateArea(col, row int, oldLength int, newText string) {

	// Очищаем старую область
	cm.MoveTo(col, row)
	fmt.Print(strings.Repeat(" ", oldLength))

	// Записываем новый текст
	cm.MoveTo(col, row)
	fmt.Print(newText)

}

// ===============================
// УТИЛИТЫ
// ===============================

// HideCursor скрывает курсор
func (cm *CursorManager) HideCursor() {
	fmt.Print("\033[?25l")
}

// ShowCursor показывает курсор
func (cm *CursorManager) ShowCursor() {
	fmt.Print("\033[?25h")
}
