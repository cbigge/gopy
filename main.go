package main

import (
	"fmt"
	"gopy/evaluator"
	"gopy/interpreter"
	"gopy/parser"
)

func main() {
	stmts := parser.StartParse("parser/test.py")
	env := interpreter.NewEnv()
	for _, stmt := range stmts {
		item := evaluator.Evaluate(stmt, env)
		fmt.Println(item.Visit())
	}

	//w := bufio.NewWriter(os.Stdout)
	//r := bufio.NewReader(os.Stdin)
	//repl.Run(w, r)
}