// Package iterator provides a simple iterator for Go slices, making use of generics for a fluent and type-safe API.
//
// Ideally, only slices of scalar types or slices of structs composed of scalar types should be used with this package. If
// you need to iterate over a slice of channels, for example, you should probably just use a for loop. This package is not
// intended to be a replacement for a for loop in every situation, but it can be useful for certain cases where multiple
// operations need to be performed on each element of a slice, but you want to express those operations in a more
// declarative way.
package iterator

import (
	"reflect"
	"sync"
)

type maybe[T any] struct {
	ok  bool // whether the current element should be included in the result
	val T    // the value of the current element
}

type iter[T any] struct {
	mu         sync.Mutex               // mutex to synchronize access to the iterator when the ThreadSafe option is used
	nextFunc   func(*iter[T]) (T, bool) // the function to be used when calling the Next method. This is set to next or synchronizedNext depending on the options used when creating the iterator.
	nextIndex  int                      // the index of the next element to be returned by the Next method
	source     []T                      // the source slice. Could be the original slice or a copy, depending on the options used when creating the iterator.
	operations []func(*maybe[T])        // the operations to be performed on each element of the source slice
}

// From returns a new iterator for the given source. There are several options that can be used to configure the
// behavior of the iterator. See the documentation for the FromOption type for more information.
func From[T any](source []T, opts ...FromOption) Of[T] {
	it := &iter[T]{
		source: source,
	}
	options := new(fromOptions)
	options.bufferLen = 64
	for _, opt := range opts {
		opt(options)
	}
	it.nextFunc = next[T]
	if options.copySource {
		it.source = make([]T, len(source))
		copy(it.source, source)
	}
	if options.threadSafe {
		it.nextFunc = synchronizedNext[T]
	}
	it.operations = make([]func(*maybe[T]), 0, options.bufferLen)
	return it
}

func next[T any](it *iter[T]) (T, bool) {
	if it.nextIndex >= len(it.source) {
		return *new(T), false
	}
	defer func() { it.nextIndex++ }()
	return it.source[it.nextIndex], true
}

func synchronizedNext[T any](it *iter[T]) (T, bool) {
	it.mu.Lock()
	defer it.mu.Unlock()
	return next(it)
}

func (it *iter[T]) Next() (T, bool) {
	return it.nextFunc(it)
}

func (it *iter[T]) ForEach(fn func(T)) {
	for {
		val, ok := it.Next()
		if !ok {
			break
		}
		fn(val)
	}
}

func (it *iter[T]) Map(fn func(T) T) Of[T] {
	it.operations = append(it.operations, func(m *maybe[T]) {
		m.val = fn(m.val)
	})
	return it
}

func (it *iter[T]) Filter(fn func(T) bool) Of[T] {
	it.operations = append(it.operations, func(m *maybe[T]) {
		m.ok = fn(m.val)
	})
	return it
}

func (it *iter[T]) Unique(opts ...UniqueOption) Of[T] {
	options := new(uniqueOptions)
	for _, opt := range opts {
		opt(options)
	}
	seen := make(map[any]struct{}, len(it.source)) // pre-allocate a map with the same size as the source slice to avoid reallocations
	filterFn := func(val T) bool {
		if _, ok := seen[val]; ok {
			return false
		}
		seen[val] = struct{}{}
		return true
	}
	if options.deref && reflect.TypeOf(*new(T)).Kind() == reflect.Ptr { // if we're dereferencing pointers AND the type of T is a pointer
		filterFn = func(val T) bool { // redefine the filterFn to dereference the pointer before checking for uniqueness
			v := reflect.ValueOf(val).Elem().Interface()
			if _, ok := seen[v]; ok {
				return false
			}
			seen[v] = struct{}{}
			return true
		}
	}
	return it.Filter(filterFn)
}

func (it *iter[T]) Collect() []T {
	result := make([]T, 0, len(it.source))
	mb := new(maybe[T]) // create a single maybe object to be reused for each iteration, preventing unnecessary allocations
	it.ForEach(func(val T) {
		mb.val = val
		mb.ok = true
		for _, op := range it.operations {
			op(mb)
			if !mb.ok {
				break
			}
		}
		if mb.ok {
			result = append(result, mb.val)
		}
		mb.ok = false
	})
	return result
}

func (it *iter[T]) Channel() <-chan T {
	ch := make(chan T, len(it.source))
	it.IntoChannel(ch, CloseChannel(true))
	return ch
}

func (it *iter[T]) IntoChannel(ch chan<- T, opts ...IntoChannelOption) {
	icos := new(intoChannelOptions)
	for _, opt := range opts {
		opt(icos)
	}
	go func() {
		if icos.closeChannel {
			defer close(ch)
		}
		it.ForEach(func(val T) {
			ch <- val
		})
	}()
}

func (it *iter[T]) CollectChannel() <-chan T {
	ch := make(chan T, len(it.source))
	it.CollectIntoChannel(ch, CloseChannel(true))
	return ch
}

func (it *iter[T]) CollectIntoChannel(ch chan<- T, opts ...IntoChannelOption) {
	icos := new(intoChannelOptions)
	for _, opt := range opts {
		opt(icos)
	}
	go func() {
		if icos.closeChannel {
			defer close(ch)
		}
		for _, val := range it.Collect() {
			ch <- val
		}
	}()
}

func (it *iter[T]) Reduce(fn func(acc T, next T) T, initial T) T {
	result := initial
	it.ForEach(func(val T) {
		result = fn(result, val)
	})
	return result
}

func (it *iter[T]) Reset() {
	it.nextIndex = 0
}
