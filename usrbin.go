package usrbin

import "io"

type Reader interface {
	Read(io.Reader) io.Reader
}

type Execer interface {
	Exec(params []string) io.Reader
}

func Pipe(in io.Reader, pipes ...Reader) io.Reader {
	out, w := io.Pipe()

	go func() {
		for _, pipe := range pipes {
			in = pipe.Read(in)
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
