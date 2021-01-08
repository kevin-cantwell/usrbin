package getopt_test

import (
	"testing"

	"github.com/kevin-cantwell/usrbin/getopt"
)

func TestGetopt(t *testing.T) {
	tests := []struct {
		name    string
		inOpts    []getopt.Opt
		in      []string
		outOpts []getopt.Option
		outArgs []string
		outErr  *getopt.Error
	}{
		{
			name: "a",
			in:   "foo bar baz",
			outOpts:  "",
		},
		{
			name: "b",
			in:   "-abc",
			out:  "-a -b -c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := getopt.New(
				"-a","--aye",
				"b",
				"abc", "aye,bee,cee", tt.inOpts
				).Parse("foo", "bar", "baz")
			if err != nil {
				t.Fatalf("got err: %+v", err)
			}

			for _, opt := range output.Options {

			}
			if string(b) != tt.out {
				t.Errorf("got %q want %q", b, tt.out)
			}
		})
	}

}
