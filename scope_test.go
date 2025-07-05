package cargo

import (
	"testing"
)

// Test basic scope creation
func TestNewScope(t *testing.T) {
	scope := NewScope()

	if scope == nil {
		t.Fatal("Expected NewScope() to return a non-nil scope")
	}

	if scope.Instances == nil {
		t.Fatal("Expected scope.Instances to be initialized")
	}
}

// Benchmark NewScope creation
func BenchmarkNewScope(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewScope()
	}
}
