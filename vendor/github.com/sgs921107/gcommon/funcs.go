/*************************************************************************
	> File Name: funcs.go
	> Author: xiangcai
	> Mail: xiangcai@gmail.com
	> Created Time: 2020年12月09日 星期三 19时39分01秒
 ************************************************************************/

package gcommon

import (
	"bytes"
	"io"
	"net/url"
	"strings"
	"time"
)

// TimeStampFlag 定义类型为int
type TimeStampFlag int

const (
	// SECOND 0: 秒
	SECOND TimeStampFlag = iota
	// MILLISECOND 1: 毫秒
	MILLISECOND
	// MICROSECOND 2: 微秒
	MICROSECOND
	// NANOSECOND 3: 纳秒
	NANOSECOND
)

/*
MapToBytes convert type map to []byte
*/
func MapToBytes(data map[string]string) []byte {
	reader := MapToReader(data)
	return ReaderToBytes(reader)
}

/*
MapToReader convert type map to io.reader
*/
func MapToReader(data map[string]string) io.Reader {
	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	return strings.NewReader(form.Encode())
}

/*
ReaderToBytes convert type io.reader to []byte
*/
func ReaderToBytes(reader io.Reader) []byte {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(reader)
	return buffer.Bytes()
}

/*
ReaderToString convert type io.reader to string
*/
func ReaderToString(reqBody io.Reader) string {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(reqBody)
	return buffer.String()
}

/*
DurationToIntSecond 将time.Duration类型转换为值为秒数的int类型
*/
func DurationToIntSecond(duration time.Duration) int {
	return int(duration) / 1e9
}

/*
TimeStamp 获取时间戳
接收一个整形 0-3
0-秒, 1-毫秒, 2-微妙, 3-纳秒
*/
func TimeStamp(flag int) int64 {
	now := time.Now()
	switch TimeStampFlag(flag) {
	case SECOND:
		return now.Unix()
	case MILLISECOND:
		return now.UnixNano() / 1e6
	case MICROSECOND:
		return now.UnixNano() / 1e3
	case NANOSECOND:
		return now.UnixNano()
	default:
		return 0
	}
}

/*
FetchURLHost 提取url的host
*/
func FetchURLHost(URL string) (string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	return u.Host, nil
}
