package runtime

import (
	"strings"

	"github.com/nicolito128/envo-lang/object"
)

func attachBuiltins(env *Environment) {
	for _, b := range builtinList {
		env.Set(b.Name, b)
	}
}

func compareValues(left, right object.Object) (int, bool) {
	switch l := left.(type) {
	case *object.Integer:
		if r, ok := right.(*object.Integer); ok {
			return int(l.Value - r.Value), true
		}
	case *object.Float:
		if r, ok := right.(*object.Float); ok {
			if l.Value < r.Value {
				return -1, true
			}
			if l.Value > r.Value {
				return 1, true
			}
			return 0, true
		}
		if r, ok := right.(*object.Integer); ok {
			if l.Value < float64(r.Value) {
				return -1, true
			}
			if l.Value > float64(r.Value) {
				return 1, true
			}
			return 0, true
		}
	case *object.String:
		if r, ok := right.(*object.String); ok {
			return strings.Compare(l.Value, r.Value), true
		}
	case *object.RawString:
		if r, ok := right.(*object.RawString); ok {
			return strings.Compare(l.Value, r.Value), true
		}
	case *object.Char:
		if r, ok := right.(*object.Char); ok {
			if len(l.Value) == 0 || len(r.Value) == 0 {
				return 0, true
			}
			return int(l.Value[0] - r.Value[0]), true
		}
	case *object.Symbol:
		if r, ok := right.(*object.Symbol); ok {
			return strings.Compare(l.H.Value(), r.H.Value()), true
		}
	}

	return 0, false
}

func unwrapSingleWordResult(obj object.Object) object.Object {
	switch raw := obj.(type) {
	case *object.RawMessage:
		if len(raw.Words) == 1 {
			return raw.Words[0]
		}
	case *object.Message:
		if raw.Body != nil && len(raw.Body.Words) == 1 {
			return raw.Body.Words[0]
		}
	}
	return obj
}
