package main

import (
	"Guess/internal/monitor"
	"Guess/internal/proxy"
	"Guess/internal/reactivity"
	"Guess/internal/ui/components"
	"Guess/internal/ui/terminal"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unicode/utf8"
)

type MemoryMonitorReport struct {
	AllocMB     string
	SysMB       string
	NumGC       string
	Goroutines  string
	HeapObjects string
}

func (m *MemoryMonitorReport) fieldNamesMemoryMonitor() []string {
	return []string{
		"AllocMB",
		"SysMB",
		"NumGC",
		"Goroutines",
		"HeapObjects",
	}
}

func createWatcher(entity *proxy.ReactiveProxy, fieldName string) {
	entity.Watch(fieldName, "Get", func(fieldName string, oldValue, newValue interface{}) {
		reactivity.Track(entity.Original(), "Get")
	})
	entity.Watch(fieldName, "Set", func(fieldName string, oldValue, newValue interface{}) {
		reactivity.Trigger(entity.Original(), "Get")
	})
}

func NewMemoryMonitorReport() *MemoryMonitorReport {
	return &MemoryMonitorReport{
		AllocMB:     "0.00 MB",
		SysMB:       "0.00 MB",
		NumGC:       "0",
		Goroutines:  "0",
		HeapObjects: "0",
	}
}

func initUI(canvas *terminal.Terminal, width, height int) *components.Box {
	canvasWidth, canvasHeight := canvas.GetSize()
	canvasPaddingV := 0
	canvasPaddingH := 0
	return components.NewBox().
		SetBorder("normal").
		SetWidth(canvasWidth).
		SetHeight(canvasHeight).
		SetPadding(canvasPaddingV, canvasPaddingH)
}

func renderUI(rootComponent *components.Box) {
	fmt.Print(rootComponent.Draw())
}

//Отчет о потреблении памяти:
//Текущая:            0.00 MB
//Системная:          0.00 MB
//Сборок мусора:      0
//Горутин:            0
//Объектов в куче:    0

func main() {
	mm := monitor.NewMemoryMonitor(1000, 100)
	mm.Start()
	defer mm.Stop()

	// Инициализация терминала (холст)
	canvas := terminal.NewTerminal()
	canvas.OnResize(func(width, height int) {
		terminal.Clear()
		fmt.Println("Размер был изменен! Жопа!")
	})
	// Генерация root-компонента терминала
	terminal.Clear()
	rootComponent := initUI(canvas, 80, 24)
	rootComponent.AddChild("Отчет о потреблении памяти:")
	// Создаю массив значений которые требуются в отчете
	listItems := map[string]string{
		"Текущая:":         "0.00 MB",
		"Системная:":       "0.00 MB",
		"Сборок мусора:":   "0",
		"Горутин:":         "0",
		"Объектов в куче:": "0",
	}
	maxWidthKey := 0
	maxWidthValue := 0
	for key, value := range listItems {
		widthKey := utf8.RuneCountInString(key)
		if widthKey > maxWidthKey {
			maxWidthKey = widthKey
		}
		widthValue := utf8.RuneCountInString(value)
		if widthValue > maxWidthValue {
			maxWidthValue = widthValue
		}
	}

	for key, value := range listItems {
		rootComponent.AddChild(
			components.NewBox().
				SetWidth(rootComponent.GetTotalWidth()).
				SetLayout("row").
				SetAlign("left").
				SetGap(1).
				AddChildren(
					components.NewBox().
						SetWidth(maxWidthKey).
						SetAlign("left").
						AddChild(key),
					components.NewBox().
						SetID(key).
						SetWidth(maxWidthValue).
						SetAlign("left").
						AddChild(value),
				))
	}
	renderUI(rootComponent)
	rootComponent.CalculateAbsolutePositions(0, 0)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	stopChan := make(chan struct{})

	// Навешиваю прокси на отчет о потреблении памяти
	report := NewMemoryMonitorReport()
	proxyMemoryMonitorReport := proxy.NewReactiveProxy(report)

	// Создаем курсор для обеспечения точечного ререндера
	cursor := terminal.NewCursorManager()

	// Настраиваем наблюдатели
	createWatcher(proxyMemoryMonitorReport, "AllocMB")
	createWatcher(proxyMemoryMonitorReport, "SysMB")
	createWatcher(proxyMemoryMonitorReport, "NumGC")
	createWatcher(proxyMemoryMonitorReport, "Goroutines")
	createWatcher(proxyMemoryMonitorReport, "HeapObjects")

	reactivity.WatchEffect(func() {
		idMaps := map[string]string{
			"AllocMB":     "Текущая:",
			"SysMB":       "Системная:",
			"NumGC":       "Сборок мусора:",
			"Goroutines":  "Горутин:",
			"HeapObjects": "Объектов в куче:",
		}
		cursor.ShowCursor()
		for _, fieldName := range report.fieldNamesMemoryMonitor() {
			value := proxyMemoryMonitorReport.Get(fieldName)
			node := rootComponent.FindByID(idMaps[fieldName])
			if node != nil {
				col, row := node.GetAbsolutePosition()
				cursor.WriteAt(row, col, fmt.Sprintf("%v", value))
			}
		}
		cursor.HideCursor()
	})

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				a, s, n, g, h := mm.PrintCurrent()
				values := map[string]interface{}{
					"AllocMB":     a,
					"SysMB":       s,
					"NumGC":       n,
					"Goroutines":  g,
					"HeapObjects": h,
				}
				for key, value := range values {
					proxyMemoryMonitorReport.Set(key, value)
				}
			case <-stopChan:
				return
			}
		}
	}()

	<-sigChan
	close(stopChan)
	terminal.Clear()
	fmt.Println("Программа завершена.")
}
