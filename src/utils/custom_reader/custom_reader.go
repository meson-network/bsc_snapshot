package custom_reader

import (
	"io"
	"sync/atomic"

	"github.com/cheggaaa/pb/v3"
	"github.com/vbauerster/mpb/v8"
)

type CustomReader struct {
	Reader      io.Reader
	Size        int64
	Have_read   int64
	UploadBar   *mpb.Bar
	DownloadBar *pb.ProgressBar
}

func (r *CustomReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	atomic.AddInt64(&r.Have_read, int64(n))

	if r.UploadBar != nil {
		r.UploadBar.SetCurrent(int64(float32(r.Have_read*100) / float32(r.Size)))
	}
	if r.DownloadBar != nil {
		r.DownloadBar.Add64(int64(n))
	}
	return n, err
}
