package cargo

import "reflect"

func From[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func RegisterT[T any](c *Container, builder func(ctx *BuilderContext) T) {
	c.Register(From[T](), From[T](), func(ctx *BuilderContext) any {
		return builder(ctx)
	})
}

func RegisterKV[K any, V any](c *Container, builder func(ctx *BuilderContext) V) {
	c.Register(From[K](), From[V](), func(ctx *BuilderContext) any {
		return builder(ctx)
	})
}

func Build[T any](c *Container, ctx *BuilderContext) T {
	return c.Build(From[T](), ctx).(T)
}

func Get[T any](c *Container, scope string, ctx *BuilderContext) T {
	return c.Get(From[T](), scope, ctx).(T)
}

func MustGet[T any](c *Container, scope string, ctx *BuilderContext) T {
	return c.MustGet(From[T](), scope, ctx).(T)
}

func Inspect(c *Container) {
	c.Inspect()
}
