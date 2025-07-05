package main

import (
	"fmt"

	"github.com/ugozlave/cargo"
)

func main() {
	container := cargo.New()
	container.CreateScope("default")
	container.CreateScope("test")
	cargo.RegisterKV[IService](container, NewServiceA)
	cargo.RegisterKV[IService](container, NewServiceB)
	cargo.RegisterKV[*ServiceA](container, NewServiceA)
	cargo.RegisterKV[ILogger](container, NewLogger)
	cargo.RegisterT(container, NewLogger)
	cargo.Get[IService](container, "default", nil)
	cargo.Get[IService](container, "default", nil)
	s := cargo.Get[IService](container, "test", nil)
	fmt.Println(s.DoSomething())
	container.DeleteScope("test")
	cargo.Get[ILogger](container, "default", nil)
	cargo.Get[*Logger](container, "default", nil)
	cargo.Get[*ServiceA](container, "default", nil)
	services := cargo.All[IService](container, "default", nil)
	for _, service := range services {
		fmt.Println(service.DoSomething())
	}
	cargo.Inspect(container)
}

type IService interface {
	DoSomething() string
}

type ServiceA struct{}

func NewServiceA(_ cargo.BuilderContext) *ServiceA {
	return &ServiceA{}
}

func (s *ServiceA) DoSomething() string {
	return "ServiceA is doing something"
}

type ServiceB struct{}

func NewServiceB(_ cargo.BuilderContext) *ServiceB {
	return &ServiceB{}
}

func (s *ServiceB) DoSomething() string {
	return "ServiceB is doing something"
}

type ILogger interface {
	Log(message string)
}

type Logger struct{}

func NewLogger(_ cargo.BuilderContext) *Logger {
	return &Logger{}
}

func (l *Logger) Log(message string) {
	println("Log:", message)
}
