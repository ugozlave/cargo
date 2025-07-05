package cargo

import "reflect"

type Service struct {
	Build func(BuilderContext) any
	Type  reflect.Type
}

func NewService(build func(BuilderContext) any, t reflect.Type) *Service {
	return &Service{
		Build: build,
		Type:  t,
	}
}
