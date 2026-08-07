package runtime

import (
	"fmt"
	"strings"
	"unique"

	"n128.xyz/n128/envo/ast"
	"n128.xyz/n128/envo/object"
)

var NIL = &object.Nil{}

func Eval(node ast.Node, env *Environment) object.Object {
	rt := New(node)
	rt.Use(env)
	return rt.Eval()
}

type Runtime struct {
	root ast.Node
	env  *Environment
}

func New(program ast.Node) *Runtime {
	rt := new(Runtime)
	rt.root = program
	return rt
}

func (rt *Runtime) Use(e *Environment) {
	rt.env = e
}

func (rt *Runtime) Eval() object.Object {
	if rt.env == nil {
		rt.env = NewEnvironment()
	}

	switch n := rt.root.(type) {
	case *ast.Program:
		return rt.evalProgram(n, rt.env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: n.Value}

	case *ast.ImaginaryLiteral:
		return &object.Complex{Value: n.Value}

	case *ast.BoolLiteral:
		return &object.Bool{Value: n.Value}

	case *ast.StringLiteral:
		return &object.String{Value: n.Value}

	case *ast.RawStringLiteral:
		return &object.RawString{Value: n.Value}

	case *ast.SymbolLiteral:
		return &object.Symbol{H: unique.Make(n.Name)}

	case *ast.CharLiteral:
		return &object.Char{Value: n.Value}

	case *ast.Identifier:
		return rt.evalIdentifier(n, rt.env)

	case *ast.Define:
		return rt.evalDefine(n, rt.env)

	case *ast.RawMessage:
		return rt.evalRawMessage(n, rt.env)

	case *ast.Message:
		return rt.evalMessage(n, rt.env)

	case *ast.BinaryOp:
		left := Eval(n.Left, rt.env)
		if object.IsError(left) {
			return left
		}
		right := Eval(n.Right, rt.env)
		if object.IsError(right) {
			return right
		}
		return rt.evalBinaryOp(n.Op, left, right, rt.env)

	case *ast.UnaryOp:
		operand := Eval(n.Operand, rt.env)
		if object.IsError(operand) {
			return operand
		}
		return rt.evalUnaryOp(n.Op, operand, rt.env)
	}

	return NIL
}

func (rt *Runtime) evalProgram(program *ast.Program, env *Environment) object.Object {
	var result object.Object = NIL
	for _, stmt := range program.Statements {
		result = Eval(stmt, env)
		if err, ok := result.(*object.Error); ok {
			return err
		}
	}
	return result
}

func (rt *Runtime) evalIdentifier(node *ast.Identifier, env *Environment) object.Object {
	name := node.Name

	def, ok := env.GetDef(name)
	if ok {
		return def
	}

	obj, ok := env.Get(name)
	if ok {
		return obj
	}

	ident := &object.Identifier{Name: node.Name}
	return ident
}

func (rt *Runtime) evalDefine(node *ast.Define, env *Environment) object.Object {
	receiver := node.Receiver.Name
	obj := &object.Define{
		Receiver: receiver,
		Pattern:  Eval(node.Pattern, env),
		Response: node.Response,
	}

	if _, ok := env.GetDef(obj.String()); !ok {
		env.Def(receiver, obj)
	}

	return obj
}

func (rt *Runtime) evalMessage(node *ast.Message, env *Environment) object.Object {
	if node.Body != nil && len(node.Body.Words) == 0 {
		if val, ok := env.Get(node.Label.Name); ok {
			return val
		}
	}
	label := node.Label.Name

	if v, ok := env.GetDef(label); ok {
		argVal := rt.evalRawMessage(node.Body, env)
		if object.IsError(argVal) {
			return argVal
		}
		argVal = unwrapSingleWordResult(argVal)

		for _, elem := range v.Elems {
			if def, ok := elem.(*object.Define); ok {
				execEnv := NewEnclosedEnvironment(env)
				if rt.matchPattern(def.Pattern, argVal, execEnv) {
					result := Eval(def.Response, execEnv)
					if object.IsError(result) {
						return result
					}
					return unwrapSingleWordResult(result)
				}
			}
		}
	}

	obj := &object.Message{
		Label: node.Label.Name,
		Body: &object.RawMessage{
			Words: make([]object.Object, 0),
		},
	}

	for _, word := range node.Body.Words {
		obj.Body.Words = append(obj.Body.Words, Eval(word, env))
	}

	return obj
}

