package runtime

import (
	"fmt"
	"unique"

	"github.com/nicolito128/envo-lang/ast"
	"github.com/nicolito128/envo-lang/object"
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
		if object.IsError(result) {
			return result
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
	patt := Eval(node.Pattern, env)
	if object.IsError(patt) {
		return patt
	}
	obj := &object.Define{
		Receiver: receiver,
		Pattern:  patt,
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

	if defList, ok := env.GetDef(label); ok {
		argVal := rt.evalRawMessage(node.Body, env)
		if object.IsError(argVal) {
			return argVal
		}
		argVal = unwrapSingleWordResult(argVal)

		for _, elem := range defList.Elems {
			if object.IsError(elem) {
				return elem
			}

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

		return &object.Error{Message: fmt.Sprintf("unknown pattern to match: %s", argVal)}
	}

	if val, ok := env.Get(label); ok {
		if builtin, ok := val.(*object.Builtin); ok {
			argVal := rt.evalRawMessage(node.Body, env)
			if object.IsError(argVal) {
				return argVal
			}
			if rawArgs, ok := argVal.(*object.RawMessage); ok {
				return builtin.Fn(rawArgs)
			}
			return builtin.Fn(&object.RawMessage{Words: []object.Object{argVal}})
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
		elem := Eval(word, env)
		if object.IsError(elem) {
			return elem
		}
		obj.Words = append(obj.Words, elem)
	}

	return obj
}

func (rt *Runtime) evalBinaryOp(op string, left, right object.Object, env *Environment) object.Object {
	if object.IsError(left) {
		return left
	}
	if object.IsError(right) {
		return right
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
	case "&&":
		return &object.Bool{Value: object.IsTruthy(left) && object.IsTruthy(right)}
	case "||":
		return &object.Bool{Value: object.IsTruthy(left) || object.IsTruthy(right)}
	}

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
		if object.IsError(saved) {
			return saved
		}
		if ok {
			return rt.evalUnaryOp(op, saved, env)
		}
	case *object.Message:
		return rt.evalUnaryOp(op, o.Body, env)
	case *object.RawMessage:
		result := make([]object.Object, 0)
		for _, word := range o.Words {
			if object.IsError(word) {
				return word
			}

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

	switch p := pattern.(type) {
	case *object.Identifier:
		if val, ok := env.GetLocal(p.Name); ok {
			if object.IsError(val) {
				return false
			}
			return object.Equal(val, arg)
		}
		env.Set(p.Name, arg)
		return true

	case *object.RawMessage:
		if len(p.Words) != len(argElems) {
			return false
		}

		count := 0
		for i := range len(p.Words) {
			word := p.Words[i]
			if object.IsError(word) {
				return false
			}

			arg := argElems[i]
			if object.IsError(arg) {
				return false
			}

			if rt.matchPattern(word, arg, env) {
				count++
			}
		}

		if count != len(argElems) {
			return false
		}
		return true
	}

	if !object.Equal(pattern, arg) {
		return false
	}

	return true
}
