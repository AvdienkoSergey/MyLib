package reactivity

import (
	"reflect"
)

func getObjectID(obj interface{}) uintptr {
	return reflect.ValueOf(obj).Pointer()
}

func Track(target interface{}, key string) {
	if activeEffect == nil {
		return
	}

	objectID := getObjectID(target)

	depsMap, exists := targetMap[objectID]
	if !exists {
		depsMap = make(map[string]*Dep)
		targetMap[objectID] = depsMap
	}

	dep, exists := depsMap[key]
	if !exists {
		dep = NewDep()
		depsMap[key] = dep
	}

	dep.track()
}

func Trigger(target interface{}, key string) {
	objectID := getObjectID(target)

	depsMap, exists := targetMap[objectID]
	if !exists {
		return
	}

	dep, exists := depsMap[key]
	if !exists {
		return
	}

	dep.trigger()
}
