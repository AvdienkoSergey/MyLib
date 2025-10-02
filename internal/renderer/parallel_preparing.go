package renderer

import (
	"runtime"
	"strings"
	"sync"
)

// PreparedCommand - подготовленная команда для отрисовки
type PreparedCommand struct {
	X, Y int
	Text string
}

// ParallelProcessor - процессор для параллельной подготовки данных
type ParallelProcessor struct {
	workers int
}

var (
	globalProcessor     *ParallelProcessor
	globalProcessorOnce sync.Once
)

// GetGlobalProcessor - получить глобальный процессор
func GetGlobalProcessor() *ParallelProcessor {
	globalProcessorOnce.Do(func() {
		globalProcessor = NewParallelProcessor(runtime.NumCPU())
	})
	return globalProcessor
}

// NewParallelProcessor - создать новый процессор
func NewParallelProcessor(workers int) *ParallelProcessor {
	if workers < 1 {
		workers = 1
	}
	return &ParallelProcessor{workers: workers}
}

// ProcessTask - обработать одну задачу параллельно (парсинг строк)
func (p *ParallelProcessor) ProcessTask(task *DrawTask) []PreparedCommand {
	contentLen := len(task.Content)
	
	// Если строк мало - обрабатываем последовательно
	if contentLen < 3 {
		return p.processSequential(task)
	}
	
	// Параллельная обработка строк
	type lineResult struct {
		index    int
		commands []PreparedCommand
	}
	
	resultsChan := make(chan lineResult, contentLen)
	var wg sync.WaitGroup
	
	// Обрабатываем каждую строку в отдельной горутине
	for i, str := range task.Content {
		wg.Add(1)
		go func(lineIndex int, line string) {
			defer wg.Done()
			
			y := task.Position.Y + lineIndex
			commands := p.processLine(line, task.Position.X, y)
			
			resultsChan <- lineResult{
				index:    lineIndex,
				commands: commands,
			}
		}(i, str)
	}
	
	// Закрываем канал после обработки всех строк
	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	
	// Собираем результаты в правильном порядке
	results := make([][]PreparedCommand, contentLen)
	for result := range resultsChan {
		results[result.index] = result.commands
	}
	
	// Объединяем все команды в один слайс
	var allCommands []PreparedCommand
	for _, commands := range results {
		allCommands = append(allCommands, commands...)
	}
	
	return allCommands
}

// processSequential - последовательная обработка (для малых задач)
func (p *ParallelProcessor) processSequential(task *DrawTask) []PreparedCommand {
	var allCommands []PreparedCommand
	
	for i, str := range task.Content {
		y := task.Position.Y + i
		commands := p.processLine(str, task.Position.X, y)
		allCommands = append(allCommands, commands...)
	}
	
	return allCommands
}

// processLine - обработать одну строку (извлечь токены)
func (p *ParallelProcessor) processLine(line string, baseX, y int) []PreparedCommand {
	var commands []PreparedCommand
	runes := []rune(line)
	
	var token strings.Builder
	tokenStartX := -1
	
	for i, r := range runes {
		if r == ' ' {
			if token.Len() > 0 {
				commands = append(commands, PreparedCommand{
					X:    baseX + tokenStartX,
					Y:    y,
					Text: token.String(),
				})
				token.Reset()
				tokenStartX = -1
			}
		} else {
			if tokenStartX == -1 {
				tokenStartX = i
			}
			token.WriteRune(r)
		}
	}
	
	// Не забываем последний токен
	if token.Len() > 0 {
		commands = append(commands, PreparedCommand{
			X:    baseX + tokenStartX,
			Y:    y,
			Text: token.String(),
		})
	}
	
	return commands
}

// ProcessBatch - обработать батч задач параллельно
func (p *ParallelProcessor) ProcessBatch(tasks []*DrawTask) []PreparedCommand {
	if len(tasks) == 0 {
		return nil
	}
	
	// Если задач мало - обрабатываем последовательно
	if len(tasks) == 1 {
		return p.ProcessTask(tasks[0])
	}
	
	type taskResult struct {
		index    int
		commands []PreparedCommand
	}
	
	resultsChan := make(chan taskResult, len(tasks))
	var wg sync.WaitGroup
	
	// Обрабатываем каждую задачу параллельно
	for i, task := range tasks {
		wg.Add(1)
		go func(taskIndex int, t *DrawTask) {
			defer wg.Done()
			
			commands := p.ProcessTask(t)
			resultsChan <- taskResult{
				index:    taskIndex,
				commands: commands,
			}
		}(i, task)
	}
	
	// Закрываем канал после обработки
	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	
	// Собираем результаты в правильном порядке
	results := make([][]PreparedCommand, len(tasks))
	for result := range resultsChan {
		results[result.index] = result.commands
	}
	
	// Объединяем все команды
	var allCommands []PreparedCommand
	for _, commands := range results {
		allCommands = append(allCommands, commands...)
	}
	
	return allCommands
}