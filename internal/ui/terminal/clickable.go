package terminal

// ClickableArea - область, которая может обрабатывать клики
type ClickableArea struct {
	X, Y          int
	Width, Height int
	OnClick       func()
	ID            string
}

// ClickManager - менеджер для отслеживания кликабельных областей
type ClickManager struct {
	areas []*ClickableArea
	term  Terminal
}

var globalManager *ClickManager

// Init инициализирует глобальный менеджер
func ClickableAreasInit() {
	if globalManager == nil {
		globalManager = &ClickManager{
			areas: make([]*ClickableArea, 0),
		}
	}
}

// Register регистрирует кликабельную область
func ClickableAreaRegister(x, y, width, height int, id string, onClick func()) *ClickableArea {
	if globalManager == nil {
		ClickableAreasInit()
	}

	area := &ClickableArea{
		X:       x,
		Y:       y,
		Width:   width,
		Height:  height,
		OnClick: onClick,
		ID:      id,
	}

	globalManager.areas = append(globalManager.areas, area)
	return area
}

// Unregister удаляет область по ID
func Unregister(id string) {
	if globalManager == nil {
		return
	}

	for i, area := range globalManager.areas {
		if area.ID == id {
			globalManager.areas = append(globalManager.areas[:i], globalManager.areas[i+1:]...)
			break
		}
	}
}

// HandleClick обрабатывает клик по координатам
func HandleClick(x, y int) bool {
	if globalManager == nil {
		return false
	}

	// Проверяем области в обратном порядке (последние зарегистрированные имеют приоритет)
	for i := len(globalManager.areas) - 1; i >= 0; i-- {
		area := globalManager.areas[i]
		if isInside(x, y, area) {
			if area.OnClick != nil {
				area.OnClick()
			}
			return true
		}
	}
	return false
}

// isInside проверяет, находится ли точка внутри области
func isInside(x, y int, area *ClickableArea) bool {
	return x >= area.X && x < area.X+area.Width &&
		y >= area.Y && y < area.Y+area.Height
}

// Clear очищает все зарегистрированные области
func ClearClickableArea() {
	if globalManager != nil {
		globalManager.areas = make([]*ClickableArea, 0)
	}
}

// GetAreas возвращает все зарегистрированные области (для отладки)
func GetAreas() []*ClickableArea {
	if globalManager == nil {
		return []*ClickableArea{}
	}
	return globalManager.areas
}
