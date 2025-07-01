package cargo

import "context"

type BuilderContext interface {
	context.Context
	C() *Container
}
