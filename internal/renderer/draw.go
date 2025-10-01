package renderer

import (
	"Guess/internal/ui/terminal"
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

type Position struct {
	X, Y int
}

type ColorSchema struct {
	FG string
	BG string
}

type DrawTask struct {
	Position      Position
	Width, Height int
	Content       []string
	ColorSchema   ColorSchema
}

func NewDrawTask() *DrawTask {
	return &DrawTask{}
}

func (t *DrawTask) SetContent(content []string) *DrawTask {
	t.Content = content
	return t
}

func (t *DrawTask) SetSize(width, height int) *DrawTask {
	t.Width = width
	t.Height = height

	if t.Width < 1 || t.Height < 1 {
		t.SetAutoSize()
	}
	return t
}

func (t *DrawTask) SetAutoSize() *DrawTask {
	if t.Content == nil {
		panic("when you work with DrawTask you should set content before set auto size")
	}

	if t.Width == 0 {
		maxLen := 0
		for _, line := range t.Content {
			lineLen := runewidth.StringWidth(line)
			if lineLen > maxLen {
				maxLen = lineLen
			}
		}
		t.Width = maxLen
	}

	if t.Height == 0 {
		t.Height = len(t.Content)
	}
	return t
}

func (t *DrawTask) SetWidth(w int) *DrawTask {
	t.Width = w
	return t
}

func (t *DrawTask) SetHeight(h int) *DrawTask {
	t.Height = h
	return t
}

func (t *DrawTask) SetPosition(x, y int) *DrawTask {
	t.Position = Position{
		X: x,
		Y: y,
	}
	return t
}

func (t *DrawTask) SetColorSchema(fg, bg string) *DrawTask {
	t.ColorSchema.FG = fg
	t.ColorSchema.BG = bg
	return t
}

func (t *DrawTask) Draw() {
	cursor := terminal.NewCursorManager()
	cursor.HideCursor()

	// Применяем цвета если заданы
	if t.ColorSchema.FG != "" || t.ColorSchema.BG != "" {
		applyColors(t.ColorSchema)
	}

	for i, str := range t.Content {
		y := t.Position.Y + i
		runes := []rune(str)

		var token strings.Builder
		tokenStartX := -1

		for ii, r := range runes {
			if r == ' ' {
				if token.Len() > 0 {
					cursor.WriteAt(t.Position.X+tokenStartX, y, token.String())
					token.Reset()
					tokenStartX = -1
				}
			} else {
				if tokenStartX == -1 {
					tokenStartX = ii
				}
				token.WriteRune(r)
			}

			if token.Len() > 0 {
				cursor.WriteAt(t.Position.X+tokenStartX, y, token.String())
			}
		}
	}

	// Сброс цветов
	fmt.Print("\033[0m")
}

func applyColors(schema ColorSchema) {
	if schema.FG != "" {
		fmt.Printf("\033[%sm", colorToANSI(schema.FG))
	}
	if schema.BG != "" {
		fmt.Printf("\033[%sm", colorToBgANSI(schema.BG))
	}
}

func colorToANSI(color string) string {
	colors := map[string]string{
		"black": "30", "red": "31", "green": "32", "yellow": "33",
		"blue": "34", "magenta": "35", "cyan": "36", "white": "37",
	}
	if code, ok := colors[color]; ok {
		return code
	}
	return "37"
}

func colorToBgANSI(color string) string {
	colors := map[string]string{
		"black": "40", "red": "41", "green": "42", "yellow": "43",
		"blue": "44", "magenta": "45", "cyan": "46", "white": "47",
	}
	if code, ok := colors[color]; ok {
		return code
	}
	return "40"
}
