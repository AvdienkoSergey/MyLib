package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Пример использования компонента Box с древовидной структурой
//
// Простой пример:
//	root := NewBox().
//		SetWidth(60).
//		SetPadding(1, 2).
//		SetBorder("rounded").
//		SetBorderColor("blue")
//
//	root.AddChild("Заголовок приложения")
//	root.AddChild(NewBox().
//		SetBorder("normal").
//		SetPadding(1, 2).
//		AddChild("Вложенный контейнер"))
//
//	fmt.Println(root.Render())
//
// Сложный пример с позиционированием:
//	page := NewBox().SetWidth(80).SetPadding(1, 2).SetBorder("double")
//
//	header := NewBox().SetBorder("rounded").SetAlign("center").SetTextBold(true)
//	header.AddChild("Мое приложение")
//
//	// Контент со смещением
//	content := NewBox().SetPadding(1, 1).SetOffset(5, 2)
//	content.AddChild("Строка 1")
//	content.AddChild(NewBox().SetBorder("normal").SetOffset(10, 0).AddChild("Смещенный блок"))
//
//	page.AddChild(header)
//	page.AddChild(content)
//
//	// Вычисляем абсолютные позиции (начало в 1,1 терминала)
//	page.CalculateAbsolutePositions(1, 1)
//
//	// Теперь можно использовать абсолютные позиции для cursor
//	page.Walk(func(b *Box, depth int) {
//		x, y := b.GetAbsolutePosition()
//		fmt.Printf("\033[%d;%dH", y, x)  // Перемещаем курсор
//		fmt.Print(b.Render())
//	})

type Box struct {
	id          string // Уникальный идентификатор
	width       int
	height      int
	paddingV    int
	paddingH    int
	marginV     int
	marginH     int
	border      string
	borderColor string
	background  string
	foreground  string
	align       string
	bold        bool
	children    []interface{} // Дочерние элементы: string или *Box
	offsetX     int           // Смещение относительно родителя по X
	offsetY     int           // Смещение относительно родителя по Y
	absoluteX   int           // Абсолютная позиция в терминале (вычисляется)
	absoluteY   int           // Абсолютная позиция в терминале (вычисляется)
	layout      string        // "column" (по умолчанию) или "row"
	gap         int           // Расстояние между дочерними элементами
}

func NewBox() *Box {
	return &Box{
		id:          "",
		width:       0,
		height:      0,
		paddingV:    0,
		paddingH:    0,
		marginV:     0,
		marginH:     0,
		border:      "",
		borderColor: "white",
		background:  "",
		foreground:  "white",
		align:       "left",
		bold:        false,
		children:    make([]interface{}, 0),
		offsetX:     0,
		offsetY:     0,
		absoluteX:   0,
		absoluteY:   0,
		layout:      "column",
		gap:         0,
	}
}

// SetID установить уникальный идентификатор
func (c *Box) SetID(id string) *Box {
	c.id = id
	return c
}

// SetWidth управляет шириной
func (c *Box) SetWidth(width int) *Box {
	c.width = width
	return c
}

// SetHeight управляет высотой
func (c *Box) SetHeight(height int) *Box {
	c.height = height
	return c
}

// SetPadding управляет внутренними отступами
func (c *Box) SetPadding(vertical, horizontal int) *Box {
	c.paddingV = vertical
	c.paddingH = horizontal
	return c
}

// SetMargin управляет внешними отступами
func (c *Box) SetMargin(vertical, horizontal int) *Box {
	c.marginV = vertical
	c.marginH = horizontal
	return c
}

// SetBorder устанавливает тип обводки
func (c *Box) SetBorder(borderType string) *Box {
	c.border = borderType
	return c
}

// SetBorderColor устанавливает цвет обводки
func (c *Box) SetBorderColor(color string) *Box {
	c.borderColor = color
	return c
}

// SetBackground устанавливает цвет фона
func (c *Box) SetBackground(color string) *Box {
	c.background = color
	return c
}

