package custom_reader

import (
	"io"
	"sync/atomic"

	"github.com/vbauerster/mpb/v8"
)

type CustomReader struct {
	Reader io.Reader
	Size   int64
	read   int64
	*mpb.Bar
}

func (r *CustomReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	atomic.AddInt64(&r.read, int64(n))

	if r.Bar != nil {
		r.Bar.SetCurrent(int64(float32(r.read*100) / float32(r.Size)))
	}
	return n, err
}