func (rt *Runtime) evalRawMessage(raw *ast.RawMessage, env *Environment) object.Object {
	if raw == nil {
		return NIL
	}

	obj := &object.RawMessage{
		Words: make([]object.Object, 0),
	}

	for _, word := range raw.Words {
		obj.Words = append(obj.Words, Eval(word, env))
	}

	return obj
}

func (rt *Runtime) evalBinaryOp(op string, left, right object.Object, env *Environment) object.Object {
	switch vLeft := left.(type) {
	case *object.Integer:
		if vRight, ok := right.(*object.Float); ok {
			return rt.evalBinaryOp(op, &object.Float{Value: float64(vLeft.Value)}, vRight, env)
		}

		if vRight, ok := right.(*object.Integer); ok {
			switch op {
			case "+":
				return &object.Integer{Value: vLeft.Value + vRight.Value}
			case "-":
				return &object.Integer{Value: vLeft.Value - vRight.Value}
			case "*":
				return &object.Integer{Value: vLeft.Value * vRight.Value}
			case "/":
				return &object.Integer{Value: vLeft.Value / vRight.Value}
			case "%":
				return &object.Integer{Value: vLeft.Value % vRight.Value}
			}
		}

		if vRight, ok := right.(*object.RawMessage); ok {
			return rt.applyBinaryOpToMessage(op, left, vRight, env)
		}

		if vRight, ok := right.(*object.Message); ok {
			return rt.applyBinaryOpToMessage(op, left, vRight.Body, env)
		}

	case *object.Float:
		if vRight, ok := right.(*object.Integer); ok {
			return rt.evalBinaryOp(op, vLeft, &object.Float{Value: float64(vRight.Value)}, env)
		}

		if vRight, ok := right.(*object.Float); ok {
			switch op {
			case "+":
				return &object.Float{Value: vLeft.Value + vRight.Value}
			case "-":
				return &object.Float{Value: vLeft.Value - vRight.Value}
			case "*":
				return &object.Float{Value: vLeft.Value * vRight.Value}
			case "/":
				return &object.Float{Value: vLeft.Value / vRight.Value}
			}
		}

		if vRight, ok := right.(*object.RawMessage); ok {
			return rt.applyBinaryOpToMessage(op, left, vRight, env)
		}

		if vRight, ok := right.(*object.Message); ok {
			return rt.applyBinaryOpToMessage(op, left, vRight.Body, env)
		}

	case *object.String:
		litRight, ok := right.Literal().(string)
		if ok {
			switch op {
			case "+":
				return &object.String{Value: vLeft.Value + litRight}
			}
		}

	case *object.RawString:
		if vRight, ok := right.(*object.RawString); ok {
			return &object.RawString{Value: vLeft.Value + vRight.Value}
		} else {
			litRight, ok := right.Literal().(string)
			if ok {
				return &object.String{Value: vLeft.Value + litRight}
			}
		}

	case *object.Identifier:
		saved, ok := env.Get(vLeft.Name)
		if ok {
			return rt.evalBinaryOp(op, saved, right, env)
		}

	case *object.Message:
		if def, ok := env.GetDef(vLeft.Label); ok {
			for _, elem := range def.Elems {
				if def, ok := elem.(*object.Define); ok {
					execEnv := NewEnclosedEnvironment(env)
					if rt.matchPattern(def.Pattern, right, execEnv) {
						return Eval(def.Response, execEnv)
					}
				}
			}
		}

		if vRight, ok := right.(*object.Message); ok {
			return rt.evalBinaryOp(op, vLeft.Body, vRight.Body, env)
		}

	case *object.RawMessage:
		if vRight, ok := right.(*object.RawMessage); ok {
			if len(vLeft.Words) != len(vRight.Words) {
				return &object.Error{Message: "invalid length between messages"}
			}

			result := make([]object.Object, 0)
			for i, wl := range vLeft.Words {
				wr := vRight.Words[i]

				elem := rt.evalBinaryOp(op, wl, wr, env)
				if object.IsError(elem) {
					return elem
				}
				result = append(result, elem)
			}

			return &object.RawMessage{Words: result}
		}

		if vRight, ok := right.(*object.Integer); ok {
			return rt.applyBinaryOpToMessage(op, left, &object.RawMessage{Words: []object.Object{vRight}}, env)
		}

		if vRight, ok := right.(*object.Float); ok {
			return rt.applyBinaryOpToMessage(op, left, &object.RawMessage{Words: []object.Object{vRight}}, env)
		}
	}

	switch op {
	case "==":
		return &object.Bool{Value: object.Equal(left, right)}
	case "!=":
		return &object.Bool{Value: !object.Equal(left, right)}
	case "<":
		if cmp, ok := compareValues(left, right); ok {
			return &object.Bool{Value: cmp < 0}
		}
	case "<=":
		if cmp, ok := compareValues(left, right); ok {
			return &object.Bool{Value: cmp <= 0}
		}
	case ">":
		if cmp, ok := compareValues(left, right); ok {
			return &object.Bool{Value: cmp > 0}
		}
	case ">=":
		if cmp, ok := compareValues(left, right); ok {
			return &object.Bool{Value: cmp >= 0}
		}
	}

	return &object.Error{Message: fmt.Sprintf("operation not supported: %s %s %s", left.Type(), op, right.Type())}
}

