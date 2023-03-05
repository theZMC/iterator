package iterator

// Of provides a high-level interface for iterating over a slice.
type Of[T any] interface {
	// Next returns the next value in the iterator, consuming it in the process, as well as a boolean indicating whether
	// there was a value to return. If there was no value to return, the returned value will be the zero value for the type.
	Next() (T, bool)
	// ForEach iterates over the iterator, calling the given function for each value and consuming the iterator. Calls to
	// ForEach are thread-safe as long as the provided function is thread-safe, but it's not recommended to call ForEach
	// from multiple goroutines, especially if the provided function is different for each goroutine.
	ForEach(fn func(T))
	// Map returns a new iterator that applies the given function to each value in the iterator upon collection. The function
	// is lazily evaluated, so it is not applied until the iterator is collected.
	Map(fn func(T) T) Of[T]
	// Filter returns a new iterator that keeps only the values in the iterator that return true when passed to the given
	// function upon collection. The function is lazily evaluated, so it is not applied until the iterator is collected.
	Filter(fn func(T) bool) Of[T]
	// Unique returns a new iterator that filters out duplicate values in the iterator upon collection. The function is
	// lazily evaluated, so it is not applied until the iterator is collected. This is a convenience method that is equivalent
	// to calling Filter with a function that keeps track of the values it has seen. If the iterator contains pointers, the
	// DerefPointers option can be used to dereference the pointers before evaluating uniqueness.
	Unique(opts ...UniqueOption) Of[T]
	// Collect applies all of the chained map and filter operations to the iterator and returns the resulting slice.
	Collect() []T
	// Channel returns a channel that will be populated with the values in the iterator. The channel will be closed when
	// there are no more values, indicating that the iterator has been consumed. This is not the same as collecting, as
	// this does not apply the chained map and filter operations to each element. If you want a channel that applies the
	// chained map and filter operations, use CollectChannel.
	Channel() <-chan T
	// IntoChannel populates the given channel with the values in the iterator. If shouldClose is true, the channel will be
	// closed when there are no more values, indicating that the iterator has been consumed. This is not the same as
	// collecting, as this does not apply the chained map and filter operations to each element. If you want the channel to
	// be populated with the values after applying the chained map and filter operations, use CollectIntoChannel.
	IntoChannel(ch chan<- T, opts ...IntoChannelOption)
	// CollectChannel returns a channel that will be populated with the values in the iterator. The channel will be closed when
	// there are no more values, indicating that the iterator has been consumed. This method does apply the chained map and
	// filter operations, so it is equivalent to calling Collect and then sending the resulting slice to a channel.
	CollectChannel() <-chan T
	// CollectIntoChannel populates the given channel with the values in the iterator. If shouldClose is true, the channel will be
	// closed when there are no more values, indicating that the iterator has been consumed. This method does apply the chained
	// map and filter operations, so it is equivalent to calling Collect and then sending the resulting slice to a channel.
	CollectIntoChannel(ch chan<- T, opts ...IntoChannelOption)
	// Reduce applies the given function to each value in the iterator, passing the result of the previous function call as the
	// first argument and the next value as the second argument until there are no more values. The initial value is passed to
	// the anonymous function as the first argument on the first iteration.
	Reduce(fn func(accumulator T, next T) T, initial T) T
	// Reset resets the iterator to the beginning of the source slice. This is useful if you want to iterate over the same
	// slice multiple times. Note that this does not reset the chained map and filter operations. If you want to reset those,
	// you should create a new iterator using the From function.
	Reset()
}
