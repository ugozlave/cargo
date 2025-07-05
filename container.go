package cargo

import (
	"fmt"
	"reflect"
	"slices"
)

type Container struct {
	Services KeyValue[reflect.Type, []*Service]
	Scopes   KeyValue[string, *Scope]
}

func New() *Container {
	return &Container{
		Services: NewCollection[reflect.Type, []*Service](slices.Clone),
		Scopes:   NewCollection[string, *Scope](nil),
	}
}

func (c *Container) Register(key reflect.Type, value reflect.Type, builder func(ctx BuilderContext) any) {
	if builder == nil {
		panic("builder function cannot be nil")
	}
	if !value.AssignableTo(key) {
		panic(fmt.Sprintf("type %v is not assignable to %v", value, key))
	}
	if reflect.TypeOf(builder(nil)) != value {
		panic(fmt.Sprintf("builder function return type %v does not match registered type %v", reflect.TypeOf(builder(nil)), value))
	}
	if !c.Services.Has(key) {
		c.Services.Set(key, make([]*Service, 0, 1))
	}
	services, _ := c.Services.Get(key)
	c.Services.Set(key, append(services, &Service{Build: builder, Type: value}))
}

func (c *Container) Build(key reflect.Type, ctx BuilderContext) any {
	services, ok := c.Services.Get(key)
	if !ok || len(services) < 1 {
		return nil
	}
	return services[len(services)-1].Build(ctx)
}

func (c *Container) MustBuild(key reflect.Type, ctx BuilderContext) any {
	services, ok := c.Services.Get(key)
	if !ok || len(services) < 1 {
		panic(fmt.Sprintf("service %v not found", key))
	}
	return services[len(services)-1].Build(ctx)
}

func (c *Container) Get(key reflect.Type, name string, ctx BuilderContext) any {
	scope, ok := c.Scopes.Get(name)
	if !ok {
		return nil
	}
	instance, ok := scope.Instances.Get(key)
	if !ok {
		instance = c.Build(key, ctx)
		scope.Instances.Set(key, instance)
	}
	return instance
}

func (c *Container) MustGet(key reflect.Type, name string, ctx BuilderContext) any {
	scope, ok := c.Scopes.Get(name)
	if !ok {
		panic(fmt.Sprintf("scope %s not found", name))
	}
	instance, ok := scope.Instances.Get(key)
	if !ok {
		instance = c.MustBuild(key, ctx)
		scope.Instances.Set(key, instance)
	}
	return instance
}

func (c *Container) Inspect() {
	fmt.Println("Services:")
	for key, services := range c.Services.Map() {
		fmt.Printf("%v:\n", key)
		for _, service := range services {
			fmt.Printf(".   %v\n", service.Type)
		}
	}
	fmt.Println("Scopes:")
	for name, scope := range c.Scopes.Map() {
		fmt.Printf("%v:\n", name)
		for key := range scope.Instances.Map() {
			fmt.Printf(".   %v\n", key)
		}
	}
}

func (c *Container) CreateScope(name string) {
	if !c.Scopes.Has(name) {
		c.Scopes.Set(name, NewScope())
	}
}

func (c *Container) DeleteScope(name string) {
	if scope, ok := c.Scopes.Get(name); ok {
		scope.Instances.Clear()
	}
	c.Scopes.Del(name)
}
