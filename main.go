package main

import (
	"fmt"
	"io"
	"os"
)

var inputInteractive io.Reader = os.Stdin
var outStream io.Writer = os.Stdout

func main() {
	args := os.Args
	argsLen := len(args)

	fmt.Println("args length:", argsLen)
}
