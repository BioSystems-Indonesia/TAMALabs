package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnique(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "Unique integers",
			input:    []int{1, 2, 2, 3, 4, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "Unique strings",
			input:    []string{"apple", "banana", "apple", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "Empty slice",
			input:    []int{},
			expected: []int(nil),
		},
		{
			name:     "No duplicates",
			input:    []string{"alpha", "beta", "gamma"},
			expected: []string{"alpha", "beta", "gamma"},
		},
		{
			name:     "Single element",
			input:    []int{42},
			expected: []int{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case []int:
				result := Unique(v)
				assert.Equal(t, tt.expected, result)
			case []string:
				result := Unique(v)
				assert.Equal(t, tt.expected, result)
			default:
				t.Fatalf("unsupported type %T", v)
			}
		})
	}
}

func TestCompareSlices(t *testing.T) {
	tests := []struct {
		name     string
		old      []interface{}
		new      []interface{}
		toDelete []interface{}
		toCreate []interface{}
	}{
		{
			name:     "Empty slices",
			old:      []interface{}{},
			new:      []interface{}{},
			toDelete: []interface{}{},
			toCreate: []interface{}{},
		},
		{
			name:     "Delete all elements",
			old:      []interface{}{1, 2, 3},
			new:      []interface{}{},
			toDelete: []interface{}{1, 2, 3},
			toCreate: []interface{}{},
		},
		{
			name:     "Create all elements",
			old:      []interface{}{},
			new:      []interface{}{1, 2, 3},
			toDelete: []interface{}{},
			toCreate: []interface{}{1, 2, 3},
		},
		{
			name:     "Same elements different order",
			old:      []interface{}{1, 2, 3},
			new:      []interface{}{3, 1, 2},
			toDelete: []interface{}{},
			toCreate: []interface{}{},
		},
		{
			name:     "Different frequencies",
			old:      []interface{}{1, 1, 2, 2, 3},
			new:      []interface{}{1, 2, 2, 2, 3, 3},
			toDelete: []interface{}{1},
			toCreate: []interface{}{2, 3},
		},
		{
			name:     "Mixed types",
			old:      []interface{}{"apple", 1, "banana", "banana"},
			new:      []interface{}{"apple", "banana", 2, 3},
			toDelete: []interface{}{"banana", 1},
			toCreate: []interface{}{2, 3},
		},
		{
			name:     "Completely different slices",
			old:      []interface{}{1, 2, 3},
			new:      []interface{}{4, 5, 6},
			toDelete: []interface{}{1, 2, 3},
			toCreate: []interface{}{4, 5, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toDelete, toCreate := CompareSlices(tt.old, tt.new)

			// Check if lengths match
			assert.Equal(t, len(tt.toDelete), len(toDelete),
				"toDelete length mismatch for test: %s", tt.name)
			assert.Equal(t, len(tt.toCreate), len(toCreate),
				"toCreate length mismatch for test: %s", tt.name)

			// Check if elements match (regardless of order)
			assert.ElementsMatch(t, tt.toDelete, toDelete,
				"toDelete elements mismatch for test: %s", tt.name)
			assert.ElementsMatch(t, tt.toCreate, toCreate,
				"toCreate elements mismatch for test: %s", tt.name)
		})
	}
}

// TestCompareSlicesTypes tests the function with specific types
func TestCompareSlicesTypes(t *testing.T) {
	t.Run("Integer slices", func(t *testing.T) {
		old := []int{1, 2, 2, 3}
		new := []int{2, 3, 3, 4}
		toDelete, toCreate := CompareSlices(old, new)

		assert.ElementsMatch(t, []int{1, 2}, toDelete)
		assert.ElementsMatch(t, []int{3, 4}, toCreate)
	})

	t.Run("String slices", func(t *testing.T) {
		old := []string{"apple", "banana", "banana", "cherry"}
		new := []string{"apple", "date", "banana", "elderberry"}
		toDelete, toCreate := CompareSlices(old, new)

		assert.ElementsMatch(t, []string{"banana", "cherry"}, toDelete)
		assert.ElementsMatch(t, []string{"date", "elderberry"}, toCreate)
	})

	t.Run("Float slices", func(t *testing.T) {
		old := []float64{1.1, 2.2, 2.2, 3.3}
		new := []float64{2.2, 3.3, 3.3, 4.4}
		toDelete, toCreate := CompareSlices(old, new)

		assert.ElementsMatch(t, []float64{1.1, 2.2}, toDelete)
		assert.ElementsMatch(t, []float64{3.3, 4.4}, toCreate)
	})
}

