package grep_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/kevin-cantwell/usrbin/pkg/grep"
)

func TestGrep(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		opts    []grep.Option
		in      string
		out     string
	}{
		{
			name:    "literal",
			pattern: "foo",
			in:      "foo\nbar\nbaz",
			out:     "foo\n",
		},
		{
			name:    "literal/newlines",
			pattern: "foo\nbar",
			in:      "foo\nbar\nbaz",
			out:     "foo\nbar\n",
		},
		{
			name:    "WithRegexps/single",
			pattern: "",
			opts:    []grep.Option{grep.WithRegexps("foo")},
			in:      "foo\nbar\nbaz\nfoobaz",
			out:     "foo\nfoobaz\n",
		},
		{
			name:    "WithRegexps/newlines",
			pattern: "",
			opts:    []grep.Option{grep.WithRegexps("foo\nbar")},
			in:      "foo\nbar\nbaz\nfoobaz",
			out:     "foo\nbar\nfoobaz\n",
		},
		{
			name:    "WithRegexps/multi",
			pattern: "",
			opts:    []grep.Option{grep.WithRegexps("foo", "bar\nbaz")},
			in:      "foo\nbar\nbaz\nfoobaz",
			out:     "foo\nbar\nbaz\nfoobaz\n",
		},
		{
			name:    "WithIgnoreCase",
			pattern: "FOO",
			opts:    []grep.Option{grep.WithIgnoreCase()},
			in:      "foo\nbar\nbaz",
			out:     "foo\n",
		},
		{
			name:    "WithIgnoreCase/fold-case",
			pattern: "(?i)FOO",
			opts:    []grep.Option{grep.WithIgnoreCase()},
			in:      "foo\nbar\nbaz",
			out:     "foo\n",
		},
		{
			name:    "WithIgnoreCase+WithInvertMatch",
			pattern: "FOO",
			opts:    []grep.Option{grep.WithIgnoreCase(), grep.WithInvertMatch()},
			in:      "foo\nbar\nbaz",
			out:     "bar\nbaz\n",
		},
		{
			name:    "WithIgnoreCase+WithWordRegexp",
			pattern: "FOO",
			opts:    []grep.Option{grep.WithIgnoreCase(), grep.WithWordRegexp()},
			in:      "foo\nbar\nbaz\nfoobar",
			out:     "foo\n",
		},
		{
			name:    "WithIgnoreCase+WithWordRegexp+WithInvertMatch",
			pattern: "FOO",
			opts:    []grep.Option{grep.WithIgnoreCase(), grep.WithWordRegexp(), grep.WithInvertMatch()},
			in:      "foo\nbar\nbaz\nfoobar",
			out:     "bar\nbaz\nfoobar\n",
		},
		{
			name:    "WithInvertMatch",
			pattern: "foo",
			opts:    []grep.Option{grep.WithInvertMatch()},
			in:      "foo\nbar\nbaz",
			out:     "bar\nbaz\n",
		},
		{
			name:    "WithWordRegexp",
			pattern: "foo",
			opts:    []grep.Option{grep.WithWordRegexp()},
			in:      "foo\nfoo bar\nbaz foo\nbar_foo_baz\nfoo-bar\nbar0foo",
			out:     "foo\nfoo bar\nbaz foo\nfoo-bar\n",
		},
		{
			name:    "WithLineRegexp",
			pattern: "foo|baz",
			opts:    []grep.Option{grep.WithLineRegexp()},
			in:      "foo\nbar\nbaz\nfoobaz",
			out:     "foo\nbaz\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := strings.NewReader(tt.in)

			out := grep.New(tt.pattern, tt.opts...).Exec(in)

			if body, err := ioutil.ReadAll(out); err != nil {
				t.Fatalf("got err: %#v", err)
			} else if string(body) != tt.out {
				t.Fatalf("got %q want %q", string(body), tt.out)
			}
		})
	}
}
