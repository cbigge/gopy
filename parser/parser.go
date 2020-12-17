package parser

import (
	"fmt"
	"gopy/ast"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gopy/lexer"
)

const (
	_ int = iota
	LOWEST
	EQUALS
	GTLT
	SUM
	PRODUCT
	PREFIX
	CALL
	AND
)
var precedence = map[lexer.TokenType]int{
	lexer.EQ: EQUALS,
	lexer.NOTEQ: EQUALS,
	lexer.LESS: GTLT,
	lexer.LESSEQ: GTLT,
	lexer.GREAT: GTLT,
	lexer.GREATEQ: GTLT,
	lexer.AND: GTLT,
	lexer.OR: GTLT,
	lexer.ADD: SUM,
	lexer.ADDEQ: SUM,
	lexer.SUB: SUM,
	lexer.SUBEQ: SUM,
	lexer.DIV: PRODUCT,
	lexer.DIVEQ: PRODUCT,
	lexer.MULT: PRODUCT,
	lexer.MULTEQ: PRODUCT,
	lexer.LEFTPAREN: CALL,
}

type Parser struct {
	tokens         []lexer.Token
	index          int
	statements     []ast.Stmt
	errors         []string
	prefixParseFns map[lexer.TokenType]prefixParseFn
	infixParseFns  map[lexer.TokenType]infixParseFn
	indentLevel    int
}

type (
	prefixParseFn func() ast.Expr
	infixParseFn func(expr ast.Expr) ast.Expr
)

func StartParse(path string) []ast.Stmt {
	file, _ := os.Open(path)
	defer file.Close()
	fileContents, _ := ioutil.ReadFile(path)
	formattedContents := strings.ReplaceAll(string(fileContents), "    ", "\t")
	tokens := lexer.StartLex(formattedContents)
	if len(tokens) == 0 {
		panic("no tokens to parse")
	}

	p := Parser{
		tokens:     tokens,
		index:      0,
		statements: []ast.Stmt{},
		errors:     []string{},
		indentLevel: 0,
	}
	p.registerFixes()
	stmts := parse(&p)
	return stmts
}

func StartParseRepl(input string) (*Parser, ast.Program) {
	formatted := strings.ReplaceAll(input, "    ", "\t")
	tokens := lexer.StartLex(formatted)
	if len(tokens) == 0 {
		panic("no tokens to parse")
	}
	p := Parser{
		tokens:     tokens,
		index:      0,
		statements: []ast.Stmt{},
		errors:     []string{},
	}
	p.registerFixes()
	parse(&p)
	return &p, ast.Program{Stmts: p.statements}
}

func (p *Parser) registerFixes() {
	p.prefixParseFns = make(map[lexer.TokenType]prefixParseFn)
	p.registerPrefix(lexer.IDENT, p.parseIdent)
	p.registerPrefix(lexer.NUM, p.parseIntLiteral)
	p.registerPrefix(lexer.STRING, p.parseStrLiteral)
	p.registerPrefix(lexer.SUB, p.parsePrefixExpr)
	p.registerPrefix(lexer.LEFTPAREN, p.parseGroupingExpr)
	p.registerPrefix(lexer.IF, p.parseIfExpr)
	p.registerPrefix(lexer.WHILE, p.parseWhileExpr)
	p.registerPrefix(lexer.STR, p.parseStrLiteral)
	p.registerPrefix(lexer.INT, p.parseIntLiteral)

	p.infixParseFns = make(map[lexer.TokenType]infixParseFn)
	p.registerInfix(lexer.EQ, p.parseInfixExpr)
	p.registerInfix(lexer.NOTEQ, p.parseInfixExpr)
	p.registerInfix(lexer.ADD, p.parseInfixExpr)
	p.registerInfix(lexer.ADDEQ, p.parseInfixExpr)
	p.registerInfix(lexer.SUB, p.parseInfixExpr)
	p.registerInfix(lexer.SUBEQ, p.parseInfixExpr)
	p.registerInfix(lexer.DIV, p.parseInfixExpr)
	p.registerInfix(lexer.DIVEQ, p.parseInfixExpr)
	p.registerInfix(lexer.MULT, p.parseInfixExpr)
	p.registerInfix(lexer.MULTEQ, p.parseInfixExpr)
	p.registerInfix(lexer.GREAT, p.parseInfixExpr)
	p.registerInfix(lexer.GREATEQ, p.parseInfixExpr)
	p.registerInfix(lexer.LESSEQ, p.parseInfixExpr)
	p.registerInfix(lexer.AND, p.parseInfixExpr)
	p.registerInfix(lexer.LEFTPAREN, p.parseCallExpr)
}

