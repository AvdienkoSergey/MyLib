package renderer

import (
	"fmt"
	"strings"
)

// LayoutDirection направление размещения детей
type LayoutDirection int8

const (
	Column LayoutDirection = iota // Вертикально (по умолчанию)
	Row                           // Горизонтально
)

type HasBorder int8

const (
	NoBorder HasBorder = iota
	Border
)

type IsRounded int8

const (
	NoRound IsRounded = iota
	Round
)

type Gap struct {
	Vertical   int8
	Horizontal int8
}

// DOMNode представляет узел DOM дерева
type DOMNode struct {
	Content   string
	Children  []*DOMNode
	X         int16
	Y         int16
	Direction LayoutDirection // Как размещать детей: row или column
	Width     int16           // 0 = auto (вычислить из content/children)
	Height    int16           // 0 = auto (вычислить из content/children)
	HasBorder HasBorder
	IsRounded IsRounded
	Padding   Gap
	Margin    Gap
}

// LayoutConfig конфигурация для позиционирования (глобальные настройки по умолчанию)
type LayoutConfig struct {
	DefaultGap Gap // Отступы по умолчанию между элементами
}

// LayoutTree позиционирует дерево с учетом всех параметров
func LayoutTree(root *DOMNode, config LayoutConfig) {
	layoutNode(root, 0, 0, config)
}

// layoutNode рекурсивно размещает узлы
// Возвращает (ширина, высота) включая margin
func layoutNode(node *DOMNode, x, y int16, config LayoutConfig) (int16, int16) {
	if node == nil {
		return 0, 0
	}

	// Применяем margin - узел начинается после margin
	node.X = x + int16(node.Margin.Horizontal)
	node.Y = y + int16(node.Margin.Vertical)

	// Листовой узел с контентом
	if len(node.Children) == 0 && node.Content != "" {
		width, height := calculateContentSize(node)
		node.Width = width
		node.Height = height

		// Возвращаем размер с margin
		totalWidth := width + int16(node.Margin.Horizontal)*2
		totalHeight := height + int16(node.Margin.Vertical)*2
		return totalWidth, totalHeight
	}

	// Контейнер с детьми
	if len(node.Children) > 0 {
		width, height := layoutChildren(node, config)

		// Если размеры не заданы явно, используем вычисленные
		if node.Width == 0 {
			node.Width = width
		}
		if node.Height == 0 {
			node.Height = height
		}

		// Возвращаем размер с margin
		totalWidth := node.Width + int16(node.Margin.Horizontal)*2
		totalHeight := node.Height + int16(node.Margin.Vertical)*2
		return totalWidth, totalHeight
	}

	// Пустой узел без контента и детей - минимальный размер
	if node.Width == 0 {
		node.Width = 3
	}
	if node.Height == 0 {
		node.Height = 3
	}

	totalWidth := node.Width + int16(node.Margin.Horizontal)*2
	totalHeight := node.Height + int16(node.Margin.Vertical)*2
	return totalWidth, totalHeight
}

// calculateContentSize вычисляет размер узла с текстовым контентом
func calculateContentSize(node *DOMNode) (int16, int16) {
	if node.Content == "" {
		return 3, 3 // минимальный размер
	}

	// Разбиваем контент по переносам строк
	lines := strings.Split(node.Content, "\n")

	// Находим самую длинную строку
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// Ширина = контент + padding*2 + border*2 (если есть)
	width := int16(maxLen)
	width += int16(node.Padding.Horizontal) * 2

	if node.HasBorder == Border {
		width += 2 // бордер слева и справа
	}

	// Минимум 3 по ширине
	if width < 3 {
		width = 3
	}

	// Высота = количество строк + padding*2 + border*2 (если есть)
	height := int16(len(lines))
	height += int16(node.Padding.Vertical) * 2

	if node.HasBorder == Border {
		height += 2 // бордер сверху и снизу
	}

	// Минимум 3 по высоте
	if height < 3 {
		height = 3
	}

	return width, height
}

