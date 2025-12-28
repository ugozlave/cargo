package cargo

import (
	"fmt"
	"slices"
)

type Container struct {
	services KeyValue[string, []*Service]
	scopes   KeyValue[string, *Scope]
}

func New() *Container {
	return &Container{
		services: NewCollection[string, []*Service](slices.Clone),
		scopes:   NewCollection[string, *Scope](nil),
	}
}

func (c *Container) Register(key string, value string, builder Builder[any]) {
	if builder == nil {
		panic("builder function cannot be nil")
	}
	if !c.services.Has(key) {
		c.services.Set(key, make([]*Service, 0, 1))
	}
	services, _ := c.services.Get(key)
	c.services.Set(key, append(services, &Service{Build: builder, Type: value}))
}

func (c *Container) Build(key string, ctx BuilderContext) any {
	services, ok := c.services.Get(key)
	count := len(services)
	if !ok || count < 1 {
		return nil
	}
	return services[count-1].Build(ctx)
}

func (c *Container) MustBuild(key string, ctx BuilderContext) any {
	services, ok := c.services.Get(key)
	count := len(services)
	if !ok || count < 1 {
		panic(fmt.Sprintf("service %v not found", key))
	}
	return services[count-1].Build(ctx)
}

func (c *Container) Builds(key string, ctx BuilderContext) []any {
	services, ok := c.services.Get(key)
	if !ok {
		return []any{}
	}
	all := make([]any, 0, len(services))
	for _, service := range services {
		all = append(all, service.Build(ctx))
	}
	return all
}

func (c *Container) Get(key string, name string, ctx BuilderContext) any {
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

func (c *Container) MustGet(key string, name string, ctx BuilderContext) any {
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

func (c *Container) Gets(key string, name string, ctx BuilderContext) []any {
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

func (c *Container) Close() {
	for scope := range c.scopes.Map() {
		c.DeleteScope(scope)
	}
	c.scopes.Clr()
	c.services.Clr()
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
		scope.Close()
		scope = nil
	}
	c.scopes.Del(name)
}
