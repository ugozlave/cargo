package cargo

import (
	"sync"
	"testing"
)

// Test basic collection operations with default copy function
func TestCollectionBasicOperations(t *testing.T) {
	// Test with string keys and int values using default copy function
	collection := NewCollection[string, int](nil)

	// Test initial state
	if collection.Has("key1") {
		t.Error("Expected empty collection to not have key1")
	}

	value, ok := collection.Get("key1")
	if ok {
		t.Error("Expected Get to return false for non-existent key")
	}
	if value != 0 {
		t.Errorf("Expected zero value, got %v", value)
	}

	// Test Set and Get
	collection.Set("key1", 42)
	if !collection.Has("key1") {
		t.Error("Expected collection to have key1 after Set")
	}

	value, ok = collection.Get("key1")
	if !ok {
		t.Error("Expected Get to return true for existing key")
	}
	if value != 42 {
		t.Errorf("Expected value 42, got %v", value)
	}

	// Test overwriting existing key
	collection.Set("key1", 100)
	value, ok = collection.Get("key1")
	if !ok {
		t.Error("Expected Get to return true for existing key")
	}
	if value != 100 {
		t.Errorf("Expected value 100, got %v", value)
	}

	// Test Del
	collection.Del("key1")
	if collection.Has("key1") {
		t.Error("Expected collection to not have key1 after Del")
	}

	value, ok = collection.Get("key1")
	if ok {
		t.Error("Expected Get to return false for deleted key")
	}
	if value != 0 {
		t.Errorf("Expected zero value, got %v", value)
	}
}

// Test collection with custom copy function
func TestCollectionWithCustomCopyFunction(t *testing.T) {
	// Create a copy function for slices
	copySlice := func(s []int) []int {
		if s == nil {
			return nil
		}
		result := make([]int, len(s))
		copy(result, s)
		return result
	}

	collection := NewCollection[string](copySlice)

	// Set a slice
	original := []int{1, 2, 3}
	collection.Set("slice", original)

	// Get the slice - should be a copy
	retrieved, ok := collection.Get("slice")
	if !ok {
		t.Error("Expected to retrieve slice")
	}

	// Verify it's a copy (different memory address)
	if &retrieved[0] == &original[0] {
		t.Error("Expected retrieved slice to be a copy, not the same reference")
	}

	// Verify contents are the same
	if len(retrieved) != len(original) {
		t.Errorf("Expected same length, got %d vs %d", len(retrieved), len(original))
	}
	for i, v := range retrieved {
		if v != original[i] {
			t.Errorf("Expected same values at index %d: %d vs %d", i, v, original[i])
		}
	}

	// Verify multiple Get() calls return independent copies
	retrieved2, _ := collection.Get("slice")
	if &retrieved[0] == &retrieved2[0] {
		t.Error("Expected each Get() call to return a fresh copy")
	}

	// But both should have the same content
	for i, v := range retrieved2 {
		if v != retrieved[i] {
			t.Errorf("Expected same values in both retrieved copies at index %d: %d vs %d", i, v, retrieved[i])
		}
	}
}

// Test Map() method
func TestCollectionMap(t *testing.T) {
	collection := NewCollection[string, int](nil)

	// Add some items
	collection.Set("a", 1)
	collection.Set("b", 2)
	collection.Set("c", 3)

	// Get map representation
	mapResult := collection.Map()

	// Verify all items are present
	if len(mapResult) != 3 {
		t.Errorf("Expected map length 3, got %d", len(mapResult))
	}

	expectedValues := map[string]int{"a": 1, "b": 2, "c": 3}
	for key, expectedValue := range expectedValues {
		if value, exists := mapResult[key]; !exists {
			t.Errorf("Expected key '%s' to exist in map", key)
		} else if value != expectedValue {
			t.Errorf("Expected value %d for key '%s', got %d", expectedValue, key, value)
		}
	}

	// Verify modifying returned map doesn't affect collection
	mapResult["d"] = 4
	if collection.Has("d") {
		t.Error("Expected collection to be unaffected by modifications to returned map")
	}
}

