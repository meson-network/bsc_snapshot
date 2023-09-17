package download

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/meson-network/bsc_snapshot/src/config"
	"github.com/meson-network/bsc_snapshot/src/model"
	"github.com/meson-network/bsc_snapshot/src/utils/file_config"

	"github.com/cheggaaa/pb/v3"
)

const (
	DEFAULT_RETRY_TIMES     = 8
	RETRY_WAIT_SECS         = 5
	DEFAULT_THREAD          = 64
	DEFAULT_REQUEST_TIMEOUT = time.Second * 7
)

func Download(configFilePath string, destDir string, thread int, retryNum int, noResume bool) error {
	if configFilePath == "" {
		fmt.Println("[ERROR] json config error, please input correct address or file path")
		return errors.New("json config error")
	}
	if thread <= 0 {
		thread = DEFAULT_THREAD
	}
	if retryNum <= 0 {
		retryNum = DEFAULT_RETRY_TIMES
	}

	conf, err := config.LoadFile4Download(configFilePath)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return err
	}

	fmt.Printf("[INFO] File: < %s > \n", conf.RawFile.FileName)
	fmt.Println("[INFO] download thread number:", thread)

	existEndPoints, err := file_config.FormatEndpoints(conf.EndPoint)
	if err != nil {
		fmt.Println("[ERROR] json config endpoint format error")
		return err
	}

	// check endpoint
	var endPoints []string
	if len(existEndPoints) == 0 {
		// if no end point info in json, default use json config download path
		endPointsFromFile, err := config.ExtractEndPointFromConfig(configFilePath)
		if err != nil {
			fmt.Println("[ERROR] ", err.Error())
			return err
		}
		// fmt.Println("[INFO] use some endpoint with json config file")

		endPoints = endPointsFromFile
	} else {
		endPoints = existEndPoints
	}

	if destDir == "" {
		destDir = "./"
	}

	// gen raw file
	rawFilePath := filepath.Join(destDir, conf.RawFile.FileName)
	rawFileDir := filepath.Dir(rawFilePath)
	err = os.MkdirAll(rawFileDir, os.ModePerm)
	if err != nil {
		fmt.Println("[ERROR] ", err.Error())
		return err
	}

	if fileStat, _ := os.Stat(rawFilePath); fileStat != nil {
		fmt.Println("[ERROR]", "<", conf.RawFile.FileName, ">", "already exist")
		return errors.New("file exist")
	}

	downloadingFileName := conf.RawFile.FileName + ".downloading"
	downloadingFilePath := filepath.Join(destDir, downloadingFileName)

	chunkMetaName := "." + conf.RawFile.FileName + ".downloaded"
	chunkMetaPath := filepath.Join(destDir, chunkMetaName)

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

	if err := dFile.Truncate(conf.RawFile.Size); err != nil {
		dFile.Close()
		fmt.Println("[ERROR] handle file err:", err.Error())
		return err
	}
	dFile.Close()

	fmt.Println("[INFO] start download...")

	threadChan := make(chan struct{}, thread)
	errorFiles := []*model.ChunkedFileInfo{}
	var errorFilesLock sync.Mutex

	var wg sync.WaitGroup

	pBar := pb.New64(conf.RawFile.Size)
	pBar.SetRefreshRate(time.Second)
	pBar.Set(pb.Bytes, true)
	pBar.Set(pb.SIBytesPrefix, true)
	pBar.Start()

	// counter := int64(0)
	for _, v := range conf.ChunkedFileList {
		chunkInfo := v
		if chunkFetchStat.IsDone(chunkInfo.FileName) {
			// if already downloaded, skip it
			pBar.Add64(chunkInfo.Size)
			continue
		}

		threadChan <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				<-threadChan
				wg.Done()
			}()

			fetcher := NewChunkFetcher(downloadingFilePath, &chunkInfo, pBar)

			if isDone := fetcher.Download(endPoints, retryNum, func(*model.ChunkedFileInfo) {
				errorFilesLock.Lock()
				errorFiles = append(errorFiles, &chunkInfo)
				errorFilesLock.Unlock()
			}); isDone {
				chunkFetchStat.SetDone(chunkInfo.FileName)
			}
		}()
	}
	// must wait wg first
	wg.Wait()
	if len(errorFiles) == 0 {
		pBar.SetCurrent(conf.RawFile.Size)
	}
	pBar.Finish()

	chunkFetchStat.Close()

	if len(errorFiles) > 0 {
		fmt.Println("[ERROR] the following files download failed, please try again.")
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
