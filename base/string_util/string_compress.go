package string_util

import (
	"bytes"
	"compress/flate"
	"io/ioutil"
)

const (
	kCompressLevel = flate.BestCompression
)

func Compress(s string) (string, error) {
	var buf bytes.Buffer
	bc, err1 := flate.NewWriter(&buf, kCompressLevel)
	if err1 != nil {
		return s, err1
	}
	_, err2 := bc.Write([]byte(s))
	if err2 != nil {
		return s, err2
	}
	err3 := bc.Flush()
	if err3 != nil {
		return s, err3
	}
	return buf.String(), nil
}
func Uncompress(s string) (string, error) {
	buf := bytes.NewBuffer([]byte(s))
	bufr := flate.NewReader(buf)
	str, _ := ioutil.ReadAll(bufr)
	return string(str), nil
}
