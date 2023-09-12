package cmd_download

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"

	"github.com/meson-network/bsc-data-file-utils/src/common/custom_reader"
	"github.com/meson-network/bsc-data-file-utils/src/file_config"
	"github.com/urfave/cli/v2"
)

const default_retry_times = 3
const default_thread = 5

func Download(clictx *cli.Context) error {

	jsonConfigAddress := clictx.String("file_config")
	thread := clictx.Int("thread")
	retry_times := clictx.Int("retry_times")

	if jsonConfigAddress == "" {
		fmt.Println("[ERROR] json config error, please input correct address or file path")
		return errors.New("json config error")
	}

	if thread <= 0 {
		thread = default_thread
	}
	threadChan := make(chan struct{}, thread)

	if retry_times <= 0 {
		retry_times = default_retry_times
	}

	// download or read jsonConfig
	config := file_config.FileConfig{}
	if strings.HasPrefix(jsonConfigAddress, "http") {
		// download json
		client := &http.Client{Timeout: 20 * time.Second}
		resp, err := client.Get(jsonConfigAddress)
		if err != nil {
			fmt.Println("[ERROR] get json config error:", err.Error())
			return err
		}
		defer resp.Body.Close()
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("[ERROR] get json config error:", err.Error())
			return err
		}
		err = json.Unmarshal(content, &config)
		if err != nil {
			fmt.Println("[ERROR] json config unmarshal error:", err.Error())
			return err
		}
	} else {
		// read json file
		content, err := os.ReadFile(jsonConfigAddress)
		if err != nil {
			fmt.Println("[ERROR] read json config error:", err.Error())
			return err
		}
		err = json.Unmarshal(content, &config)
		if err != nil {
			fmt.Println("[ERROR] json config unmarshal error:", err.Error())
			return err
		}
	}

	// check endpoint
	endPoint := ""
	if len(config.EndPoint) == 0 {
		// if no end point info in json, default use json config download path
		if strings.HasPrefix(jsonConfigAddress, "http") {
			i := strings.LastIndex(jsonConfigAddress, "/")
			if i < 0 {
				fmt.Println("[ERROR] download endpoint error")
				return errors.New("download endpoint error")
			}
			fmt.Println("[INFO] use some endpoint with json config file")
			endPoint = jsonConfigAddress[:i]
		} else {
			fmt.Println("[ERROR] download endpoint not exist")
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
		fmt.Println("[ERROR]", "<", config.RawFile.FileName, ">", "already exist")
		return errors.New("file exist")
	}

	downloadingFileName := config.RawFile.FileName + ".downloading"
	downloadingFilePath := filepath.Join("./", downloadingFileName)
	dFile, err := os.OpenFile(downloadingFilePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println("[ERROR] open file err:", err.Error())
		return err
	}

	err = dFile.Truncate(config.RawFile.Size)
	if err != nil {
		dFile.Close()
		fmt.Println("[ERROR] handle file err:", err.Error())
		return err
	}
	dFile.Close()

	fmt.Println("[INFO] start download...")

	errorFiles := []*file_config.ChunkedFileInfo{}
	var errorFilesLock sync.Mutex
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithAutoRefresh())
	wg.Add(len(config.ChunkedFileList))
	counter := int64(0)
	for _, v := range config.ChunkedFileList {
		chunkInfo := v
		threadChan <- struct{}{}

		c := atomic.AddInt64(&counter, 1)
		bar := p.AddBar(int64(100),
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(fmt.Sprintf(" %d / %d %s ", c, len(config.ChunkedFileList), chunkInfo.FileName)),
				// decor.DSyncWidth bit enables column width synchronization
				decor.Percentage(decor.WCSyncSpace),
			),
			mpb.AppendDecorators(
				decor.OnComplete(
					decor.Name(""), "SUCCESS ",
				),
				decor.OnAbort(
					decor.Elapsed(decor.ET_STYLE_GO), "FAILED ",
				),
				
			),
		)
		bar.SetPriority(int(c))

		go func() {
			defer func() {
				<-threadChan
				wg.Done()
			}()

			bar.SetPriority(math.MaxInt - int(c))

			// try some times if download failed
			for try := 0; try < retry_times; try++ {
				bar.SetCurrent(0)

				downloadUrl := endPoint + "/" + chunkInfo.FileName
				err := downloadPart(downloadUrl, downloadingFilePath, chunkInfo.Size, chunkInfo.Offset, chunkInfo.Md5, bar)
				if err != nil {
					if try < retry_times-1 {
						time.Sleep(3 * time.Second)
						continue
					}
					errorFilesLock.Lock()
					defer errorFilesLock.Unlock()
					errorFiles = append(errorFiles, &chunkInfo)
					bar.Abort(false)
					bar.SetPriority(math.MaxInt - int(c) - len(config.ChunkedFileList))
				} else {
					if !bar.Completed() {
						bar.SetCurrent(100)
					}
					bar.SetPriority(int(c))
					break
				}
			}
		}()
	}
	p.Wait()
	wg.Wait()
	if len(errorFiles) > 0 {

		fmt.Println("[ERROR] the following files download failed, please try again:")
		for _, v := range errorFiles {
			fmt.Println(v.FileName)
		}
		return errors.New("download error")
	}

	err = os.Rename(downloadingFilePath, rawFilePath)
	if err != nil {
		fmt.Println("[ERROR] rename download file err:",err.Error())
		return err
	}

	fmt.Println("[INFO] download finish")

	return nil
}

func downloadPart(downloadUrl string, downloadFilePath string, chunkSize int64, chunkOffset int64, chunkMd5 string, bar *mpb.Bar) error {
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