func parse(p *Parser) []ast.Stmt {
	for !p.end() {
		//statement, err := p.declaration()
		stmt := p.parseStmt()
		if stmt != nil {
			p.statements = append(p.statements, stmt)
		}
		p.next()
	}
	return p.statements
}

func (p *Parser) parseStmt() ast.Stmt {
	switch p.current().Name {
	case lexer.IDENT:
		if p.peek().Name == lexer.EQUALS {
			return p.parseVarStmt()
		} else {
			return p.parseExprStmt()
		}
	case lexer.NL:
		return nil
	default:
		return p.parseExprStmt()
	}
}

func (p *Parser) parseVarStmt() *ast.VarStmt {
	stmt := &ast.VarStmt{Token: p.current()}
	stmt.Ident = &ast.Identifier{Token: p.current(), Val: p.current().Val}
	if !p.expectPeek(lexer.EQUALS) {
		return nil
	}
	p.next()
	stmt.Value = p.parseExpr(LOWEST)
	for !p.checkCurrent(lexer.NL) && !p.end() {
		p.next()
	}
	return stmt
}

func (p *Parser) parseExprStmt() *ast.ExprStmt {
	stmt := &ast.ExprStmt{Token: p.current()}
	stmt.Expr = p.parseExpr(LOWEST)
	return stmt
}

func (p *Parser) parseExpr(precedence int) ast.Expr {
	pre := p.prefixParseFns[p.current().Name]
	if pre == nil {
		return nil
	}
	left := pre()
	for !p.checkCurrent(lexer.NL) && precedence < p.peekPrec() {
		in := p.infixParseFns[p.peek().Name]
		if in == nil {
			return left
		}
		p.next()
		left = in(left)
	}
	return left
}

func (p *Parser) parseIdent() ast.Expr {
	return &ast.Identifier{Token: p.current(), Val: p.current().Val}
}

func (p *Parser) parseIntLiteral() ast.Expr {
	il := &ast.IntLiteral{Token: p.current()}
	val, err := strconv.ParseInt(p.current().Val, 0, 64)
	if err != nil {
		err := fmt.Sprintf("could not parse %q as int", p.current().Val)
		p.errors = append(p.errors, err)
		return nil
	}
	il.Value = val
	return il
}

func (p *Parser) parseStrLiteral() ast.Expr {
	return &ast.StrLiteral{Token: p.current(), Value: p.current().Val}
}

func (p *Parser) parsePrefixExpr() ast.Expr {
	expr := &ast.PrefixExpr{Token: p.current(), Op: p.current().Val}
	p.next()
	expr.Expr = p.parseExpr(PREFIX)
	return expr
}

func (p *Parser) parseInfixExpr(l ast.Expr) ast.Expr {
	expr := &ast.InfixExpr{Token: p.current(), Op: p.current().Val, Left: l}
	prec := p.currentPrec()
	p.next()
	expr.Right = p.parseExpr(prec)
	return expr
}

