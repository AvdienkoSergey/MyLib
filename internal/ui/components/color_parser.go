package components

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrInvalidColorFormat возвращается когда формат цвета некорректен
	ErrInvalidColorFormat = errors.New("invalid color format")
	// ErrInvalidRGBValue возвращается когда RGB значения вне диапазона 0-255
	ErrInvalidRGBValue = errors.New("RGB values must be in range 0-255")
	// ErrInvalidHexFormat возвращается когда hex формат некорректен
	ErrInvalidHexFormat = errors.New("hex color must be 6 characters")
)

// ColorParser отвечает за парсинг цветов
type ColorParser struct{}

// NewColorParser создает новый парсер цветов
func NewColorParser() *ColorParser {
	return &ColorParser{}
}

// ParseToANSI парсит цвет и возвращает ANSI код
// Возвращает пустую строку и nil для пустого цвета (используется дефолтный цвет терминала)
// Возвращает error для невалидного формата цвета
func (p *ColorParser) ParseToANSI(color string, isBackground bool) (string, error) {
	if color == "" {
		return "", nil
	}

	// RGB формат: rgb(255,0,0)
	if strings.HasPrefix(color, "rgb(") && strings.HasSuffix(color, ")") {
		r, g, b, err := p.parseRGB(color)
		if err != nil {
			return "", fmt.Errorf("parsing RGB color %q: %w", color, err)
		}
		return p.rgbToANSI(r, g, b, isBackground), nil
	}

	// Hex формат: #FF0000
	if strings.HasPrefix(color, "#") {
		r, g, b, err := p.parseHex(color)
		if err != nil {
			return "", fmt.Errorf("parsing hex color %q: %w", color, err)
		}
		return p.rgbToANSI(r, g, b, isBackground), nil
	}

	// 256 формат: число от 0 до 255
	if num, err := strconv.Atoi(color); err == nil {
		if num < 0 || num > 255 {
			return "", fmt.Errorf("256-color value %d out of range (must be 0-255)", num)
		}
		return p.color256ToANSI(color, isBackground), nil
	}

	// Базовые именованные цвета
	result := p.basicColorToANSI(color, isBackground)
	if result == "" {
		return "", fmt.Errorf("unknown color name %q: %w", color, ErrInvalidColorFormat)
	}
	return result, nil
}

// parseRGB парсит RGB значение из строки rgb(255,0,0)
func (p *ColorParser) parseRGB(color string) (r, g, b int, err error) {
	color = strings.TrimPrefix(color, "rgb(")
	color = strings.TrimSuffix(color, ")")
	parts := strings.Split(color, ",")
	if len(parts) != 3 {
		return 0, 0, 0, fmt.Errorf("%w: expected 3 values, got %d", ErrInvalidColorFormat, len(parts))
	}

	r, err = strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("%w: invalid red value: %v", ErrInvalidColorFormat, err)
	}
	if r < 0 || r > 255 {
		return 0, 0, 0, fmt.Errorf("%w: red=%d", ErrInvalidRGBValue, r)
	}

	g, err = strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("%w: invalid green value: %v", ErrInvalidColorFormat, err)
	}
	if g < 0 || g > 255 {
		return 0, 0, 0, fmt.Errorf("%w: green=%d", ErrInvalidRGBValue, g)
	}

	b, err = strconv.Atoi(strings.TrimSpace(parts[2]))
	if err != nil {
		return 0, 0, 0, fmt.Errorf("%w: invalid blue value: %v", ErrInvalidColorFormat, err)
	}
	if b < 0 || b > 255 {
		return 0, 0, 0, fmt.Errorf("%w: blue=%d", ErrInvalidRGBValue, b)
	}

	return r, g, b, nil
}

// parseHex парсит hex значение из строки #FF0000
func (p *ColorParser) parseHex(color string) (r, g, b int, err error) {
	color = strings.TrimPrefix(color, "#")
	if len(color) != 6 {
		return 0, 0, 0, fmt.Errorf("%w: got %d characters", ErrInvalidHexFormat, len(color))
	}
	n, scanErr := fmt.Sscanf(color, "%02x%02x%02x", &r, &g, &b)
	if scanErr != nil {
		return 0, 0, 0, fmt.Errorf("%w: %v", ErrInvalidHexFormat, scanErr)
	}
	if n != 3 {
		return 0, 0, 0, fmt.Errorf("%w: scanned %d values instead of 3", ErrInvalidHexFormat, n)
	}
	return r, g, b, nil
}

// rgbToANSI конвертирует RGB в ANSI код (24-bit true color)
func (p *ColorParser) rgbToANSI(r, g, b int, isBackground bool) string {
	if isBackground {
		return fmt.Sprintf("48;2;%d;%d;%d", r, g, b)
	}
	return fmt.Sprintf("38;2;%d;%d;%d", r, g, b)
}

// color256ToANSI конвертирует 256-color в ANSI код
func (p *ColorParser) color256ToANSI(color string, isBackground bool) string {
	if isBackground {
		return fmt.Sprintf("48;5;%s", color)
	}
	return fmt.Sprintf("38;5;%s", color)
}

// basicColorToANSI конвертирует базовый именованный цвет в ANSI код
func (p *ColorParser) basicColorToANSI(color string, isBackground bool) string {
	fgColors := map[string]string{
		"black": "30", "red": "31", "green": "32", "yellow": "33",
		"blue": "34", "magenta": "35", "cyan": "36", "white": "37",
		"bright_black": "90", "bright_red": "91", "bright_green": "92", "bright_yellow": "93",
		"bright_blue": "94", "bright_magenta": "95", "bright_cyan": "96", "bright_white": "97",
	}

	bgColors := map[string]string{
		"black": "40", "red": "41", "green": "42", "yellow": "43",
		"blue": "44", "magenta": "45", "cyan": "46", "white": "47",
		"bright_black": "100", "bright_red": "101", "bright_green": "102", "bright_yellow": "103",
		"bright_blue": "104", "bright_magenta": "105", "bright_cyan": "106", "bright_white": "107",
	}

	if isBackground {
		if code, ok := bgColors[color]; ok {
			return code
		}
	} else {
		if code, ok := fgColors[color]; ok {
			return code
		}
	}

	return ""
}
