package main

import (
	"io"
	"os"

	"github.com/kevin-cantwell/usrbin/pkg/grep"
	getopt "github.com/pborman/getopt"
)

var (
	regexp      string
	invertMatch bool
)

func init() {
	// getopt
	// getopt.FlagLong(&regexp, "regexp", 'e', "use PATTERN for matching")
	// getopt.FlagLong(&invertMatch, "invert-match", 'v', "select non-matching lines")
}

func main() {
	getopt.Parse()

	var pattern string

	var opts []grep.Opt
	if regexp != "" {
		opts = append(opts, grep.WithRegexps(regexp))
	} else {
		pattern = os.Args[1]
	}
	if invertMatch {
		opts = append(opts, grep.WithInvertMatch())
	}

	output := grep.New(pattern, opts...).Exec(os.Stdin)
	io.Copy(os.Stdout, output)
}
