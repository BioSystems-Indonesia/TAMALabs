package util

import (
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
