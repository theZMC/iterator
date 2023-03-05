package iterator_test

import (
	"reflect"
	"testing"

	"github.com/thezmc/iterator"
)

type test[T any] struct {
	source   []T
	configFn func(iterator.Of[T])
	expected []T
}

func runCollect[T any](t *testing.T, test test[T]) {
	t.Helper()
	it := iterator.From(test.source)
	if test.configFn != nil {
		test.configFn(it)
	}
	result := it.Collect()
	if !reflect.DeepEqual(result, test.expected) {
		t.Errorf("expected %+v, got %+v", test.expected, result)
	}
}

func Test_Iterator_Collect_Ints(t *testing.T) {
	tests := map[string]test[int]{
		"empty": {
			source:   []int{},
			configFn: nil,
			expected: []int{},
		},
		"single": {
			source:   []int{1},
			configFn: nil,
			expected: []int{1},
		},
		"multiple": {
			source:   []int{1, 2, 3},
			configFn: nil,
			expected: []int{1, 2, 3},
		},
		"map": {
			source: []int{1, 2, 3},
			configFn: func(it iterator.Of[int]) {
				it.Map(func(val int) int {
					return val * 2
				})
			},
			expected: []int{2, 4, 6},
		},
		"filter": {
			source: []int{1, 2, 3},
			configFn: func(it iterator.Of[int]) {
				it.Filter(func(val int) bool {
					return val%2 == 0
				})
			},
			expected: []int{2},
		},
		"unique": {
			source: []int{1, 2, 3, 2, 1},
			configFn: func(it iterator.Of[int]) {
				it.Unique()
			},
			expected: []int{1, 2, 3},
		},
		"map, filter, unique": {
			source: []int{1, 2, 3, 2, 1},
			configFn: func(it iterator.Of[int]) {
				it.Map(func(val int) int {
					return val * 2 // 2, 4, 6, 4, 2
				}).Filter(func(val int) bool {
					return val%2 == 0
				}).Unique() // 2, 4, 6
			},
			expected: []int{2, 4, 6},
		},
		"multiple maps": {
			source: []int{1, 2, 3},
			configFn: func(it iterator.Of[int]) {
				it.Map(func(val int) int {
					return val * 2 // 2, 4, 6
				}).Map(func(val int) int {
					return val * 3 // 6, 12, 18
				})
			},
			expected: []int{6, 12, 18},
		},
		"multiple filters": {
			source: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
			configFn: func(it iterator.Of[int]) {
				it.Filter(func(val int) bool {
					return val%2 == 0 // 2, 4, 6, 8
				}).Filter(func(val int) bool {
					return val%3 == 0 // 6
				})
			},
			expected: []int{6},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			runCollect(t, test)
		})
	}
}

func Test_Iterator_Collect_Strings(t *testing.T) {
	tests := map[string]test[string]{
		"empty": {
			source:   []string{},
			configFn: nil,
			expected: []string{},
		},
		"single": {
			source:   []string{"a"},
			configFn: nil,
			expected: []string{"a"},
		},
		"multiple": {
			source:   []string{"a", "b", "c"},
			configFn: nil,
			expected: []string{"a", "b", "c"},
		},
		"map": {
			source: []string{"a", "b", "c"},
			configFn: func(it iterator.Of[string]) {
				it.Map(func(val string) string {
					return val + val
				})
			},
			expected: []string{"aa", "bb", "cc"},
		},
		"filter": {
			source: []string{"a", "b", "c"},
			configFn: func(it iterator.Of[string]) {
				it.Filter(func(val string) bool {
					return val == "b"
				})
			},
			expected: []string{"b"},
		},
		"unique": {
			source: []string{"a", "b", "c", "b", "a"},
			configFn: func(it iterator.Of[string]) {
				it.Unique()
			},
			expected: []string{"a", "b", "c"},
		},
		"map, filter, unique": {
			source: []string{"a", "b", "c", "b", "a"},
			configFn: func(it iterator.Of[string]) {
				it.Map(func(val string) string {
					return val + val
				}).Filter(func(val string) bool {
					return val == "aa"
				}).Unique()
			},
			expected: []string{"aa"},
		},
		"multiple maps": {
			source: []string{"a", "b", "c"},
			configFn: func(it iterator.Of[string]) {
				it.Map(func(val string) string {
					return val + val
				}).Map(func(val string) string {
					return val + val
				})
			},
			expected: []string{"aaaa", "bbbb", "cccc"},
		},
		"multiple filters": {
			source: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i"},
			configFn: func(it iterator.Of[string]) {
				it.Filter(func(val string) bool {
					return val == "a" || val == "b" || val == "c"
				}).Filter(func(val string) bool {
					return val == "b"
				})
			},
			expected: []string{"b"},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			runCollect(t, test)
		})
	}
}