// Test Map() method with custom copy function
func TestCollectionMapWithCopyFunction(t *testing.T) {
	// Test with pointer values to ensure copy function is used
	type TestStruct struct {
		Value int
	}

	copyFunc := func(ts *TestStruct) *TestStruct {
		if ts == nil {
			return nil
		}
		return &TestStruct{Value: ts.Value}
	}

	collection := NewCollection[string](copyFunc)

	// Add items
	item1 := &TestStruct{Value: 10}
	item2 := &TestStruct{Value: 20}
	collection.Set("item1", item1)
	collection.Set("item2", item2)

	// Get map
	mapResult := collection.Map()

	// Verify items are copied
	if mapResult["item1"] == item1 {
		t.Error("Expected map to contain copies, not original references")
	}
	if mapResult["item1"].Value != item1.Value {
		t.Error("Expected copied item to have same value")
	}
}

// Test Clone() method
func TestCollectionClone(t *testing.T) {
	original := NewCollection[string, int](nil)

	// Add some items
	original.Set("a", 1)
	original.Set("b", 2)
	original.Set("c", 3)

	// Clone the collection
	cloned := original.Clone()

	// Verify clone is a different instance
	if original == cloned {
		t.Error("Expected clone to be a different instance")
	}

	// Verify clone has same content
	for _, key := range []string{"a", "b", "c"} {
		if !cloned.Has(key) {
			t.Errorf("Expected cloned collection to have key '%s'", key)
		}

		originalValue, _ := original.Get(key)
		clonedValue, _ := cloned.Get(key)
		if originalValue != clonedValue {
			t.Errorf("Expected same value for key '%s': %d vs %d", key, originalValue, clonedValue)
		}
	}

	// Verify modifications to clone don't affect original
	cloned.Set("d", 4)
	if original.Has("d") {
		t.Error("Expected original collection to be unaffected by clone modifications")
	}

	// Verify modifications to original don't affect clone
	original.Set("e", 5)
	if cloned.Has("e") {
		t.Error("Expected cloned collection to be unaffected by original modifications")
	}
}

// Test Clone() method with custom copy function
func TestCollectionCloneWithCopyFunction(t *testing.T) {
	type TestStruct struct {
		Value int
	}

	copyFunc := func(ts *TestStruct) *TestStruct {
		if ts == nil {
			return nil
		}
		return &TestStruct{Value: ts.Value}
	}

	original := NewCollection[string](copyFunc)

	// Add items
	item1 := &TestStruct{Value: 10}
	original.Set("item1", item1)

	// Clone
	cloned := original.Clone()

	// Get items from both collections
	originalItem, _ := original.Get("item1")
	clonedItem, _ := cloned.Get("item1")

	// Verify they are different instances but same values
	if originalItem == clonedItem {
		t.Error("Expected cloned item to be a different instance")
	}
	if originalItem.Value != clonedItem.Value {
		t.Error("Expected cloned item to have same value")
	}

	// Modify original item's value
	originalItem.Value = 999
	clonedItem2, _ := cloned.Get("item1")
	if clonedItem2.Value == 999 {
		t.Error("Expected cloned item to be independent of original")
	}
}

// Test Clear operation
func TestCollectionClear(t *testing.T) {
	collection := NewCollection[string, int](nil)

	// Add some items
	collection.Set("a", 1)
	collection.Set("b", 2)
	collection.Set("c", 3)

	// Verify items exist
	if !collection.Has("a") || !collection.Has("b") || !collection.Has("c") {
		t.Error("Expected all items to exist before clear")
	}

	// Clear collection
	collection.Clear()

	// Verify items are gone
	if collection.Has("a") || collection.Has("b") || collection.Has("c") {
		t.Error("Expected no items to exist after clear")
	}

	// Verify Map() returns empty map
	mapResult := collection.Map()
	if len(mapResult) != 0 {
		t.Errorf("Expected empty map after clear, got %d items", len(mapResult))
	}

	// Verify we can still add items after clear
	collection.Set("new", 100)
	if !collection.Has("new") {
		t.Error("Expected to be able to add items after clear")
	}
}

// Test deleting non-existent key
func TestCollectionDelNonExistent(t *testing.T) {
	collection := NewCollection[string, int](nil)

	// Delete from empty collection - should not panic
	collection.Del("nonexistent")

	// Add an item and delete a different key
	collection.Set("exists", 42)
	collection.Del("nonexistent")

	// Verify existing item is still there
	if !collection.Has("exists") {
		t.Error("Expected existing item to remain after deleting non-existent key")
	}
}

