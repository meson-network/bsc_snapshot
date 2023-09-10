package cmd_download

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/meson-network/bsc-data-file-utils/src/split_file"
	"github.com/urfave/cli/v2"
)

func Download(clictx *cli.Context) error {

	jsonConfig := clictx.String("file_config")
	thread := clictx.Int("thread")

	if jsonConfig == "" {
		fmt.Println("[ERROR] json config error")
		return errors.New("json config error")
	}

	if thread == 0 {
		thread = 3
	}
	threadChan := make(chan struct{}, thread)

	// download or read jsonConfig
	config := split_file.FileSplitConfig{}
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
	if len(config.EndPoint) == 0 {
		return errors.New("download endpoint not exist")
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
	dFile, err := os.OpenFile(downloadingFilePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	err = dFile.Truncate(config.RawFile.Size)
	if err != nil {
		return err
	}
	dFile.Close()

	errorFiles := []*split_file.ChunkedFileInfo{}
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

			downloadUrl := config.EndPoint[0] + "/" + chunkInfo.FileName
			err := downloadPart(downloadUrl, downloadingFilePath, chunkInfo.Size, chunkInfo.Offset, chunkInfo.Md5)
			c := atomic.AddInt64(&counter, 1)
			if err != nil {
				log.Println(c, "/", len(config.ChunkedFileList), chunkInfo.FileName, "download err:", err)
				errorFilesLock.Lock()
				defer errorFilesLock.Unlock()
				errorFiles = append(errorFiles, &chunkInfo)
			} else {
				log.Println(c, "/", len(config.ChunkedFileList), chunkInfo.FileName, "downloaded")
			}
		}()

	}

	wg.Wait()

	if len(errorFiles) > 0 {
		log.Println("the following files download failed, please try again:")
		for _, v := range errorFiles {
			log.Println(v.FileName)
		}
		return errors.New("download error")
	}

	err=os.Rename(downloadingFilePath,rawFilePath)
	if err!=nil{
		return err
	}

	log.Println("download finish")

	return nil
}

func downloadPart(downloadUrl string, downloadFilePath string, chunkSize int64, chunkOffset int64, chunkMd5 string) error {
	// read local md5
	file, err := os.OpenFile(downloadFilePath, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Seek(chunkOffset, 0)
	if err != nil {
		log.Println("seek err:", err)
		return err
	}
	h := md5.New()
	_, err = io.CopyN(h, file, chunkSize)
	if err != nil {
		log.Println("read exist file chunk err:", err)
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
	file.Seek(chunkOffset, 0)
	h = md5.New()
	target := io.MultiWriter(file, h)

	_, err = io.CopyN(target, resp.Body,chunkSize)
	if err != nil {
		log.Println("read body err:", err)
		return err
	}
	md5Str := hex.EncodeToString(h.Sum(nil))
	if md5Str == chunkMd5 {
		return nil
	}

	return errors.New("md5 not equal")

}
