package cmd_split

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/meson-network/bsc-data-file-utils/src/split_file"
	"github.com/urfave/cli/v2"
)

func Split(clictx *cli.Context) error {

	//read params
	originFile := clictx.String("file")
	destDir := clictx.String("dest")
	sizeStr := clictx.String("size")
	thread := clictx.Int("thread")
	if thread == 0 {
		thread = runtime.NumCPU()
	}

	if originFile == "" {
		fmt.Println("[ERROR] Invalid file. Please input the path of the source file that you want to split with param '--file=path to file'")
		return errors.New("source file error")
	}

	if destDir == "" {
		fmt.Println("[INFO] No destination folder is entered, the default folder './dest' will be used")
		destDir = "./dest"
	}

	if sizeStr == "" {
		fmt.Println("[ERROR] Invalid size. Please input the size of each shard '--size=shard size', ex. --size=100M")
		return errors.New("size error")
	}

	chunkSize, err := sizeStrToNum(sizeStr)
	if err != nil {
		fmt.Println("[ERROR] Invalid size. Please input the size of each shard '--size=shard size', ex. --size=100M")
		return errors.New("size error")
	}

	return splitFile(originFile, destDir, chunkSize, thread)
}

func splitFile(originFile string, destDir string, chunkSize int64, thread int) error {
	// read origin file
	fileInfo, err := os.Stat(originFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("[ERROR] source file not exist")
			return errors.New("source file not exist")
		}
		return errors.New("source file not exist")
	}
	// fileName := fileInfo.Name()

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
	splitConfig := split_file.FileSplitConfig{
		RawFile: split_file.RawFileInfo{
			FileName: fileInfo.Name(),
			Size:     fileInfo.Size(),
		},
		EndPoint:        []string{},
		ChunkedFileList: []split_file.ChunkedFileInfo{},
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
			name, fileSize, md5Str, err := genSplitFile(fileInfo, originFile, index, chunkSize, destDir)
			if err != nil {
				fmt.Println("err:", err)
				errChan <- err
				return
			}
			c := atomic.AddInt64(&counter, 1)
			fmt.Println(c, "/", num, "fileName:",name, "fileSize:", fileSize, "md5:", md5Str)
			jsonLock.Lock()
			defer jsonLock.Unlock()
			splitConfig.ChunkedFileList = append(splitConfig.ChunkedFileList, split_file.ChunkedFileInfo{
				FileName: name,
				Md5:      md5Str,
				Size:     fileSize,
				Offset:   (index - 1) * chunkSize,
			})
		}()

	}

	wg.Wait()

	sort.Slice(splitConfig.ChunkedFileList,func(i, j int) bool {
		return splitConfig.ChunkedFileList[i].Offset<splitConfig.ChunkedFileList[j].Offset
	})

	configFilePath := filepath.Join(destDir, split_file.FILES_CONFIG_JSON_NAME)
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

func genSplitFile(fileInfo fs.FileInfo, originFilePath string, i int64, chunkSize int64, destDir string) (name string, size int64, md5Str string, err error) {
	fi, err := os.OpenFile(originFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return "", 0, "", err
	}
	defer fi.Close()

	fileName := fileInfo.Name()
	partSize := chunkSize
	_, err = fi.Seek((i-1)*chunkSize, 0)
	if err != nil {
		return "", 0, "", err
	}
	if partSize > int64(fileInfo.Size()-(i-1)*chunkSize) {
		partSize = fileInfo.Size() - (i-1)*chunkSize
	}

	name = fmt.Sprintf("%s.%d", fileName, i)
	dfName := fmt.Sprintf("./%s.%d", fileName, i)
	dfName = filepath.Join(destDir, dfName)
	
	df, err := os.OpenFile(dfName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
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

	return name, partSize, md5Str, nil
}

func sizeStrToNum(sizeStr string) (int64, error) {
	sizeStr = strings.TrimSpace(sizeStr)
	sizeStr = strings.ToLower(sizeStr)
	sizeStr = strings.TrimSuffix(sizeStr, "b")

	var err error
	rate := int64(1)
	size := int64(0)
	if strings.HasSuffix(sizeStr, "g") {
		size, err = strconv.ParseInt(sizeStr[:len(sizeStr)-1], 10, 64)
		if err != nil {
			return 0, err
		}
		rate = 1024 * 1024 * 1024
	} else if strings.HasSuffix(sizeStr, "m") {
		size, err = strconv.ParseInt(sizeStr[:len(sizeStr)-1], 10, 64)
		if err != nil {
			return 0, err
		}
		rate = 1024 * 1024
	} else if strings.HasSuffix(sizeStr, "k") {
		size, err = strconv.ParseInt(sizeStr[:len(sizeStr)-1], 10, 64)
		if err != nil {
			return 0, err
		}
		rate = 1024
	} else {
		size, err = strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			return 0, err
		}
	}

	if size <= 0 {
		return 0, errors.New("size error")
	}

	return size * rate, nil
}
