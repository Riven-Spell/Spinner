package common

import (
	"io"
)

type RedirectWriter struct {
	UnderlyingWriters []io.Writer
}

func (r *RedirectWriter) Write(b []byte) (n int, err error) {
	for _, v := range r.UnderlyingWriters {
		_, err := v.Write(b)

		if err != nil {
			return len(b), err
		}
	}

	return len(b), nil
}
