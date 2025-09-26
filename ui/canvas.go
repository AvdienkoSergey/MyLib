package ui

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// Canvas - современная реализация работы с терминалом
type Canvas struct {
	width  int
	height int
	cached bool
}

// NewCanvas создает новый объект для работы с терминалом
func NewCanvas() *Canvas {
	return &Canvas{}
}

// GetSize возвращает размеры терминала (width, height)
func (t *Canvas) GetSize() (int, int) {
	if !t.cached {
		t.updateSize()
	}
	return t.width, t.height
}

// GetWidth возвращает ширину терминала
func (t *Canvas) GetWidth() int {
	width, _ := t.GetSize()
	return width
}

// GetHeight возвращает высоту терминала
func (t *Canvas) GetHeight() int {
	_, height := t.GetSize()
	return height
}

// Refresh обновляет кешированные размеры
func (t *Canvas) Refresh() {
	t.cached = false
	t.updateSize()
}

// updateSize обновляет размеры терминала
func (t *Canvas) updateSize() {
	width, height := getCanvasSize()
	t.width = width
	t.height = height
	t.cached = true
}

// getCanvasSize - современная версия определения размера терминала
func getCanvasSize() (width, height int) {
	// Способ 1: Для Unix-систем используем golang.org/x/term
	if width, height := getTerminalSizeXTerm(); width > 0 && height > 0 {
		return width, height
	}

	// Способ 3: Windows API
	if runtime.GOOS == "windows" {
		if width, height := getTerminalSizeWindows(); width > 0 && height > 0 {
			return width, height
		}
	}

	// Способ 6: Значения по умолчанию
	return 80, 24
}

// getTerminalSizeXTerm использует golang.org/x/term (рекомендуемый способ)
func getTerminalSizeXTerm() (width, height int) {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		width, height, err := term.GetSize(int(os.Stdin.Fd()))
		if err == nil {
			return width, height
		}
	}

	return 0, 0 // Пока не используем, чтобы избежать зависимости
}

// getTerminalSizeWindows использует Windows Console API
func getTerminalSizeWindows() (width, height int) {
	// Для Windows можно использовать PowerShell или WinAPI
	// Простое решение через PowerShell:
	cmd := exec.Command("powershell", "-Command",
		"$host.UI.RawUI.WindowSize.Width; $host.UI.RawUI.WindowSize.Height")

	output, err := cmd.Output()
	if err != nil {
		return 0, 0
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) >= 2 {
		if w, err := strconv.Atoi(strings.TrimSpace(lines[0])); err == nil {
			width = w
		}
		if h, err := strconv.Atoi(strings.TrimSpace(lines[1])); err == nil {
			height = h
		}
	}

	return width, height
}
