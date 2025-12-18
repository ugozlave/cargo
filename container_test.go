package cargo

import (
	"fmt"
	"reflect"
	"testing"
)

// Test basic container creation
func TestContainerNew(t *testing.T) {
	container := New()

	if container == nil {
		t.Fatal("Expected New() to return a non-nil container")
	}

	if container.services == nil {
		t.Fatal("Expected container.Services to be initialized")
	}

	if container.scopes == nil {
		t.Fatal("Expected container.Scopes to be initialized")
	}
}

// Test container creation independence
func TestContainerNewIndependence(t *testing.T) {
	container1 := New()
	container2 := New()

	if container1 == container2 {
		t.Error("Expected New() to return different container instances")
	}

	if container1.services == container2.services {
		t.Error("Expected different containers to have independent Services")
	}

	if container1.scopes == container2.scopes {
		t.Error("Expected different containers to have independent Scopes")
	}
}

// Test Register function
func TestContainerRegister(t *testing.T) {
	container := New()
	stringPtrType := reflect.TypeOf((*string)(nil))

	// Test registering a valid service
	builder := func(ctx BuilderContext) any {
		s := "test service"
		return &s
	}

	container.Register(stringPtrType, stringPtrType, builder)

	// Verify service was registered
	if !container.services.Has(stringPtrType) {
		t.Error("Expected service to be registered")
	}

	services, _ := container.services.Get(stringPtrType)
	if len(services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(services))
	}

	if services[0].Type != stringPtrType {
		t.Errorf("Expected service type %v, got %v", stringPtrType, services[0].Type)
	}
}

// Test Register with multiple services for same key
func TestContainerRegisterMultiple(t *testing.T) {
	container := New()
	stringType := reflect.TypeOf("")

	// Register first service
	builder1 := func(ctx BuilderContext) any {
		return "service1"
	}
	container.Register(stringType, stringType, builder1)

	// Register second service for same key
	builder2 := func(ctx BuilderContext) any {
		return "service2"
	}
	container.Register(stringType, stringType, builder2)

	// Verify both services are registered
	services, _ := container.services.Get(stringType)
	if len(services) != 2 {
		t.Errorf("Expected 2 services, got %d", len(services))
	}
}

// Test Register panic conditions
func TestContainerRegisterPanics(t *testing.T) {
	container := New()
	stringType := reflect.TypeOf("")
	intType := reflect.TypeOf(0)

	// Test nil builder panic
	t.Run("nil builder", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for nil builder")
			}
		}()
		container.Register(stringType, stringType, nil)
	})

	// Test non-assignable type panic
	t.Run("non-assignable type", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for non-assignable type")
			}
		}()
		builder := func(ctx BuilderContext) any {
			return "string"
		}
		container.Register(intType, stringType, builder)
	})
}

// Test Build function
func TestContainerBuild(t *testing.T) {
	container := New()
	ctx := t.Context()
	stringType := reflect.TypeOf("")

	// Test building non-existent service
	result := container.Build(stringType, ctx)
	if result != nil {
		t.Errorf("Expected nil for non-existent service, got %v", result)
	}

	// Register a service
	builder := func(ctx BuilderContext) any {
		return "test service"
	}
	container.Register(stringType, stringType, builder)

	// Test building existing service
	result = container.Build(stringType, ctx)
	if result != "test service" {
		t.Errorf("Expected 'test service', got %v", result)
	}
}

// Test Build returns latest registered service
func TestContainerBuildLatest(t *testing.T) {
	container := New()
	ctx := t.Context()
	stringType := reflect.TypeOf("")

	// Register first service
	builder1 := func(ctx BuilderContext) any {
		return "service1"
	}
	container.Register(stringType, stringType, builder1)

	// Register second service
	builder2 := func(ctx BuilderContext) any {
		return "service2"
	}
	container.Register(stringType, stringType, builder2)

	// Build should return the latest registered service
	result := container.Build(stringType, ctx)
	if result != "service2" {
		t.Errorf("Expected 'service2' (latest), got %v", result)
	}
}

// Test MustBuild function
func TestContainerMustBuild(t *testing.T) {
	container := New()
	ctx := t.Context()
	stringType := reflect.TypeOf("")

	// Test MustBuild with non-existent service (should panic)
	t.Run("non-existent service", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for non-existent service")
			}
		}()
		container.MustBuild(stringType, ctx)
	})

	// Register a service
	builder := func(ctx BuilderContext) any {
		return "test service"
	}
	container.Register(stringType, stringType, builder)

	// Test MustBuild with existing service
	result := container.MustBuild(stringType, ctx)
	if result != "test service" {
		t.Errorf("Expected 'test service', got %v", result)
	}
}

// Test CreateScope function
func TestContainerCreateScope(t *testing.T) {
	container := New()

	// Test creating new scope
	container.CreateScope("test-scope")

	if !container.scopes.Has("test-scope") {
		t.Error("Expected scope to be created")
	}

	scope, _ := container.scopes.Get("test-scope")
	if scope == nil {
		t.Error("Expected scope to be non-nil")
	}

	// Test creating scope that already exists (should not create duplicate)
	container.CreateScope("test-scope")

	scopes := container.scopes.Map()
	count := 0
	for name := range scopes {
		if name == "test-scope" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("Expected 1 scope named 'test-scope', got %d", count)
	}
}

