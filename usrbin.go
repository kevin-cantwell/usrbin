package usrbin

import "io"

type Cmd interface {
	Exec(io.Reader) io.Reader
}

func Pipe(in io.Reader, cmds ...Cmd) io.Reader {
	out, w := io.Pipe()

	go func() {
		for _, cmd := range cmds {
			in = cmd.Exec(in)
		}
		_, err := io.Copy(w, in)
		if err == io.EOF {
			w.Close()
			return
		}
		w.CloseWithError(err)
	}()

	return out
}