// Test concurrent access (thread safety)
func TestCollectionConcurrentAccess(t *testing.T) {
	collection := NewCollection[int, int](nil)
	var wg sync.WaitGroup

	// Number of goroutines and operations per goroutine
	numGoroutines := 10
	operationsPerGoroutine := 100

	// Concurrent writes
	for i := range numGoroutines {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := goroutineID*operationsPerGoroutine + j
				collection.Set(key, key*2)
			}
		}(i)
	}

	// Concurrent reads
	for i := range numGoroutines {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				key := goroutineID*operationsPerGoroutine + j
				collection.Get(key)
				collection.Has(key)
			}
		}(i)
	}

	// Concurrent Map() and Clone() operations
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				collection.Map()
				collection.Clone()
			}
		}()
	}

	// Concurrent deletes
	for i := range numGoroutines / 2 {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine/2; j++ {
				key := goroutineID*operationsPerGoroutine + j
				collection.Del(key)
			}
		}(i)
	}

	wg.Wait()

	// Test should complete without race conditions or panics
}

// Test with different types
func TestCollectionDifferentTypes(t *testing.T) {
	// Test with int keys and string values
	intCollection := NewCollection[int, string](nil)
	intCollection.Set(1, "one")
	intCollection.Set(2, "two")

	if !intCollection.Has(1) {
		t.Error("Expected collection to have key 1")
	}

	value, ok := intCollection.Get(2)
	if !ok || value != "two" {
		t.Errorf("Expected 'two', got %v", value)
	}

	// Test with string keys and pointer values
	type TestStruct struct {
		Name string
		ID   int
	}

	ptrCollection := NewCollection[string, *TestStruct](nil)
	testObj := &TestStruct{Name: "test", ID: 1}
	ptrCollection.Set("obj1", testObj)

	retrievedObj, ok := ptrCollection.Get("obj1")
	if !ok {
		t.Error("Expected to retrieve test object")
	}
	if retrievedObj != testObj {
		t.Error("Expected same pointer reference (default copy)")
	}
	if retrievedObj.Name != "test" || retrievedObj.ID != 1 {
		t.Errorf("Expected Name='test' ID=1, got Name='%s' ID=%d", retrievedObj.Name, retrievedObj.ID)
	}
}

// Test KeyValue interface compliance
func TestCollectionKeyValueInterface(t *testing.T) {
	// Verify Collection implements KeyValue interface
	var _ KeyValue[string, int] = NewCollection[string, int](nil)

	// This test will fail at compile time if Collection doesn't implement KeyValue
}

// Benchmark tests
func BenchmarkCollectionSet(b *testing.B) {
	collection := NewCollection[int, int](nil)

	for i := 0; b.Loop(); i++ {
		collection.Set(i, i*2)
	}
}

func BenchmarkCollectionGet(b *testing.B) {
	collection := NewCollection[int, int](nil)

	// Pre-populate
	for i := range 1000 {
		collection.Set(i, i*2)
	}

	for i := 0; b.Loop(); i++ {
		collection.Get(i % 1000)
	}
}

func BenchmarkCollectionMap(b *testing.B) {
	collection := NewCollection[int, int](nil)

	// Pre-populate
	for i := range 100 {
		collection.Set(i, i*2)
	}

	for b.Loop() {
		collection.Map()
	}
}

func BenchmarkCollectionClone(b *testing.B) {
	collection := NewCollection[int, int](nil)

	// Pre-populate
	for i := range 100 {
		collection.Set(i, i*2)
	}

	for b.Loop() {
		collection.Clone()
	}
}

func BenchmarkCollectionConcurrentRead(b *testing.B) {
	collection := NewCollection[int, int](nil)

	// Pre-populate
	for i := range 1000 {
		collection.Set(i, i*2)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			collection.Get(i % 1000)
			i++
		}
	})
}

func BenchmarkCollectionConcurrentWrite(b *testing.B) {
	collection := NewCollection[int, int](nil)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			collection.Set(i, i*2)
			i++
		}
	})
}
