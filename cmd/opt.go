package cmd

/*
	Notes on opt behavior in gnu programs:

	# grep --w
	grep: option '--w' is ambiguous; possibilities: '--with-filename' '--word-regexp'

	# grep --with-filename=foo
	grep: option '--with-filename' doesn't allow an argument

	# grep --with-filenamef
	grep: unrecognized option '--with-filenamef'

	# grep -t
	grep: invalid option -- 't'

	# grep --t # ok because it unambiguously matches '--text'

	# grep -x=1
	grep: invalid option -- '='

	# ls -wo
	ls: invalid line width: 'o'

	# ls -w1 # ok
	# ls -w 1 # ok
*/

// var (
// 	optset = map[string]opt{}
// )

// type opt struct {
// 	short string
// 	long  string
// 	value interface{}
// }

// func (o opt) Bool() bool {
// 	getopt.Parse()
// 	b, ok := o.value.(bool)
// 	if !ok {
// 		fmt.Println()
// 	}
// }

// type args []string

// func (a args) opts(alias string) ([]opt, bool) {

// }
