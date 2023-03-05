package iterator

// fromOptions is a struct that holds the options for creating an iterator using the From function.
type fromOptions struct {
	copySource bool // whether to copy the source slice when creating the iterator
	threadSafe bool // whether to use a mutex when making calls to the Next method
	bufferLen  int  // the initial capacity of the operations buffer
}

// FromOption is a function that configures the parameters when creating an iterator using the From function.
type FromOption func(*fromOptions)

// CopySource returns an option that specifies whether the source slice should be copied when creating the iterator.
func CopySource(shouldCopy bool) FromOption {
	return func(opts *fromOptions) {
		opts.copySource = shouldCopy
	}
}

// ThreadSafe returns an option that specifies whether the iterator should be thread-safe when making calls to the Next method.
// Note that this option incurs a performance penalty, as it requires the use of a mutex.
func ThreadSafe(shouldLock bool) FromOption {
	return func(opts *fromOptions) {
		opts.threadSafe = shouldLock
	}
}

// BufferLen returns an option that specifies the initial capacity of the operations (like filter, map) buffer.
// This option is useful if you know in advance how many operations you'll be performing on the iterator.
// The default value is 64.
func BufferLen(bufferLen int) FromOption {
	return func(opts *fromOptions) {
		opts.bufferLen = bufferLen
	}
}

// uniqueOptions is a struct that holds the conditions for the Unique method.
type uniqueOptions struct {
	deref bool // whether to dereference pointers before evaluating uniqueness
}

// UniqueOption is a function that configures the conditions for the Unique method.
type UniqueOption func(*uniqueOptions)

// DerefPointers returns a UniqueOption that specifies whether the iterator should dereference pointers before evaluating uniqueness.
// Note that this option incurs a performance penalty, as it requires the use of reflection. If performance is critical, but
// you still need to evaluate uniqueness on the dereferenced values, you should consider writing your own Filter function
// instead of using this option as you'll have concrete types to work with.
func DerefPointers(shouldDeref bool) UniqueOption {
	return func(opts *uniqueOptions) {
		opts.deref = shouldDeref
	}
}

type intoChannelOptions struct {
	closeChannel bool // whether to close the channel when the iterator is exhausted
}

// IntoChannelOption is a function that configures the conditions for the IntoChannel method.
type IntoChannelOption func(*intoChannelOptions)

// CloseChannel returns an IntoChannelOption that specifies whether the channel should be closed when the iterator is exhausted.
func CloseChannel(shouldClose bool) IntoChannelOption {
	return func(opts *intoChannelOptions) {
		opts.closeChannel = shouldClose
	}
}
