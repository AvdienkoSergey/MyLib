package terminal

import (
	"fmt"
	"strings"
	"time"
)

// CursorManager управляет позиционированием курсора в терминале
type CursorManager struct {
	currentRow int
	currentCol int
	saved      bool
	savedRow   int
	savedCol   int
}

// NewCursorManager создает новый менеджер курсора
func NewCursorManager() *CursorManager {
	return &CursorManager{
		currentRow: 1,
		currentCol: 1,
	}
}

// Position представляет позицию в терминале
type Position struct {
	Row int
	Col int
}

// ===============================
// ОСНОВНЫЕ ОПЕРАЦИИ С КУРСОРОМ
// ===============================

// MoveTo перемещает курсор на указанную позицию (строка, столбец)
func (cm *CursorManager) MoveTo(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
	cm.currentRow = row
	cm.currentCol = col
}

// SavePosition сохраняет текущую позицию курсора
func (cm *CursorManager) SavePosition() {
	fmt.Print("\033[s") // Сохранить позицию курсора
	cm.saved = true
	cm.savedRow = cm.currentRow
	cm.savedCol = cm.currentCol
}

// RestorePosition восстанавливает сохраненную позицию курсора
func (cm *CursorManager) RestorePosition() {
	if cm.saved {
		fmt.Print("\033[u") // Восстановить позицию курсора
		cm.currentRow = cm.savedRow
		cm.currentCol = cm.savedCol
	}
}

// GetPosition возвращает текущую позицию курсора
func (cm *CursorManager) GetPosition() Position {
	return Position{Row: cm.currentRow, Col: cm.currentCol}
}

// ===============================
// ОПЕРАЦИИ ОЧИСТКИ
// ===============================

// ClearLine очищает указанную строку
func (cm *CursorManager) ClearLine(row int) {
	cm.SavePosition()
	cm.MoveTo(row, 1)
	fmt.Print("\033[K") // Очистить от курсора до конца строки
	cm.RestorePosition()
}

// ClearArea очищает область от указанной позиции на заданную длину
func (cm *CursorManager) ClearArea(row, col, length int) {
	cm.SavePosition()
	cm.MoveTo(row, col)
	fmt.Print(strings.Repeat(" ", length))
	cm.RestorePosition()
}

// ClearRect очищает прямоугольную область
func (cm *CursorManager) ClearRect(startRow, startCol, endRow, endCol int) {
	cm.SavePosition()
	for row := startRow; row <= endRow; row++ {
		cm.MoveTo(row, startCol)
		width := endCol - startCol + 1
		fmt.Print(strings.Repeat(" ", width))
	}
	cm.RestorePosition()
}

// ===============================
// ОПЕРАЦИИ ВЫВОДА
// ===============================

// WriteAt записывает текст в указанной позиции
func (cm *CursorManager) WriteAt(row, col int, text string) {
	cm.SavePosition()
	cm.MoveTo(row, col)
	fmt.Print(text)
	cm.RestorePosition()
}

// WriteAtWithStyle записывает стилизованный текст в указанной позиции
func (cm *CursorManager) WriteAtWithStyle(row, col int, text string, style Style) {
	cm.SavePosition()
	cm.MoveTo(row, col)
	fmt.Print(style.Apply(text))
	cm.RestorePosition()
}

// UpdateArea заменяет содержимое области новым текстом
func (cm *CursorManager) UpdateArea(row, col int, oldLength int, newText string) {
	cm.SavePosition()

	// Очищаем старую область
	cm.MoveTo(row, col)
	fmt.Print(strings.Repeat(" ", oldLength))

	// Записываем новый текст
	cm.MoveTo(row, col)
	fmt.Print(newText)

	cm.RestorePosition()
}

// ===============================
// СТИЛИ И АНИМАЦИИ
// ===============================

// Style представляет стиль текста
type Style struct {
	FgColor   int  // Цвет текста (30-37)
	BgColor   int  // Цвет фона (40-47)
	Bold      bool // Жирный
	Italic    bool // Курсив
	Underline bool // Подчеркивание
	Blink     bool // Мигание
	Reverse   bool // Инверсия
}

