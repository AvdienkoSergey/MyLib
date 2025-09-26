package components

import (
	"Guess/ui"

	"github.com/charmbracelet/lipgloss"
)

// Пример использования компонента Box
// Базовый контейнер для группировки элементов
//
//	box := NewBox().
//		SetWidth(50).
//		SetPadding(2, 4).
//		SetBorder("rounded").
//		SetBorderColor("blue").
//		SetBackground("black")
//
//	fmt.Println(box.Render("Содержимое контейнера"))

type Box struct {
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
}

func NewBox() *Box {
	return &Box{
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
	}
}

func (c *Box) SetWidth(width int) *Box {
	c.width = width
	return c
}

func (c *Box) SetWidthPercent(percent int) *Box {
	width := ui.NewCanvas().GetWidth()
	c.width = width * percent / 100
	return c
}

func (c *Box) SetHeight(height int) *Box {
	c.height = height
	return c
}

func (c *Box) SetHeightPercent(percent int) *Box {
	height := ui.NewCanvas().GetHeight()
	c.height = height * percent / 100
	return c
}

func (c *Box) SetPadding(vertical, horizontal int) *Box {
	c.paddingV = vertical
	c.paddingH = horizontal
	return c
}

func (c *Box) SetMargin(vertical, horizontal int) *Box {
	c.marginV = vertical
	c.marginH = horizontal
	return c
}

func (c *Box) SetBorder(borderType string) *Box {
	c.border = borderType
	return c
}

func (c *Box) SetBorderColor(color string) *Box {
	c.borderColor = color
	return c
}

func (c *Box) SetBackground(color string) *Box {
	c.background = color
	return c
}

func (c *Box) SetTextColor(color string) *Box {
	c.foreground = color
	return c
}

func (c *Box) SetAlign(align string) *Box {
	c.align = align
	return c
}

func (c *Box) SetTextBold(bold bool) *Box {
	c.bold = bold
	return c
}

func (c *Box) Render(content string) string {
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
