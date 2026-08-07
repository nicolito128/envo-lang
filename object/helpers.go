package object

import "slices"

func IsTruthy(obj Object) bool {
	if obj == nil {
		return false
	}

	switch o := obj.(type) {
	case *Bool:
		return o.Value
	case *Integer:
		return o.Value != 0
	case *Nil:
		return false
	default:
		return true
	}
}

func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}

func Equal(a, b Object) bool {
	if a == nil || b == nil {
		return a == b
	}

	if a.Type() != b.Type() {
		if unwrappedA := unwrapComparableValue(a); unwrappedA != nil {
			return Equal(unwrappedA, b)
		}
		if unwrappedB := unwrapComparableValue(b); unwrappedB != nil {
			return Equal(a, unwrappedB)
		}
		return false
	}

	switch valA := a.(type) {
	case *Integer:
		if valB, ok := b.(*Integer); ok {
			return valA.Value == valB.Value
		}
	case *Float:
		if valB, ok := b.(*Float); ok {
			return valA.Value == valB.Value
		}
	case *Complex:
		if valB, ok := b.(*Complex); ok {
			return valA.Value == valB.Value
		}
	case *String:
		if valB, ok := b.(*String); ok {
			return valA.Value == valB.Value
		}
	case *RawString:
		if valB, ok := b.(*RawString); ok {
			return valA.Value == valB.Value
		}
	case *Char:
		if valB, ok := b.(*Char); ok {
			return slices.Compare(valA.Value, valB.Value) == 0
		}
	case *Symbol:
		if valB, ok := b.(*Symbol); ok {
			return valA.H.Value() == valB.H.Value()
		}
	case *Bool:
		if valB, ok := b.(*Bool); ok {
			return valA.Value == valB.Value
		}
	case *Identifier:
		if valB, ok := b.(*Identifier); ok {
			return valA.Name == valB.Name
		}
	case *RawMessage:
		if valB, ok := b.(*RawMessage); ok {
			if len(valA.Words) != len(valB.Words) {
				return false
			}
			for i, wordA := range valA.Words {
				wordB := valB.Words[i]
				if !Equal(wordA, wordB) {
					return false
				}
			}
			return true
		}
	case *Message:
		if valB, ok := b.(*Message); ok {
			return Equal(valA.Body, valB.Body)
		}
	case *Nil:
		return b.Type() == NIL_OBJ
	}
	return false
}

func unwrapComparableValue(obj Object) Object {
	switch o := obj.(type) {
	case *RawMessage:
		if len(o.Words) == 1 {
			return o.Words[0]
		}
	case *Message:
		if o.Body != nil && len(o.Body.Words) == 1 {
			return o.Body.Words[0]
		}
	}
	return nil
}
