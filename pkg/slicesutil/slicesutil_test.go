package slicesutil

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var in []int = nil
		got, found := Find(in, nil)
		assert.False(t, found)
		assert.Zero(t, got)
	})

	t.Run("empty", func(t *testing.T) {
		got, found := Find([]int{}, func(element int) bool { return element > 0 })
		assert.False(t, found)
		assert.Zero(t, got)
	})

	t.Run("simple", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		result, found := Find([]Person{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 30},
			{Name: "Charlie", Age: 20},
		}, func(element Person) bool { return element.Age == 30 })
		assert.True(t, found)
		assert.Equal(t, Person{Name: "Bob", Age: 30}, result)
	})
}

func TestFilter(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var in []int = nil

		assert.Equal(t, []int{}, Filter(in, nil))
	})

	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, []int{}, Filter([]int{}, func(element int) bool { return element > 0 }))
	})

	t.Run("simple", func(t *testing.T) {
		assert.Equal(t, []int{2, 4}, Filter([]int{1, 2, 3, 4, 5, 7, 9}, func(element int) bool { return element%2 == 0 }))
	})
}

func TestIntersect(t *testing.T) {
	t.Run("Struct Slices with Custom Compare", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		inputSliceA := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
			{Name: "Charlie", Age: 30},
		}
		inputSliceB := []Person{
			{Name: "Bob", Age: 35},
			{Name: "Alice", Age: 30},
			{Name: "David", Age: 25},
		}
		expectedSlice := []Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
		}
		compare := func(p1, p2 Person) bool {
			return p1.Name == p2.Name // Compare based on Name only
		}
		resultSlice := IntersectFunc(inputSliceA, inputSliceB, compare)

		// We need to compare structs field by field because assert.Equal might not work directly for struct slices with custom compare.
		assert.Equal(t, len(expectedSlice), len(resultSlice))

		for i := range expectedSlice {
			assert.True(t, compare(expectedSlice[i], resultSlice[i]))
		}
	})

	t.Run("Nil Slice A", func(t *testing.T) {
		var inputSliceA []int = nil

		inputSliceB := []int{1, 2, 3}
		expectedSlice := []int{}
		compare := func(a, b int) bool { return a == b }
		resultSlice := IntersectFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("Nil Slice B", func(t *testing.T) {
		inputSliceA := []int{1, 2, 3}

		var inputSliceB []int = nil

		expectedSlice := []int{}
		compare := func(a, b int) bool { return a == b }
		resultSlice := IntersectFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("A is Subset of B", func(t *testing.T) {
		inputSliceA := []int{1, 2, 3}
		inputSliceB := []int{0, 1, 2, 3, 4}
		expectedSlice := []int{1, 2, 3}
		compare := func(a, b int) bool { return a == b }
		resultSlice := IntersectFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("B is Subset of A", func(t *testing.T) {
		inputSliceA := []int{0, 1, 2, 3, 4}
		inputSliceB := []int{1, 2, 3}
		expectedSlice := []int{1, 2, 3}
		compare := func(a, b int) bool { return a == b }
		resultSlice := IntersectFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("Identical Slices", func(t *testing.T) {
		inputSliceA := []int{1, 2, 3}
		inputSliceB := []int{1, 2, 3}
		expectedSlice := []int{1, 2, 3}
		compare := func(a, b int) bool { return a == b }
		resultSlice := IntersectFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("Empty Slice A", func(t *testing.T) {
		inputSliceA := []int{}
		inputSliceB := []int{1, 2, 3}
		expectedSlice := []int{}
		compare := func(a, b int) bool { return a == b }
		resultSlice := IntersectFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("Empty Slice B", func(t *testing.T) {
		inputSliceA := []int{1, 2, 3}
		inputSliceB := []int{}
		expectedSlice := []int{}
		compare := func(a, b int) bool { return a == b }
		resultSlice := IntersectFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})
}

func TestMap(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, []string{}, Map([]int{}, strconv.Itoa))
	})

	t.Run("int to string", func(t *testing.T) {
		assert.Equal(t, []string{"1", "2", "3", "4", "5"}, Map([]int{1, 2, 3, 4, 5}, strconv.Itoa))
	})

	t.Run("struct to string", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		assert.Equal(t, []string{"Alice", "Bob"}, Map([]Person{
			{Name: "Alice", Age: 30},
			{Name: "Bob", Age: 25},
		}, func(p Person) string { return p.Name }))
	})
}

func TestPartition(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var in []int = nil

		assert.Equal(t, [][]int{}, Partition(in, 3))
	})

	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, [][]int{}, Partition([]int{}, 3))
	})

	t.Run("size is zero", func(t *testing.T) {
		assert.Equal(t, [][]int{}, Partition([]int{1, 2, 3, 4, 5}, 0))
	})

	t.Run("size is negative", func(t *testing.T) {
		assert.Equal(t, [][]int{}, Partition([]int{1, 2, 3, 4, 5}, -2))
	})

	t.Run("size is one", func(t *testing.T) {
		assert.Equal(t, [][]int{{1}, {2}, {3}}, Partition([]int{1, 2, 3}, 1))
	})

	t.Run("size is equal to slice length", func(t *testing.T) {
		assert.Equal(t, [][]int{{1, 2, 3}}, Partition([]int{1, 2, 3}, 3))
	})

	t.Run("size is larger than slice length", func(t *testing.T) {
		assert.Equal(t, [][]int{{1, 2, 3}}, Partition([]int{1, 2, 3}, 5))
	})

	t.Run("size does not divide slice length exactly", func(t *testing.T) {
		assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7}}, Partition([]int{1, 2, 3, 4, 5, 6, 7}, 3))
	})
}

