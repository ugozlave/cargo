package cargo

import (
	"fmt"
	"reflect"
)

type Container struct {
	Services Services
	Scopes   Scopes
}

func New() *Container {
	return &Container{
		Services: make(Services),
		Scopes:   make(Scopes),
	}
}

func (c *Container) Register(key reflect.Type, value reflect.Type, builder func(ctx BuilderContext) any) {
	if !value.AssignableTo(key) {
		panic(fmt.Sprintf("type %v is not assignable to %v", value, key))
	}
	c.Services[key] = &Service{Build: builder}
}

func (c *Container) Build(key reflect.Type, ctx BuilderContext) any {
	service, ok := c.Services[key]
	if !ok {
		return nil
	}
	return service.Build(ctx)
}

func (c *Container) MustBuild(key reflect.Type, ctx BuilderContext) any {
	service, ok := c.Services[key]
	if !ok {
		panic(fmt.Sprintf("service %v not found", key))
	}
	return service.Build(ctx)
}

func (c *Container) Get(key reflect.Type, name string, ctx BuilderContext) any {
	scope, ok := c.Scopes[name]
	if !ok {
		return nil
	}
	instance, ok := scope.Instances[key]
	if !ok {
		instance = c.Build(key, ctx)
		scope.Instances[key] = instance
	}
	return instance
}

func (c *Container) MustGet(key reflect.Type, name string, ctx BuilderContext) any {
	scope, ok := c.Scopes[name]
	if !ok {
		panic(fmt.Sprintf("scope %s not found", name))
	}
	instance, ok := scope.Instances[key]
	if !ok {
		instance = c.MustBuild(key, ctx)
		scope.Instances[key] = instance
	}
	return instance
}

func (c *Container) Inspect() {
	fmt.Printf("%v\n", c.Services)
}
