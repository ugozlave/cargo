package cargo

import "context"

type BuilderContext struct {
	context.Context
	*Container
}
