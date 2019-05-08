/*
	Package grep exposes grep-like functionality. A best effort is made to
	mirror GNU grep 3.3 (https://www.gnu.org/software/grep/manual/grep.html).
*/
package grep

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"regexp"
	"regexp/syntax"
	"strings"
)

type Option func(*options)

// WithRegexps uses one or more patterns; newlines within patterns
// separate each pattern from the next. If this option is used multiple times
// or is combined with the WithFiles option, search for all patterns given.
func WithRegexps(patterns ...string) Option {
	return func(opts *options) {
		opts.e = append(opts.e, patterns...)
	}
}

// WithFiles obtain patterns from files, one per line. If this option is
// combined with the WithRegexps option, search for all patterns given.
// The empty file contains zero patterns, and therefore matches nothing.
func WithFiles(files ...*os.File) Option {
	return func(opts *options) {
		opts.f = append(opts.f, files...)
	}
}

// WithIgnoreCase ignores case distinctions, so that characters that differ
// only in case match each other. Setting this option is identical to specifying
// a case-insensitive flag in pattern.
func WithIgnoreCase() Option {
	return func(opts *options) {
		opts.i = true
	}
}

// WithInvertMatch inverts the sense of matching, to select non-matching lines.
func WithInvertMatch() Option {
	return func(opts *options) {
		opts.v = true
	}
}

// WithWordRegexp selects only those lines containing matches that form whole
// words. The test is that the matching substring must either be at the beginning
// of the line, or preceded by a non-word constituent character. Similarly, it
// must be either at the end of the line or followed by a non-word constituent
// character. Word constituent characters are letters, digits, and the underscore.
// This option has no effect if WithLineRegexp is also specified.
func WithWordRegexp() Option {
	return func(opts *options) {
		if opts.x {
			return
		}
		opts.w = true
	}
}

// WithLineRegexp selects only those matches that exactly match the whole line.
// For regular expression patterns, this is like parenthesizing each pattern and
// then surrounding it with ‘^’ and ‘$’.
func WithLineRegexp() Option {
	return func(opts *options) {
		opts.x = true
		opts.w = false
	}
}

type options struct {
	// Matching Control
	// https://www.gnu.org/software/grep/manual/grep.html#Matching-Control

	//   -e, --regexp=PATTERN      use PATTERN for matching
	e []string
	//   -f, --file=FILE           obtain PATTERN from FILE
	f []*os.File
	//   -i, --ignore-case         ignore case distinctions
	i bool
	//   -v, --invert-match        select non-matching lines
	v bool
	//   -w, --word-regexp         force PATTERN to match only whole words
	w bool
	//   -x, --line-regexp         force PATTERN to match only whole lines
	x bool

	// General Output control
	// https://www.gnu.org/software/grep/manual/grep.html#General-Output-Control

	// Other Options
	//   -z, --null-data           a data line ends in 0 byte, not newline
	z bool

	//   -m, --max-count=NUM       stop after NUM selected lines
	m int
	//   -b, --byte-offset         print the byte offset with output lines
	b bool
	//   -n, --line-number         print line number with output lines
	n bool
	//       --line-buffered       flush output on every line
	lineBuffered bool
	//   -H, --with-filename       print file name with output lines
	H bool
	//   -h, --no-filename         suppress the file name prefix on output
	h bool
	//       --label=LABEL         use LABEL as the standard input file name prefix
	label string
	//   -o, --only-matching       show only the part of a line matching PATTERN
	o bool
	//   -q, --quiet, --silent     suppress all normal output
	//       --binary-files=TYPE   assume that binary files are TYPE;
	//                             TYPE is 'binary', 'text', or 'without-match'
	//   -a, --text                equivalent to --binary-files=text
	//   -I                        equivalent to --binary-files=without-match
	//   -d, --directories=ACTION  how to handle directories;
	//                             ACTION is 'read', 'recurse', or 'skip'
	//   -D, --devices=ACTION      how to handle devices, FIFOs and sockets;
	//                             ACTION is 'read' or 'skip'
	//   -r, --recursive           like --directories=recurse
	//   -R, --dereference-recursive  likewise, but follow all symlinks
	//       --include=FILE_PATTERN  search only files that match FILE_PATTERN
	//       --exclude=FILE_PATTERN  skip files and directories matching FILE_PATTERN
	//       --exclude-from=FILE   skip files matching any file pattern from FILE
	//       --exclude-dir=PATTERN  directories that match PATTERN will be skipped.
	//   -L, --files-without-match  print only names of FILEs with no selected lines
	//   -l, --files-with-matches  print only names of FILEs with selected lines
	//   -c, --count               print only a count of selected lines per FILE
	//   -T, --initial-tab         make tabs line up (if needed)
	//   -Z, --null                print 0 byte after FILE name

	// Context control:
	//   -B, --before-context=NUM  print NUM lines of leading context
	//   -A, --after-context=NUM   print NUM lines of trailing context
	//   -C, --context=NUM         print NUM lines of output context
	//   -NUM                      same as --context=NUM
	//       --color[=WHEN],
	//       --colour[=WHEN]       use markers to highlight the matching strings;
	//                             WHEN is 'always', 'never', or 'auto'
	//   -U, --binary              do not strip CR characters at EOL (MSDOS/Windows)

	// Programs:
	// https://www.gnu.org/software/grep/manual/grep.html#grep-Programs
}