// Apply применяет стиль к тексту
func (s Style) Apply(text string) string {
	var codes []string

	if s.Bold {
		codes = append(codes, "1")
	}
	if s.Italic {
		codes = append(codes, "3")
	}
	if s.Underline {
		codes = append(codes, "4")
	}
	if s.Blink {
		codes = append(codes, "5")
	}
	if s.Reverse {
		codes = append(codes, "7")
	}
	if s.FgColor != 0 {
		codes = append(codes, fmt.Sprintf("%d", s.FgColor))
	}
	if s.BgColor != 0 {
		codes = append(codes, fmt.Sprintf("%d", s.BgColor))
	}

	if len(codes) == 0 {
		return text
	}

	return fmt.Sprintf("\033[%sm%s\033[0m", strings.Join(codes, ";"), text)
}

// Предопределенные стили
var (
	StyleRed     = Style{FgColor: 31}
	StyleGreen   = Style{FgColor: 32}
	StyleYellow  = Style{FgColor: 33}
	StyleBlue    = Style{FgColor: 34}
	StyleMagenta = Style{FgColor: 35}
	StyleCyan    = Style{FgColor: 36}
	StyleWhite   = Style{FgColor: 37}

	StyleBold      = Style{Bold: true}
	StyleBlink     = Style{Blink: true}
	StyleUnderline = Style{Underline: true}

	StyleRedBold   = Style{FgColor: 31, Bold: true}
	StyleGreenBold = Style{FgColor: 32, Bold: true}
	StyleBlueBold  = Style{FgColor: 34, Bold: true}
)

// BlinkAt мигает текстом в указанной позиции
func (cm *CursorManager) BlinkAt(row, col int, text string, style Style, duration time.Duration, interval time.Duration) {
	style.Blink = true
	cm.WriteAtWithStyle(row, col, text, style)

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		timer := time.After(duration)
		visible := true

		for {
			select {
			case <-timer:
				// Убираем мигание в конце
				style.Blink = false
				cm.WriteAtWithStyle(row, col, text, style)
				return

			case <-ticker.C:
				if visible {
					cm.ClearArea(row, col, len(text))
				} else {
					cm.WriteAtWithStyle(row, col, text, style)
				}
				visible = !visible
			}
		}
	}()
}

// ===============================
// ВЫСОКОУРОВНЕВЫЕ ОПЕРАЦИИ
// ===============================

// ProgressBar отображает прогресс-бар в указанной позиции
type ProgressBar struct {
	Row      int
	Col      int
	Width    int
	Progress float64 // 0.0 - 1.0
	Style    Style
	cm       *CursorManager
}

// NewProgressBar создает новый прогресс-бар
func (cm *CursorManager) NewProgressBar(row, col, width int, style Style) *ProgressBar {
	return &ProgressBar{
		Row:   row,
		Col:   col,
		Width: width,
		Style: style,
		cm:    cm,
	}
}

// Update обновляет прогресс-бар
func (pb *ProgressBar) Update(progress float64) {
	pb.Progress = progress
	if progress > 1.0 {
		progress = 1.0
	}
	if progress < 0.0 {
		progress = 0.0
	}

	filled := int(float64(pb.Width) * progress)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", pb.Width-filled)

	pb.cm.WriteAtWithStyle(pb.Row, pb.Col, bar, pb.Style)
}

// Counter отображает счетчик в указанной позиции
type Counter struct {
	Row   int
	Col   int
	Value int
	Style Style
	cm    *CursorManager
}

// NewCounter создает новый счетчик
func (cm *CursorManager) NewCounter(row, col int, style Style) *Counter {
	return &Counter{
		Row:   row,
		Col:   col,
		Style: style,
		cm:    cm,
	}
}

// Update обновляет значение счетчика
func (c *Counter) Update(value int) {
	oldText := fmt.Sprintf("%d", c.Value)
	newText := fmt.Sprintf("%d", value)

	c.cm.UpdateArea(c.Row, c.Col, len(oldText), newText)
	c.Value = value
}

// Increment увеличивает счетчик на 1
func (c *Counter) Increment() {
	c.Update(c.Value + 1)
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
