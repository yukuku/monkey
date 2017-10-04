package repl

import (
	"bufio"
	"fmt"
	"io"
	"lexer"
	"parser"
)

const PROMPT = "\n🐵> "

func Start(in io.Reader, out io.Writer) {
	sc := bufio.NewScanner(in)

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
			for _, s := range prog.Statements {
				fmt.Print(s.String())
			}
		}
	}
}
