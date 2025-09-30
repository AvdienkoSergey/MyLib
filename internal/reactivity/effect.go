package reactivity

type ReactiveEffect struct {
	ID     int
	Fn     EffectFunc
	Deps   []*Dep
	Active bool
	OnStop []func()
}

func NewReactiveEffect(fn EffectFunc) *ReactiveEffect {
	effectIdSeq++
	return &ReactiveEffect{
		ID:     effectIdSeq,
		Fn:     fn,
		Deps:   make([]*Dep, 0),
		Active: true,
	}
}

func (e *ReactiveEffect) Run() {
	if !e.Active {
		return
	}

	prevEffect := activeEffect
	effectStack = append(effectStack, e)

	defer func() {
		effectStack = effectStack[:len(effectStack)-1]
		activeEffect = prevEffect
	}()

	activeEffect = e
	e.cleanupDeps()
	e.Fn()
}

func (e *ReactiveEffect) cleanupDeps() {
	for _, dep := range e.Deps {
		dep.removeSub(e)
	}
	e.Deps = e.Deps[:0]
}

func (e *ReactiveEffect) Stop() {
	if e.Active {
		e.cleanupDeps()
		e.Active = false
		for _, onStop := range e.OnStop {
			onStop()
		}
	}
}

func WatchEffect(fn EffectFunc) *ReactiveEffect {
	effect := NewReactiveEffect(fn)
	effect.Run()
	return effect
}
