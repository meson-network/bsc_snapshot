package download

import (
	"errors"
	"fmt"
	"math/rand"
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
	DEFAULT_RETRY_TIMES     = 5
	DEFAULT_THREAD          = 128
	DEFAULT_REQUEST_TIMEOUT = time.Second * 7
)

func Download(configFilePath string, thread int, retryNum int, noResume bool) error {
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

	// gen raw file
	rawFilePath := filepath.Join("./", conf.RawFile.FileName)
	if fileStat, _ := os.Stat(rawFilePath); fileStat != nil {
		fmt.Println("[ERROR]", "<", conf.RawFile.FileName, ">", "already exist")
		return errors.New("file exist")
	}

	downloadingFileName := conf.RawFile.FileName + ".downloading"
	downloadingFilePath := filepath.Join("./", downloadingFileName)

	chunkMetaName := "." + conf.RawFile.FileName + ".downloaded"
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

			fetcher := NewChunkFetcher(downloadingFilePath, chunkInfo.Size, chunkInfo.Offset, chunkInfo.Md5, pBar)

			// try some times if download failed
			for try := 0; try < retryNum; try++ {
				
				// pick endpoint random
				currentEndpoint := endPoints[rand.Intn(len(endPoints))]
				downloadUrl := currentEndpoint + "/" + chunkInfo.FileName

				downloadSize, err := fetcher.Download(downloadUrl)
				if err != nil {
					if try < retryNum-1 {
						pBar.Add64(-downloadSize)
						time.Sleep(3 * time.Second)
						continue
					}
					pBar.Add64(-downloadSize)
					errorFilesLock.Lock()
					errorFiles = append(errorFiles, &chunkInfo)
					errorFilesLock.Unlock()
				} else {
					chunkFetchStat.SetDone(chunkInfo.FileName)
					return
				}
			}
		}()
	}
	// must wait wg first
	wg.Wait()
	if len(errorFiles) == 0{
		pBar.SetCurrent(conf.RawFile.Size)
	}
	pBar.Finish()

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
