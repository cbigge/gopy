package parser

import (
	"testing"
)

func TestStartParse(t *testing.T) {
	d := StartParse("test2.py")

	for _, stmt := range d {
		t.Logf("%s: %s\n", stmt.TokenLiteral(), stmt.String())
	}
	return
}
