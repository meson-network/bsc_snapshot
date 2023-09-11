package cmd_upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/meson-network/bsc-data-file-utils/src/file_config"
	"github.com/meson-network/bsc-data-file-utils/src/uploader/uploader_r2"

	"github.com/gosuri/uiprogress"
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
	if thread == 0 {
		thread = default_thread
	}
	threadChan := make(chan struct{}, thread)

	if retry_times == 0 {
		retry_times = default_retry_times
	}

	// read json from originDir
	configFilePath := filepath.Join(originDir, file_config.FILES_CONFIG_JSON_NAME)
	jsonContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	fileConfig := file_config.FileConfig{}
	err = json.Unmarshal(jsonContent, &fileConfig)
	if err != nil {
		return err
	}

	client, err := uploader_r2.GenR2Client(accountId, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	fmt.Println("start upload")

	uiprogress.Start()

	fileList := fileConfig.ChunkedFileList
	errorFiles := []*file_config.ChunkedFileInfo{}
	var errorFilesLock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(fileList))
	counter := int64(0)
	for _, v := range fileList {
		fileInfo := v
		threadChan <- struct{}{}
		go func() {
			defer func() {
				<-threadChan
				wg.Done()
			}()
			c := atomic.AddInt64(&counter, 1)
			bar := uiprogress.AddBar(100).AppendCompleted().PrependElapsed()
			bar.PrependFunc(func(b *uiprogress.Bar) string {
				return fmt.Sprintf(" %d / %d %s ", c, len(fileList), fileInfo.FileName)
			})

			// try some times if upload failed
			for try := 0; try < retry_times; try++ {
				bar.Set(0)

				localFilePath := filepath.Join(originDir, fileInfo.FileName)
				err := uploader_r2.UploadFile(client, bucketName, additional_path, fileInfo.FileName, localFilePath, fileInfo.Md5, bar)

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
					errorFiles = append(errorFiles, &fileInfo)
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

	if len(errorFiles) > 0 {
		uiprogress.Stop()
		fmt.Println("the following files upload failed, please try again:")
		for _, v := range errorFiles {
			fmt.Println(v.FileName)
		}
		return errors.New("upload error")
	}

	//upload config
	localFilePath := filepath.Join(originDir, file_config.FILES_CONFIG_JSON_NAME)
	bar := uiprogress.AddBar(100).AppendCompleted().PrependElapsed()
	bar.Set(0)
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf(" %s ", file_config.FILES_CONFIG_JSON_NAME)
	})

	err = uploader_r2.UploadFile(client, bucketName, additional_path, file_config.FILES_CONFIG_JSON_NAME, localFilePath, "", bar)
	if err != nil {
		bar.AppendFunc(func(b *uiprogress.Bar) string {
			return "FAILED"
		})
		uiprogress.Stop()
		fmt.Println("upload json file error")
	} else {
		bar.Set(100)
		bar.AppendFunc(func(b *uiprogress.Bar) string {
			return "SUCCESS"
		})
		uiprogress.Stop()
		fmt.Println("upload job finish")
	}

	return nil
}
