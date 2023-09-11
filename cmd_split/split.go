package cmd_split

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/meson-network/bsc-data-file-utils/src/common/parse_size"
	"github.com/meson-network/bsc-data-file-utils/src/file_config"
	"github.com/meson-network/bsc-data-file-utils/src/split_file"
	"github.com/urfave/cli/v2"
)

const default_dest = "./dest"

func Split(clictx *cli.Context) error {

	//read params
	originFilePath := clictx.String("file")
	destDir := clictx.String("dest")
	sizeStr := clictx.String("size")
	thread := clictx.Int("thread")
	if thread == 0 {
		thread = runtime.NumCPU()
	}

	if originFilePath == "" {
		fmt.Println("[ERROR] Invalid file. Please input the path of the source file that you want to split with param '--file=path to file'")
		return errors.New("source file error")
	}

	if destDir == "" {
		fmt.Println("[INFO] No destination folder is entered, the default folder '" + default_dest + "' will be used")
		destDir = default_dest
	}

	if sizeStr == "" {
		fmt.Println("[ERROR] Invalid size. Please input the size of each shard '--size=shard size', ex. --size=100M")
		return errors.New("size error")
	}

	chunkSize, err := parse_size.ParseToByte(sizeStr)
	if err != nil {
		fmt.Println("[ERROR] Invalid size. Please input the size of each shard '--size=shard size', ex. --size=100M")
		return errors.New("size error")
	}

	return splitFile(originFilePath, destDir, chunkSize, thread)
}

func splitFile(originFilePath string, destDir string, chunkSize int64, thread int) error {
	// read origin file
	fileInfo, err := os.Stat(originFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("[ERROR] source file not exist")
			return errors.New("source file not exist")
		}
		return errors.New("source file not exist")
	}

	num := math.Ceil(float64(fileInfo.Size()) / float64(chunkSize))
	if num <= 1 {
		fmt.Println("[ERROR] source file is smaller than chunk size")
		return errors.New("chunk size error")
	}

	// if dest dir exist
	dirInfo, err := os.Stat(destDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if dirInfo != nil {
		// todo warn user

		os.RemoveAll(destDir)
	}

	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	// json file
	splitConfig := file_config.FileConfig{
		RawFile: file_config.RawFileInfo{
			FileName: fileInfo.Name(),
			Size:     fileInfo.Size(),
		},
		EndPoint:        []string{},
		ChunkedFileList: []file_config.ChunkedFileInfo{},
	}
	var jsonLock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(int(num))
	counter := int64(0)

	threadChan := make(chan struct{}, thread)
	errChan := make(chan error)

	fmt.Println("start to split file", "chunk size:", chunkSize, "byte")
	var i int64 = 1
	for ; i <= int64(num); i++ {
		index := i

		select {
		case threadChan <- struct{}{}:
		case err := <-errChan:
			return err
		}

		go func() {
			defer func() {
				<-threadChan
				wg.Done()
			}()
			name, fileSize, md5Str, err := split_file.GenFileChunk(fileInfo.Name(), fileInfo.Size(), originFilePath, index, chunkSize, destDir)
			if err != nil {
				fmt.Println("err:", err)
				errChan <- err
				return
			}
			c := atomic.AddInt64(&counter, 1)
			fmt.Println(c, "/", num, "fileName:", name, "fileSize:", fileSize, "md5:", md5Str)
			jsonLock.Lock()
			defer jsonLock.Unlock()
			splitConfig.ChunkedFileList = append(splitConfig.ChunkedFileList, file_config.ChunkedFileInfo{
				FileName: name,
				Md5:      md5Str,
				Size:     fileSize,
				Offset:   (index - 1) * chunkSize,
			})
		}()

	}

	wg.Wait()

	sort.Slice(splitConfig.ChunkedFileList, func(i, j int) bool {
		return splitConfig.ChunkedFileList[i].Offset < splitConfig.ChunkedFileList[j].Offset
	})

	configFilePath := filepath.Join(destDir, file_config.FILES_CONFIG_JSON_NAME)
	configJson, err := json.Marshal(splitConfig)
	if err != nil {
		return err
	}
	err = os.WriteFile(configFilePath, configJson, os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Println("split finish")

	return nil
}
