package runtime

import (
	"github.com/nicolito128/envo-lang/object"
)

type Environment struct {
	bindings map[string]object.Object
	rules    map[string]*object.DefineList
	outer    *Environment
}

func NewEnvironment() *Environment {
	env := &Environment{
		bindings: make(map[string]object.Object),
		rules:    make(map[string]*object.DefineList),
	}
	attachBuiltins(env)
	return env
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (object.Object, bool) {
	val, ok := e.bindings[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}
	return val, ok
}

func (e *Environment) GetLocal(name string) (object.Object, bool) {
	val, ok := e.bindings[name]
	return val, ok
}

func (e *Environment) Set(name string, val object.Object) object.Object {
	e.bindings[name] = val
	return val
}

func (e *Environment) Def(receiver string, def *object.Define) {
	if _, ok := e.rules[receiver]; !ok {
		e.rules[receiver] = &object.DefineList{Elems: make([]object.Object, 0)}
	}

	defList, ok := e.rules[receiver]
	if ok {
		defList.Elems = append(defList.Elems, def)
	}
}

func (e *Environment) GetDef(receiver string) (*object.DefineList, bool) {
	rules, ok := e.rules[receiver]
	if !ok && e.outer != nil {
		return e.outer.GetDef(receiver)
	}
	return rules, ok
}
