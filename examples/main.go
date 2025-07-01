package main

import (
	"github.com/ugozlave/cargo"
)

func main() {
	container := cargo.New()
	cargo.RegisterT(container, NewA)
	cargo.RegisterT(container, NewB)
	cargo.Inspect(container)
}

type A struct {
	Counter int
}

func NewA(ctx *cargo.BuilderContext) *A {
	return &A{}
}

type B struct {
	A *A
}

func NewB(ctx *cargo.BuilderContext) *B {
	return &B{
		A: cargo.MustGet[*A](ctx.Container, "scope", ctx),
	}
}
