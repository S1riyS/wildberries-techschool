package datastructure

// SlidingWindow represents a sliding window over a slice with iteration support.
type SlidingWindow[T any] struct {
	data  []T
	start int
	size  int
	count int // current number of elements
}

// NewSlidingWindow creates a new sliding window with the specified size.
func NewSlidingWindow[T any](size int) *SlidingWindow[T] {
	if size < 0 {
		panic("size cannot be negative")
	}

	// Create a nil slice for zero-sized windows
	var data []T
	if size > 0 {
		data = make([]T, size)
	}

	return &SlidingWindow[T]{
		data:  data,
		size:  size,
		start: 0,
		count: 0,
	}
}

// Add adds an element to the window, sliding if necessary.
func (sw *SlidingWindow[T]) Add(item T) {
	if sw.size == 0 {
		// Window of size 0 cannot hold any elements
		return
	}

	if sw.count < sw.size {
		// Window not full yet, add to next position
		sw.data[sw.count] = item
		sw.count++
	} else {
		// Window is full, slide by replacing the oldest element
		sw.data[sw.start] = item
		sw.start = (sw.start + 1) % sw.size
	}
}

// At returns the element at the specified index (0-based from oldest to newest).
func (sw *SlidingWindow[T]) At(index int) T {
	if sw.size == 0 {
		panic("window of size 0 has no elements")
	}
	if index < 0 || index >= sw.count {
		panic("index out of range")
	}
	return sw.data[(sw.start+index)%sw.size]
}

// Iterator returns a channel for iterating over window elements.
func (sw *SlidingWindow[T]) Iterator() <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for i := range sw.count {
			ch <- sw.At(i)
		}
	}()
	return ch
}
