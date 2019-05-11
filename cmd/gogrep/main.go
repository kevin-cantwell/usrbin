package main

import (
	"fmt"
	"io"
	"os"

	"github.com/kevin-cantwell/usrbin/pkg/grep"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	flags = []flag{
		{"regexp", "e", []string{}, "use PATTERN for matching"},
		{"file", "f", []string{}, "obtain PATTERN from FILE"},
		{"ignore-case", "i", false, "ignore case distinctions"},
		{"invert-match", "v", false, "select non-matching lines"},
		{"word-regexp", "w", false, "force PATTERN to match only whole words"},
		{"line-regexp", "x", false, "force PATTERN to match only whole lines"},
	}
)

type flag struct {
	name  string
	short string
	val   interface{}
	use   string
}

func setFlags(flagset *pflag.FlagSet) {
	for _, f := range flags {
		switch val := f.val.(type) {
		case []string:
			flagset.StringArrayP(f.name, f.short, val, f.use)
		case bool:
			flagset.BoolP(f.name, f.short, val, f.use)
		case string:
			flagset.StringP(f.name, f.short, val, f.use)
		}
	}
}

func main() {
	var exitCode int

	cmd := &cobra.Command{}
	cmd.SetUsageTemplate(usage)
	cmd.SetHelpTemplate(help)

	setFlags(cmd.Flags())

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(os.Args) == 1 {
			return cmd.Usage()
		}

		flagset := cmd.Flags()

		var (
			pattern string
			files   = args
		)

		if len(args) > 0 {
			pattern = args[0]
			files = args[1:]
		}

		var opts []grep.Opt

		for _, flag := range flags {
			if !flagset.Lookup(flag.name).Changed {
				continue
			}

			switch name := flag.name; name {
			case "regexp":
				e, err := flagset.GetStringArray(name)
				if err != nil {
					return err
				}
				for _, pattern := range e {
					opts = append(opts, grep.WithRegexps(pattern))
				}
				pattern = ""
				files = args
			case "file":
				f, err := flagset.GetStringArray(name)
				if err != nil {
					return err
				}
				for _, filename := range f {
					file, err := os.Open(filename)
					if err != nil {
						return err
					}
					opts = append(opts, grep.WithFiles(file))
				}
				pattern = ""
				files = args
			case "ignore-case":
				i, err := flagset.GetBool(name)
				if err != nil {
					return err
				}
				if i {
					opts = append(opts, grep.WithIgnoreCase())
				}

				// TODO: if pattern == "" show usage
			case "invert-match":
				v, err := flagset.GetBool(name)
				if err != nil {
					return err
				}
				if v {
					opts = append(opts, grep.WithInvertMatch())
				}
			case "word-regexp":
				w, err := flagset.GetBool(name)
				if err != nil {
					return err
				}
				if w {
					opts = append(opts, grep.WithWordRegexp())
				}
			case "line-regexp":
				x, err := flagset.GetBool(name)
				if err != nil {
					return err
				}
				if x {
					opts = append(opts, grep.WithLineRegexp())
				}
			}
		}

		var inputs []io.Reader
		for _, filename := range files {
			if filename == "-" {
				inputs = append(inputs, os.Stdin)
				continue
			}
			file, err := os.Open(filename)
			if err != nil {
				fmt.Printf("gogrep: %s: No such file or directory\n", filename)
				exitCode = 1
				continue
			}
			inputs = append(inputs, file)
			defer file.Close()
		}

		var input io.Reader
		if len(inputs) > 0 {
			input = io.MultiReader(inputs...)
		} else {
			input = os.Stdin
		}

		output := grep.New(pattern, opts...).Exec(input)
		_, err := io.Copy(os.Stdout, output)
		return err
	}

	cmd.Execute()
	if exitCode == 0 {
		exitCode = 1
	}
	os.Exit(exitCode)
}

const usage = `Usage: gogrep [OPTION]... PATTERN [FILE]...
Try 'gogrep --help' for more information.
`

const help = `Usage: gogrep [OPTION]... PATTERN [FILE]...
Search for PATTERN in each FILE.
Example: gogrep -i 'hello world' menu.h main.c

Pattern selection and interpretation:
  -E, --extended-regexp     PATTERN is an extended regular expression
  -F, --fixed-strings       PATTERN is a set of newline-separated strings
  -G, --basic-regexp        PATTERN is a basic regular expression (default)
  -P, --perl-regexp         PATTERN is a Perl regular expression
  -e, --regexp=PATTERN      use PATTERN for matching
  -f, --file=FILE           obtain PATTERN from FILE
  -i, --ignore-case         ignore case distinctions
  -w, --word-regexp         force PATTERN to match only whole words
  -x, --line-regexp         force PATTERN to match only whole lines
  -z, --null-data           a data line ends in 0 byte, not newline

Miscellaneous:
  -s, --no-messages         suppress error messages
  -v, --invert-match        select non-matching lines
  -V, --version             display version information and exit
      --help                display this help text and exit

Output control:
  -m, --max-count=NUM       stop after NUM selected lines
  -b, --byte-offset         print the byte offset with output lines
  -n, --line-number         print line number with output lines
      --line-buffered       flush output on every line
  -H, --with-filename       print file name with output lines
  -h, --no-filename         suppress the file name prefix on output
      --label=LABEL         use LABEL as the standard input file name prefix
  -o, --only-matching       show only the part of a line matching PATTERN
  -q, --quiet, --silent     suppress all normal output
      --binary-files=TYPE   assume that binary files are TYPE;
                            TYPE is 'binary', 'text', or 'without-match'
  -a, --text                equivalent to --binary-files=text
  -I                        equivalent to --binary-files=without-match
  -d, --directories=ACTION  how to handle directories;
                            ACTION is 'read', 'recurse', or 'skip'
  -D, --devices=ACTION      how to handle devices, FIFOs and sockets;
                            ACTION is 'read' or 'skip'
  -r, --recursive           like --directories=recurse
  -R, --dereference-recursive  likewise, but follow all symlinks
      --include=FILE_PATTERN  search only files that match FILE_PATTERN
      --exclude=FILE_PATTERN  skip files and directories matching FILE_PATTERN
      --exclude-from=FILE   skip files matching any file pattern from FILE
      --exclude-dir=PATTERN  directories that match PATTERN will be skipped.
  -L, --files-without-match  print only names of FILEs with no selected lines
  -l, --files-with-matches  print only names of FILEs with selected lines
  -c, --count               print only a count of selected lines per FILE
  -T, --initial-tab         make tabs line up (if needed)
  -Z, --null                print 0 byte after FILE name

Context control:
  -B, --before-context=NUM  print NUM lines of leading context
  -A, --after-context=NUM   print NUM lines of trailing context
  -C, --context=NUM         print NUM lines of output context
  -NUM                      same as --context=NUM
      --color[=WHEN],
      --colour[=WHEN]       use markers to highlight the matching strings;
                            WHEN is 'always', 'never', or 'auto'
  -U, --binary              do not strip CR characters at EOL (MSDOS/Windows)

When FILE is '-', read standard input.  With no FILE, read '.' if
recursive, '-' otherwise.  With fewer than two FILEs, assume -h.
Exit status is 0 if any line is selected, 1 otherwise;
if any error occurs and -q is not given, the exit status is 2.
`
