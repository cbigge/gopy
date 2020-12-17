package interpreter

import "fmt"

type Item interface {
	Type() ItemType
	Visit() string
}
type ItemType string

const (
	ERR = "ERR"
	INT = "INT"
	STR = "STR"
	BOOL = "BOOL"
	BUILTIN = "BUILTIN"
)

type Error struct {
	Err string
}

func (e *Error) Type() ItemType { return ERR }
func (e *Error) Visit() string { return e.Err }

type Int struct {
	Val int64
}

func (i *Int) Type() ItemType { return INT }
func (i *Int) Visit() string { return fmt.Sprintf("%d", i.Val) }

type Str struct {
	Val string
}

func (s *Str) Type() ItemType { return STR }
func (s *Str) Visit() string { return s.Val }

type Bool struct {
	Val bool
}

func (b *Bool) Type() ItemType { return BOOL }
func (b *Bool) Visit() string { return fmt.Sprintf("%t", b.Val) }

type BuiltinFunction func(args ...Item) Item
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ItemType { return BUILTIN }
func (b *Builtin) Visit() string { return "builtin function" }