// layoutChildren размещает детей внутри контейнера
func layoutChildren(node *DOMNode, config LayoutConfig) (int16, int16) {
	if len(node.Children) == 0 {
		return 3, 3
	}

	// Начальная позиция для детей (внутри контейнера)
	// Учитываем border (если есть) + padding
	offsetX := int16(0)
	offsetY := int16(0)

	if node.HasBorder == Border {
		offsetX += 1
		offsetY += 1
	}
	offsetX += int16(node.Padding.Horizontal)
	offsetY += int16(node.Padding.Vertical)

	childX := node.X + offsetX
	childY := node.Y + offsetY

	var contentWidth, contentHeight int16

	// Определяем gap между детьми
	gap := config.DefaultGap
	// Можно переопределить gap для конкретного узла (если нужно в будущем)

	if node.Direction == Row {
		// ГОРИЗОНТАЛЬНОЕ размещение (row)
		maxHeight := int16(0)
		currentX := childX

		for i, child := range node.Children {
			childW, childH := layoutNode(child, currentX, childY, config)

			// Следующий ребенок справа (с учетом gap)
			currentX += childW
			if i < len(node.Children)-1 {
				currentX += int16(gap.Horizontal)
			}

			contentWidth += childW
			if childH > maxHeight {
				maxHeight = childH
			}
		}

		// Добавляем gap между детьми
		if len(node.Children) > 1 {
			contentWidth += int16(len(node.Children)-1) * int16(gap.Horizontal)
		}

		contentHeight = maxHeight

	} else {
		// ВЕРТИКАЛЬНОЕ размещение (column)
		maxWidth := int16(0)
		currentY := childY

		for i, child := range node.Children {
			childW, childH := layoutNode(child, childX, currentY, config)

			// Следующий ребенок ниже (с учетом gap)
			currentY += childH
			if i < len(node.Children)-1 {
				currentY += int16(gap.Vertical)
			}

			contentHeight += childH
			if childW > maxWidth {
				maxWidth = childW
			}
		}

		// Добавляем gap между детьми
		if len(node.Children) > 1 {
			contentHeight += int16(len(node.Children)-1) * int16(gap.Vertical)
		}

		contentWidth = maxWidth
	}

	// Итоговый размер контейнера = содержимое + padding*2 + border*2
	totalWidth := contentWidth + int16(node.Padding.Horizontal)*2
	totalHeight := contentHeight + int16(node.Padding.Vertical)*2

	if node.HasBorder == Border {
		totalWidth += 2
		totalHeight += 2
	}

	// Минимум 3x3
	if totalWidth < 3 {
		totalWidth = 3
	}
	if totalHeight < 3 {
		totalHeight = 3
	}

	return totalWidth, totalHeight
}

// RenderTree рисует дерево в терминале
func RenderTree(root *DOMNode, config LayoutConfig) {
	fmt.Println("=== Визуализация дерева ===\n")

	// Вычисляем позиции
	LayoutTree(root, config)

	// Находим границы
	maxX, maxY := getTreeBounds(root)

	// Создаем canvas
	width := int(maxX) + 2
	height := int(maxY) + 2

	if width > 120 {
		width = 120
		fmt.Println("⚠️  Дерево слишком широкое, обрезано до 120 символов")
	}
	if height > 50 {
		height = 50
		fmt.Println("⚠️  Дерево слишком высокое, обрезано до 50 строк")
	}

	canvas := make([][]rune, height)
	for i := range canvas {
		canvas[i] = make([]rune, width)
		for j := range canvas[i] {
			canvas[i][j] = ' '
		}
	}

	// Рисуем узлы
	drawNode(root, canvas)

	// Выводим
	for _, row := range canvas {
		fmt.Println(string(row))
	}

	fmt.Printf("\nРазмер: %d×%d\n", width, height)
}

// getTreeBounds находит максимальные координаты
func getTreeBounds(node *DOMNode) (int16, int16) {
	if node == nil {
		return 0, 0
	}

	maxX := node.X + node.Width + int16(node.Margin.Horizontal)
	maxY := node.Y + node.Height + int16(node.Margin.Vertical)

	for _, child := range node.Children {
		childMaxX, childMaxY := getTreeBounds(child)
		if childMaxX > maxX {
			maxX = childMaxX
		}
		if childMaxY > maxY {
			maxY = childMaxY
		}
	}

	return maxX, maxY
}

// drawNode рисует один узел
func drawNode(node *DOMNode, canvas [][]rune) {
	if node == nil {
		return
	}

	x := int(node.X)
	y := int(node.Y)
	w := int(node.Width)
	h := int(node.Height)

	// Проверка границ
	if y >= len(canvas) || x >= len(canvas[0]) || w < 3 || h < 3 {
		return
	}

	// Рисуем border (если есть)
	if node.HasBorder == Border {
		drawBorder(canvas, x, y, w, h, node.IsRounded == Round)
	}

	// Рисуем content (если есть и это листовой узел)
	if len(node.Children) == 0 && node.Content != "" {
		drawContent(canvas, node, x, y, w, h)
	}

	// Рекурсивно для детей
	for _, child := range node.Children {
		drawNode(child, canvas)
	}
}

// drawBorder рисует рамку
func drawBorder(canvas [][]rune, x, y, w, h int, rounded bool) {
	if y >= len(canvas) || x >= len(canvas[0]) {
		return
	}

	// Символы для рамки
	var topLeft, topRight, bottomLeft, bottomRight rune
	if rounded {
		topLeft, topRight = '╭', '╮'
		bottomLeft, bottomRight = '╰', '╯'
	} else {
		topLeft, topRight = '┌', '┐'
		bottomLeft, bottomRight = '└', '┘'
	}

	// Верх
	if y < len(canvas) {
		for i := 0; i < w && x+i < len(canvas[0]); i++ {
			if i == 0 {
				canvas[y][x+i] = topLeft
			} else if i == w-1 {
				canvas[y][x+i] = topRight
			} else {
				canvas[y][x+i] = '─'
			}
		}
	}

	// Бока
	for j := 1; j < h-1 && y+j < len(canvas); j++ {
		if x < len(canvas[0]) {
			canvas[y+j][x] = '│'
		}
		if x+w-1 < len(canvas[0]) {
			canvas[y+j][x+w-1] = '│'
		}
	}

	// Низ
	if y+h-1 < len(canvas) {
		for i := 0; i < w && x+i < len(canvas[0]); i++ {
			if i == 0 {
				canvas[y+h-1][x+i] = bottomLeft
			} else if i == w-1 {
				canvas[y+h-1][x+i] = bottomRight
			} else {
				canvas[y+h-1][x+i] = '─'
			}
		}
	}
}

