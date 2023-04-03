package app

import "fmt"

// Application is the interface for the application
type Application interface {
	Start() error
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
	fmt.Println("Fib(40):", a.Fib(40))

	// TODO: Add your application logic here

	return nil
}

// Fib returns nth fibonacci number
func (a *app) Fib(n int) int {
	if n <= 1 {
		return n
	}
	return a.Fib(n-1) + a.Fib(n-2)
}
