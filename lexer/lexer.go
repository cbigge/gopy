package lexer

import (
	"errors"
	"fmt"
	"unicode"
)

type tokenKey int
type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF = "EOF"
	NL = "NEWLINE"
	IDENT = "IDENT"

	// Keywords
	IF = "IF"
	ELIF = "ELIF"
	ELSE = "ELSE"
	WHILE = "WHILE"
	FOR = "FOR"
	IN = "IN"
	PRINT = "PRINT"
	INT = "INT"
	STR = "STR"
	AND = "AND"
	OR = "OR"

	// Literals
	STRING = "STRING"
	NUM = "NUM"

	// Punctuation
	LEFTPAREN = "("
	RIGHTPAREN = ")"
	COLON = ":"
	EQUALS = "="
	COMMA = ","
	INDENT = "INDENT"

	// Operations
	ADD = "+"
	SUB = "-"
	MULT = "*"
	DIV = "/"
	MOD = "%"
	POW = "^"
	NOT = "!"

	// Operation Assignment
	ADDEQ = "+="
	SUBEQ = "-="
	MULTEQ = "*="
	DIVEQ = "/="
	MODEQ = "%="
	POWEQ = "^="

	// Comparison
	LESS = "<"
	LESSEQ = "<="
	GREAT = ">"
	GREATEQ = ">="
	EQ = "=="
	NOTEQ = "!="
)

const (
	_ tokenKey = iota
	TokenKeyword
	TokenPunct
	TokenIdent
	TokenInt
	TokenString
)

type Token struct {
	Name TokenType
	Val  string
	Pos  tokenPos
}

type tokenPos struct {
	row int
	col int
}

func (t Token) GetPosition() string {
	return fmt.Sprintf("line %d, column %d", t.Pos.row, t.Pos.col)
}

func (t Token) GetCol() int {
	return t.Pos.col
}

type Lexer struct {
	index int
	input string
	line int
	column int
	current rune
	currentType tokenKey
	tokens []Token
}

func StartLex(input string) []Token {
	l := &Lexer{
		input: input,
		line: 1,
		column: 0,
		current: ' ',
		tokens: []Token{},
	}
	lex(l)
	eof := Token{
		Name: EOF,
		Val:  "",
		Pos:  tokenPos{l.line+1, 0},
	}
	l.tokens = append(l.tokens, eof)
	return l.tokens
}

func lex(l *Lexer) {
	for ok := true; ok; ok = l.index < len(l.input) {
		l.current = rune(l.input[l.index])
		switch l.current {
		case '\n':
			l.lexNL()
			l.nextLine()
		case '\t':
			l.lexTab()
		case '=':
			if nextChar, err := l.peek(); nextChar == '=' && err == nil {
				l.lexPunct(EQ, "==")
				l.index++
				l.column++
			} else {
				l.lexPunct(EQUALS, "=")
			}
		case '<':
			if nextChar, err := l.peek(); nextChar == '=' && err == nil {
				l.lexPunct(LESSEQ, "<=")
				l.index++
				l.column++
			} else {
				l.lexPunct(LESS, "<")
			}
		case '>':
			if nextChar, err := l.peek(); nextChar == '=' && err == nil {
				l.lexPunct(GREATEQ, ">=")
				l.index++
				l.column++
			} else {
				l.lexPunct(GREAT, ">")
			}
		case '!':
			if nextChar, err := l.peek(); nextChar == '=' && err == nil {
				l.lexPunct(NOTEQ, "!=")
			} else {
				l.lexPunct(NOT, "!")
			}
		case ':':
			l.lexPunct(COLON, ":")
		case '(':
			l.lexPunct(LEFTPAREN, "(")
		case ')':
			l.lexPunct(RIGHTPAREN, ")")
		case '+':
			if nextChar, err := l.peek(); nextChar == '=' && err == nil {
				l.lexPunct(ADDEQ, "+=")
				l.index++
				l.column++
			} else {
				l.lexPunct(ADD, "+")
			}
		case '-':
			if nextChar, err := l.peek(); nextChar == '=' && err == nil{
				l.lexPunct(SUBEQ, "-=")
				l.index++
				l.column++
			} else {
				l.lexPunct(SUB, "-")
			}
		case '*':
			if nextChar, err := l.peek(); nextChar == '=' && err == nil {
				l.lexPunct(MULTEQ, "*=")
				l.index++
				l.column++
			} else {
				l.lexPunct(MULT, "*")
			}
		case '/':
			if nextChar, err := l.peek(); nextChar == '=' && err == nil {
				l.lexPunct(DIVEQ, "/=")
				l.index++
				l.column++
			} else {
				l.lexPunct(DIV, "/")
			}
		case '^':
			if nextChar, err := l.peek(); nextChar == '=' && err == nil {
				l.lexPunct(POWEQ, "^=")
				l.index++
				l.column++
			} else {
				l.lexPunct(POW, "^")
			}
		case '%':
			if nextChar, err := l.peek(); nextChar == '=' && err == nil {
				l.lexPunct(MODEQ, "%=")
				l.index++
				l.column++
			} else {
				l.lexPunct(MOD, "%")
			}
		case ',':
			l.lexPunct(COMMA, ",")
		case '"':
			l.lexString()
		default:
			if unicode.IsSpace(l.current) {
				l.currentType = -1
			} else if l.current == '\n' || l.current == '#' {
				l.nextLine()
			} else if unicode.IsDigit(l.current) {
				if l.currentType == TokenIdent || l.currentType == TokenKeyword {
					l.lexText(string(l.current))
				} else {
					l.lexInt()
				}
			} else if unicode.IsLetter(l.current) || l.current == '_' {
				l.lexText(string(l.current))
			} else {
				l.tokens = append(l.tokens, Token{Name: ILLEGAL, Val: "nil"})
			}
		}
		l.index++
		l.column++
	}
}