// TestCompareSlicesEdgeCases tests edge cases and potential error conditions
func TestCompareSlicesEdgeCases(t *testing.T) {
	t.Run("Nil slices", func(t *testing.T) {
		var old, new []int
		toDelete, toCreate := CompareSlices(old, new)

		assert.Empty(t, toDelete)
		assert.Empty(t, toCreate)
	})

	t.Run("Large number of duplicates", func(t *testing.T) {
		old := []int{1, 1, 1, 1, 1} // 5 ones
		new := []int{1, 1, 1}       // 3 ones
		toDelete, toCreate := CompareSlices(old, new)

		assert.Len(t, toDelete, 2) // Should delete 2 ones
		assert.ElementsMatch(t, []int{1, 1}, toDelete)
		assert.Empty(t, toCreate)
	})

	t.Run("Single element slices", func(t *testing.T) {
		old := []int{1}
		new := []int{2}
		toDelete, toCreate := CompareSlices(old, new)

		assert.ElementsMatch(t, []int{1}, toDelete)
		assert.ElementsMatch(t, []int{2}, toCreate)
	})
}

// TestMap tests the Map function with different use cases
func TestMap(t *testing.T) {
	// Test 1: Transforming a slice of integers to their squares
	ints := []int{1, 2, 3, 4, 5}
	expectedSquares := []int{1, 4, 9, 16, 25}
	squares := Map(ints, func(n int) int {
		return n * n
	})
	for i, v := range squares {
		if v != expectedSquares[i] {
			t.Errorf("Test 1: Expected %d, got %d", expectedSquares[i], v)
		}
	}

	// Test 2: Transforming a slice of strings to uppercase
	strs := []string{"go", "is", "awesome"}
	expectedUppercase := []string{"GO", "IS", "AWESOME"}
	uppercased := Map(strs, func(s string) string {
		return strings.ToUpper(s)
	})
	for i, v := range uppercased {
		if v != expectedUppercase[i] {
			t.Errorf("Test 2: Expected %s, got %s", expectedUppercase[i], v)
		}
	}

	// Test 3: Transforming an empty slice
	empty := []int{}
	expectedEmpty := []int{}
	emptyResult := Map(empty, func(n int) int {
		return n * n
	})
	if len(emptyResult) != len(expectedEmpty) {
		t.Errorf("Test 3: Expected %v, got %v", expectedEmpty, emptyResult)
	}

	// Test 4: Transforming a slice of strings (edge case)
	strsEdge := []string{"a", "b", "c"}
	expectedEdge := []string{"A", "B", "C"}
	edgeResult := Map(strsEdge, func(s string) string {
		return strings.ToUpper(s)
	})
	for i, v := range edgeResult {
		if v != expectedEdge[i] {
			t.Errorf("Test 4: Expected %s, got %s", expectedEdge[i], v)
		}
	}
}

func TestFlatten(t *testing.T) {
	// Test 1: Flatten a 2D slice of integers
	ints := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8},
	}
	expectedInts := []int{1, 2, 3, 4, 5, 6, 7, 8}
	flattenedInts := Flatten(ints)
	for i, v := range flattenedInts {
		if v != expectedInts[i] {
			t.Errorf("Test 1: Expected %d, got %d", expectedInts[i], v)
		}
	}

	// Test 2: Flatten a 2D slice of strings
	strs := [][]string{
		{"hello", "world"},
		{"go", "is", "awesome"},
	}
	expectedStrs := []string{"hello", "world", "go", "is", "awesome"}
	flattenedStrs := Flatten(strs)
	for i, v := range flattenedStrs {
		if v != expectedStrs[i] {
			t.Errorf("Test 2: Expected %s, got %s", expectedStrs[i], v)
		}
	}

	// Test 3: Flatten an empty 2D slice
	empty := [][]int{}
	expectedEmpty := []int{}
	flattenedEmpty := Flatten(empty)
	if len(flattenedEmpty) != len(expectedEmpty) {
		t.Errorf("Test 3: Expected %v, got %v", expectedEmpty, flattenedEmpty)
	}

	// Test 4: Flatten a 2D slice with empty inner slices
	withEmptySubArrays := [][]int{
		{1, 2},
		{},
		{3, 4},
	}
	expectedWithEmptySubArrays := []int{1, 2, 3, 4}
	flattenedWithEmptySubArrays := Flatten(withEmptySubArrays)
	for i, v := range flattenedWithEmptySubArrays {
		if v != expectedWithEmptySubArrays[i] {
			t.Errorf("Test 4: Expected %d, got %d", expectedWithEmptySubArrays[i], v)
		}
	}
}