func (p *Parser) parseGroupingExpr() ast.Expr {
	p.next()
	expr := p.parseExpr(LOWEST)
	if !p.expectPeek(lexer.RIGHTPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseIfExpr() ast.Expr {
	expr := &ast.IfExpr{Token: p.current()}
	p.next()
	expr.Cond = p.parseExpr(LOWEST)
	if !p.expectPeek(lexer.COLON) {
		return nil
	}
	if !p.expectPeek(lexer.NL) {
		return nil
	}
	expr.Pass = p.parseBlockStmt()
	if !p.expectCurrent(lexer.ELSE) {
		return expr
	} else {
		expr.Fail = p.parseBlockStmt()
	}
	return expr
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	p.next()
	for p.checkPeek(lexer.INDENT) {
		p.next()
	}
	p.next()
	p.indentLevel++
	b := &ast.BlockStmt{Token: p.current()}
	b.Stmts = []ast.Stmt{}
	for 4*p.indentLevel <= p.current().GetCol() && !p.end() {
		stmt := p.parseStmt()
		if stmt != nil {
			b.Stmts = append(b.Stmts, stmt)
		}
		p.next()
		for p.checkCurrent(lexer.NL) || p.checkCurrent(lexer.INDENT) {
			p.next()
		}
	}
	p.indentLevel--
	return b
}

func (p *Parser) parseCallExpr(fn ast.Expr) ast.Expr {
	expr := &ast.CallExpr{Token: p.current(), Func: fn}
	expr.Args = p.parseCallArgs()
	return expr
}

func (p *Parser) parseCallArgs() []ast.Expr {
	var args []ast.Expr
	if p.checkPeek(lexer.RIGHTPAREN) {
		p.next()
		return args
	}
	p.next()
	args = append(args, p.parseExpr(LOWEST))
	for p.checkPeek(lexer.COMMA) {
		p.next()
		p.next()
		args = append(args, p.parseExpr(LOWEST))
	}
	if !p.expectPeek(lexer.RIGHTPAREN) {
		return nil
	}
	return args
}

func (p *Parser) parseWhileExpr() ast.Expr {
	expr := &ast.WhileExpr{Token: p.current()}
	p.next()
	expr.Cond = p.parseExpr(LOWEST)
	if !p.expectPeek(lexer.COLON) {
		return nil
	}
	if !p.expectPeek(lexer.NL) {
		return nil
	}
	expr.Body = p.parseBlockStmt()
	return expr
}

func (p *Parser) expectCurrent(t lexer.TokenType) bool {
	if p.checkCurrent(t) {
		p.next()
		return true
	} else {
		return false
	}
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.checkPeek(t) {
		p.next()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) checkCurrent(t lexer.TokenType) bool {
	if p.end() { return false }
	return p.current().Name == t
}

func (p *Parser) checkPeek(t lexer.TokenType) bool {
	if p.end() { return false }
	return p.peek().Name == t
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t lexer.TokenType) {
	err := fmt.Sprintf("error at %s: expected next token to be %s, got %s instead",
		p.peek().GetPosition(), t, p.peek().Name)
	p.errors = append(p.errors, err)
}

func (p *Parser) next() {
	if !p.end() {
		p.index++
	}
}

func (p *Parser) end() bool {
	token := p.current()
	return token.Name == lexer.EOF
}

func (p *Parser) current() lexer.Token {
	return p.tokens[p.index]
}

func (p *Parser) peek() lexer.Token {
	if !p.end() {
		return p.tokens[p.index+1]
	}
	return lexer.Token{}
}

func (p *Parser) registerPrefix(t lexer.TokenType, fn prefixParseFn) {
	p.prefixParseFns[t] = fn
}

func (p *Parser) registerInfix(t lexer.TokenType, fn infixParseFn) {
	p.infixParseFns[t] = fn
}

func (p *Parser) currentPrec() int {
	if prec, ok := precedence[p.current().Name]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) peekPrec() int {
	if prec, ok := precedence[p.peek().Name]; ok {
		return prec
	}
	return LOWEST
}