// SetTextColor устанавливает цвет текста
func (c *Box) SetTextColor(color string) *Box {
	c.foreground = color
	return c
}

// SetAlign устанавливает позиционирование контента
func (c *Box) SetAlign(align string) *Box {
	c.align = align
	return c
}

// SetTextBold делает текст жирным
func (c *Box) SetTextBold(bold bool) *Box {
	c.bold = bold
	return c
}

// SetOffset устанавливает смещение относительно родителя
func (c *Box) SetOffset(x, y int) *Box {
	c.offsetX = x
	c.offsetY = y
	return c
}

// SetLayout устанавливает направление размещения дочерних элементов
// "column" - вертикально (по умолчанию), "row" - горизонтально
func (c *Box) SetLayout(layout string) *Box {
	c.layout = layout
	return c
}

// SetGap устанавливает расстояние между дочерними элементами
func (c *Box) SetGap(gap int) *Box {
	c.gap = gap
	return c
}

// GetID получить уникальный идентификатор
func (c *Box) GetID() string {
	return c.id
}

// GetOffset возвращает смещение относительно родителя
func (c *Box) GetOffset() (x, y int) {
	return c.offsetX, c.offsetY
}

// GetAbsolutePosition возвращает абсолютную позицию в терминале
func (c *Box) GetAbsolutePosition() (x, y int) {
	return c.absoluteX, c.absoluteY
}

// AddChild добавляет дочерний элемент (string или *Box)
func (c *Box) AddChild(child interface{}) *Box {
	c.children = append(c.children, child)
	return c
}

// AddChildren добавляет несколько дочерних элементов
func (c *Box) AddChildren(children ...interface{}) *Box {
	c.children = append(c.children, children...)
	return c
}

// RemoveChild удаляет дочерний элемент по индексу
func (c *Box) RemoveChild(index int) *Box {
	if index >= 0 && index < len(c.children) {
		c.children = append(c.children[:index], c.children[index+1:]...)
	}
	return c
}

// ClearChildren удаляет все дочерние элементы
func (c *Box) ClearChildren() *Box {
	c.children = make([]interface{}, 0)
	return c
}

// GetChildren возвращает все дочерние элементы
func (c *Box) GetChildren() []interface{} {
	return c.children
}

// joinHorizontal объединяет элементы горизонтально используя lipgloss.JoinHorizontal
func (c *Box) joinHorizontal(parts []string) string {
	if len(parts) == 0 {
		return ""
	}

	// Если задан gap, добавляем пробелы между элементами
	if c.gap > 0 {
		gapStyle := lipgloss.NewStyle().Width(c.gap)
		var partsWithGaps []string
		for i, part := range parts {
			partsWithGaps = append(partsWithGaps, part)
			if i < len(parts)-1 {
				partsWithGaps = append(partsWithGaps, gapStyle.Render(""))
			}
		}
		return lipgloss.JoinHorizontal(lipgloss.Top, partsWithGaps...)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

// joinVertical объединяет элементы вертикально используя lipgloss.JoinVertical
func (c *Box) joinVertical(parts []string) string {
	if len(parts) == 0 {
		return ""
	}

	// Если задан gap, добавляем пустые строки между элементами
	if c.gap > 0 {
		gapStyle := lipgloss.NewStyle().Height(c.gap)
		var partsWithGaps []string
		for i, part := range parts {
			partsWithGaps = append(partsWithGaps, part)
			if i < len(parts)-1 {
				partsWithGaps = append(partsWithGaps, gapStyle.Render(""))
			}
		}
		return lipgloss.JoinVertical(lipgloss.Left, partsWithGaps...)
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// visualLength вычисляет визуальную длину строки (без ANSI escape-последовательностей)
func (c *Box) visualLength(s string) int {
	// Простое удаление ANSI escape-последовательностей
	inEscape := false
	length := 0

	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if s[i] == 'm' || s[i] == 'H' {
				inEscape = false
			}
			continue
		}
		length++
	}

	return length
}

// truncateToVisualLength обрезает строку до указанной визуальной длины
func (c *Box) truncateToVisualLength(s string, maxLen int) string {
	visualLen := 0
	result := strings.Builder{}
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			result.WriteByte(s[i])
			continue
		}
		if inEscape {
			result.WriteByte(s[i])
			if s[i] == 'm' || s[i] == 'H' {
				inEscape = false
			}
			continue
		}

		if visualLen >= maxLen {
			break
		}

		result.WriteByte(s[i])
		visualLen++
	}

	return result.String()
}