func (rt *Runtime) applyBinaryOpToMessage(op string, left object.Object, right *object.RawMessage, env *Environment) object.Object {
	result := make([]object.Object, 0, len(right.Words))
	for _, word := range right.Words {
		elem := rt.evalBinaryOp(op, left, word, env)
		if object.IsError(elem) {
			return elem
		}
		result = append(result, elem)
	}
	return unwrapSingleWordResult(&object.RawMessage{Words: result})
}

func (rt *Runtime) evalUnaryOp(op string, operand object.Object, env *Environment) object.Object {
	switch o := operand.(type) {
	case *object.Integer:
		switch op {
		case "+":
			return &object.Integer{Value: +o.Value}
		case "-":
			return &object.Integer{Value: -o.Value}
		}
	case *object.Float:
		switch op {
		case "+":
			return &object.Float{Value: +o.Value}
		case "-":
			return &object.Float{Value: -o.Value}
		}
	case *object.Bool:
		switch op {
		case "!":
			return &object.Bool{Value: !o.Value}
		}
	case *object.Identifier:
		saved, ok := env.Get(o.Name)
		if ok {
			return rt.evalUnaryOp(op, saved, env)
		}
	case *object.Message:
		return rt.evalUnaryOp(op, o.Body, env)
	case *object.RawMessage:
		result := make([]object.Object, 0)
		for _, word := range o.Words {
			elem := rt.evalUnaryOp(op, word, env)
			if object.IsError(elem) {
				return elem
			}

			result = append(result, elem)
		}
		return &object.RawMessage{Words: result}
	}

	return &object.Error{Message: fmt.Sprintf("unsupported unary operation: %s %s", op, operand.Type())}
}

func (rt *Runtime) matchPattern(pattern object.Object, arg object.Object, env *Environment) bool {
	arg = unwrapSingleWordResult(arg)
	if object.Equal(pattern, arg) {
		return true
	}

	var argElems []object.Object
	if rawArg, ok := arg.(*object.RawMessage); ok {
		argElems = rawArg.Words
	} else if arg == NIL {
		argElems = []object.Object{}
	} else {
		argElems = []object.Object{arg}
	}

	for _, arg := range argElems {
		switch p := pattern.(type) {
		case *object.Identifier:
			if val, ok := env.GetLocal(p.Name); ok {
				return object.Equal(val, arg)
			}
			env.Set(p.Name, arg)
			return true

		case *object.RawMessage:
			for _, word := range p.Words {
				if rt.matchPattern(word, arg, env) {
					return true
				}
			}
		}

		if !object.Equal(pattern, arg) {
			return false
		}
	}

	return true
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
