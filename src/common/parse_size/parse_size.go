package parse_size

import (
	"errors"
	"strconv"
	"strings"
)

func ParseToByte(sizeStr string) (int64, error) {
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
