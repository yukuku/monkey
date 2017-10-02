package main

import (
	"os/user"
	"fmt"
	"repl"
	"os"
)

func main() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello there %q! Welcome to 🐵.\nPlease start typing commands.\n", u.Username)
	repl.Start(os.Stdin, os.Stdout)
}