func Test_Iterator_Collect_Structs(t *testing.T) {
	type person struct {
		name     string
		age      int
		canDrink bool
	}
	tests := map[string]test[person]{
		"empty": {
			source:   []person{},
			configFn: nil,
			expected: []person{},
		},
		"single": {
			source:   []person{{"Felicita", 23, true}},
			configFn: nil,
			expected: []person{{"Felicita", 23, true}},
		},
		"multiple": {
			source:   []person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}},
			configFn: nil,
			expected: []person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}},
		},
		"map": {
			source: []person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}},
			configFn: func(it iterator.Of[person]) {
				it.Map(func(val person) person {
					val.age = val.age * 2
					return val
				})
			},
			expected: []person{{"Felicita", 46, true}, {"Luis", 48, false}, {"Juan", 50, true}},
		},
		"filter": {
			source: []person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}},
			configFn: func(it iterator.Of[person]) {
				it.Filter(func(val person) bool {
					return val.name[1] == 'u'
				})
			},
			expected: []person{{"Luis", 24, false}, {"Juan", 25, true}},
		},
		"unique": {
			source: []person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}, {"Luis", 24, false}, {"Felicita", 23, true}},
			configFn: func(it iterator.Of[person]) {
				it.Unique()
			},
			expected: []person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}},
		},
		"map, filter, unique": {
			source: []person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}, {"Luis", 24, false}, {"Felicita", 23, true}},
			configFn: func(it iterator.Of[person]) {
				it.Map(func(val person) person {
					val.age = val.age * 2
					return val
				}).Filter(func(val person) bool {
					return val.name[1] == 'u'
				}).Unique()
			},
			expected: []person{{"Luis", 48, false}, {"Juan", 50, true}},
		},
		"multiple maps": {
			source: []person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}},
			configFn: func(it iterator.Of[person]) {
				it.Map(func(val person) person {
					val.age = val.age * 2
					return val
				}).Map(func(val person) person {
					val.name = val.name + " Valdez"
					return val
				})
			},
			expected: []person{{"Felicita Valdez", 46, true}, {"Luis Valdez", 48, false}, {"Juan Valdez", 50, true}},
		},
		"multiple filters": {
			source: []person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}, {"Fernando", 26, false}, {"Fernanda", 27, true}, {"Fernando", 28, false}, {"Fernanda", 29, true}, {"Fernando", 30, false}, {"Fernanda", 31, true}},
			configFn: func(it iterator.Of[person]) {
				it.Filter(func(val person) bool {
					return val.name[1] == 'u'
				}).Filter(func(val person) bool {
					return val.name[2] == 'a'
				})
			},
			expected: []person{{"Juan", 25, true}},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			runCollect(t, test)
		})
	}
}

func Test_Iterator_Collect_Pointers(t *testing.T) {
	one := 1
	two := 2
	three := 3
	four := 4
	five := 5
	six := 6
	seven := 7
	eight := 8
	nine := 9
	ten := 10

	tests := map[string]test[*int]{
		"empty": {
			source:   []*int{},
			configFn: nil,
			expected: []*int{},
		},
		"single": {
			source:   []*int{&one},
			configFn: nil,
			expected: []*int{&one},
		},
		"multiple": {
			source:   []*int{&one, &two, &three},
			configFn: nil,
			expected: []*int{&one, &two, &three},
		},
		"filter": {
			source: []*int{&one, &two, &three, &four, &five, &six, &seven, &eight, &nine, &ten},
			configFn: func(it iterator.Of[*int]) {
				it.Filter(func(val *int) bool {
					return *val%2 == 0
				})
			},
			expected: []*int{&two, &four, &six, &eight, &ten},
		},
		"unique": {
			source: []*int{&one, &two, &three, &four, &five, &six, &seven, &eight, &nine, &ten, &one, &two, &three, &four, &five, &six, &seven, &eight, &nine, &ten},
			configFn: func(it iterator.Of[*int]) {
				it.Unique()
			},
			expected: []*int{&one, &two, &three, &four, &five, &six, &seven, &eight, &nine, &ten},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			runCollect(t, test)
		})
	}
}

