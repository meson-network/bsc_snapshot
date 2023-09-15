package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/meson-network/bsc_snapshot/src/model"
)

func LoadFile4Upload(configFilePath string) (*model.FileConfig, error) {
	jsonContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("read file config err: %s", err.Error())
	}
	fileConfig := &model.FileConfig{}
	err = json.Unmarshal(jsonContent, fileConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshal file config err: %s", err.Error())
	}

	return fileConfig, nil
}

func LoadFile4Download(configFilePath string) (*model.FileConfig, error) {
	config := &model.FileConfig{}
	if strings.HasPrefix(configFilePath, "http") {
		// download json
		client := &http.Client{Timeout: 20 * time.Second}
		resp, err := client.Get(configFilePath)
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
		content, err := os.ReadFile(configFilePath)
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

func ExtractEndPointFromConfig(jsonConfigAddress string) ([]string, error) {
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
