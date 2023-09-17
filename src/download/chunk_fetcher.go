package download

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"

	"github.com/meson-network/bsc_snapshot/src/model"
	"github.com/meson-network/bsc_snapshot/src/utils/custom_reader"
)

type ChunkFetcher struct {
	filePath    string
	chunkInfo   *model.ChunkedFileInfo
	downloadBar *pb.ProgressBar
}

func NewChunkFetcher(filePath string, chunkInfo *model.ChunkedFileInfo, bar *pb.ProgressBar) *ChunkFetcher {

	return &ChunkFetcher{
		filePath:    filePath,
		chunkInfo:   chunkInfo,
		downloadBar: bar,
	}
}

func (c *ChunkFetcher) Download(endPoints []string, retryNum int,
	errFn func(*model.ChunkedFileInfo)) bool {

	file, err := os.OpenFile(c.filePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		errFn(c.chunkInfo)
	}
	defer file.Close()

	_, err = file.Seek(c.chunkInfo.Offset, 0)
	if err != nil {
		// fmt.Println("seek err:", err)
		errFn(c.chunkInfo)
	}
	defer resp.Body.Close()

	// try some times if download failed
	downloaded := int64(0)
	md5hash := md5.New()
	for try := 0; try < retryNum; try++ {
		// pick endpoint random
		currentEndpoint := endPoints[rand.Intn(len(endPoints))]
		downloadUrl := currentEndpoint + "/" + c.chunkInfo.FileName

		downloadSize, err := c.doDownload(downloadUrl, downloaded, file, &md5hash)
		if err != nil {
			downloaded += downloadSize
			if try < retryNum-1 {
				time.Sleep(RETRY_WAIT_SECS * time.Second)
				continue
			}
			errFn(c.chunkInfo)
		} else {
			return true
		}
	}

	return false
}

func (c *ChunkFetcher) doDownload(url string, startPos int64, file *os.File, md5hash *hash.Hash) (int64, error) {

	// download
	resp, err := c.fetchChunk(url, int(startPos))
	if err != nil {
		// fmt.Println("fetch err:", err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return 0, errors.New("response status code error")
	}

	haveRead := startPos
	if startPos > 0 {
		if resp.StatusCode != http.StatusPartialContent {
			haveRead = 0
		}

		contentRange := resp.Header.Get("Content-Range")
		if contentRange == "" {
			haveRead = 0
		}
	}

	if (c.chunkInfo.Size - haveRead) != resp.ContentLength {
		return 0, errors.New("remote file size error")
	}

	// use custom reader to show upload progress
	reader := &custom_reader.CustomReader{
		Reader:      resp.Body,
		Size:        c.chunkInfo.Size,
		Pos:         haveRead,
		DownloadBar: c.downloadBar,
		UploadBar:   nil,
	}

	file.Seek(c.chunkInfo.Offset+haveRead, 0)
	return writeChunk(reader, file, c.chunkInfo.Size, c.chunkInfo.Md5, md5hash)
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

func writeChunk(src io.Reader, dst io.Writer, chunkSize int64, chunkMd5 string, md5Hash *hash.Hash) (int64, error) {

	target := io.MultiWriter(dst, *md5Hash)

	copySize := copyContent(target, src)

	md5Str := hex.EncodeToString((*md5Hash).Sum(nil))
	if strings.EqualFold(md5Str, chunkMd5) {
		return copySize, nil
	}

	return copySize, errors.New("md5 not equal")
}

func (c *ChunkFetcher) fetchChunk(url string, offset int) (*http.Response, error) {
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

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if offset > 0 {
		req.Header.Set("Range", "bytes="+strconv.Itoa(offset)+"-")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func copyContent(dst io.Writer, src io.Reader) int64 {
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

	return int64(written)
}
