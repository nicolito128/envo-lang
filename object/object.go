package object

import (
	"fmt"
	"strconv"
	"strings"
	"unique"

	"n128.xyz/n128/envo/ast"
)

type ObjectType string

const (
	IDENT_OBJ       ObjectType = "IDENT"
	INTEGER_OBJ     ObjectType = "INTEGER"
	FLOAT_OBJ       ObjectType = "FLOAT"
	COMPLEX_OBJ     ObjectType = "COMPLEX"
	STRING_OBJ      ObjectType = "STRING"
	RAW_STRING_OBJ  ObjectType = "RAW_STRING"
	CHAR_OBJ        ObjectType = "CHAR"
	SYMBOL_OBJ      ObjectType = "SYMBOL"
	BOOL_OBJ        ObjectType = "BOOL"
	ERROR_OBJ       ObjectType = "ERROR"
	RAW_MESSAGE_OBJ ObjectType = "RAW_MESSAGE"
	MESSAGE_OBJ     ObjectType = "MESSAGE"
	DEFINE_OBJ      ObjectType = "DEFINE"
	DEFINE_LIST_OBJ ObjectType = "DEFINE_LIST"
	BUILTIN_OBJ     ObjectType = "BUILTIN"
	NIL_OBJ         ObjectType = "NIL"
)

type Object interface {
	Type() ObjectType
	Literal() any
	String() string
}

// Builtin
type Builtin struct {
	Name string
	Fn   func(*RawMessage) Object
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Literal() any     { return nil }
func (b *Builtin) String() string   { return "builtin" }

// Identifier
type Identifier struct{ Name string }

func (i *Identifier) Type() ObjectType { return INTEGER_OBJ }
func (i *Identifier) Literal() any     { return i.Name }
func (i *Identifier) String() string   { return fmt.Sprintf("%s", i.Name) }

// Integer
type Integer struct{ Value int64 }

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Literal() any     { return i.Value }
func (i *Integer) String() string   { return strconv.FormatInt(i.Value, 10) }

// Float
type Float struct{ Value float64 }

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Literal() any     { return f.Value }
func (f *Float) String() string   { return strconv.FormatFloat(f.Value, 'g', -1, 64) }

// Complex
type Complex struct{ Value complex128 }

func (c *Complex) Type() ObjectType { return COMPLEX_OBJ }
func (c *Complex) Literal() any     { return c.Value }
func (c *Complex) String() string   { return fmt.Sprintf("%v", c.Value) }

// String
type String struct{ Value string }

func (s *String) Type() ObjectType { return SYMBOL_OBJ }
func (s *String) Literal() any     { return s.Value }

func (s *String) String() string {
	q, err := strconv.Unquote(`"` + s.Value + `"`)
	if err != nil {
		return fmt.Sprintf("\"%s\"", s.Value)
	}
	return fmt.Sprintf("\"%s\"", q)
}

// RawString
type RawString struct{ Value string }

func (s *RawString) Type() ObjectType { return SYMBOL_OBJ }
func (s *RawString) Literal() any     { return s.Value }
func (s *RawString) String() string   { return fmt.Sprintf("`%s`", s.Value) }

// Symbol
type Symbol struct{ H unique.Handle[string] }

func (s *Symbol) Type() ObjectType { return SYMBOL_OBJ }
func (s *Symbol) Literal() any     { return s.H.Value() }
func (s *Symbol) String() string   { return s.H.Value() }

// Symbol
type Char struct{ Value []rune }

func (c *Char) Type() ObjectType { return CHAR_OBJ }
func (c *Char) Literal() any     { return c.Value }
func (c *Char) String() string   { return fmt.Sprintf("'%s'", string(c.Value)) }

// Bool
type Bool struct{ Value bool }

func (b *Bool) Type() ObjectType { return BOOL_OBJ }
func (b *Bool) Literal() any     { return b.Value }
func (b *Bool) String() string   { return fmt.Sprintf("%t", b.Value) }

// Error
type Error struct{ Message string }

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Literal() any     { return e.Message }
func (e *Error) String() string   { return "ERROR: " + e.Message }

// Nil
type Nil struct{}

func (n *Nil) Type() ObjectType { return NIL_OBJ }
func (n *Nil) Literal() any     { return nil }
func (n *Nil) String() string   { return "nil" }

// Define : <pattern> <response>
type Define struct {
	Receiver string
	Pattern  Object
	Response *ast.RawMessage
}

func (d *Define) Type() ObjectType { return DEFINE_OBJ }
func (d *Define) Literal() any     { return d.Receiver }

func (d *Define) String() string {
	return fmt.Sprintf("%s%s%s", d.Receiver, d.Pattern, d.Response)
}

// DefineList
type DefineList struct {
	Elems []Object
}

func (d *DefineList) Type() ObjectType { return DEFINE_LIST_OBJ }
func (d *DefineList) Literal() any     { return d.Elems }

func (d *DefineList) String() string {
	var elems []string
	for _, elem := range d.Elems {
		elems = append(elems, elem.String())
	}
	return "{" + strings.Join(elems, ",") + "}"
}

// RawMessage
type RawMessage struct {
	Words []Object
}

func (r *RawMessage) Type() ObjectType { return RAW_MESSAGE_OBJ }
func (r *RawMessage) Literal() any     { return r.Words }

func (r *RawMessage) String() string {
	var elems []string
	for _, elem := range r.Words {
		elems = append(elems, elem.String())
	}
	return "{" + strings.Join(elems, ", ") + "}"
}

// Message
type Message struct {
	Label string
	Body  *RawMessage
}

func (m *Message) Type() ObjectType { return MESSAGE_OBJ }
func (m *Message) Literal() any     { return m.Body.Words }

func (m *Message) String() string {
	var elems []string
	for _, elem := range m.Body.Words {
		elems = append(elems, elem.String())
	}
	return m.Label + "{" + strings.Join(elems, ", ") + "}"
}
