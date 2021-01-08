The `userbin` project takes  aim at translating common POSIX and GNU programs as pure Go interfaces.

Every program translated must naturally consume some input and emit some output. Thus, the core of usrbin is
the `Cmd` interface:

```go
type Cmd interface {
	Exec(input io.Reader) (output io.Reader)
}
``` 

Consider this `grep` example:

```sh
grep -vi "foobar" -
```

In Go, the exact same functionality would look like:

```go
package main

import "github.com/kevin-cantwell/usrbin/grep"

func main() {
    g := grep.New("foobar", grep.WithInvertMatch(), grep.WithIgnoreCase())
    output := g.Exec(os.Stdin)
    io.Copy(os.Stdout, output)
}
```