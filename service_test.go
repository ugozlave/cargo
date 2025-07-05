package cargo

import (
	"reflect"
	"testing"
)

// Test basic service creation
func TestNewService(t *testing.T) {
	builder := func(ctx BuilderContext) any {
		return &struct{}{}
	}
	value := reflect.TypeOf(&struct{}{})
	service := NewService(builder, value)

	if service == nil {
		t.Fatal("Expected NewService to return a non-nil service")
	}

	if service.Build == nil {
		t.Error("Expected service.Build to be non-nil")
	}

	if service.Type != value {
		t.Errorf("Expected service.Type to be %v, got %v", value, service.Type)
	}
}

// Benchmark service creation
func BenchmarkNewService(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewService(nil, nil)
	}
}
