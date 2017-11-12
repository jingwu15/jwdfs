package client

import (
	"bytes"
	util "dfs/lib"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	os "os"
	"path"
)

var logger *log.Logger

func Upload(filename string, filekey string, targetUrl string) string {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("filekey", filekey)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("file_upload", filename)
	if err != nil {
		response := "101," + err.Error()
		logger.Println(response)
		return response
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		response := "101," + err.Error()
		logger.Println(response)
		return response
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		response := "101," + err.Error()
		logger.Println(response)
		return response
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		response := "101," + err.Error()
		logger.Println(response)
		return response
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response := "101," + err.Error()
		logger.Println(response)
		return response
	}
	return string(resp_body)
}

func Download(filekey string, targetUrl string) string {
	resp, err := http.Get(targetUrl + "?filekey=" + filekey)
	if err != nil {
		response := "101," + err.Error()
		logger.Println(response)
		return response
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response := "101," + err.Error()
		logger.Println(response)
		return response
	}
	if resp.StatusCode == 701 {
		response := "701," + string(resp_body)
		logger.Println(response)
		return response
	}
	destfile := viper.Get("client.filename").(string)
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
		logger.Println(response)
		return response
	}
	fmt.Println(len(resp_body))
	_, err = fh.Write(resp_body)
	if err != nil {
		response := "101," + err.Error()
		logger.Println(response)
		return response
	}
	return ""
}

func Info(filekey string, targetUrl string) string {
	url := targetUrl + "?filekey=" + filekey

	resp, err := http.Get(url)
	if err != nil {
		response := "101," + err.Error()
		return response
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		response := "101," + err.Error()
		return response
	}
	return string(resp_body)
}
