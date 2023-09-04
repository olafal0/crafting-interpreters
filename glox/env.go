package glox

import "fmt"

type Environment struct {
	enclosing *Environment
	vars      map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		vars:      map[string]any{},
	}
}

func (e *Environment) Get(name string) (any, bool) {
	v, ok := e.vars[name]
	if ok {
		return v, ok
	}
	if e.enclosing == nil {
		return v, ok
	}
	return e.enclosing.Get(name)
}

func (e *Environment) Declare(name string, val any) error {
	_, ok := e.vars[name]
	if ok {
		return fmt.Errorf("redeclaration of var %s", name)
	}
	e.vars[name] = val
	return nil
}

func (e *Environment) Set(name string, val any) error {
	_, ok := e.vars[name]
	if ok {
		e.vars[name] = val
		return nil
	}
	if e.enclosing == nil {
		return fmt.Errorf("unknown var %s", name)
	}
	return e.enclosing.Set(name, val)
}
