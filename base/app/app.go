package app

import "fmt"

// Application is the interface for the application
type Application interface {
	Start() error
	Hi() string
}

// app is the implementation of the application
type app struct {
}

// NewApplication creates a new application
func NewApplication() Application {
	return &app{}
}

func (a *app) Start() error {
	fmt.Println("Starting application")
	fmt.Println(a.Hi())

	// TODO: Add your application logic here

	return nil
}

func (a *app) Hi() string {
	return "Hi Mom!"
}
