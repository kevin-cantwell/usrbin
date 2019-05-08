package main

import (
	"io"
	"os"

	"github.com/kevin-cantwell/usrbin"
)

func main() {
	output := usrbin.Grep(os.Stdin, os.Args[1])
	io.Copy(os.Stdout, output)
}
