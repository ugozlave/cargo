package cargo

import (
	"fmt"
	"reflect"
	"slices"
)

type Container struct {
	services KeyValue[reflect.Type, []*Service]
	scopes   KeyValue[string, *Scope]
}

func New() *Container {
	return &Container{
		services: NewCollection[reflect.Type, []*Service](slices.Clone),
		scopes:   NewCollection[string, *Scope](nil),
	}
}

func (c *Container) Register(key reflect.Type, value reflect.Type, builder func(BuilderContext) any) {
	if builder == nil {
		panic("builder function cannot be nil")
	}
	if !value.AssignableTo(key) {
		panic(fmt.Sprintf("type %v is not assignable to %v", value, key))
	}
	if !c.services.Has(key) {
		c.services.Set(key, make([]*Service, 0, 1))
	}
	services, _ := c.services.Get(key)
	c.services.Set(key, append(services, &Service{Build: builder, Type: value}))
}

func (c *Container) Build(key reflect.Type, ctx BuilderContext) any {
	services, ok := c.services.Get(key)
	count := len(services)
	if !ok || count < 1 {
		return nil
	}
	return services[count-1].Build(ctx)
}

func (c *Container) MustBuild(key reflect.Type, ctx BuilderContext) any {
	services, ok := c.services.Get(key)
	count := len(services)
	if !ok || count < 1 {
		panic(fmt.Sprintf("service %v not found", key))
	}
	return services[count-1].Build(ctx)
}

func (c *Container) Get(key reflect.Type, name string, ctx BuilderContext) any {
	scope, ok := c.scopes.Get(name)
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
	scope, ok := c.scopes.Get(name)
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

func (c *Container) All(key reflect.Type, name string, ctx BuilderContext) []any {
	scope, ok := c.scopes.Get(name)
	if !ok {
		return []any{}
	}
	services, ok := c.services.Get(key)
	if !ok {
		return []any{}
	}
	all := make([]any, 0, len(services))
	for _, service := range services {
		instance, ok := scope.Instances.Get(service.Type)
		if !ok {
			instance = service.Build(ctx)
			scope.Instances.Set(service.Type, instance)
		}
		all = append(all, instance)
	}
	return all
}

func (c *Container) Inspect() {
	fmt.Println("Services:")
	for key, services := range c.services.Map() {
		fmt.Printf(".   %v:\n", key)
		for _, service := range services {
			fmt.Printf("    .   %v\n", service.Type)
		}
	}
	fmt.Println()
	fmt.Println("Scopes:")
	for name, scope := range c.scopes.Map() {
		fmt.Printf(".   %v:\n", name)
		for key := range scope.Instances.Map() {
			fmt.Printf("    .   %v\n", key)
		}
	}
}

func (c *Container) CreateScope(name string) {
	if !c.scopes.Has(name) {
		c.scopes.Set(name, NewScope())
	}
}

func (c *Container) DeleteScope(name string) {
	if scope, ok := c.scopes.Get(name); ok {
		scope.Instances.Clear()
	}
	c.scopes.Del(name)
}
