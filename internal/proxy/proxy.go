package proxy

import (
	"fmt"
	"reflect"
)

// WatcherFunc Тип для функций-наблюдателей
type WatcherFunc func(fieldName string, oldValue, newValue interface{})

// ReactiveProxy Реактивный Proxy
type ReactiveProxy struct {
	target      interface{}
	getWatchers map[string][]WatcherFunc
	setWatchers map[string][]WatcherFunc
	history     []ChangeRecord
}

// ChangeRecord Запись об изменении
type ChangeRecord struct {
	Field    string
	OldValue interface{}
	NewValue interface{}
	Time     string
}

const MaxHistorySize = 2

// NewReactiveProxy Конструктор
func NewReactiveProxy(target interface{}) *ReactiveProxy {
	return &ReactiveProxy{
		target:      target,
		getWatchers: make(map[string][]WatcherFunc),
		setWatchers: make(map[string][]WatcherFunc),
		history:     make([]ChangeRecord, 0),
	}
}

// Уведомить всех наблюдателей
func (p *ReactiveProxy) notify(fieldName string, key string, oldValue, newValue interface{}) {
	switch key {
	case "Get":
		if watchers, exists := p.getWatchers[fieldName]; exists {
			for _, watcher := range watchers {
				watcher(fieldName, oldValue, newValue)
			}
		}
		break
	case "Set":
		if watchers, exists := p.setWatchers[fieldName]; exists {
			for _, watcher := range watchers {
				watcher(fieldName, oldValue, newValue)
			}
		}
		break
	default:
		fmt.Printf("Unknown watcher key: %s\n", key)
	}

	if oldValue != newValue {
		// Записываем в историю
		p.history = append(p.history, ChangeRecord{
			Field:    fieldName,
			OldValue: oldValue,
			NewValue: newValue,
			Time:     fmt.Sprintf("%d", len(p.history)+1),
		})

		if len(p.history) > MaxHistorySize {
			// Удаляем старые записи (FIFO)
			copy(p.history, p.history[1:])
			p.history = p.history[:MaxHistorySize]
		}
	}
}

// Original Получить оригинальную структуру
func (p *ReactiveProxy) Original() interface{} {
	return p.target
}

// Watch Добавить наблюдателя за полем
func (p *ReactiveProxy) Watch(fieldName string, key string, watcher WatcherFunc) {
	switch key {
	case "Get":
		if p.getWatchers[fieldName] == nil {
			p.getWatchers[fieldName] = make([]WatcherFunc, 0)
		}
		p.getWatchers[fieldName] = append(p.getWatchers[fieldName], watcher)
		break
	case "Set":
		if p.setWatchers[fieldName] == nil {
			p.setWatchers[fieldName] = make([]WatcherFunc, 0)
		}
		p.setWatchers[fieldName] = append(p.setWatchers[fieldName], watcher)
		break
	default:
		fmt.Printf("Unknown watcher key: %s\n", key)
	}

}

// Get Получить значение
func (p *ReactiveProxy) Get(fieldName string) interface{} {
	value := reflect.ValueOf(p.target).Elem()
	field := value.FieldByName(fieldName)

	if !field.IsValid() {
		return nil
	}

	// Уведомляем наблюдателей (старое значение совпадает с новым)
	p.notify(fieldName, "Get", field.Interface(), field.Interface())

	return field.Interface()
}

// Set Установить значение с уведомлением наблюдателей
func (p *ReactiveProxy) Set(fieldName string, newValue interface{}) {
	value := reflect.ValueOf(p.target).Elem()
	field := value.FieldByName(fieldName)

	if !field.IsValid() || !field.CanSet() {
		return
	}

	oldValue := field.Interface()

	// Проверяем, действительно ли значение изменилось
	if reflect.DeepEqual(oldValue, newValue) {
		return
	}

	// Устанавливаем новое значение
	field.Set(reflect.ValueOf(newValue))

	// Уведомляем наблюдателей
	p.notify(fieldName, "Set", oldValue, newValue)
}

// GetHistory Получить историю изменений
func (p *ReactiveProxy) GetHistory() []ChangeRecord {
	return p.history
}

// ClearHistory Очистить историю
func (p *ReactiveProxy) ClearHistory() {
	p.history = make([]ChangeRecord, 0)
}
