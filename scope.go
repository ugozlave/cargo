package cargo

import "reflect"

type Scope struct {
	Instances KeyValue[reflect.Type, any]
}

func NewScope() *Scope {
	return &Scope{
		Instances: NewCollection[reflect.Type, any](nil),
	}
}
