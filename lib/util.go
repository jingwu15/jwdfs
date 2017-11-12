package lib

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

func FormatPathSuffix(pathname string) string {
	length := len(pathname)
	suffix := pathname[length-1 : length]
	if suffix != "/" {
		pathname = pathname + "/"
	}
	return pathname
}

func Byte2string(in [16]byte) []byte {
	tmp := make([]byte, 16)
	for _, value := range in {
		tmp = append(tmp, value)
	}

	return tmp[16:]
}

func Md5_sum(raw []byte) string {
	md5Sum := md5.Sum(raw)
	return hex.EncodeToString(Byte2string(md5Sum))
}

func Json_encode(jsonMap map[string]string) string {
	jsonData, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return string(jsonData)
}

func Json_decode_file(jsonfile string) (map[string]string, error) {
	var jsonMap map[string]string
	bytes, err := ioutil.ReadFile(jsonfile)
	if err != nil {
		return jsonMap, err
	}
	err = json.Unmarshal(bytes, &jsonMap)
	if err != nil {
		return jsonMap, err
	}
	return jsonMap, nil
}

func File_exists(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func IsEmpty(a interface{}) bool {
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}
