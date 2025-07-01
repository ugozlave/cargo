package cargo

import "reflect"

type Service struct {
	Build func(ctx *BuilderContext) any
}

type Services Collection[reflect.Type, *Service]