func TestToAny(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, []any{}, ToAny([]int{}))
	})

	t.Run("int", func(t *testing.T) {
		assert.Equal(t, []any{1, 2, 3, 4, 5}, ToAny([]int{1, 2, 3, 4, 5}))
	})
}

func TestRemoveEmpty(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		var in []int = nil

		assert.Equal(t, []int{}, RemoveEmpty(in))
	})

	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, []int{}, RemoveEmpty([]int{}))
	})

	t.Run("simple", func(t *testing.T) {
		assert.Equal(t, []int{1, 2, 3, 4, 5}, RemoveEmpty([]int{0, 1, 2, 3, 0, 4, 5, 0}))
	})

	t.Run("only empty elements", func(t *testing.T) {
		assert.Equal(t, []string{}, RemoveEmpty([]string{"", ""}))
	})
}

func TestSubtract(t *testing.T) {
	t.Run("Nil Slice A", func(t *testing.T) {
		var inputSliceA []int = nil

		inputSliceB := []int{1, 2, 3}
		expectedSlice := []int{}
		compare := func(a, b int) bool { return a == b }
		resultSlice := SubtractFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("Nil Slice B", func(t *testing.T) {
		inputSliceA := []int{1, 2, 3}

		var inputSliceB []int = nil

		expectedSlice := []int{1, 2, 3}
		compare := func(a, b int) bool { return a == b }
		resultSlice := SubtractFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("Empty Slice A", func(t *testing.T) {
		inputSliceA := []int{}
		inputSliceB := []int{1, 2, 3}
		expectedSlice := []int{}
		compare := func(a, b int) bool { return a == b }
		resultSlice := SubtractFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("Empty Slice B", func(t *testing.T) {
		inputSliceA := []int{1, 2, 3}
		inputSliceB := []int{}
		expectedSlice := []int{1, 2, 3}
		compare := func(a, b int) bool { return a == b }
		resultSlice := SubtractFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("Both Slices Empty", func(t *testing.T) {
		inputSliceA := []int{}
		inputSliceB := []int{}
		expectedSlice := []int{}
		compare := func(a, b int) bool { return a == b }
		resultSlice := SubtractFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("No Difference", func(t *testing.T) {
		inputSliceA := []int{1, 2, 3}
		inputSliceB := []int{1, 2, 3, 4, 5}
		expectedSlice := []int{}
		compare := func(a, b int) bool { return a == b }
		resultSlice := SubtractFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("Some Difference", func(t *testing.T) {
		inputSliceA := []int{1, 2, 3, 4, 5}
		inputSliceB := []int{3, 4, 5, 6, 7}
		expectedSlice := []int{1, 2}
		compare := func(a, b int) bool { return a == b }
		resultSlice := SubtractFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("Duplicates in A, Some in B", func(t *testing.T) {
		inputSliceA := []int{1, 2, 2, 3, 4, 4, 4, 5}
		inputSliceB := []int{2, 4, 6}
		expectedSlice := []int{1, 3, 5}
		compare := func(a, b int) bool { return a == b }
		resultSlice := SubtractFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("String Slices Case-Insensitive", func(t *testing.T) {
		inputSliceA := []string{"Apple", "Banana", "Orange", "Grape"}
		inputSliceB := []string{"banana", "GRAPE", "Kiwi"}
		expectedSlice := []string{"Apple", "Orange"}
		compare := func(a, b string) bool {
			return strings.ToLower(a) == strings.ToLower(b)
		}
		resultSlice := SubtractFunc(inputSliceA, inputSliceB, compare)
		assert.Equal(t, expectedSlice, resultSlice)
	})
}

func TestMergeFunc(t *testing.T) {
	t.Run("Empty Slices A and B", func(t *testing.T) {
		inputSliceA := []int{}
		inputSliceB := []int{}

		var expectedSlice []int = nil

		compareFunc := func(a, b int) bool { return a == b }
		mergeFunc := func(a, b int) int { return a + b }
		resultSlice := MergeFunc(inputSliceA, inputSliceB, compareFunc, mergeFunc)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("Empty Slice A, Non-Empty Slice B", func(t *testing.T) {
		inputSliceA := []int{}
		inputSliceB := []int{1, 2, 3}
		expectedSlice := []int{1, 2, 3}
		compareFunc := func(a, b int) bool { return a == b }
		mergeFunc := func(a, b int) int { return a + b }
		resultSlice := MergeFunc(inputSliceA, inputSliceB, compareFunc, mergeFunc)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("Non-Empty Slice A, Empty Slice B", func(t *testing.T) {
		inputSliceA := []int{1, 2, 3}
		inputSliceB := []int{}
		expectedSlice := []int{1, 2, 3}
		compareFunc := func(a, b int) bool { return a == b }
		mergeFunc := func(a, b int) int { return a + b }
		resultSlice := MergeFunc(inputSliceA, inputSliceB, compareFunc, mergeFunc)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("No Duplicates - Integers", func(t *testing.T) {
		inputSliceA := []int{1, 2, 3}
		inputSliceB := []int{4, 5, 6}
		expectedSlice := []int{1, 2, 3, 4, 5, 6}
		compareFunc := func(a, b int) bool { return a == b }
		mergeFunc := func(a, b int) int { return a + b }
		resultSlice := MergeFunc(inputSliceA, inputSliceB, compareFunc, mergeFunc)
		assert.Equal(t, expectedSlice, resultSlice)
	})

	t.Run("No Duplicates - Strings", func(t *testing.T) {
		inputSliceA := []string{"Apple", "Banana", "Orange"}
		inputSliceB := []string{"Grape", "Kiwi", "Mango"}
		expectedSlice := []string{"Apple", "Banana", "Orange", "Grape", "Kiwi", "Mango"}
		compareFunc := func(a, b string) bool { return a == b }
		mergeFunc := func(a, b string) string { return a + b }
		resultSlice := MergeFunc(inputSliceA, inputSliceB, compareFunc, mergeFunc)
		assert.Equal(t, expectedSlice, resultSlice)
	})
}
