package endpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/meson-network/bsc_snapshot/src/model"
	"github.com/meson-network/bsc_snapshot/src/utils/file_config"
)

func AddEndpoint(configPath string, endpoints []string) error {
	config, err := readFileConfig(configPath)
	if err != nil {
		fmt.Println("[ERROR] read config error:", err.Error())
		return err
	}
	if len(endpoints) == 0 {
		fmt.Println("[ERROR] please input endpoint with param -e=<endpoint address>")
		return errors.New("endpoints not find")
	}

	f_ep, err := file_config.FormatEndpoints(endpoints)
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
		return err
	}

	addEndpoint(config, f_ep)
	err = saveFileConfig(config, configPath)
	if err != nil {
		fmt.Println("[ERROR] save config error:", err.Error())
	} else {
		fmt.Println("[INFO] config saved")
	}
	return err
}

func RemoveEndpoint(configPath string, endpoints []string) error {
	config, err := readFileConfig(configPath)
	if err != nil {
		fmt.Println("[ERROR] read config error:", err.Error())
		return err
	}
	if len(endpoints) == 0 {
		fmt.Println("[ERROR] please input endpoint with param -e=<endpoint address>")
		return errors.New("endpoints not find")
	}

	f_ep, err := file_config.FormatEndpoints(endpoints)
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
		return err
	}

	removeEndpoint(config, f_ep)

	err = saveFileConfig(config, configPath)
	if err != nil {
		fmt.Println("[ERROR] save config error:", err.Error())
	} else {
		fmt.Println("[INFO] config saved")
	}
	return err
}

func SetEndpoint(configPath string, endpoints []string) error {
	config, err := readFileConfig(configPath)
	if err != nil {
		fmt.Println("[ERROR] read config error:", err.Error())
		return err
	}

	f_ep, err := file_config.FormatEndpoints(endpoints)
	if err != nil {
		fmt.Println("[ERROR]", err.Error())
		return err
	}
	if len(endpoints) == 0 {
		fmt.Println("[ERROR] please input endpoint with param -e=<endpoint address>")
		return errors.New("endpoints not find")
	}

	setEndpoint(config, f_ep)

	err = saveFileConfig(config, configPath)
	if err != nil {
		fmt.Println("[ERROR] save config error:", err.Error())
	} else {
		fmt.Println("[INFO] config saved")
	}
	return err
}

func ClearEndpoint(configPath string) error {
	config, err := readFileConfig(configPath)
	if err != nil {
		fmt.Println("[ERROR] read config error:", err.Error())
		return err
	}

	clearEndpoint(config)
	err = saveFileConfig(config, configPath)
	if err != nil {
		fmt.Println("[ERROR] save config error:", err.Error())
	} else {
		fmt.Println("[INFO] config saved")
	}
	return err
}

func PrintEndpoint(configPath string) error {
	config, err := readFileConfig(configPath)
	if err != nil {
		fmt.Println("[ERROR] read config error:", err.Error())
		return err
	}

	fmt.Println("current endpoint:")
	for _, v := range config.EndPoint {
		fmt.Println("\t", v)
	}

	return nil
}

func readFileConfig(filePath string) (*model.FileConfig, error) {

	jsonContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	fileConfig := model.FileConfig{}
	err = json.Unmarshal(jsonContent, &fileConfig)
	if err != nil {
		return nil, err
	}

	return &fileConfig, nil
}

func saveFileConfig(c *model.FileConfig, savePath string) error {

	configJson, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = os.WriteFile(savePath, configJson, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func addEndpoint(c *model.FileConfig, endpoints []string) {

	existEndpointMap := map[string]struct{}{}
	for _, v := range c.EndPoint {
		existEndpointMap[v] = struct{}{}
	}
	for _, v := range endpoints {
		_, exist := existEndpointMap[v]
		if !exist {
			c.EndPoint = append(c.EndPoint, v)
		}
	}
}

func removeEndpoint(c *model.FileConfig, endpoints []string) {
	toRemoveEndpointMap := map[string]struct{}{}
	for _, v := range endpoints {
		toRemoveEndpointMap[v] = struct{}{}
	}

	newEndpoints := []string{}
	for _, v := range c.EndPoint {
		_, exist := toRemoveEndpointMap[v]
		if !exist {
			newEndpoints = append(newEndpoints, v)
		}
	}

	c.EndPoint = newEndpoints
}

func setEndpoint(c *model.FileConfig, endpoints []string) {
	c.EndPoint = endpoints
}

func clearEndpoint(c *model.FileConfig) {
	c.EndPoint = []string{}
}
