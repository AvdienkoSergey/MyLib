package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// Terminal - современная реализация работы с терминалом
type Terminal struct {
	width  int
	height int
	cached bool
}

// NewTerminal создает новый объект для работы с терминалом
func NewTerminal() *Terminal {
	return &Terminal{}
}

// GetSize возвращает размеры терминала (width, height)
func (t *Terminal) GetSize() (int, int) {
	if !t.cached {
		t.updateSize()
	}
	return t.width, t.height
}

// GetWidth возвращает ширину терминала
func (t *Terminal) GetWidth() int {
	width, _ := t.GetSize()
	return width
}

// GetHeight возвращает высоту терминала
func (t *Terminal) GetHeight() int {
	_, height := t.GetSize()
	return height
}

// Refresh обновляет кешированные размеры
func (t *Terminal) Refresh() {
	t.cached = false
	t.updateSize()
}

// updateSize обновляет размеры терминала
func (t *Terminal) updateSize() {
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

// Clear Кроссплатформенная очистка экрана
func Clear() {
	// Проверяем, есть ли переменная TERM (признак настоящего терминала)
	if os.Getenv("TERM") == "" {
		// Не в терминале, используем простые переводы строк
		fmt.Println("\033[?25l")
	}

	switch runtime.GOOS {
	case "darwin", "linux":
		// Для macOS и Linux - используем комбинацию ANSI кодов
		// \033[H - курсор в (1,1)
		// \033[2J - очистить весь экран
		// \033[3J - очистить scrollback buffer (для macOS Terminal)
		fmt.Print("\033[H\033[2J\033[3J")

	case "windows":
		// Для Windows
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			// Fallback к ANSI если cls не работает
			fmt.Print("\033[2J\033[H")
		}

	default:
		fmt.Print("\033[2J\033[H")
	}
}