// Test DeleteScope function
func TestContainerDeleteScope(t *testing.T) {
	container := New()

	// Create and populate a scope
	container.CreateScope("test-scope")
	scope, _ := container.scopes.Get("test-scope")
	stringType := reflect.TypeOf("")
	scope.Instances.Set(stringType, "test instance")

	// Verify scope exists and has instances
	if !container.scopes.Has("test-scope") {
		t.Error("Expected scope to exist before deletion")
	}
	if !scope.Instances.Has(stringType) {
		t.Error("Expected scope to have instance before deletion")
	}

	// Delete the scope
	container.DeleteScope("test-scope")

	// Verify scope is deleted
	if container.scopes.Has("test-scope") {
		t.Error("Expected scope to be deleted")
	}

	// Test deleting non-existent scope (should not panic)
	container.DeleteScope("non-existent-scope")
}

// Test Get function
func TestContainerGet(t *testing.T) {
	container := New()
	ctx := t.Context()
	stringType := reflect.TypeOf("")

	// Test getting from non-existent scope
	result := container.Get(stringType, "non-existent-scope", ctx)
	if result != nil {
		t.Errorf("Expected nil for non-existent scope, got %v", result)
	}

	// Create scope and register service
	container.CreateScope("test-scope")
	builder := func(ctx BuilderContext) any {
		return "scoped service"
	}
	container.Register(stringType, stringType, builder)

	// Test getting from scope (should build and cache)
	result1 := container.Get(stringType, "test-scope", ctx)
	if result1 != "scoped service" {
		t.Errorf("Expected 'scoped service', got %v", result1)
	}

	// Test getting again (should return cached instance)
	result2 := container.Get(stringType, "test-scope", ctx)
	if result2 != result1 {
		t.Error("Expected same instance from scope cache")
	}
}

// Test MustGet function
func TestContainerMustGet(t *testing.T) {
	container := New()
	ctx := t.Context()
	stringType := reflect.TypeOf("")

	// Test MustGet with non-existent scope (should panic)
	t.Run("non-existent scope", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for non-existent scope")
			}
		}()
		container.MustGet(stringType, "non-existent-scope", ctx)
	})

	// Create scope and register service
	container.CreateScope("test-scope")
	builder := func(ctx BuilderContext) any {
		return "scoped service"
	}
	container.Register(stringType, stringType, builder)

	// Test MustGet with existing scope
	result := container.MustGet(stringType, "test-scope", ctx)
	if result != "scoped service" {
		t.Errorf("Expected 'scoped service', got %v", result)
	}
}

// Test scope isolation
func TestContainerScopeIsolation(t *testing.T) {
	container := New()
	ctx := t.Context()
	stringType := reflect.TypeOf("")

	// Register a service that returns different values each time
	var counter int
	builder := func(ctx BuilderContext) any {
		counter++
		return fmt.Sprintf("instance-%d", counter)
	}
	container.Register(stringType, stringType, builder)

	// Create two scopes
	container.CreateScope("scope1")
	container.CreateScope("scope2")

	// Get instances from different scopes
	instance1 := container.Get(stringType, "scope1", ctx)
	instance2 := container.Get(stringType, "scope2", ctx)

	// They should be different instances
	if instance1 == instance2 {
		t.Error("Expected different instances from different scopes")
	}

	// Getting from same scope should return same instance
	instance1Again := container.Get(stringType, "scope1", ctx)
	if instance1 != instance1Again {
		t.Error("Expected same instance from same scope")
	}
}

// Test dependency injection scenario
func TestContainerDependencyInjection(t *testing.T) {
	container := New()
	ctx := t.Context()

	// Define simple types for dependency injection
	type Database struct {
		Name string
	}

	type UserService struct {
		DB *Database
	}

	databasePtr := reflect.TypeOf(&Database{})
	userServicePtr := reflect.TypeOf(&UserService{})

	// Register database service
	dbBuilder := func(ctx BuilderContext) any {
		return &Database{Name: "test-db"}
	}
	container.Register(databasePtr, databasePtr, dbBuilder)

	// Register user service that depends on database
	userBuilder := func(ctx BuilderContext) any {
		if ctx == nil {
			// Return a dummy instance for type checking during registration
			return &UserService{DB: &Database{Name: "dummy"}}
		}
		db := container.Build(databasePtr, ctx)
		return &UserService{DB: db.(*Database)}
	}
	container.Register(userServicePtr, userServicePtr, userBuilder)

	// Build user service
	result := container.Build(userServicePtr, ctx)
	userService, ok := result.(*UserService)
	if !ok {
		t.Errorf("Expected *UserService, got %T", result)
	}

	// Test that dependency was injected
	if userService.DB == nil {
		t.Error("Expected database dependency to be injected")
	}

	if userService.DB.Name != "test-db" {
		t.Errorf("Expected database name 'test-db', got '%s'", userService.DB.Name)
	}
}

// Benchmark container operations
func BenchmarkContainerNew(b *testing.B) {

	for b.Loop() {
		New()
	}
}

func BenchmarkContainerRegister(b *testing.B) {
	container := New()
	stringType := reflect.TypeOf("")
	builder := func(ctx BuilderContext) any {
		return "benchmark service"
	}

	for b.Loop() {
		container.Register(stringType, stringType, builder)
	}
}

func BenchmarkContainerBuild(b *testing.B) {
	container := New()
	ctx := b.Context()
	stringType := reflect.TypeOf("")
	builder := func(ctx BuilderContext) any {
		return "benchmark service"
	}
	container.Register(stringType, stringType, builder)

	for b.Loop() {
		container.Build(stringType, ctx)
	}
}

func BenchmarkContainerScopedGet(b *testing.B) {
	container := New()
	ctx := b.Context()
	stringType := reflect.TypeOf("")
	builder := func(ctx BuilderContext) any {
		return "benchmark service"
	}
	container.Register(stringType, stringType, builder)
	container.CreateScope("benchmark-scope")

	for b.Loop() {
		container.Get(stringType, "benchmark-scope", ctx)
	}
}
