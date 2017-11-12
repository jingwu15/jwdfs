package server

import (
	dfsutil "dfs/lib"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	os "os"
	"path"
	"strconv"
	strings "strings"
)

var configMap map[string]string
var logger *log.Logger

func Upload(w http.ResponseWriter, r *http.Request) {
	filekey := string(r.FormValue("filekey"))
	var filepath string
	if strings.HasPrefix(filekey, "/") {
		filekey = strings.TrimLeftFunc(filekey, func(r rune) bool { return r == '/' })
	} else {
	}
	filepath = viper.Get("server.updir").(string) + filekey
	_, err := os.Open(filepath)
	if os.IsNotExist(err) {
		os.MkdirAll(path.Dir(filepath), 0755)
	}

	r.ParseMultipartForm(32 << 20)
	srcFile, _, err := r.FormFile("file_upload")
	if err != nil {
		logger.Println("201," + err.Error())
		fmt.Fprintln(w, "201,"+err.Error())
		return
	}
	defer srcFile.Close()

	destFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Println("201," + err.Error())
		fmt.Fprintln(w, "201,"+err.Error())
		return
	}
	defer destFile.Close()
	io.Copy(destFile, srcFile)
	fmt.Fprintln(w, "000,"+filekey)
	return
}

//文件下载
func Download(w http.ResponseWriter, req *http.Request) {
	queryForm, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		w.WriteHeader(701)
		logger.Println("201," + err.Error())
		fmt.Fprintln(w, "201,"+err.Error())
		return
	}
	if _, ok := queryForm["filekey"]; !ok {
		w.WriteHeader(701)
		logger.Println("201, filekey is requirement!")
		fmt.Fprintln(w, "201, filekey is requirement!")
		return
	}
	filekey := queryForm["filekey"][0]
	filepath := viper.Get("server.updir").(string) + filekey
	if !dfsutil.File_exists(filepath) {
		w.WriteHeader(701)
		fmt.Fprintln(w, "201,file is not exists,"+filekey)
	} else {
		w.Header().Set("Content-type", "image/jpeg")
		w.Header().Set("Content-Disposition", "attachment; filename="+path.Base(filepath))
		w.Header().Set("Filekey", filekey)
		http.ServeFile(w, req, filepath)
	}
	return
}

//文件信息
func Info(w http.ResponseWriter, req *http.Request) {
	queryForm, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		logger.Println("201," + err.Error())
		fmt.Fprintln(w, "201,"+err.Error())
		return
	}
	filekey := req.FormValue("filekey")
	filepath := viper.Get("server.updir").(string) + filekey
	file, _ := os.Open(filepath)
	body, _ := ioutil.ReadAll(file)
	md5Hash := dfsutil.Md5_sum(body)
	response := map[string]string{
		"filekey":  queryForm["filekey"][0],
		"hash":     md5Hash,
		"filesize": strconv.Itoa(len(body)),
	}
	fmt.Fprintln(w, dfsutil.Json_encode(response))
	return
}

func Start() {
	logger = dfsutil.GetLogger("/tmp/dfs_server.log")

	http.HandleFunc("/upload", Upload)
	http.HandleFunc("/download", Download)
	http.HandleFunc("/info", Info)
	address := viper.Get("server.host").(string) + ":" + viper.Get("server.port").(string)
	fmt.Println(address)
	http.ListenAndServe(address, nil)
}
