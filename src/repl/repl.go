package repl

import (
	"io"
	"bufio"
	"fmt"
	"lexer"
	"token"
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

		for tok := lx.NextToken(); tok.Type != token.EOF; tok = lx.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
