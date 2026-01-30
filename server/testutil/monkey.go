package testutil

import (
	"reflect"
	"sync"

	"github.com/agiledragon/gomonkey/v2"
)

// Drop-in compatibility layer for bou.ke/monkey using gomonkey/v2.
// gomonkey/v2 supports darwin/arm64 (Apple Silicon) whereas bou.ke/monkey does not.

var (
	mu       sync.Mutex
	patchMap = make(map[uintptr]*gomonkey.Patches)
)

func Patch(target, replacement interface{}) {
	mu.Lock()
	defer mu.Unlock()
	ptr := reflect.ValueOf(target).Pointer()
	// Reset existing patch before applying new one to avoid leaking stale patches
	if old, ok := patchMap[ptr]; ok {
		old.Reset()
	}
	p := gomonkey.ApplyFunc(target, replacement)
	patchMap[ptr] = p
}

func PatchInstanceMethod(target reflect.Type, methodName string, replacement interface{}) {
	mu.Lock()
	defer mu.Unlock()
	p := gomonkey.ApplyMethod(target, methodName, replacement)
	key := reflect.ValueOf(methodName).Pointer() ^ target.Size()
	patchMap[key] = p
}

func Unpatch(target interface{}) {
	mu.Lock()
	defer mu.Unlock()
	ptr := reflect.ValueOf(target).Pointer()
	if p, ok := patchMap[ptr]; ok {
		p.Reset()
		delete(patchMap, ptr)
	}
}

func UnpatchAll() {
	mu.Lock()
	defer mu.Unlock()
	for k, p := range patchMap {
		p.Reset()
		delete(patchMap, k)
	}
}