func (l *Lexer) peek() (rune, error) {
	if l.index+1 != len(l.input) {
		return rune(l.input[l.index+1]), nil
	}
	return ' ', errors.New("end of input")
}

func (l *Lexer) lexTab() {
	var tok Token
	tok.Name = INDENT
	tok.Val = INDENT
	tok.Pos = tokenPos{
		row: l.line,
		col: l.column,
	}
	l.column += 3
	l.tokens = append(l.tokens, tok)
}

func (l *Lexer) lexPunct(name TokenType, val string) {
	var tok Token
	l.currentType = TokenPunct
	tok.Name = name
	tok.Val = val
	tok.Pos = tokenPos{
		row: l.line,
		col: l.column,
	}
	l.tokens = append(l.tokens, tok)
}

func (l *Lexer) lexText(val string) {
	reserved := []string{
		"if",
		"elif",
		"else",
		"while",
		"for",
		"in",
		"print",
		"int",
		"str",
		"and",
		"or",
	}

	var tok Token
	if l.currentType == TokenKeyword {
		l.currentType = TokenIdent
		tok = l.tokens[len(l.tokens)-1]
		l.tokens = l.tokens[:len(l.tokens)-1]
		tok.Name = IDENT
		tok.Val += val
	} else if l.currentType == TokenIdent {
		tok = l.tokens[len(l.tokens)-1]
		l.tokens = l.tokens[:len(l.tokens)-1]
		tok.Val += val
	} else {
		l.currentType = TokenIdent
		tok.Name = IDENT
		tok.Val = val
		tok.Pos = tokenPos{
			row: l.line,
			col: l.column,
		}
	}
	for _, keyword := range reserved {
		if tok.Val == keyword {
			switch tok.Val {
			case "if":
				tok.Name = IF
			case "elif":
				tok.Name = ELIF
			case "else":
				tok.Name = ELSE
			case "while":
				tok.Name = WHILE
			case "for":
				tok.Name = FOR
			case "in":
				tok.Name = IN
			case "print":
				tok.Name = PRINT
			case "int":
				tok.Name = INT
			case "str":
				tok.Name = STR
			case "and":
				tok.Name = AND
			case "or":
				tok.Name = OR
			}
		}
	}
	l.tokens = append(l.tokens, tok)
}

func (l *Lexer) lexString() {
	var tok Token
	l.index++
	start := l.index
	l.currentType = TokenString
	next, err := l.peek()
	if err != nil {
		tok.Name = ILLEGAL
		tok.Val = ""
		tok.Pos = tokenPos{
			row: l.line,
			col: l.column,
		}
		l.tokens = append(l.tokens, tok)
		return
	}
	for next != '"' {
		l.index++
		if next, err = l.peek(); err != nil {
			break
		}
	}
	tok.Name = STRING
	tok.Val = l.input[start:l.index+1]
	tok.Pos = tokenPos{
		row: l.line,
		col: l.column,
	}
	l.tokens = append(l.tokens, tok)
	l.index++
}

func (l *Lexer) lexInt() {
	var tok Token
	start := l.index
	l.currentType = TokenInt
	next, err := l.peek()
	if err != nil {
		tok.Name = NUM
		tok.Val = string(l.input[start])
		tok.Pos = tokenPos{
			row: l.line,
			col: l.column,
		}
		l.tokens = append(l.tokens, tok)
		return
	}
	for unicode.IsDigit(next) {
		l.index++
		if next, err = l.peek(); err != nil {
			break
		}
	}
	tok.Name = NUM
	tok.Val = l.input[start:l.index+1]
	tok.Pos = tokenPos{
		row: l.line,
		col: l.column,
	}
	l.tokens = append(l.tokens, tok)
}

func (l *Lexer) lexNL() {
	tok := Token{
		Name: NL,
		Val:  NL,
		Pos:  tokenPos{l.line, l.column},
	}
	l.tokens = append(l.tokens, tok)
}

func (l *Lexer) nextLine() {
	l.line++
	l.column = 0
	l.currentType = -1
	for l.input[l.index] != '\n' && l.index < len(l.input)-1 {
		l.index++
	}
}