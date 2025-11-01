package main

import (
	"github.com/andrewneudegg/calc/pkg/display"
)

func main() {
	repl := display.NewREPL()
	repl.Run()
}
