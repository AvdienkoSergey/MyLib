package main

import (
	"Guess/internal/ui/components"
	"fmt"
)

func main() {
	// Создаем стиль
	style := components.NewUIComponentStyleBuilder().
		SetColor("red").
		SetBackground("yellow").
		SetItalic(true).
		SetFlex(components.Flex{
			Direction: components.FlexDirectionRow,
			Justify:   components.FlexJustifyCenter,
			Align:     components.FlexAlignStretch,
			Gap:       10,
		}).
		Build()

	// Создаем компонент
	component := components.NewUIComponentBuilder().
		SetText("Hello World\nIt`s me").
		SetWidth(100).
		SetHeight(50).
		SetStyle(style).
		Build()

	// Рендерим стиль
	renderer := components.NewStyleAdapter()
	styledText := renderer.ApplyStyleToText(component.Text, component.Style)

	fmt.Println(styledText)
}
