package components

import (
	"fmt"
	"strings"
)

// StyleAdapter отвечает за преобразование стилей в ANSI коды
type StyleAdapter struct {
	colorParser *ColorParser
}

// NewStyleAdapter создает новый рендерер стилей
func NewStyleAdapter() *StyleAdapter {
	return &StyleAdapter{
		colorParser: NewColorParser(),
	}
}

// AdaptTextStyle конвертирует TextStyle в ANSI коды
func (r *StyleAdapter) AdaptTextStyle(style TextStyle) string {
	var codes []string

	if style.Bold {
		codes = append(codes, "1")
	}
	if style.Dim {
		codes = append(codes, "2")
	}
	if style.Italic {
		codes = append(codes, "3")
	}
	if style.Underline {
		codes = append(codes, "4")
	}
	if style.Blink {
		codes = append(codes, "5")
	}
	if style.Reverse {
		codes = append(codes, "7")
	}
	if style.Hidden {
		codes = append(codes, "8")
	}
	if style.Strike {
		codes = append(codes, "9")
	}

	return strings.Join(codes, ";")
}

// AdaptStyle генерирует полную ANSI последовательность для стиля
// Возвращает пустую строку и error при невалидном цвете
func (r *StyleAdapter) AdaptStyle(style UIComponentStyle) (string, error) {
	var codes []string

	// Текстовые стили
	if textStyleCode := r.AdaptTextStyle(style.TextStyle); textStyleCode != "" {
		codes = append(codes, textStyleCode)
	}

	// Цвет переднего плана
	if style.Color != "" {
		fgCode, err := r.colorParser.ParseToANSI(style.Color, false)
		if err != nil {
			return "", fmt.Errorf("invalid foreground color: %w", err)
		}
		if fgCode != "" {
			codes = append(codes, fgCode)
		}
	}

	// Цвет фона
	if style.Background != "" {
		bgCode, err := r.colorParser.ParseToANSI(style.Background, true)
		if err != nil {
			return "", fmt.Errorf("invalid background color: %w", err)
		}
		if bgCode != "" {
			codes = append(codes, bgCode)
		}
	}

	if len(codes) == 0 {
		return "", nil
	}

	return fmt.Sprintf("\033[%sm", strings.Join(codes, ";")), nil
}

// ApplyStyleToText применяет стиль к тексту
// Возвращает текст без изменений если стиль невалиден (ошибку игнорируем для graceful degradation)
func (r *StyleAdapter) ApplyStyleToText(text string, style UIComponentStyle) string {
	ansi, err := r.AdaptStyle(style)
	if err != nil || ansi == "" {
		return text
	}
	return ansi + text + r.ResetStyle()
}

// ResetStyle возвращает ANSI последовательность сброса
func (r *StyleAdapter) ResetStyle() string {
	return "\033[0m"
}
