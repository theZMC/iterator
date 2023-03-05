package iterator_test

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/thezmc/iterator"
)

type TestStruct struct {
	Num int64
}

var (
	IntResult    []int64
	StructResult []*TestStruct
)

func makeRandomSlice(b *testing.B, size int) []int64 {
	source := make([]int64, size)
	rdr := rand.Reader
	for i := range source {
		randInt, err := rand.Int(rdr, big.NewInt(int64(size*10)))
		if err != nil {
			b.Fatal(err)
		}
		source[i] = randInt.Int64()
	}
	return source
}

func Benchmark_Iterator_Ints_NoCopy(b *testing.B) {
	nums := makeRandomSlice(b, 1_000_000)
	iter := iterator.From(nums)
	iter.Filter(func(i int64) bool {
		return i%2 == 0
	}).Map(func(i int64) int64 {
		return i * 2
	}).Filter(func(i int64) bool {
		return i%3 == 0
	}).Map(func(i int64) int64 {
		return i / 2
	})
	for i := 0; i < b.N; i++ {
		IntResult = iter.Collect()
		iter.Reset()
	}
}

func Benchmark_Iterator_Ints(b *testing.B) {
	nums := makeRandomSlice(b, 1_000_000)
	iter := iterator.From(nums, iterator.CopySource(true))
	iter.Filter(func(i int64) bool {
		return i%2 == 0
	}).Map(func(i int64) int64 {
		return i * 2
	}).Filter(func(i int64) bool {
		return i%3 == 0
	}).Map(func(i int64) int64 {
		return i / 2
	})
	for i := 0; i < b.N; i++ {
		IntResult = iter.Collect()
		iter.Reset()
	}
}

func Benchmark_NoIterator_Ints(b *testing.B) {
	nums := makeRandomSlice(b, 1_000_000)
	times2 := func(i int64) int64 {
		return i * 2
	}
	div2 := func(i int64) int64 {
		return i / 2
	}
	isEven := func(i int64) bool {
		return i%2 == 0
	}
	isDiv3 := func(i int64) bool {
		return i%3 == 0
	}
	for i := 0; i < b.N; i++ {
		IntResult = genericFilter(nums, isEven)
		IntResult = genericMap(IntResult, times2)
		IntResult = genericFilter(IntResult, isDiv3)
		IntResult = genericMap(IntResult, div2)
	}
}

func Benchmark_Iterator_Structs_NoCopy(b *testing.B) {
	strs := make([]*TestStruct, 1_000_000)
	for i := range strs {
		strs[i] = &TestStruct{Num: int64(i)}
	}
	iter := iterator.From(strs)
	iter.Filter(func(s *TestStruct) bool {
		return s.Num%2 == 0
	}).Map(func(s *TestStruct) *TestStruct {
		return &TestStruct{Num: s.Num * 2}
	}).Filter(func(s *TestStruct) bool {
		return s.Num%3 == 0
	}).Map(func(s *TestStruct) *TestStruct {
		return &TestStruct{Num: s.Num / 2}
	})
	for i := 0; i < b.N; i++ {
		StructResult = iter.Collect()
		iter.Reset()
	}
}

func Benchmark_Iterator_Structs(b *testing.B) {
	strs := make([]*TestStruct, 1_000_000)
	for i := range strs {
		strs[i] = &TestStruct{Num: int64(i)}
	}
	iter := iterator.From(strs, iterator.CopySource(true))
	iter.Filter(func(s *TestStruct) bool {
		return s.Num%2 == 0
	}).Map(func(s *TestStruct) *TestStruct {
		return &TestStruct{Num: s.Num * 2}
	}).Filter(func(s *TestStruct) bool {
		return s.Num%3 == 0
	}).Map(func(s *TestStruct) *TestStruct {
		return &TestStruct{Num: s.Num / 2}
	})
	for i := 0; i < b.N; i++ {
		StructResult = iter.Collect()
		iter.Reset()
	}
}

func Benchmark_NoIterator_Structs(b *testing.B) {
	strs := make([]*TestStruct, 1_000_000)
	for i := range strs {
		strs[i] = &TestStruct{Num: int64(i)}
	}
	times2 := func(s *TestStruct) *TestStruct {
		return &TestStruct{Num: s.Num * 2}
	}
	div2 := func(s *TestStruct) *TestStruct {
		return &TestStruct{Num: s.Num / 2}
	}
	isEven := func(s *TestStruct) bool {
		return s.Num%2 == 0
	}
	isDiv3 := func(s *TestStruct) bool {
		return s.Num%3 == 0
	}
	for i := 0; i < b.N; i++ {
		StructResult = genericFilter(strs, isEven)
		StructResult = genericMap(StructResult, times2)
		StructResult = genericFilter(StructResult, isDiv3)
		StructResult = genericMap(StructResult, div2)
	}
}

func genericMap[T any](source []T, fn func(T) T) []T {
	result := make([]T, len(source))
	for i := range source {
		result[i] = fn(source[i])
	}
	return result
}

func genericFilter[T any](source []T, fn func(T) bool) []T {
	result := make([]T, 0, len(source))
	for i := range source {
		if fn(source[i]) {
			result = append(result, source[i])
		}
	}
	return result
}
