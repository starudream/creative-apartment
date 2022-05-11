package iio

import (
	"bytes"
	"io"
)

func ReadBody(r io.ReadCloser) (io.ReadCloser, []byte) {
	bs, err := io.ReadAll(r)
	if err != nil {
		return r, nil
	}
	return io.NopCloser(bytes.NewReader(bs)), bs
}
