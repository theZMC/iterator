# `iterator`
[![codecov](https://codecov.io/github/theZMC/iterator/branch/main/graph/badge.svg?token=F3HIWA9OGD)](https://codecov.io/github/theZMC/iterator)
[![Go Report Card](https://goreportcard.com/badge/github.com/theZMC/iterator)](https://goreportcard.com/report/github.com/theZMC/iterator)
[![CI](https://github.com/theZMC/iterator/actions/workflows/ci.yml/badge.svg)](https://github.com/theZMC/iterator/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/theZMC/iterator.svg)](https://pkg.go.dev/github.com/theZMC/iterator)

The iterator package provides a simple iterator interface for Go. As go's generics don't allow the use of generic type
parameters on methods that differ from the types on the receiver, this iterator implementation isn't quite as powerful
as you might find in other languages. That said, I still think there are many situations where this package can be useful.

## Usage
The package provides a single interface, `Of`, which provides all of the high-level functionality of existing
and possible future iterator implementations.

### Creating an iterator
To create an iterator, you can use `iterator.From` which takes a slice of any type and returns an iterator over that
slice. For example, to create an iterator over a slice of integers:
```go
it := iterator.From([]int{1, 2, 3})
```

### Simple iteration
To iterate over an iterator, you can use the `Next` method. This method will return the next value (maybe) and a boolean
indicating whether or not there was a next value. For example, to iterate over the iterator created above:
```go
for {
  val, ok := it.Next()
  if !ok {
    break
  }
  fmt.Println(val)
}
// Output:
// 1
// 2
// 3
```

### Functional iteration
This iterator implementation also provides `Filter`, `Map`, `Reduce`, and `ForEach` methods which allow you to perform
common functional operations on an iterator. The iterator can then be collected into a slice using the `Collect` method.
For example, to filter an iterator to only even numbers and then collect the results into a slice:
```go
result := iterator.From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
  Filter(func(val int) bool {
    return val % 2 == 0
  }).
  Collect()

fmt.Println(result)
// Output:
// [2 4 6 8 10]
```

### Chaining
The `Filter` and `Map` methods return the iterator itself, allowing you to chain these methods together. For example, to
filter an iterator to only even numbers, double each value, and then collect the results into a slice:
```go
result := iterator.From([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).
  Filter(func(val int) bool {
    return val % 2 == 0
  }).
  Map(func(val int) int {
    return val * 2
  }).
  Collect()

fmt.Println(result)
// Output:
// [4 8 12 16 20]
```

### Using `ForEach`
The `ForEach` method is similar to the `Next` method, but it doesn't return a value. Instead, it takes a function which
is called for each value in the iterator, performing some side effect. For example, to print each value in an iterator:
```go
iterator.From([]int{1, 2, 3}).ForEach(func(val int) {
  fmt.Println(val)
})
// Output:
// 1
// 2
// 3
```

### Using `Reduce`
You can use the `Reduce` method when you only want one value from an iterator. This method takes a function which is
called for each value and accumulates the result. For example, to compute the sum of all values in an iterator:
```go
sum := iterator.From([]int{1, 2, 3}).Reduce(func(acc, val int) int {
  return acc + val
}, 0)

fmt.Println(sum)
// Output:
// 6
```

### Iterating into channels
The `Of` interface also provides some convenient channel methods. The `Channel` method returns a channel which will
receive all of the values in the iterator and will be closed when the iterator is exhausted. Example:
```go
ch := iterator.From([]int{1, 2, 3}).Channel()
for val := range ch {
  fmt.Println(val)
}
// Output:
// 1
// 2
// 3
```

The `IntoChannel` method does the same thing, but with a user-provided channel. The sends happen in a separate goroutine,
so this method is non-blocking, even if the channel is full. Example:
```go
ch := make(chan int)
iterator.From([]int{1, 2, 3}).IntoChannel(ch)
for val := range ch {
  fmt.Println(val)
}
// Output:
// 1
// 2
// 3
```
The `IntoChannel` method is useful when you want to iterate over multiple iterators in parallel, but want to send the
results to a single channel.

By default, the `IntoChannel` method will not close the channel when the iterator is exhausted, leaving that up to the
caller. However, if you want the channel to be closed when the iterator is exhausted, you can pass a functional option
to the `IntoChannel` method to do so. Example:
```go
ch := make(chan int)
iterator.From([]int{1, 2, 3}).IntoChannel(ch, iterator.CloseChannel(true))
```

Both the `Channel` and `IntoChannel` methods iterate without applying any functional operations. If you want to apply the
chained `Filter` and `Map` operations, you can use the `CollectChannel` or `CollectIntoChannel` methods. Other than
that, these methods work the same as the `Channel` and `IntoChannel` methods.

## Performance
Because go lacks tail call optimization, the `Collect` method does cause quite a few allocations. Despite this, benchmarks
do show that this implementation is still quite fast. Take a look at the benchmarks in the package and compare the results
yourself. This package may not be the best choice in all situations, and I leave it up to you to benchmark your particular
use case to determine if it's right for you.

## Contributing
Contributions are welcome! If you find a bug or have a feature request, please open an issue. If you'd like to contribute
code, please open an issue stating the problem/feature along with a pull request that you think could resolve it. Any new
implementations must satisfy the `Of` interface and their constructor should return that interface rather than the
concrete type.

## License
This project was created by [Zach Callahan](https://zmc.dev) and is licensed under the Apache License, Version 2.0. See
the [LICENSE](LICENSE) file for more details.