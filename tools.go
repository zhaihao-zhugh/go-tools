package tools

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const time_layout = "2006-01-02 15:04:05"

func TimeFormate(origin_time string) string {
	// const time_layout = "2006-01-02 15:04:05"
	// default_time := time.Date(2006, 1, 2, 15, 04, 05, 0, time.Local)
	default_str := time.Now().Format(time_layout)

	if origin_time == "" {
		return default_str
	}

	{
		var parse_str_time time.Time
		parse_str_time, err := time.Parse(time_layout, origin_time)
		if err == nil {
			return parse_str_time.Format(time_layout)
		}
	}

	{
		var parse_str_time time.Time
		parse_str_time, err := time.Parse("20060102 15:04:05", origin_time)
		if err == nil {
			return parse_str_time.Format(time_layout)
		}
	}

	{
		var parse_str_time time.Time
		parse_str_time, err := time.Parse("2006-01-02T15:04:05", origin_time)
		if err == nil {
			return parse_str_time.Format(time_layout)
		}
	}

	{
		var parse_str_time time.Time
		parse_str_time, err := time.Parse(time.RFC3339, origin_time)
		if err == nil {
			return parse_str_time.Format(time_layout)
		}
	}

	return default_str
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func SumCheck(data []byte) (sum byte) {
	for _, v := range data {
		sum += v
	}
	return sum
}

func ByteSliceToHexString(src []byte) (dest string) {
	for _, v := range src {
		dest += fmt.Sprintf("%02X", v)
	}
	return
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

func IsFileExit(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func Ticker(t time.Duration, handler func()) {
	timeTickerChan := time.Tick(t)
	for {
		handler()
		<-timeTickerChan
	}
}

func SplitString(s string, sep []rune) []string {
	Split := func(r rune) bool {
		for _, v := range sep {
			if v == r {
				return true
			}
		}
		return false
	}
	return strings.FieldsFunc(s, Split)
}
