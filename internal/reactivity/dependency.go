package reactivity

type Dep struct {
	ID          int
	Subscribers []*ReactiveEffect
	subsMap     map[*ReactiveEffect]bool
}

func NewDep() *Dep {
	depIdSeq++
	return &Dep{
		ID:          depIdSeq,
		Subscribers: make([]*ReactiveEffect, 0),
		subsMap:     make(map[*ReactiveEffect]bool),
	}
}

func (d *Dep) addSub(effect *ReactiveEffect) {
	if !d.subsMap[effect] {
		d.Subscribers = append(d.Subscribers, effect)
		d.subsMap[effect] = true
		effect.Deps = append(effect.Deps, d)
	}
}

func (d *Dep) removeSub(effect *ReactiveEffect) {
	if d.subsMap[effect] {
		delete(d.subsMap, effect)
		for i, sub := range d.Subscribers {
			if sub == effect {
				d.Subscribers = append(d.Subscribers[:i], d.Subscribers[i+1:]...)
				break
			}
		}
	}
}

func (d *Dep) notify() {
	for _, effect := range d.Subscribers {
		if effect.Active {
			effect.Run()
		}
	}
}

func (d *Dep) track() {
	if activeEffect != nil {
		d.addSub(activeEffect)
	}
}

func (d *Dep) trigger() {
	d.notify()
}
