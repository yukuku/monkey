package repl

import (
	"bufio"
	"fmt"
	"io"
	"lexer"
	"parser"
	"evaluator"
	"object"
)

const PROMPT = "\n🐵> "

func Start(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Print(PROMPT)
		if !sc.Scan() {
			return
		}
		line := sc.Text()
		lx := lexer.New(line)
		p := parser.New(lx)
		prog := p.Parse()

		errors := p.Errors()
		if len(errors) > 0 {
			fmt.Printf("Found %d error(s):\n", len(errors))
			for _, msg := range errors {
				fmt.Printf("- %s\n", msg)
			}

		} else {
			out := evaluator.Eval(prog, env)
			fmt.Print(out.Inspect())
			fmt.Println()
		}
	}
}
