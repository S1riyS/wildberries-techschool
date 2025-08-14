package main

import (
	"fmt"
	"sync"
)

// SafeMap is a thread-safe map with generic types
type SafeMap[U comparable, V any] struct {
	mu   sync.RWMutex
	data map[U]V
}

func NewSafeMap[U comparable, V any]() *SafeMap[U, V] {
	return &SafeMap[U, V]{
		data: make(map[U]V),
	}
}

func (m *SafeMap[U, V]) Set(key U, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *SafeMap[U, V]) Get(key U) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	return val, ok
}

func main() {
	safeMap := NewSafeMap[string, int]()
	var wg sync.WaitGroup

	valuesCount := 100
	for i := range valuesCount {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", i)
			fmt.Printf("Write %s: %v\n", key, i)
			safeMap.Set(key, i)
		}(i)
	}

	// ! NOTE:
	// Could have used wg.Wait() here in order to read from fully populated map.
	// However for the sake of example I will leave a possibility to "race condition".

	for i := range valuesCount {
		key := fmt.Sprintf("key_%d", i)
		if val, ok := safeMap.Get(key); ok {
			fmt.Printf("Got %s: %v\n", key, val)
		} else {
			fmt.Printf("%s not found\n", key)
		}
	}

	wg.Wait()
}
