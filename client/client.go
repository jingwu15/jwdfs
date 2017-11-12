package client

import (
	"bytes"
	util "dfs/lib"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	os "os"
	"path"
	"strconv"
	"strings"
)

func Upload(filename string, filekey string, targetUrl string) string {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("filekey", filekey)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("file_upload", filename)
	if err != nil {
		return "101," + err.Error()
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		return "101," + err.Error()
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return "101," + err.Error()
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return strconv.Itoa(resp.StatusCode) + "," + err.Error()
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return strconv.Itoa(resp.StatusCode) + "," + err.Error()
	}
	return strconv.Itoa(resp.StatusCode) + "," + strings.Replace(string(resp_body), "\n", "", -1)
}

func Download(filekey string, targetUrl string) string {
	resp, err := http.Get(targetUrl + "?filekey=" + filekey)
	if err != nil {
		return strconv.Itoa(resp.StatusCode) + "," + err.Error()
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return strconv.Itoa(resp.StatusCode) + "," + err.Error()
	}
	if resp.StatusCode != 200 {
		return strconv.Itoa(resp.StatusCode) + "," + strings.Replace(string(resp_body), "\n", "", -1)
	}
	destfile := viper.Get("client.downfile").(string)
	if destfile == "" {
		destfile = viper.Get("client.downdir").(string) + "/" + filekey
	}
	destdir := path.Dir(destfile)
	if !util.File_exists(destdir) {
		os.MkdirAll(destdir, 0755)
	}
	fh, err := os.OpenFile(destfile, os.O_CREATE|os.O_WRONLY, 0644)
	defer fh.Close()
	if err != nil {
		response := "101," + err.Error()
		return response
	}
	_, err = fh.Write(resp_body)
	if err != nil {
		response := "101," + err.Error()
		return response
	}
	return "200,"
}

func Info(filekey string, targetUrl string) string {
	url := targetUrl + "?filekey=" + filekey

	resp, err := http.Get(url)
	if err != nil {
		return strconv.Itoa(resp.StatusCode) + "," + err.Error()
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return strconv.Itoa(resp.StatusCode) + "," + err.Error()
	}
	return strconv.Itoa(resp.StatusCode) + "," + strings.Replace(string(resp_body), "\n", "", -1)
}
