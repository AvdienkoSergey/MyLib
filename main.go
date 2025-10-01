package main

import (
	"Guess/internal/monitor"
	"Guess/internal/proxy"
	"Guess/internal/reactivity"
	"Guess/internal/renderer"
	"Guess/internal/ui/terminal"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
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

func NewMemoryMonitorReport() *MemoryMonitorReport {
	return &MemoryMonitorReport{
		AllocMB:     "0.00 MB",
		SysMB:       "0.00 MB",
		NumGC:       "0",
		Goroutines:  "0",
		HeapObjects: "0",
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

func main() {
	mm := monitor.NewMemoryMonitor(1000, 100)
	mm.Start()
	defer mm.Stop()

	terminal.Clear() // Очищаем экран
	renderer.
		NewDrawTask().
		SetContent([]string{
			"┌──────────────────────────────┐",
			"│ Отчет о потреблении памяти:  │",
			"└──────────────────────────────┘",
		}).
		SetPosition(1, 1).
		SetColorSchema("yellow", "").
		SetAutoSize().
		Draw()
	renderer.
		NewDrawTask().
		SetContent([]string{
			"┌───────────────────┐",
			"│ Текущая:          │",
			"└───────────────────┘",
		}).
		SetPosition(1, 3).
		SetColorSchema("yellow", "").
		SetAutoSize().
		Draw()
	renderer.
		NewDrawTask().
		SetContent([]string{
			"┌──────────┐",
			"│ 0.00 MB  │",
			"└──────────┘",
		}).
		SetPosition(21, 3).
		SetColorSchema("yellow", "").
		SetAutoSize().
		Draw()
	renderer.
		NewDrawTask().
		SetContent([]string{
			"┌───────────────────┐",
			"│ Системная:        │",
			"└───────────────────┘",
		}).
		SetPosition(1, 5).
		SetColorSchema("yellow", "").
		SetAutoSize().
		Draw()
	renderer.
		NewDrawTask().
		SetContent([]string{
			"┌──────────┐",
			"│ 0.00 MB  │",
			"└──────────┘",
		}).
		SetPosition(21, 5).
		SetColorSchema("yellow", "").
		SetAutoSize().
		Draw()
	renderer.
		NewDrawTask().
		SetContent([]string{
			"┌───────────────────┐",
			"│ Сборок мусора:    │",
			"└───────────────────┘",
		}).
		SetPosition(1, 7).
		SetColorSchema("yellow", "").
		SetAutoSize().
		Draw()
	renderer.
		NewDrawTask().
		SetContent([]string{
			"┌──────────┐",
			"│ 0        │",
			"└──────────┘",
		}).
		SetPosition(21, 7).
		SetColorSchema("yellow", "").
		SetAutoSize().
		Draw()
	renderer.
		NewDrawTask().
		SetContent([]string{
			"┌───────────────────┐",
			"│ Горутин:          │",
			"└───────────────────┘",
		}).
		SetPosition(1, 9).
		SetColorSchema("yellow", "").
		SetAutoSize().
		Draw()
	renderer.
		NewDrawTask().
		SetContent([]string{
			"┌──────────┐",
			"│ 0        │",
			"└──────────┘",
		}).
		SetPosition(21, 9).
		SetColorSchema("yellow", "").
		SetAutoSize().
		Draw()
	renderer.
		NewDrawTask().
		SetContent([]string{
			"┌───────────────────┐",
			"│ Объектов в куче:  │",
			"└───────────────────┘",
		}).
		SetPosition(1, 11).
		SetColorSchema("yellow", "").
		SetAutoSize().
		Draw()
	renderer.
		NewDrawTask().
		SetContent([]string{
			"┌──────────┐",
			"│ 0        │",
			"└──────────┘",
		}).
		SetPosition(21, 11).
		SetColorSchema("yellow", "").
		SetAutoSize().
		Draw()

	// Создаем дерево для шапки сайта
	//renderer.DemoMain()

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
		posMaps := map[string]renderer.Position{
			"AllocMB":     {X: 23, Y: 4},
			"SysMB":       {X: 23, Y: 6},
			"NumGC":       {X: 23, Y: 8},
			"Goroutines":  {X: 23, Y: 10},
			"HeapObjects": {X: 23, Y: 12},
		}
		cursor.ShowCursor()
		for _, fieldName := range report.fieldNamesMemoryMonitor() {
			value := proxyMemoryMonitorReport.Get(fieldName)
			pos := posMaps[fieldName]
			cursor.WriteAt(pos.X, pos.Y, fmt.Sprintf("%v", value))
		}
		cursor.HideCursor()
	})

	go func() {
		ticker := time.NewTicker(1 * time.Second)
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