// Draw рендерит содержимое (устаревший метод, оставлен для совместимости)
func (c *Box) Draw() string {
	borderSize := GetBorderSize(c.border)
	if borderSize.Width != 0 {
		c.width -= borderSize.Width
		c.height -= borderSize.Height
	}
	c.SetOffset(CalculateOffset(c.width,
		c.border,
		c.layout,
		c.align,
		c.paddingH,
		c.paddingV,
		c.marginH,
		c.marginV,
		c.children,
	))
	content := c.drawChildren()
	return c.drawWithContent(content)
}

// drawChildren рекурсивно рендерит все дочерние элементы
func (c *Box) drawChildren() string {
	if len(c.children) == 0 {
		return ""
	}

	var parts []string
	for _, child := range c.children {
		switch v := child.(type) {
		case string:
			parts = append(parts, v)
		case *Box:
			// Рекурсивный рендеринг вложенного Box
			parts = append(parts, v.Draw())
		default:
			// Для других типов используем fmt.Sprint
			parts = append(parts, strings.TrimSpace(strings.Join(strings.Fields(string(rune(0))), " ")))
		}
	}

	// Объединяем в зависимости от layout
	if c.layout == "row" {
		return c.joinHorizontal(parts)
	}

	// Вертикальное размещение (по умолчанию)
	return c.joinVertical(parts)
}

// drawWithContent применяет стили к контенту
func (c *Box) drawWithContent(content string) string {
	style := lipgloss.NewStyle()

	// Размеры
	if c.width > 0 {
		style = style.Width(c.width)
	}
	if c.height > 0 {
		style = style.Height(c.height)
	}

	// Отступы
	if c.paddingV > 0 || c.paddingH > 0 {
		style = style.Padding(c.paddingV, c.paddingH)
	}
	if c.marginV > 0 || c.marginH > 0 {
		style = style.Margin(c.marginV, c.marginH)
	}

	// Рамка
	if c.border != "" {
		switch c.border {
		case "normal":
			style = style.Border(lipgloss.NormalBorder())
			c.width = 80
			c.height = 30
		case "rounded":
			style = style.Border(lipgloss.RoundedBorder())
		case "thick":
			style = style.Border(lipgloss.ThickBorder())
		case "double":
			style = style.Border(lipgloss.DoubleBorder())
		case "hidden":
			style = style.Border(lipgloss.HiddenBorder())
		}
		style = style.BorderForeground(lipgloss.Color(c.borderColor))
	}

	// Цвета
	if c.background != "" {
		style = style.Background(lipgloss.Color(c.background))
	}
	if c.foreground != "" {
		style = style.Foreground(lipgloss.Color(c.foreground))
	}

	// Выравнивание
	switch c.align {
	case "center":
		style = style.Align(lipgloss.Center)
	case "right":
		style = style.Align(lipgloss.Right)
	default:
		style = style.Align(lipgloss.Left)
	}

	// Жирный текст
	if c.bold {
		style = style.Bold(true)
	}

	return style.Render(content)
}

// State возвращает текущее состояние Box
func (c *Box) State() (width, height, paddingV, paddingH, marginV, marginH int) {
	return c.width, c.height, c.paddingV, c.paddingH, c.marginV, c.marginH
}

// GetChildCount возвращает количество дочерних элементов
func (c *Box) GetChildCount() int {
	return len(c.children)
}

// Walk обходит дерево элементов с помощью callback функции
func (c *Box) Walk(callback func(*Box, int)) {
	c.walk(callback, 0)
}

