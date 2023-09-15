package download

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"

	"github.com/meson-network/bsc_snapshot/src/file_config"
)

const (
	DEFAULT_RETRY_TIMES     = 5
	DEFAULT_THREAD          = 5
	DEFAULT_REQUEST_TIMEOUT = time.Second * 7
)

func Download(jsonConfigAddress string, thread int, retryNum int, noResume bool) error {
	if jsonConfigAddress == "" {
		fmt.Println("[ERROR] json config error, please input correct address or file path")
		return errors.New("json config error")
	}
	if thread <= 0 {
		thread = DEFAULT_THREAD
	}
	if retryNum <= 0 {
		retryNum = DEFAULT_RETRY_TIMES
	}

	config, err := loadFileConfig(jsonConfigAddress)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return err
	}

	existEndPoints, err := file_config.FormatEndpoints(config.EndPoint)
	if err != nil {
		fmt.Println("[ERROR] json config endpoint format error")
		return err
	}

	// check endpoint
	var endPoints []string
	if len(existEndPoints) == 0 {
		// if no end point info in json, default use json config download path
		endPointsFromFile, err := parserEndPointFromFileConfig(jsonConfigAddress)
		if err != nil {
			fmt.Println("[ERROR] ", err.Error())
			return err
		}
		fmt.Println("[INFO] use some endpoint with json config file")

		endPoints = endPointsFromFile
	} else {
		endPoints = existEndPoints
	}

	// gen raw file
	rawFilePath := filepath.Join("./", config.RawFile.FileName)

	if fileStat, _ := os.Stat(rawFilePath); fileStat != nil {
		fmt.Println("[ERROR]", "<", config.RawFile.FileName, ">", "already exist")
		return errors.New("file exist")
	}

	downloadingFileName := config.RawFile.FileName + ".downloading"
	downloadingFilePath := filepath.Join("./", downloadingFileName)

	chunkMetaName := "." + config.RawFile.FileName + ".downloaded"
	chunkMetaPath := filepath.Join("./", chunkMetaName)

	if noResume {
		// clean up existing block
		os.Remove(downloadingFilePath)
		os.Remove(chunkMetaPath)
	}

	chunkFetchStat := NewChunkFetchStat()
	err = chunkFetchStat.LoadChunkMeta(chunkMetaPath)
	if err != nil {
		fmt.Println("[ERROR] finished job read error")
		return err
	}
	defer chunkFetchStat.Close()

	dFile, err := os.OpenFile(downloadingFilePath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println("[ERROR] open file err:", err.Error())
		return err
	}

	if err := dFile.Truncate(config.RawFile.Size); err != nil {
		dFile.Close()
		fmt.Println("[ERROR] handle file err:", err.Error())
		return err
	}
	dFile.Close()

	fmt.Println("[INFO] start download...")

	threadChan := make(chan struct{}, thread)
	errorFiles := []*file_config.ChunkedFileInfo{}
	var errorFilesLock sync.Mutex

	var wg sync.WaitGroup
	progressBar := mpb.New(mpb.WithAutoRefresh())

	counter := int64(0)
	for _, v := range config.ChunkedFileList {
		chunkInfo := v
		if chunkFetchStat.IsDone(chunkInfo.FileName) {
			// if already downloaded, skip it
			continue
		}

		threadChan <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				<-threadChan
				wg.Done()
			}()

			c := atomic.AddInt64(&counter, 1)
			bar := progressBar.AddBar(
				int64(100),
				mpb.BarRemoveOnComplete(),
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

			bar.SetPriority(math.MaxInt - len(config.ChunkedFileList) + int(c))

			fetcher := NewChunkFetcher(downloadingFilePath, chunkInfo.Size, chunkInfo.Offset, chunkInfo.Md5, bar)

			// try some times if download failed
			for try := 0; try < retryNum; try++ {
				bar.SetCurrent(0)

				// pick endpoint random
				currentEndpoint := endPoints[rand.Intn(len(endPoints))]
				downloadUrl := currentEndpoint + "/" + chunkInfo.FileName

				err := fetcher.Download(downloadUrl)
				if err != nil {
					if try < retryNum-1 {
						time.Sleep(3 * time.Second)
						continue
					}
					errorFilesLock.Lock()
					errorFiles = append(errorFiles, &chunkInfo)
					errorFilesLock.Unlock()

					bar.Abort(false)
					bar.SetPriority(math.MaxInt - int(c) - len(config.ChunkedFileList))
				} else {
					if !bar.Completed() {
						bar.SetCurrent(100)
					}
					bar.SetPriority(int(c))
					chunkFetchStat.SetDone(chunkInfo.FileName)
					return
				}
			}
		}()
	}
	progressBar.Wait()
	wg.Wait()

	chunkFetchStat.Close()

	if len(errorFiles) > 0 {
		fmt.Println("[ERROR] the following files download failed, please try again:")
		for _, v := range errorFiles {
			fmt.Println(v.FileName)
		}
		return errors.New("download error")
	}

	err = os.Rename(downloadingFilePath, rawFilePath)
	if err != nil {
		fmt.Println("[ERROR] rename download file err:", err.Error())
		return err
	}
	os.Remove(chunkMetaPath)

	fmt.Println("[INFO] download finish")

	return nil
}

func loadFileConfig(jsonConfigAddress string) (*file_config.FileConfig, error) {
	config := &file_config.FileConfig{}
	if strings.HasPrefix(jsonConfigAddress, "http") {
		// download json
		client := &http.Client{Timeout: 20 * time.Second}
		resp, err := client.Get(jsonConfigAddress)
		if err != nil {
			return nil, fmt.Errorf("get json config error: %s", err.Error())
		}
		defer resp.Body.Close()
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("get json config error: %s", err.Error())
		}
		err = json.Unmarshal(content, config)
		if err != nil {
			return nil, fmt.Errorf("json config unmarshal error: %s", err.Error())
		}
	} else {
		// read json file
		content, err := os.ReadFile(jsonConfigAddress)
		if err != nil {
			return nil, fmt.Errorf("read json config error: %s", err.Error())
		}
		err = json.Unmarshal(content, config)
		if err != nil {
			return nil, fmt.Errorf("json config unmarshal error: %s", err.Error())
		}
	}

	return config, nil
}

func parserEndPointFromFileConfig(jsonConfigAddress string) ([]string, error) {
	if strings.HasPrefix(jsonConfigAddress, "http") {
		i := strings.LastIndex(jsonConfigAddress, "/")
		if i < 0 {
			return nil, errors.New("download endpoint error")
		}

		endPoints := []string{}
		endPoints = append(endPoints, jsonConfigAddress[:i])
		return endPoints, nil
	} else {
		return nil, errors.New("download endpoint not exist")
	}
}
