package cmd_upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/meson-network/bsc-data-file-utils/src/file_config"
	"github.com/meson-network/bsc-data-file-utils/src/uploader/uploader_r2"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"

	"github.com/urfave/cli/v2"
)

const default_retry_times = 3
const default_thread = 5

func Upload_r2(clictx *cli.Context) error {

	//read params
	originDir := clictx.String("dir")
	thread := clictx.Int("thread")
	bucketName := clictx.String("bucket_name")
	additional_path := clictx.String("additional_path")
	accountId := clictx.String("account_id")
	accessKeyId := clictx.String("access_key_id")
	accessKeySecret := clictx.String("access_key_secret")
	retry_times := clictx.Int("retry_times")

	// default use thread
	if thread <= 0 {
		thread = default_thread
	}
	threadChan := make(chan struct{}, thread)

	if retry_times <= 0 {
		retry_times = default_retry_times
	}

	// read json from originDir
	configFilePath := filepath.Join(originDir, file_config.FILES_CONFIG_JSON_NAME)
	jsonContent, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("[ERROR] read file config err:",err)
		return err
	}
	fileConfig := file_config.FileConfig{}
	err = json.Unmarshal(jsonContent, &fileConfig)
	if err != nil {
		fmt.Println("[ERROR] unmarshal file config err:",err)
		return err
	}

	client, err := uploader_r2.GenR2Client(accountId, accessKeyId, accessKeySecret)
	if err != nil {
		fmt.Println("[ERROR] gen r2 client err:",err)
		return err
	}

	fmt.Println("[INFO] start upload...")

	fileList := fileConfig.ChunkedFileList
	errorFiles := []*file_config.ChunkedFileInfo{}
	var errorFilesLock sync.Mutex
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg), mpb.WithAutoRefresh())
	wg.Add(len(fileList))
	counter := int64(0)
	for _, v := range fileList {
		fileInfo := v
		threadChan <- struct{}{}

		c := atomic.AddInt64(&counter, 1)
		bar := p.AddBar(int64(100),
			mpb.BarFillerClearOnComplete(),
			mpb.PrependDecorators(
				// simple name decorator
				decor.Name(fmt.Sprintf(" %d / %d %s ", c, len(fileList), fileInfo.FileName)),
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

			// try some times if upload failed
			for try := 0; try < retry_times; try++ {
				bar.SetCurrent(0)

				localFilePath := filepath.Join(originDir, fileInfo.FileName)
				err := uploader_r2.UploadFile(client, bucketName, additional_path, fileInfo.FileName, localFilePath, fileInfo.Md5, bar)

				if err != nil {
					if try < retry_times-1 {
						time.Sleep(3 * time.Second)
						continue
					}

					errorFilesLock.Lock()
					defer errorFilesLock.Unlock()
					errorFiles = append(errorFiles, &fileInfo)
					bar.Abort(false)
					bar.SetPriority(math.MaxInt - int(c) - len(fileList))
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
		fmt.Println("[ERROR] the following files upload failed, please try again:")
		for _, v := range errorFiles {
			fmt.Println(v.FileName)
		}
		return errors.New("upload error")
	}

	//upload config
	localFilePath := filepath.Join(originDir, file_config.FILES_CONFIG_JSON_NAME)
	u_p:=mpb.New(mpb.WithAutoRefresh())
	bar :=u_p.New(int64(100),
		mpb.BarStyle(),
		mpb.PrependDecorators(
			// simple name decorator
			decor.Name(fmt.Sprintf(" %s ", file_config.FILES_CONFIG_JSON_NAME)),
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

	err = uploader_r2.UploadFile(client, bucketName, additional_path, file_config.FILES_CONFIG_JSON_NAME, localFilePath, "", bar)
	p.Wait()
	if err != nil {
		bar.Abort(false)
		fmt.Println("[ERROR] upload json file error")
	} else {
		if !bar.Completed() {
			bar.SetCurrent(100)
		}
		fmt.Println("[INFO] upload job finish")
	}

	return nil
}
