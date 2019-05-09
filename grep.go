package usrbin

import (
	"io"

	"github.com/kevin-cantwell/usrbin/pkg/grep"
)

func Grep(pattern string, input io.Reader, opts ...grep.Option) io.Reader {
	return grep.New(pattern, opts...).Exec(input)
}
