package ast

import (
	"bytes"
	"gopy/lexer"
	"strconv"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Stmt interface {
	Node
	statementNode()
}

type Expr interface {
	Node
	expressionNode()
}

type Program struct {
	Stmts []Stmt
}

func (p *Program) TokenLiteral() string {
	if len(p.Stmts) > 0 {
		return p.Stmts[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var result bytes.Buffer
	for _, stmt := range p.Stmts {
		result.WriteString(stmt.String())
	}
	return result.String()
}

type VarStmt struct {
	Token lexer.Token
	Ident *Identifier
	Value Expr
}

func (vs *VarStmt) statementNode() {}
func (vs *VarStmt) TokenLiteral() string { return vs.Token.Val }
func (vs *VarStmt) String() string {
	var result bytes.Buffer
	result.WriteString(vs.TokenLiteral() + " ")
	result.WriteString(vs.Ident.String())
	result.WriteString(" = ")
	if vs.Value != nil {
		result.WriteString(vs.Value.String())
	}
	return result.String()
}

type Identifier struct {
	Token lexer.Token
	Val string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Val }
func (i *Identifier) String() string { return i.Val }

type ExprStmt struct {
	Token lexer.Token
	Expr Expr
}

func (es *ExprStmt) statementNode() {}
func (es *ExprStmt) TokenLiteral() string { return es.Token.Val }
func (es *ExprStmt) String() string {
	if es.Expr != nil {
		return es.Expr.String()
	}
	return ""
}

type IntLiteral struct {
	Token lexer.Token
	Value int64
}

func (il *IntLiteral) expressionNode() {}
func (il *IntLiteral) TokenLiteral() string { return il.Token.Val }
func (il *IntLiteral) String() string { return strconv.Itoa(int(il.Value)) }

type StrLiteral struct {
	Token lexer.Token
	Value string
}

func (sl *StrLiteral) expressionNode() {}
func (sl *StrLiteral) TokenLiteral() string { return sl.Token.Val }
func (sl *StrLiteral) String() string { return sl.Value }

type PrefixExpr struct {
	Token lexer.Token
	Op string
	Expr Expr
}

func (pe *PrefixExpr) expressionNode() {}
func (pe *PrefixExpr) TokenLiteral() string { return pe.Token.Val }
func (pe *PrefixExpr) String() string {
	var result bytes.Buffer
	result.WriteString("(")
	result.WriteString(pe.Op)
	result.WriteString(pe.Expr.String())
	result.WriteString(")")
	return result.String()
}

type InfixExpr struct {
	Token lexer.Token
	Left Expr
	Op string
	Right Expr
}

func (ie *InfixExpr) expressionNode() {}
func (ie *InfixExpr) TokenLiteral() string { return ie.Token.Val }
func (ie *InfixExpr) String() string {
	var result bytes.Buffer
	result.WriteString("(")
	result.WriteString(ie.Left.String())
	result.WriteString(" " + ie.Op + " ")
	result.WriteString(ie.Right.String())
	result.WriteString(")")
	return result.String()
}

type IfExpr struct {
	Token lexer.Token
	Cond Expr
	Pass *BlockStmt
	Fail *BlockStmt
}

func (ie *IfExpr) expressionNode()      {}
func (ie *IfExpr) TokenLiteral() string { return ie.Token.Val }
func (ie *IfExpr) String() string {
	var result bytes.Buffer
	result.WriteString("if")
	result.WriteString(ie.Cond.String())
	result.WriteString(" : ")
	result.WriteString(ie.Pass.String())
	if ie.Fail != nil {
		result.WriteString(" else: ")
		result.WriteString(ie.Fail.String())
	}
	return result.String()
}

type BlockStmt struct {
	Token lexer.Token
	Stmts []Stmt
}

func (bs *BlockStmt) expressionNode()      {}
func (bs *BlockStmt) TokenLiteral() string { return bs.Token.Val }
func (bs *BlockStmt) String() string {
	var result bytes.Buffer
	for _, stmt := range bs.Stmts {
		result.WriteString(stmt.String())
	}
	return result.String()
}

type CallExpr struct {
	Token lexer.Token
	Func Expr
	Args []Expr
}

func (ce *CallExpr) expressionNode() {}
func (ce *CallExpr) TokenLiteral() string { return ce.Token.Val }
func (ce *CallExpr) String() string {
	var result bytes.Buffer
	var args []string
	for _, arg := range ce.Args {
		args = append(args, arg.String())
	}
	result.WriteString(ce.Func.String())
	result.WriteString("(")
	result.WriteString(strings.Join(args, ", "))
	result.WriteString(")")
	return result.String()
}

type WhileExpr struct {
	Token lexer.Token
	Cond Expr
	Body *BlockStmt
}

func (we *WhileExpr) expressionNode() {}
func (we *WhileExpr) TokenLiteral() string { return we.Token.Val }
func (we *WhileExpr) String() string {
	var result bytes.Buffer
	result.WriteString(we.Token.Val)
	result.WriteString("(")
	result.WriteString(we.Cond.String())
	result.WriteString(")")
	return result.String()
}