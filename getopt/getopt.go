/*
	Package getopt implements getopt(1) according to the manpage:

		https://www.linux.org/docs/man1/getopt.html

	The logic, output, and return codes attempt to mirror getopt(1) behavior
	exactly, including known bugs.
*/
package getopt

import (
	"strings"
)

// Error provides information about a getopt failure.
type Error struct {
	// Msgs provides user messages prefixed with either '<progname>: ' or
	// 'getopts: ' if the name option was left unset or the return code is
	// anything but 1.
	Msgs []string
	// ReturnCode indicates the error code of getopt. Possible values
	// are:
	//     1 if parameter parsing returns errors
	//     2 if it does not understand its own parameters
	//     3 if an internal error occurs like out-of-memory
	ReturnCode int
}

func (err *Error) Error() string {
	return strings.Join(err.Msgs, "; ")
}

var (
	UnparsableCode    int = 1
	UnknownParamsCode int = 2
	InternalErrorCode int = 3
)

type Opt func(*Opts)

func WithShortOpts(shortopts string) Opt {
	return func(opts *Opts) {
		if len(shortopts) == 0 {
			return
		}

		// The first character of shortopts may be '+' or '-' to influence the
		// way options are parsed and output is generated.
		switch shortopts[0] {
		case ':':
			// this appears to be a getopts (builtin) feature, but is still available in GNU getopt???
			WithSilentErrors()(opts)
			shortopts = shortopts[1:]
		case '+':
			WithScanPosixlyCorrect()(opts)
			shortopts = shortopts[1:]
		case '-':
			WithScanInPlace()(opts)
			shortopts = shortopts[1:]
		}

		if len(shortopts) == 0 {
			return
		}

		opts.shortopts = shortopts
	}
}

func WithLongOpts(longopts ...string) Opt {
	return func(opts *Opts) {
		if len(longopts) == 0 {
			return
		}

	}
}

func WithSilentErrors() Opt {
	return func(opts *Opts) {
		opgs.silentErrors = true
	}
}

func WithAlternative() Opt {
	return func(opts *Opts) {
		opts.alternative = true
	}
}

func WithName(name string) Opt {
	return func(opts *Opts) {
		opts.name = name
	}
}

func WithScanPosixlyCorrect() Opt {
	return func(opts *Opts) {
		opts.scanMode = '+'
	}
}

func WithScanInPlace() Opt {
	return func(opts *Opts) {
		opts.scanMode = '-'
	}
}

type opt struct {
	name     string
	argument bool
	optional bool
}

type Opts struct {
	shortopts    string
	longopts     []string
	name         string
	alternative  bool
	silentErrors bool

	//  0 : default
	// '+': POSIXLY_CORRECT
	// '-': in-place
	scanMode rune
}

type Output struct {
	Options []Option
	Args    []string
	Err     *Error
}

type Option struct {
	Name  string
	Value string
}

type Getopt struct {
	opts *Opts
}

func New(options ...Opt) *Getopt {
	opts := &Opts{
		name:        "getopt",
		alternative: false,
		scanMode:    0,
	}
	for _, opt := range options {
		opt(opts)
	}
	return &Getopt{
		opts: opts,
	}
}

var (
	unignoredShtOptChars = runeSet("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ !\"'#$%&()*+-./<=>@[\\]^_`{|}~,:")
	// same as validShortOptChars plus '?', minus ','
	unignoredLngOptChars = runeSet("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ !\"'#$%&()*+-./<=>@[\\]^_`{|}~?:")
)

func runeSet(s string) map[rune]bool {
	m := map[rune]bool{}
	for _, r := range s {
		m[r] = true
	}
	return m
}

func (cmd *Getopt) Parse(parameters ...string) (*Output, error) {
	var getoptErrs []string

	var shortopts string

	var s []opt
	var prev rune
	for i, curr := range shortopts {
		switch curr {
		case ':':

			continue
		}
		if shortopts[i] != ':' {
			s[len(s)-1].argument = true
		}

		prev = curr
	}

	// var output Output

	// for i := 0; i < len(parameters); i++ {
	// 	p := parameters[i]
	// 	switch {
	// 	case p == "--":
	// 		output.Args = append(output.Args, parameters[i:]...)
	// 		return &output, nil
	// 	case len(p) < 2:
	// 		output.Args = append(output.Args, p)
	// 	case p[:2] == "--":
	// 		name := p[2:]
	// 		opt, ok := cmd.opts.longs[name]
	// 		if !ok {
	// 			return nil, errors.New("getopt: unrecognized option '" + p + "'")
	// 		}
	// 		if opt.argument {
	// 			if !opt.optional {
	// 				dr
	// 			}
	// 		}
	// 	case p[0] == '-':
	// 	default:
	// 		output.Args = append(output.Args, p)
	// 	}
	// }

	return &output, nil
}
