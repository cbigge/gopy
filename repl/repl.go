package repl

import (
	"bufio"
	"fmt"
	"gopy/evaluator"
	"gopy/interpreter"
	"io"
	"gopy/parser"
)

func Run(w *bufio.Writer, r *bufio.Reader) {
	scanner := bufio.NewScanner(r)
	environment := interpreter.NewEnv()
	for {
		fmt.Printf("REPL> ")
		input := scanner.Scan()
		if !input {
			return
		}
		line := scanner.Text()
		p, program := parser.StartParseRepl(line)
		if len(p.Errors()) != 0 {
			printParserErrors(w, p.Errors())
			continue
		}
		io.WriteString(w, program.String())
		eval := evaluator.Evaluate(&program, environment)
		if eval != nil {
			fmt.Printf("%v\n", eval.Visit())
		}
	}
}

func printParserErrors(w io.Writer, errors []string) {
	for _, err := range errors {
		io.WriteString(w, "\t"+err+"\n")
		fmt.Println(err)
	}
}