package usrbin

import (
	"io"

	"github.com/kevin-cantwell/usrbin/pkg/grep"
)

func Grep(input io.Reader, pattern string, opts ...grep.Opt) io.Reader {
	return grep.New(pattern, opts...).Exec(input)
}
