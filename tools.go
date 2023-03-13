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

const TimeLayout = "2006-01-02 15:04:05"

func TimeFormate(origin_time string) string {
	default_str := time.Now().Format(TimeLayout)

	if origin_time == "" {
		return default_str
	}

	{
		var parse_str_time time.Time
		parse_str_time, err := time.Parse(TimeLayout, origin_time)
		if err == nil {
			return parse_str_time.Format(TimeLayout)
		}
	}

	{
		var parse_str_time time.Time
		parse_str_time, err := time.Parse("20060102 15:04:05", origin_time)
		if err == nil {
			return parse_str_time.Format(TimeLayout)
		}
	}

	{
		var parse_str_time time.Time
		parse_str_time, err := time.Parse("2006-01-02T15:04:05", origin_time)
		if err == nil {
			return parse_str_time.Format(TimeLayout)
		}
	}

	{
		var parse_str_time time.Time
		parse_str_time, err := time.Parse(time.RFC3339, origin_time)
		if err == nil {
			return parse_str_time.Format(TimeLayout)
		}
	}

	return default_str
}
func StringToTime(origin_time string) time.Time {
	default_time := time.Now()
	if origin_time == "" {
		return time.Time{}
	}

	{
		var parse_str_time time.Time
		parse_str_time, err := time.ParseInLocation(TimeLayout, origin_time, time.Local)
		if err == nil {
			return parse_str_time
		}
	}

	{
		var parse_str_time time.Time
		parse_str_time, err := time.ParseInLocation("20060102 15:04:05", origin_time, time.Local)
		if err == nil {
			return parse_str_time
		}
	}

	{
		var parse_str_time time.Time
		parse_str_time, err := time.ParseInLocation("2006-01-02T15:04:05", origin_time, time.Local)
		if err == nil {
			return parse_str_time
		}
	}

	{
		var parse_str_time time.Time
		parse_str_time, err := time.ParseInLocation(time.RFC3339, origin_time, time.Local)
		if err == nil {
			return parse_str_time
		}
	}

	return default_time
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

func Ticker(t time.Duration, handler func()) {
	timeTickerChan := time.Tick(t)
	for {
		handler()
		<-timeTickerChan
	}
}

func Timeout(t_c chan time.Duration, t time.Duration, handler func()) {
	select {
	case v := <-t_c:
		go Timeout(t_c, v, handler)
	case <-time.After(t):
		handler()
	}
}

func BuilderString(str ...string) string {
	var builder strings.Builder
	for _, v := range str {
		builder.WriteString(v)
	}
	return builder.String()
}
