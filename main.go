package main

import (
	"Guess/ui"
	"Guess/ui/components"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

func setWindowsEnvironmentVars() {
	fmt.Println("📝 Устанавливаем переменные окружения...")

	// Устанавливаем TERM
	if os.Getenv("TERM") == "" {
		os.Setenv("TERM", "xterm-256color")
		fmt.Println("   ✅ TERM установлена: xterm-256color")
	}

	// Устанавливаем COLORTERM
	if os.Getenv("COLORTERM") == "" {
		os.Setenv("COLORTERM", "truecolor")
		fmt.Println("   ✅ COLORTERM установлена: truecolor")
	}

	// Принудительно включаем цвета
	os.Setenv("FORCE_COLOR", "1")
	fmt.Println("   ✅ FORCE_COLOR установлена: 1")

	// Убираем отключение цветов
	if os.Getenv("NO_COLOR") != "" {
		os.Unsetenv("NO_COLOR")
		fmt.Println("   ✅ NO_COLOR удалена")
	}

	fmt.Println()
}

func testBasicColors() {
	fmt.Println("🎨 БАЗОВЫЕ ЦВЕТА (16 цветов):")

	colors := []struct {
		name  string
		color lipgloss.Color
	}{
		{"Черный", lipgloss.Color("0")},
		{"Красный", lipgloss.Color("1")},
		{"Зеленый", lipgloss.Color("2")},
		{"Желтый", lipgloss.Color("3")},
		{"Синий", lipgloss.Color("4")},
		{"Пурпурный", lipgloss.Color("5")},
		{"Циан", lipgloss.Color("6")},
		{"Белый", lipgloss.Color("7")},
		{"Серый", lipgloss.Color("8")},
		{"Яркий красный", lipgloss.Color("9")},
		{"Яркий зеленый", lipgloss.Color("10")},
		{"Яркий желтый", lipgloss.Color("11")},
		{"Яркий синий", lipgloss.Color("12")},
		{"Яркий пурпурный", lipgloss.Color("13")},
		{"Яркий циан", lipgloss.Color("14")},
		{"Яркий белый", lipgloss.Color("15")},
	}

	for _, c := range colors {
		style := lipgloss.NewStyle().
			Background(c.color).
			Foreground(lipgloss.Color("15")).
			Padding(0, 1).
			Bold(true)

		fmt.Printf("   %s %s\n", style.Render("████"), c.name)
	}
	fmt.Println()
}

func main() {

	//setWindowsEnvironmentVars()

	box := components.NewBox().
		SetWidthPercent(90).
		SetPadding(2, 4).
		SetBorder("rounded").
		SetTextColor("9").
		SetBorderColor("10").
		SetBackground("11")

	fmt.Println(box.Render("Содержимое контейнера"))
	fmt.Println(box)

	canvas := ui.NewCanvas()

	fmt.Println(canvas.GetSize())

	//fmt.Printf("🖥️  ОС: %s\n", runtime.GOOS)
	//fmt.Printf("📺 TERM: %s\n", os.Getenv("TERM"))
	//fmt.Printf("🎯 COLORTERM: %s\n", os.Getenv("COLORTERM"))
	//fmt.Printf("🔧 TERM_PROGRAM: %s\n", os.Getenv("TERM_PROGRAM"))
	//fmt.Println()

	//testBasicColors()
}
