package cargo

import "reflect"

type Scope struct {
	Instances Collection[reflect.Type, any]
}

type Scopes Collection[string, *Scope]

func (scopes Scopes) Create(name string) {
	if _, ok := scopes[name]; !ok {
		scopes[name] = &Scope{Instances: make(Collection[reflect.Type, any])}
	}
}

func (scopes Scopes) Delete(name string) {
	if scope, ok := scopes[name]; ok {
		for t := range scope.Instances {
			delete(scope.Instances, t)
		}

	}
	delete(scopes, name)
}
