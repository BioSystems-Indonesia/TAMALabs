package util

// Unique filters unique items from a slice
func Unique[T comparable](items []T) []T {
	seen := make(map[T]struct{}) // Use a map to track seen items
	var result []T

	for _, item := range items {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func Contains[T comparable](ss []T, v T) bool {
	for _, s := range ss {
		if s == v {
			return true
		}
	}
	return false
}

func Filter[T any](ss []T, filterFunc func(T) bool) []T {
	var ret []T
	for _, s := range ss {
		if filterFunc(s) {
			ret = append(ret, s)
		}
	}
	return ret
}

func CompareSlices[T comparable](old, new []T) (toDelete []T, toCreate []T) {
	// Create maps to track element frequencies
	oldMap := make(map[T]int)
	newMap := make(map[T]int)

	// Count frequencies in old slice
	for _, item := range old {
		oldMap[item]++
	}

	// Count frequencies in new slice
	for _, item := range new {
		newMap[item]++
	}

	// Find elements to delete (present in old but not in new, or have higher frequency in old)
	for item, oldCount := range oldMap {
		newCount := newMap[item]
		if oldCount > newCount {
			// Add the item to toDelete as many times as it appears more frequently in old
			for i := 0; i < oldCount-newCount; i++ {
				toDelete = append(toDelete, item)
			}
		}
	}

	// Find elements to create (present in new but not in old, or have higher frequency in new)
	for item, newCount := range newMap {
		oldCount := oldMap[item]
		if newCount > oldCount {
			// Add the item to toCreate as many times as it appears more frequently in new
			for i := 0; i < newCount-oldCount; i++ {
				toCreate = append(toCreate, item)
			}
		}
	}

	return toDelete, toCreate
}

func Map[T any, U any](input []T, transform func(T) U) []U {
	var result []U
	for _, v := range input {
		result = append(result, transform(v))
	}
	return result
}

func Flatten[T any](input [][]T) []T {
	var result []T
	for _, subArray := range input {
		result = append(result, subArray...)
	}
	return result
}

// RemoveDuplicates removes duplicates from a slice of type T based on a key field.
// The keyFunc extracts a comparable key of type K from each item.
//
//	people := []Person{
//			{ID: 1, Name: "Alice"},
//			{ID: 2, Name: "Bob"},
//			{ID: 1, Name: "Alice"},
//			{ID: 3, Name: "Charlie"},
//			{ID: 2, Name: "Bob"},
//		}
//
//		uniquePeople := RemoveDuplicatesFromStruct(people, func(p Person) int {
//			return p.ID
//		})
func RemoveDuplicatesFromStruct[T any, K comparable](items []T, keyFunc func(T) K) []T {
	seen := make(map[K]struct{})
	result := make([]T, 0, len(items))
	for _, item := range items {
		key := keyFunc(item)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