func (c *Box) walk(callback func(*Box, int), depth int) {
	callback(c, depth)
	for _, child := range c.children {
		if box, ok := child.(*Box); ok {
			box.walk(callback, depth+1)
		}
	}
}

// CalculateAbsolutePositions вычисляет абсолютные позиции для всех элементов дерева
// startX, startY - начальная позиция корневого элемента в терминале
func (c *Box) CalculateAbsolutePositions(startX, startY int) {
	c.calculatePositions(startX, startY, 0, 0)
}

func (c *Box) calculatePositions(parentAbsX, parentAbsY, parentContentX, parentContentY int) {
	// Вычисляем абсолютную позицию текущего Box
	c.absoluteX = parentAbsX + c.offsetX
	c.absoluteY = parentAbsY + c.offsetY

	// Вычисляем позицию контента внутри текущего Box
	contentX := c.absoluteX + c.marginH
	contentY := c.absoluteY + c.marginV

	// Если есть рамка, добавляем 1 символ
	if c.border != "" && c.border != "hidden" {
		contentX += 1
		contentY += 1
	}

	// Добавляем padding
	contentX += c.paddingH
	contentY += c.paddingV

	// Обрабатываем дочерние элементы в зависимости от layout
	if c.layout == "row" {
		// Горизонтальное размещение
		currentX := contentX
		for _, child := range c.children {
			if box, ok := child.(*Box); ok {
				box.calculatePositions(currentX, contentY, currentX, contentY)
				currentX += box.GetTotalWidth()

				// Добавляем gap между элементами
				if c.gap > 0 {
					currentX += c.gap
				}
			} else {
				// Для строк - сложнее, нужно знать их ширину
				currentX += 10 // примерная ширина, или передавайте реальную
			}
		}
	} else {
		// Вертикальное размещение (column)
		currentY := contentY
		for _, child := range c.children {
			if box, ok := child.(*Box); ok {
				box.calculatePositions(contentX, currentY, contentX, currentY)
				currentY += box.GetTotalHeight()

				// Добавляем gap между элементами
				if c.gap > 0 {
					currentY += c.gap
				}
			} else {
				currentY += 1
			}
		}
	}
}

// GetTotalHeight возвращает полную высоту элемента с учетом всех отступов
func (c *Box) GetTotalHeight() int {
	height := 0

	// Margin
	height += c.marginV * 2

	// Border
	if c.border != "" && c.border != "hidden" {
		height += 2
	}

	// Padding
	height += c.paddingV * 2

	// Контент
	if c.height > 0 {
		height += c.height
	} else {
		// Если высота не задана, считаем по количеству дочерних элементов
		for _, child := range c.children {
			if box, ok := child.(*Box); ok {
				height += box.GetTotalHeight()
			} else {
				height += 1
			}
		}
		if len(c.children) == 0 {
			height += 1 // Минимум 1 строка
		}
	}

	return height
}

// GetTotalWidth возвращает полную ширину элемента с учетом всех отступов
func (c *Box) GetTotalWidth() int {
	width := 0

	// Margin
	width += c.marginH * 2

	// Border
	if c.border != "" && c.border != "hidden" {
		width += 2
	}

	// Padding
	width += c.paddingH * 2

	// Контент
	if c.width > 0 {
		width += c.width
	} else {
		width += 20 // Ширина по умолчанию
	}

	return width
}

// FindByID ищет элемент по ID в дереве
func (c *Box) FindByID(id string) *Box {
	if c.id == id {
		return c
	}

	for _, child := range c.children {
		if box, ok := child.(*Box); ok {
			if found := box.FindByID(id); found != nil {
				return found
			}
		}
	}

	return nil
}

// GetParent ищет родителя для указанного элемента
func (c *Box) GetParent(target *Box) *Box {
	for _, child := range c.children {
		if box, ok := child.(*Box); ok {
			if box == target {
				return c
			}
			if parent := box.GetParent(target); parent != nil {
				return parent
			}
		}
	}
	return nil
}
