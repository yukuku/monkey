package main

import (
	"fmt"
	"os"
	"os/user"
	"repl"
)

func main() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello there %q! Welcome to 🐵.\nPlease start typing commands.\n", u.Username)
	repl.Start(os.Stdin, os.Stdout)
}
