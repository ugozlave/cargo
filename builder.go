package cargo

type Builder[T any] func(BuilderContext) T
