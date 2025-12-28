package cargo

type Scope struct {
	Instances KeyValue[string, any]
}

func NewScope() *Scope {
	return &Scope{
		Instances: NewCollection[string, any](nil),
	}
}

func (s *Scope) Close() {
	for _, instance := range s.Instances.Map() {
		if closer, ok := instance.(Closer); ok {
			closer.Close()
		}
		instance = nil
	}
	s.Instances.Clr()
}
