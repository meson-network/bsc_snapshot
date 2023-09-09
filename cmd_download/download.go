package cmd_download

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/meson-network/bsc-data-file-utils/src/split_file"
	"github.com/urfave/cli/v2"
)

func Download(clictx *cli.Context) error {

	jsonConfig := clictx.String("json")
	threadCount := clictx.Int("thread")

	if jsonConfig == "" {
		fmt.Println("[ERROR] json config error")
		return errors.New("json config error")
	}

	if threadCount == 0 {
		threadCount = runtime.NumCPU()
	}

	// download or read jsonConfig
	config:=split_file.FileSplitConfig{}
	if strings.HasPrefix(jsonConfig, "http"){
		// download json
	}else{
		// read json file
		content,err:=os.ReadFile(jsonConfig)
		if err!=nil{
			return err
		}
		json.Unmarshal(content,&config)
	}

	return nil
}
