package split_file

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func GenFileChunk(originFileName string, originFileSize int64, originFilePath string, chunkNum int64, chunkSize int64, destDir string) (name string, size int64, md5Str string, err error) {
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
