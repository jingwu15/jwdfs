package server

import (
	dfsutil "jwdfs/lib"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	os "os"
	"path"
	"strconv"
	strings "strings"
)

var configMap map[string]string

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
		w.WriteHeader(701)
		fmt.Fprintln(w, err.Error())
		return
	}
	defer srcFile.Close()

	destFile, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		w.WriteHeader(701)
		fmt.Fprintln(w, err.Error())
		return
	}
	defer destFile.Close()
	io.Copy(destFile, srcFile)
	fmt.Fprintln(w, filekey)
	return
}

//文件下载
func Download(w http.ResponseWriter, req *http.Request) {
	queryForm, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		w.WriteHeader(701)
		fmt.Fprintln(w, err.Error())
		return
	}
	if _, ok := queryForm["filekey"]; !ok {
		w.WriteHeader(701)
		fmt.Fprintln(w, "filekey is requirement!")
		return
	}
	filekey := queryForm["filekey"][0]
	filepath := viper.Get("server.updir").(string) + filekey
	if !dfsutil.File_exists(filepath) {
		w.WriteHeader(701)
		fmt.Fprintln(w, "file is not exists,"+filekey)
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
		w.WriteHeader(701)
		fmt.Fprintln(w, err.Error())
		return
	}
	filekey := req.FormValue("filekey")
	filepath := viper.Get("server.updir").(string) + filekey
	file, err := os.Open(filepath)
	if err == nil {
		body, _ := ioutil.ReadAll(file)
		md5Hash := dfsutil.Md5_sum(body)
		response := map[string]string{
			"filekey":  queryForm["filekey"][0],
			"hash":     md5Hash,
			"filesize": strconv.Itoa(len(body)),
		}
		fmt.Fprintln(w, dfsutil.Json_encode(response))
	} else {
		w.WriteHeader(701)
		fmt.Fprintln(w, filepath+" not exists!")
	}
	return
}

func Start() {
	http.HandleFunc("/upload", Upload)
	http.HandleFunc("/download", Download)
	http.HandleFunc("/info", Info)
	address := viper.Get("server.host").(string) + ":" + viper.Get("server.port").(string)
	fmt.Println(address)
	http.ListenAndServe(address, nil)
}
