package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func Post(url string, param interface{}) (res *http.Response, err error) {
	cli := &http.Client{Timeout: 1 * time.Minute}
	data, err := json.Marshal(param)
	if err != nil {
		return res, err
	}
	fmt.Println("post data:", string(data))
	res, err = cli.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return res, err
	}
	return res, err
}

func Get(url string) (res *http.Response, err error) {
	cli := &http.Client{Timeout: 1 * time.Minute}
	res, err = cli.Get(url)
	if err != nil {
		return res, err
	}
	return res, err
}
