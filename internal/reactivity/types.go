package reactivity

// Общие типы и переменные
type EffectFunc func()

// Глобальное состояние
var (
	activeEffect *ReactiveEffect
	shouldTrack  = true
	effectIdSeq  = 0
	depIdSeq     = 0
	effectStack  []*ReactiveEffect
	targetMap    = make(TargetMap)
)

// Карта зависимостей
type TargetMap map[uintptr]map[string]*Dep