// drawContent рисует текстовое содержимое
func drawContent(canvas [][]rune, node *DOMNode, x, y, w, h int) {
	lines := strings.Split(node.Content, "\n")

	// Начальная позиция для текста (с учетом border и padding)
	textX := x
	textY := y

	if node.HasBorder == Border {
		textX += 1
		textY += 1
	}
	textX += int(node.Padding.Horizontal)
	textY += int(node.Padding.Vertical)

	// Рисуем каждую строку
	for lineIdx, line := range lines {
		currentY := textY + lineIdx
		if currentY >= len(canvas) {
			break
		}

		for charIdx, ch := range line {
			currentX := textX + charIdx
			if currentX >= len(canvas[0]) {
				break
			}
			canvas[currentY][currentX] = ch
		}
	}
}

// PrintTree выводит дерево текстом с параметрами
func PrintTree(node *DOMNode, indent string) {
	if node == nil {
		return
	}

	dir := "col"
	if node.Direction == Row {
		dir = "row"
	}

	border := ""
	if node.HasBorder == Border {
		border = " [border]"
		if node.IsRounded == Round {
			border = " [rounded]"
		}
	}

	content := ""
	if node.Content != "" {
		content = fmt.Sprintf(" '%s'", node.Content)
	}

	fmt.Printf("%s%s%s [%s] (%d,%d) %dx%d p:%d,%d m:%d,%d\n",
		indent, content, border, dir,
		node.X, node.Y, node.Width, node.Height,
		node.Padding.Horizontal, node.Padding.Vertical,
		node.Margin.Horizontal, node.Margin.Vertical)

	for _, child := range node.Children {
		PrintTree(child, indent+"  ")
	}
}

// DemoBasic базовое демо
func DemoBasic() {
	// Простая карточка с текстом
	card := &DOMNode{
		Content:   "Hello\nWorld!",
		HasBorder: Border,
		IsRounded: Round,
		Padding:   Gap{Vertical: 1, Horizontal: 2},
		Margin:    Gap{Vertical: 1, Horizontal: 1},
	}

	config := LayoutConfig{
		DefaultGap: Gap{Vertical: 1, Horizontal: 2},
	}

	fmt.Println("=== Простая карточка ===")
	LayoutTree(card, config)
	PrintTree(card, "")
	fmt.Println()
	RenderTree(card, config)
}

// DemoHeader демо шапки сайта
func DemoHeader() {
	// Header с горизонтальным размещением
	header := &DOMNode{
		Direction: Row,
		HasBorder: Border,
		Padding:   Gap{Vertical: 1, Horizontal: 1},
	}

	logo := &DOMNode{
		Content:   "LOGO",
		HasBorder: Border,
		Padding:   Gap{Vertical: 0, Horizontal: 1},
	}

	nav := &DOMNode{
		Direction: Row,
		HasBorder: Border,
		Padding:   Gap{Vertical: 0, Horizontal: 1},
	}

	home := &DOMNode{Content: "Home", Margin: Gap{Horizontal: 1}}
	about := &DOMNode{Content: "About", Margin: Gap{Horizontal: 1}}
	contact := &DOMNode{Content: "Contact", Margin: Gap{Horizontal: 1}}
	nav.Children = []*DOMNode{home, about, contact}

	actions := &DOMNode{
		Direction: Row,
		HasBorder: Border,
		Padding:   Gap{Vertical: 0, Horizontal: 1},
	}

	login := &DOMNode{Content: "Login", Margin: Gap{Horizontal: 1}}
	signup := &DOMNode{Content: "SignUp", HasBorder: Border, IsRounded: Round, Padding: Gap{Horizontal: 1}}
	actions.Children = []*DOMNode{login, signup}

	header.Children = []*DOMNode{logo, nav, actions}

	config := LayoutConfig{
		DefaultGap: Gap{Vertical: 1, Horizontal: 2},
	}

	fmt.Println("=== Header сайта ===")
	LayoutTree(header, config)
	PrintTree(header, "")
	fmt.Println()
	RenderTree(header, config)
}

// DemoMain запуск всех демо
func DemoMain() {
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("ДЕМОНСТРАЦИЯ LAYOUT СИСТЕМЫ")
	fmt.Println(strings.Repeat("=", 70))

	fmt.Println("\n1. Базовая карточка")
	fmt.Println(strings.Repeat("-", 70))
	DemoBasic()

	fmt.Println("\n\n2. Header сайта")
	fmt.Println(strings.Repeat("-", 70))
	DemoHeader()
}
