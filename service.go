package cargo

type Service struct {
	Build Builder[any]
	Type  string
}

func NewService(build Builder[any], t string) *Service {
	return &Service{
		Build: build,
		Type:  t,
	}
}
