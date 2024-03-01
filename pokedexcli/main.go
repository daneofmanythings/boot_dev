package main

import (
	"os"

	"gitlab.com/daneofmanythings/pokedex/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
