package download

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/vbauerster/mpb/v8"

	"github.com/meson-network/bsc-snapshot/src/utils/custom_reader"
)

type ChunkFetcher struct {
	filePath    string
	chunkSize   int64
	chunkOffset int64
	chunkMd5    string
	bar         *mpb.Bar
}

func NewChunkFetcher(filePath string,
	chunkSize int64, chunkOffset int64, chunkMd5 string, bar *mpb.Bar) *ChunkFetcher {

	return &ChunkFetcher{
		filePath:    filePath,
		chunkSize:   chunkSize,
		chunkOffset: chunkOffset,
		chunkMd5:    chunkMd5,
		bar:         bar,
	}
}

func (c *ChunkFetcher) Download(downloadUrl string) error {
	file, err := os.OpenFile(c.filePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Seek(c.chunkOffset, 0)
	if err != nil {
		// fmt.Println("seek err:", err)
		return err
	}

	if exists := validateChunk(file, c.chunkSize, c.chunkMd5); exists {
		return nil
	}

	// download
	resp, err := c.fetchChunk(downloadUrl)
	if err != nil {
		// fmt.Println("fetch err:", err)
		return err
	}
	defer resp.Body.Close()
	if c.chunkSize != resp.ContentLength {
		return errors.New("remote file size error")
	}

	// use custom reader to show upload progress
	reader := &custom_reader.CustomReader{
		Reader: resp.Body,
		Size:   c.chunkSize,
		Bar:    c.bar,
	}

	file.Seek(c.chunkOffset, 0)
	return writeChunk(reader, file, c.chunkSize, c.chunkMd5)
}

func validateChunk(src io.Reader, chunkSize int64, chunkMd5 string) bool {

	md5hash := md5.New()
	_, err := io.CopyN(md5hash, src, chunkSize)
	if err != nil {
		// fmt.Println("read exist file chunk err:", err)
		return false
	} else {
		md5Str := hex.EncodeToString(md5hash.Sum(nil))
		if strings.EqualFold(md5Str, chunkMd5) {
			return true
		}
	}

	return false
}

func writeChunk(src io.Reader, dst io.Writer, chunkSize int64, chunkMd5 string) error {

	md5hash := md5.New()
	target := io.MultiWriter(dst, md5hash)

	copyContent(target, src)

	md5Str := hex.EncodeToString(md5hash.Sum(nil))
	if strings.EqualFold(md5Str, chunkMd5) {
		return nil
	}

	return errors.New("md5 not equal")
}

func (c *ChunkFetcher) fetchChunk(downloadUrl string) (*http.Response, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(network, addr, DEFAULT_REQUEST_TIMEOUT)
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
			ResponseHeaderTimeout: DEFAULT_REQUEST_TIMEOUT,
		},
	}
	resp, err := client.Get(downloadUrl)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func copyContent(dst io.Writer, src io.Reader) {
	// need check download speed
	buff := make([]byte, 32*1024)
	written := 0
	finishChan := make(chan struct{})
	go func() {
		defer func() {
			finishChan <- struct{}{}
		}()
		for {
			nr, er := src.Read(buff)
			if nr > 0 {
				nw, ew := dst.Write(buff[0:nr])
				if nw > 0 {
					written += nw
				}
				if ew != nil {
					// err = ew
					break
				}
				if nr != nw {
					// err = io.ErrShortWrite
					break
				}
			}
			if er != nil {
				if er != io.EOF {
					// err = er
				}
				break
			}
		}
	}()

	spaceTime := time.Second * 10
	ticker := time.NewTicker(spaceTime)
	defer ticker.Stop()
	lastWtn := 0
	stop := false

outLoop:
	for {
		select {
		case <-finishChan:
			break outLoop
		case <-ticker.C:
			// if no data transfer in 10 seconds
			if written-lastWtn == 0 {
				stop = true
				break
			}
			lastWtn = written
		}
		if stop {
			break
		}
	}
}