func Test_Iterator_Collect_Pointers_Deref(t *testing.T) {
	type person struct {
		name     string
		age      int
		canDrink bool
	}
	tests := map[string]test[*person]{
		"empty": {
			source:   []*person{},
			configFn: nil,
			expected: []*person{},
		},
		"single": {
			source:   []*person{{"Felicita", 23, true}},
			configFn: nil,
			expected: []*person{{"Felicita", 23, true}},
		},
		"multiple": {
			source:   []*person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}},
			configFn: nil,
			expected: []*person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}},
		},
		"map": {
			source: []*person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}},
			configFn: func(it iterator.Of[*person]) {
				it.Map(func(val *person) *person {
					val.age = val.age * 2
					return val
				})
			},
			expected: []*person{{"Felicita", 46, true}, {"Luis", 48, false}, {"Juan", 50, true}},
		},
		"filter": {
			source: []*person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}},
			configFn: func(it iterator.Of[*person]) {
				it.Filter(func(val *person) bool {
					return val.name[1] == 'u'
				})
			},
			expected: []*person{{"Luis", 24, false}, {"Juan", 25, true}},
		},
		"unique": {
			source: []*person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}, {"Luis", 24, false}, {"Felicita", 23, true}},
			configFn: func(it iterator.Of[*person]) {
				it.Unique(iterator.DerefPointers(true))
			},
			expected: []*person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}},
		},
		"map, filter, unique": {
			source: []*person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}, {"Luis", 24, false}, {"Felicita", 23, true}},
			configFn: func(it iterator.Of[*person]) {
				it.Map(func(val *person) *person {
					val.age = val.age * 2
					return val
				}).Filter(func(val *person) bool {
					return val.name[1] == 'u'
				}).Unique(iterator.DerefPointers(true))
			},
			expected: []*person{{"Luis", 48, false}, {"Juan", 50, true}},
		},
		"multiple maps": {
			source: []*person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}},
			configFn: func(it iterator.Of[*person]) {
				it.Map(func(val *person) *person {
					val.age = val.age * 2
					return val
				}).Map(func(val *person) *person {
					val.name = val.name + " Valdez"
					return val
				})
			},
			expected: []*person{{"Felicita Valdez", 46, true}, {"Luis Valdez", 48, false}, {"Juan Valdez", 50, true}},
		},
		"multiple filters": {
			source: []*person{{"Felicita", 23, true}, {"Luis", 24, false}, {"Juan", 25, true}, {"Fernando", 26, false}, {"Fernanda", 27, true}, {"Fernando", 28, false}, {"Fernanda", 29, true}, {"Fernando", 30, false}, {"Fernanda", 31, true}},
			configFn: func(it iterator.Of[*person]) {
				it.Filter(func(val *person) bool {
					return val.name[1] == 'u'
				}).Filter(func(val *person) bool {
					return val.name[2] == 'a'
				})
			},
			expected: []*person{{"Juan", 25, true}},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			runCollect(t, test)
		})
	}
}

func Test_Iterator_Reduce(t *testing.T) {
	iter := iterator.From([]int{1, 2, 3, 4, 5})
	sum := iter.Reduce(func(acc, val int) int {
		return acc + val
	}, 0)
	if sum != 15 {
		t.Errorf("Expected 15, got %d", sum)
	}
}

func Test_Iterator_Reset(t *testing.T) {
	iter := iterator.From([]int{1, 2, 3, 4, 5})
	iter.Next()
	iter.Reset()
	if num, ok := iter.Next(); !ok || num != 1 {
		t.Errorf("Expected 1, got %d", num)
	}
	iter.Map(func(val int) int {
		return val * 2
	})
	iter.Next()
	iter.Reset()
	if nums := iter.Collect(); !reflect.DeepEqual(nums, []int{2, 4, 6, 8, 10}) { // proving that the operations were not reset
		t.Errorf("Expected [2 4 6 8 10], got %v", nums)
	}
}

func Test_Iterator_Channel(t *testing.T) {
	iter := iterator.From([]int{1, 2, 3, 4, 5})
	ch := iter.Channel()
	num := 1
	for i := range ch {
		if i != num {
			t.Errorf("Expected %d, got %d", num, i)
		}
		num++
	}
}

func Test_Iterator_CollectChannel(t *testing.T) {
	iter := iterator.From([]int{1, 2, 3, 4, 5})
	iter.Map(func(val int) int {
		return val * 2
	})
	ch := iter.CollectChannel()
	num := 2
	for i := range ch {
		if i != num {
			t.Errorf("Expected %d, got %d", num, i)
		}
		num += 2
	}
}
