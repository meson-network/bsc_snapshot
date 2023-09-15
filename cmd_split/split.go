package cmd_split

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/meson-network/bsc-data-file-utils/src/common/parse_size"
	"github.com/meson-network/bsc-data-file-utils/src/file_config"
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
		fmt.Println("[ERROR] read source file err:", err.Error())
		return errors.New("source file error")
	}

	num := math.Ceil(float64(fileInfo.Size()) / float64(chunkSize))
	if num <= 1 {
		fmt.Println("[ERROR] source file is smaller than chunk size")
		return errors.New("chunk size error")
	}

	// if dest dir exist
	dirInfo, err := os.Stat(destDir)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("[ERROR] dest dir err:", err.Error())
		return err
	}
	if dirInfo != nil {
		fmt.Print("[WARN] dest dir already exist, all content will be overwrite (yes or no ?):")
		var input string
		fmt.Scanln(&input)
		input = strings.TrimSpace(input)
		input = strings.ToLower(input)
		if input != "y" && input != "yes" {
			return nil
		}
		os.RemoveAll(destDir)
	}

	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		fmt.Println("[ERROR] build dest dir err:", err.Error())
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

	fmt.Println("[INFO] start to split file", "chunk size:", chunkSize, "byte")

	for i := int64(1); i <= int64(num); i++ {
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
			name, fileSize, md5Str, err := genFileChunk(fileInfo.Name(), fileInfo.Size(), originFilePath, index, chunkSize, destDir)
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
		fmt.Println("[ERROR] marshal file config err:", err.Error())
		return err
	}
	err = os.WriteFile(configFilePath, configJson, os.ModePerm)
	if err != nil {
		fmt.Println("[ERROR] save file config err:", err.Error())
		return err
	}

	fmt.Println("[INFO] split finish")

	return nil
}

func genFileChunk(originFileName string, originFileSize int64, originFilePath string, chunkNum int64, chunkSize int64, destDir string) (name string, size int64, md5Str string, err error) {
	fi, err := os.OpenFile(originFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", 0, "", err
	}
	defer fi.Close()

	_, err = fi.Seek((chunkNum-1)*chunkSize, 0)
	if err != nil {
		return "", 0, "", err
	}
	partSize := chunkSize
	if partSize > int64(originFileSize-(chunkNum-1)*chunkSize) {
		partSize = originFileSize - (chunkNum-1)*chunkSize
	}

	chunkFileName := fmt.Sprintf("%s.%d", originFileName, chunkNum)
	chunkFilePath := filepath.Join(destDir, chunkFileName)

	df, err := os.OpenFile(chunkFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return "", 0, "", err
	}
	defer df.Close()

	h := md5.New()
	target := io.MultiWriter(df, h)

	_, err = io.CopyN(target, fi, partSize)
	if err != nil {
		return "", 0, "", err
	}
	md5Str = hex.EncodeToString(h.Sum(nil))

	return chunkFileName, partSize, md5Str, nil
}
