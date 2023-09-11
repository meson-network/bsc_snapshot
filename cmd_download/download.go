package cmd_download

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gosuri/uiprogress"
	"github.com/meson-network/bsc-data-file-utils/src/common/custom_reader"
	"github.com/meson-network/bsc-data-file-utils/src/file_config"
	"github.com/urfave/cli/v2"
)

const default_retry_times = 3
const default_thread = 5

func Download(clictx *cli.Context) error {

	jsonConfig := clictx.String("file_config")
	thread := clictx.Int("thread")
	retry_times := clictx.Int("retry_times")

	if jsonConfig == "" {
		fmt.Println("[ERROR] json config error")
		return errors.New("json config error")
	}

	if thread == 0 {
		thread = default_thread
	}
	threadChan := make(chan struct{}, thread)

	if retry_times == 0 {
		retry_times = default_retry_times
	}

	// download or read jsonConfig
	config := file_config.FileConfig{}
	if strings.HasPrefix(jsonConfig, "http") {
		// download json
		client := &http.Client{Timeout: 20 * time.Second}
		resp, err := client.Get(jsonConfig)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		json.Unmarshal(content, &config)
	} else {
		// read json file
		content, err := os.ReadFile(jsonConfig)
		if err != nil {
			return err
		}
		json.Unmarshal(content, &config)
	}

	// check endpoint
	endPoint := ""
	if len(config.EndPoint) == 0 {
		if strings.HasPrefix(jsonConfig, "http") {
			i := strings.LastIndex(jsonConfig, "/")
			if i < 0 {
				return errors.New("download end point error")
			}
			endPoint = jsonConfig[:i]
		} else {
			return errors.New("download endpoint not exist")
		}
	} else {
		endPoint = config.EndPoint[0]
	}

	// gen raw file
	rawFilePath := filepath.Join("./", config.RawFile.FileName)
	fileStat, _ := os.Stat(rawFilePath)
	// file already exist
	if fileStat != nil {
		fmt.Println("file already exist")
		return errors.New("file exist")
	}

	downloadingFileName := config.RawFile.FileName + ".downloading"
	downloadingFilePath := filepath.Join("./", downloadingFileName)
	dFile, err := os.OpenFile(downloadingFilePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}

	err = dFile.Truncate(config.RawFile.Size)
	if err != nil {
		dFile.Close()
		return err
	}
	dFile.Close()

	uiprogress.Start()

	errorFiles := []*file_config.ChunkedFileInfo{}
	var errorFilesLock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(config.ChunkedFileList))
	counter := int64(0)
	for _, v := range config.ChunkedFileList {
		chunkInfo := v
		threadChan <- struct{}{}
		go func() {
			defer func() {
				<-threadChan
				wg.Done()
			}()

			c := atomic.AddInt64(&counter, 1)
			bar := uiprogress.AddBar(100).AppendCompleted().PrependElapsed()
			bar.PrependFunc(func(b *uiprogress.Bar) string {
				return fmt.Sprintf(" %d / %d %s ", c, len(config.ChunkedFileList), chunkInfo.FileName)
			})

			// try some times if download failed
			for try := 0; try < retry_times; try++ {
				bar.Set(0)

				downloadUrl := endPoint + "/" + chunkInfo.FileName
				err := downloadPart(downloadUrl, downloadingFilePath, chunkInfo.Size, chunkInfo.Offset, chunkInfo.Md5, bar)
				if err != nil {
					if try < retry_times-1 {
						time.Sleep(3 * time.Second)
						continue
					}
					bar.AppendFunc(func(b *uiprogress.Bar) string {
						return "FAILED"
					})
					errorFilesLock.Lock()
					defer errorFilesLock.Unlock()
					errorFiles = append(errorFiles, &chunkInfo)
				} else {
					bar.Set(100)
					bar.AppendFunc(func(b *uiprogress.Bar) string {
						return "SUCCESS"
					})
					break
				}
			}
		}()
	}

	wg.Wait()
	uiprogress.Stop()
	if len(errorFiles) > 0 {

		fmt.Println("the following files download failed, please try again:")
		for _, v := range errorFiles {
			fmt.Println(v.FileName)
		}
		return errors.New("download error")
	}

	err = os.Rename(downloadingFilePath, rawFilePath)
	if err != nil {
		return err
	}

	fmt.Println("download finish")

	return nil
}

func downloadPart(downloadUrl string, downloadFilePath string, chunkSize int64, chunkOffset int64, chunkMd5 string, bar *uiprogress.Bar) error {
	// read local md5
	file, err := os.OpenFile(downloadFilePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Seek(chunkOffset, 0)
	if err != nil {
		fmt.Println("seek err:", err)
		return err
	}
	h := md5.New()
	_, err = io.CopyN(h, file, chunkSize)
	if err != nil {
		fmt.Println("read exist file chunk err:", err)
	} else {
		md5Str := hex.EncodeToString(h.Sum(nil))
		if md5Str == chunkMd5 {
			return nil
		}
	}

	// download
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if chunkSize != resp.ContentLength {
		return errors.New("remote file size error")
	}

	file.Seek(chunkOffset, 0)
	h = md5.New()
	target := io.MultiWriter(file, h)

	// use custom reader to show upload progress
	reader := &custom_reader.CustomReader{
		Reader: resp.Body,
		Size:   chunkSize,
		Bar:    bar,
	}

	_, err = io.CopyN(target, reader, chunkSize)
	if err != nil {
		fmt.Println("read body err:", err)
		return err
	}
	md5Str := hex.EncodeToString(h.Sum(nil))
	if md5Str == chunkMd5 {
		return nil
	}

	return errors.New("md5 not equal")

}
