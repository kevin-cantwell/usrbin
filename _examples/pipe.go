package main

import (
	"io"
	"os"

	"github.com/kevin-cantwell/usrbin"
	"github.com/kevin-cantwell/usrbin/grep"
)

func main() {
	pipe := usrbin.Pipe(os.Stdin, grep.New("foo"), grep.New("bar"))
	io.Copy(os.Stdout, pipe)
}
