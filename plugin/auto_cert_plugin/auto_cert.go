package auto_cert_plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"net/http"
	"os"
	"time"

	"github.com/coreservice-io/job"
	"github.com/meson-network/bsc-data-file-utils/basic"
)

var instanceMap = map[string]*Cert{}

type Cert struct {
	Download_url        string
	Local_crt_path      string
	Local_key_path      string
	Auto_updating       bool
	Check_interval_secs int
}

type Config struct {
	Download_url        string
	Local_crt_path      string
	Local_key_path      string
	Check_interval_secs int
}

func GetInstance() *Cert {
	return GetInstance_("default")
}

func GetInstance_(name string) *Cert {
	cert := instanceMap[name]
	if cert == nil {
		basic.Logger.Errorln(name + " auto_cert plugin null")
	}
	return cert
}

// update_change_callback func(new_crt_content, new_key_content)
func (cert *Cert) AutoUpdate(update_change_callback func(string, string)) {
	if !cert.Auto_updating {
		cert.Auto_updating = true

		job.Start(context.Background(), job.JobConfig{
			Name:          "cert_auto_update_job",
			Job_type:      job.TYPE_PANIC_REDO,
			Interval_secs: int64(cert.Check_interval_secs),
			Process_fn: func(j *job.Job) {
				cert.update_(update_change_callback)
			},
			On_panic: func(job *job.Job, panic_err interface{}) {
				basic.Logger.Errorln(panic_err)
			},
		}, nil)

	}
}

type RemoteRespCert struct {
	Crt_content string `json:"crt_content"`
	Key_content string `json:"key_content"`
}

type RemoteResp struct {
	Cert         RemoteRespCert `json:"cert"`
	Meta_status  int64          `json:"meta_status"`
	Meta_message string         `json:"meta_message"`
}

// update from remote url
func (cert *Cert) update_(Update_change_callback func(string, string)) error {

	downloadClient := http.Client{
		Timeout: time.Second * 60, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, cert.Download_url, nil)
	if err != nil {
		return err
	}

	res, getErr := downloadClient.Do(req)
	if getErr != nil {
		return getErr
	}

	if res.Body == nil {
		return errors.New("no body response")
	}

	defer res.Body.Close()

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return readErr
	}

	rp := RemoteResp{}
	jsonErr := json.Unmarshal(body, &rp)
	if jsonErr != nil {
		return jsonErr
	}

	if rp.Meta_status <= 0 {
		return errors.New(rp.Meta_message)
	}

	// /////////////////////////
	change := false
	// read old .crt
	old_crt_content, read_err := os.ReadFile(cert.Local_crt_path)
	if read_err != nil {
		change = true
	} else {
		if string(old_crt_content) != rp.Cert.Crt_content {
			change = true
		}
	}

	// read old .key
	old_key_content, read_err := os.ReadFile(cert.Local_key_path)
	if read_err != nil {
		change = true
	} else {
		if string(old_key_content) != rp.Cert.Key_content {
			change = true
		}
	}

	// ////save .crt and .key/////
	if change {
		crt_file_err := file_overwrite(cert.Local_crt_path, rp.Cert.Crt_content)
		if crt_file_err != nil {
			return crt_file_err
		}

		key_file_err := file_overwrite(cert.Local_key_path, rp.Cert.Key_content)
		if key_file_err != nil {
			return key_file_err
		}

		if Update_change_callback != nil {
			Update_change_callback(rp.Cert.Crt_content, rp.Cert.Key_content)
		}
	}

	return nil
}

func file_overwrite(path string, content string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, werr := f.WriteString(content)
	if werr != nil {
		return werr
	}
	return nil
}

// init_download ==true ,it will update from remote url in init process
func Init(conf *Config, init_download bool) error {
	return Init_("default", conf, init_download)
}

func Init_(name string, conf *Config, init_download bool) error {
	if name == "" {
		name = "default"
	}

	_, exist := instanceMap[name]
	if exist {
		return fmt.Errorf("cert instance <%s> has already been initialized", name)
	}

	if conf.Download_url == "" || conf.Local_crt_path == "" || conf.Local_key_path == "" {
		return errors.New("params Download_url|Local_crt_path|Local_key_path must all be set ")
	}

	cert := &Cert{
		conf.Download_url,
		conf.Local_crt_path,
		conf.Local_key_path,
		false,
		conf.Check_interval_secs,
	}

	if init_download {
		first_update_err := cert.update_(nil)
		if first_update_err != nil {
			// will try again after 5 second
			time.Sleep(5 * time.Second)
			second_update_err := cert.update_(nil)
			if second_update_err != nil {
				return errors.New("cert init failed," + first_update_err.Error())
			}
		}
	}

	instanceMap[name] = cert
	return nil
}
