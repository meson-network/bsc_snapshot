package api

import (
	"errors"
	"time"

	"github.com/imroc/req"
)

const post = "POST"
const get = "GET"
const API_TIMEOUT_SECS = 10

func Get(url string, token string, result interface{}) error {
	return request(get, url, token, nil, API_TIMEOUT_SECS, result)
}

func Get_(url string, token string, timeOutSec int, result interface{}) error {
	return request(get, url, token, nil, timeOutSec, result)
}

func POST(url string, token string, postData interface{}, result interface{}) error {
	return request(post, url, token, postData, API_TIMEOUT_SECS, result)
}

func POST_(url string, token string, postData interface{}, timeOutSec int, result interface{}) error {
	return request(post, url, token, postData, timeOutSec, result)
}

func request(method string, url string, token string, postData interface{}, timeOutSec int, result interface{}) error {
	//optional
	// if result != nil {
	// 	t := reflect.TypeOf(result).Kind()
	// 	if t != reflect.Ptr && t != reflect.Slice && t != reflect.Map {
	// 		return errors.New("value only support Pointer Slice and Map")
	// 	}
	// }

	authHeader := req.Header{
		"Accept": "application/json",
	}

	if token != "" {
		authHeader["Authorization"] = "Bearer " + token
	}

	req.SetTimeout(time.Duration(timeOutSec) * time.Second)

	var resp *req.Resp
	var err error

	switch method {
	case get:
		resp, err = req.Get(url, authHeader)
	case post:
		resp, err = req.Post(url, authHeader, req.BodyJSON(postData))
	default:
		// imposssible
	}

	if err != nil {
		return err
	}

	if resp.Response().StatusCode >= 200 && resp.Response().StatusCode < 300 {
	} else {
		return errors.New("network error")
	}

	err = resp.ToJSON(result)
	if err != nil {
		return errors.New("json error")
	}

	return nil
}
