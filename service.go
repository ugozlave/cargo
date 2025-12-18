package cargo

import "reflect"

type Service struct {
	Build Builder[any]
	Type  reflect.Type
}

func NewService(build Builder[any], t reflect.Type) *Service {
	return &Service{
		Build: build,
		Type:  t,
	}
}
