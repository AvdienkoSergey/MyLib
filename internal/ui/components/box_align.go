package components

import (
	"unicode/utf8"
)

func CalculateAlignOffset(cWidth int, cBorder string, cPaddingH int, cLayout string, cAlign string, cChildren []interface{}) int {
	// Вычисляем доступную ширину для контента
	availableWidth := cWidth

	// Вычитаем border
	borderSize := GetBorderSize(cBorder)
	if borderSize.Width > 0 {
		availableWidth -= borderSize.Width
	}

	// Вычитаем padding
	availableWidth -= cPaddingH * 2

	contentWidth := 0
	isRowLayout := cLayout == "row"
	// Для получения ширины контента требуется обработать вложенные элементы
	for _, child := range cChildren {
		switch v := child.(type) {
		case string:
			strWidth := utf8.RuneCountInString(v)
			if strWidth > contentWidth && isRowLayout == false {
				contentWidth = strWidth
			}
			if isRowLayout == true {
				contentWidth += strWidth
			}
		case *Box:
			childWidth := v.width
			if childWidth > contentWidth && isRowLayout == false {
				contentWidth = childWidth
			}
			if isRowLayout == true {
				contentWidth += childWidth
			}
		default:
			// Просто ничего не делаем
		}
	}
	// Если контент шире доступного места - выравнивание не применяется
	if contentWidth >= availableWidth {
		return 0
	}

	// Вычисляем смещение в зависимости от align
	switch cAlign {
	case "center":
		return (availableWidth - contentWidth) / 2
	case "right":
		return availableWidth - contentWidth
	case "left":
	default:
		return 0
	}
	return 0
}
