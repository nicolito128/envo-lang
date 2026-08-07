package runtime

import (
	"fmt"
	"strings"

	"n128.xyz/n128/envo/object"
)

// builtin functions
var (
	print = &object.Builtin{
		Name: "print",
		Fn: func(rm *object.RawMessage) object.Object {
			var s strings.Builder
			for _, word := range rm.Words {
				switch w := word.(type) {
				case *object.Integer,
					*object.Float,
					*object.Complex,
					*object.Char,
					*object.String,
					*object.Bool,
					*object.Symbol:
					s.WriteString(fmt.Sprintf("%v", w.Literal()))
				default:
					s.WriteString(w.String())
				}
			}
			fmt.Print(s.String())
			return NIL
		},
	}

	println = &object.Builtin{
		Name: "println",
		Fn: func(rm *object.RawMessage) object.Object {
			var s strings.Builder
			for _, word := range rm.Words {
				switch w := word.(type) {
				case *object.Integer,
					*object.Float,
					*object.Complex,
					*object.Char,
					*object.String,
					*object.Bool,
					*object.Symbol:
					s.WriteString(fmt.Sprintf("%v", w.Literal()))
				default:
					s.WriteString(w.String())
				}
			}
			fmt.Println(s.String())
			return NIL
		},
	}

	scan = &object.Builtin{
		Name: "scan",
		Fn: func(rm *object.RawMessage) object.Object {
			var input string
			_, err := fmt.Scan(&input)
			if err != nil {
				return &object.Error{Message: err.Error()}
			}
			return &object.String{Value: input}
		},
	}

	scanln = &object.Builtin{
		Name: "scanln",
		Fn: func(rm *object.RawMessage) object.Object {
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				return &object.Error{Message: err.Error()}
			}
			return &object.String{Value: input}
		},
	}

	lenFn = &object.Builtin{
		Name: "len",
		Fn: func(rm *object.RawMessage) object.Object {
			res := &object.Integer{}

			if len(rm.Words) > 1 {
				return &object.Error{Message: "len builtin do not allows more than 1 argument"}
			}

			switch w := rm.Words[0].(type) {
			case *object.RawMessage:
				res.Value = int64(len(w.Words))
			case *object.Message:
				res.Value = int64(len(w.Body.Words))
			default:
				res.Value++
			}

			return res
		},
	}

	typeFn = &object.Builtin{
		Name: "type",
		Fn: func(rm *object.RawMessage) object.Object {
			if len(rm.Words) == 0 {
				return &object.Error{Message: "type builtin requires at least one argument"}
			}
			if len(rm.Words) > 1 {
				return &object.Error{Message: "type builtin do not allows more than 1 argument"}
			}

			word := rm.Words[0]
			return &object.String{Value: string(word.Type())}
		},
	}

	headFn = &object.Builtin{
		Name: "head",
		Fn: func(rm *object.RawMessage) object.Object {
			if len(rm.Words) == 0 {
				return &object.Error{Message: "head builtin requires at least one argument"}
			}
			if len(rm.Words) > 1 {
				return &object.Error{Message: "head builtin do not allows more than 1 argument"}
			}

			w, ok := rm.Words[0].(*object.RawMessage)
			if !ok {
				msg, ok := rm.Words[0].(*object.Message)
				if !ok {
					return &object.Error{Message: "head builtin requires a raw message or a message as argument"}
				}
				w = msg.Body
			}

			if len(w.Words) == 0 {
				return &object.RawMessage{Words: make([]object.Object, 0)}
			}

			return w.Words[0]
		},
	}

	tailFn = &object.Builtin{
		Name: "tail",
		Fn: func(rm *object.RawMessage) object.Object {
			if len(rm.Words) == 0 {
				return &object.Error{Message: "tail builtin requires at least one argument"}
			}
			if len(rm.Words) > 1 {
				return &object.Error{Message: "tail builtin do not allows more than 1 argument"}
			}

			w, ok := rm.Words[0].(*object.RawMessage)
			if !ok {
				msg, ok := rm.Words[0].(*object.Message)
				if !ok {
					return &object.Error{Message: "tail builtin requires a raw message or a message as argument"}
				}
				w = msg.Body
			}

			if len(w.Words) == 0 {
				return &object.RawMessage{Words: make([]object.Object, 0)}
			}

			return &object.RawMessage{Words: w.Words[1:]}
		},
	}
)

var builtinList = [...]*object.Builtin{
	print,
	println,
	scan,
	scanln,
	lenFn,
	typeFn,
	headFn,
	tailFn,
}
