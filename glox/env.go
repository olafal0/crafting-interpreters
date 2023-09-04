package glox

import "fmt"

type Environment struct {
	vars map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		vars: map[string]any{},
	}
}

func (e *Environment) Get(name string) (any, bool) {
	v, ok := e.vars[name]
	return v, ok
}

func (e *Environment) Declare(name string, val any) error {
	_, ok := e.vars[name]
	if ok {
		return fmt.Errorf("redeclaration of var %s", name)
	}
	e.vars[name] = val
	return nil
}
