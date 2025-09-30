package proxy

import (
	"testing"
)

type TestPerson struct {
	Name string
	Age  int
}

// TestNewReactiveProxy проверяем, что конструктор работает правильно
func TestNewReactiveProxy(t *testing.T) {
	person := &TestPerson{Name: "Вася", Age: 25}
	proxy := NewReactiveProxy(person)

	// Проверяем, что proxy создался не пустой
	if proxy == nil {
		t.Error("Proxy не должен быть nil")
		return
	}

	// Проверяем, что наблюдатели инициализированы
	if proxy.getWatchers == nil {
		t.Error("getWatchers должны быть инициализированы")
	}

	if proxy.setWatchers == nil {
		t.Error("setWatchers должны быть инициализированы")
	}

	// Проверяем, что история пустая
	if len(proxy.history) != 0 {
		t.Error("История должна быть пустой при создании")
	}
}

// TestGet проверяем получение значений
func TestGet(t *testing.T) {
	person := &TestPerson{Name: "Маша", Age: 30}
	proxy := NewReactiveProxy(person)

	// Проверяем получение имени
	name := proxy.Get("Name")
	if name != "Маша" {
		t.Errorf("Ожидали 'Маша', получили '%v'", name)
	}

	// Проверяем получение возраста
	age := proxy.Get("Age")
	if age != 30 {
		t.Errorf("Ожидали 30, получили %v", age)
	}

	// Проверяем несуществующее поле
	nothing := proxy.Get("NotExists")
	if nothing != nil {
		t.Error("Несуществующее поле должно возвращать nil")
	}
}

// TestSet проверяем установку значений
func TestSet(t *testing.T) {
	person := &TestPerson{Name: "Петя", Age: 20}
	proxy := NewReactiveProxy(person)

	// Меняем имя
	proxy.Set("Name", "Коля")

	// Проверяем, что значение изменилось
	newName := proxy.Get("Name")
	if newName != "Коля" {
		t.Errorf("Ожидали 'Коля', получили '%v'", newName)
	}

	// Проверяем, что в исходной структуре тоже изменилось
	if person.Name != "Коля" {
		t.Error("Значение в исходной структуре должно измениться")
	}
}

// TestWatch проверяем работу наблюдателей
func TestWatch(t *testing.T) {
	person := &TestPerson{Name: "Аня", Age: 28}
	proxy := NewReactiveProxy(person)

	// Переменные для проверки вызова наблюдателя
	var watchedField string
	var watchedOldValue interface{}
	var watchedNewValue interface{}
	wasCalled := false

	// Добавляем наблюдателя
	proxy.Watch("Name", "Set", func(fieldName string, oldValue, newValue interface{}) {
		watchedField = fieldName
		watchedOldValue = oldValue
		watchedNewValue = newValue
		wasCalled = true
	})

	// Меняем значение
	proxy.Set("Name", "Катя")

	// Проверяем, что наблюдатель был вызван
	if !wasCalled {
		t.Error("Наблюдатель должен был быть вызван")
	}

	if watchedField != "Name" {
		t.Errorf("Ожидали поле 'Name', получили '%s'", watchedField)
	}

	if watchedOldValue != "Аня" {
		t.Errorf("Ожидали старое значение 'Аня', получили '%v'", watchedOldValue)
	}

	if watchedNewValue != "Катя" {
		t.Errorf("Ожидали новое значение 'Катя', получили '%v'", watchedNewValue)
	}
}

// TestHistory проверяем историю изменений
func TestHistory(t *testing.T) {
	person := &TestPerson{Name: "Лена", Age: 22}
	proxy := NewReactiveProxy(person)

	// Делаем несколько изменений
	proxy.Set("Name", "Оля")
	proxy.Set("Age", 25)
	proxy.Set("Name", "Света")
	proxy.Get("Name")
	proxy.Get("Age")

	// Проверяем историю
	history := proxy.GetHistory()
	if len(history) != 5 {
		t.Errorf("Ожидали 5 записи в истории, получили %d", len(history))
	}

	// Проверяем первую запись
	if history[0].Field != "Name" || history[0].OldValue != "Лена" || history[0].NewValue != "Оля" {
		t.Error("Первая запись в истории неверная")
	}

	// Проверяем последнюю запись
	if history[4].Field != "Age" || history[4].OldValue != 25 {
		t.Error("Последняя запись в истории неверная")
	}

	// Очищаем историю
	proxy.ClearHistory()
	history = proxy.GetHistory()
	if len(history) != 0 {
		t.Error("История должна быть пустой после очистки")
	}
}

// TestGetSameValue проверяем, что при попытке Get будет вызвано уведомление
func TestGetSameValue(t *testing.T) {
	person := &TestPerson{Name: "Рома", Age: 33}
	proxy := NewReactiveProxy(person)

	wasCalled := false
	proxy.Watch("Name", "Get", func(fieldName string, oldValue, newValue interface{}) {
		wasCalled = true
	})

	proxy.Get("Name")

	if !wasCalled {
		t.Error("Наблюдатель должен вызываться при наличии поля")
	}

	// История должна быть пустой
	history := proxy.GetHistory()
	if len(history) == 0 {
		t.Error("История не должна быть пустой")
	}
}

// TestSetSameValue проверяем, что одинаковые значения не вызывают уведомления при попытке Set
func TestSetSameValue(t *testing.T) {
	person := &TestPerson{Name: "Рома", Age: 33}
	proxy := NewReactiveProxy(person)

	wasCalled := false
	proxy.Watch("Name", "Set", func(fieldName string, oldValue, newValue interface{}) {
		wasCalled = true
	})

	// Устанавливаем то же самое значение
	proxy.Set("Name", "Рома")

	// Наблюдатель не должен быть вызван
	if wasCalled {
		t.Error("Наблюдатель не должен вызываться при установке того же значения")
	}

	// История должна быть пустой
	history := proxy.GetHistory()
	if len(history) != 0 {
		t.Error("История должна быть пустой при установке того же значения")
	}
}
