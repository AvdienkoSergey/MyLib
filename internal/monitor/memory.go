package monitor

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type MemoryStats struct {
	Timestamp    time.Time
	AllocMB      float64 // Текущая аллоцированная память
	TotalAllocMB float64 // Общая аллоцированная память за время работы
	SysMB        float64 // Системная память
	NumGC        uint32  // Количество сборок мусора
	Goroutines   int     // Количество горутин
	HeapObjects  uint64  // Количество объектов в куче
}

type MemoryMonitor struct {
	stats      []MemoryStats
	interval   time.Duration
	running    bool
	stopChan   chan bool
	mutex      sync.RWMutex
	maxRecords int
}

func NewMemoryMonitor(interval time.Duration, maxRecords int) *MemoryMonitor {
	return &MemoryMonitor{
		stats:      make([]MemoryStats, 0),
		interval:   interval,
		stopChan:   make(chan bool),
		maxRecords: maxRecords,
	}
}

// Start Запуск мониторинга
func (m *MemoryMonitor) Start() {
	if m.running {
		return
	}

	m.running = true

	go func() {
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.collectStats()
			case <-m.stopChan:
				return
			}
		}
	}()
}

// Stop Остановка мониторинга
func (m *MemoryMonitor) Stop() {
	if !m.running {
		return
	}

	m.running = false
	m.stopChan <- true
}

// Сбор статистики
func (m *MemoryMonitor) collectStats() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	stats := MemoryStats{
		Timestamp:    time.Now(),
		AllocMB:      float64(memStats.Alloc) / 1024 / 1024,
		TotalAllocMB: float64(memStats.TotalAlloc) / 1024 / 1024,
		SysMB:        float64(memStats.Sys) / 1024 / 1024,
		NumGC:        memStats.NumGC,
		Goroutines:   runtime.NumGoroutine(),
		HeapObjects:  memStats.HeapObjects,
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.stats = append(m.stats, stats)

	// Ограничиваем количество записей
	if len(m.stats) > m.maxRecords {
		m.stats = m.stats[1:]
	}
}

// GetCurrentStats Получить текущую статистику
func (m *MemoryMonitor) GetCurrentStats() MemoryStats {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.stats) == 0 {
		return MemoryStats{}
	}

	return m.stats[len(m.stats)-1]
}

// GetAllStats Получить все записи
func (m *MemoryMonitor) GetAllStats() []MemoryStats {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Возвращаем копию
	result := make([]MemoryStats, len(m.stats))
	copy(result, m.stats)
	return result
}

// PrintCurrent Показать текущее состояние
func (m *MemoryMonitor) PrintCurrent() (AllocMB, SysMB, NumGC, Goroutines, HeapObjects string) {
	current := m.GetCurrentStats()
	if current.Timestamp.IsZero() {
		return "0.00 MB", "0.00 MB", "0", "0", "0"
	}

	allocMB := fmt.Sprintf(`%.2f MB`, current.AllocMB)
	sysMB := fmt.Sprintf(`%.2f MB`, current.SysMB)
	numGC := fmt.Sprintf(`%d`, current.NumGC)
	goroutines := fmt.Sprintf(`%d`, current.Goroutines)
	heapObjects := fmt.Sprintf(`%d`, current.HeapObjects)

	return allocMB, sysMB, numGC, goroutines, heapObjects
}

// PrintDelta Показать изменения за период
func (m *MemoryMonitor) PrintDelta(seconds int) {
	all := m.GetAllStats()
	if len(all) < 2 {
		fmt.Println("📊 Недостаточно данных для сравнения")
		return
	}

	now := time.Now()
	var start, end MemoryStats

	// Находим записи за указанный период
	for i := len(all) - 1; i >= 0; i-- {
		if now.Sub(all[i].Timestamp).Seconds() <= float64(seconds) {
			if end.Timestamp.IsZero() {
				end = all[i]
			}
			start = all[i]
		}
	}

	if start.Timestamp.IsZero() || end.Timestamp.IsZero() {
		fmt.Printf("📊 Недостаточно данных за %d секунд\n", seconds)
		return
	}

	allocDelta := end.AllocMB - start.AllocMB
	gcDelta := int(end.NumGC) - int(start.NumGC)
	goroutineDelta := end.Goroutines - start.Goroutines
	objectsDelta := int64(end.HeapObjects) - int64(start.HeapObjects)

	fmt.Printf(`📊 Изменения за %d сек:
  💾 Память: %.2f MB (%+.2f)
  🗑️ Сборок мусора: +%d
  🔄 Горутин: %d (%+d)
  📦 Объектов: %d (%+d)
`,
		seconds, end.AllocMB, allocDelta,
		gcDelta, end.Goroutines, goroutineDelta,
		end.HeapObjects, objectsDelta)
}

// ForceGCAndMeasure Принудительная сборка мусора с измерением
func (m *MemoryMonitor) ForceGCAndMeasure() (before, after MemoryStats) {
	// Замер до GC
	m.collectStats()
	before = m.GetCurrentStats()

	fmt.Println("🗑️ Принудительная сборка мусора...")
	runtime.GC()
	runtime.GC() // Двойной вызов для полной очистки

	time.Sleep(100 * time.Millisecond) // Даём время на очистку

	// Замер после GC
	m.collectStats()
	after = m.GetCurrentStats()

	freed := before.AllocMB - after.AllocMB
	fmt.Printf("🧹 Освобождено: %.2f MB (было: %.2f MB, стало: %.2f MB)\n",
		freed, before.AllocMB, after.AllocMB)

	return before, after
}