// Grep searches input files for matches to patterns. When it finds a match in
// a line, it copies the line to the output.
//
// Though Grep expects to do the matching on text, it has no limits on input
// line length other than available memory, and it can match arbitrary
// characters within a line. If the final byte of an input file is not a
// newline, grep silently supplies one. Since newline is also a separator for
// the list of patterns, there is no way to match newline characters in a text.
type Grep struct {
	pattern string
	opts    *options
}

// New returns a Grep that matches pattern with opts set. The pattern argument
// contains one or more patterns separated by newlines. Each resulting pattern is
// interpreted according to the regexp package.
func New(pattern string, opts ...Option) *Grep {
	options := &options{}
	for _, opt := range opts {
		opt(options)
	}
	return &Grep{
		pattern: pattern,
		opts:    options,
	}
}

func (cmd *Grep) Exec(input io.Reader) io.Reader {
	r, w := io.Pipe()

	matchers, err := cmd.matchers()
	if err != nil {
		w.CloseWithError(err)
		return r
	}

	go func() {

		s := bufio.NewScanner(input)
		for s.Scan() {
			for _, m := range matchers {
				line := s.Bytes()
				if m.Match(line) {
					_, err := w.Write(append(line, '\n'))
					if err != nil {
						w.CloseWithError(err)
						return
					}
					break
				}
			}
		}
		w.CloseWithError(s.Err())

	}()

	return r
}

type matcher struct {
	regexp *regexp.Regexp
	opts   *options
}

func (m *matcher) Match(line []byte) bool {
	matches := func() bool {
		if !m.regexp.Match(line) {
			return false
		}

		// match lines only
		if m.opts.x {
			match := m.regexp.Find(line)
			if m.opts.i {
				return bytes.EqualFold(match, line)
			}
			return bytes.Equal(match, line)
		}

		// match whole words only
		if m.opts.w {
			indexes := m.regexp.FindAllIndex(line, -1)
			for _, i := range indexes {
				begin, end := i[0], i[1]
				switch {
				case begin == 0 && end == len(line):
					return true
				case begin == 0 && !syntax.IsWordChar(rune(line[end])):
					return true
				case end == len(line) && !syntax.IsWordChar(rune(line[begin-1])):
					return true
				}
			}
			return false
		}

		return true
	}()

	// invert match if necessary
	return matches != m.opts.v // xor
}

func (cmd *Grep) matchers() ([]*matcher, error) {
	var matchers []*matcher

	appendRegex := func(expr string) error {
		xflags := syntax.Perl // -p, --perl-regexp
		if cmd.opts.i {
			xflags |= syntax.FoldCase // -i, --ignore-case
		}
		parsed, err := syntax.Parse(expr, syntax.Perl|syntax.FoldCase)
		if err != nil {
			return err
		}
		regex, err := regexp.Compile(parsed.String())
		if err != nil {
			return err
		}
		matchers = append(matchers, &matcher{regexp: regex, opts: cmd.opts})
		return nil
	}

	// obtain patterns from input, split on newlines. But only if regexps and files are unset.
	if len(cmd.opts.e) == 0 && len(cmd.opts.f) == 0 {
		for _, expr := range strings.Split(cmd.pattern, "\n") {
			if err := appendRegex(expr); err != nil {
				return nil, err
			}
		}
	}

	// obtain patterns from regexp option, split on newlines
	for _, pattern := range cmd.opts.e {
		for _, expr := range strings.Split(pattern, "\n") {
			if err := appendRegex(expr); err != nil {
				return nil, err
			}
		}
	}

	// obtain patterns from files, one per line
	for _, file := range cmd.opts.f {
		s := bufio.NewScanner(file)
		for s.Scan() {
			if err := appendRegex(s.Text()); err != nil {
				return nil, err
			}
		}
		if err := s.Err(); err != nil {
			return nil, err
		}
	}

	return matchers, nil
}
