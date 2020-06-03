package handler

import (
	"encoding/json"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"filestore-server/meta"
)

// 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "GET"{
		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	}else if r.Method == "POST" {
		// 接受文件流及存出到本地目录
		// 返回：文件句柄，文件头，错误信息
		file, head, err := r.FormFile("file")
		if err != nil{
			fmt.Printf("Fail to get data, err: %s", err.Error())
		}

		// 创建新文件的元信息
		fileMeta := meta.FielMeta{
			FileName: head.Filename,
			Location: "/home/yanchaowei/tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		// 关掉句柄
		defer file.Close()
		// 创建本地文件接受文件流, 返回新创建的文件的句柄
		newFile, err := os.Create(fileMeta.Location)
		if err != nil{
			fmt.Printf("Fail to create file, err:%s\n", err.Error())
			return
		}
		defer newFile.Close()

		// 将文件数据流拷贝到新文件
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil{
			fmt.Printf("Fail to save data into file, err: %s\n", err.Error())
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		//meta.UpdateFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)
		fmt.Printf("fileMeta: %+v", fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request)  {
	io.WriteString(w, "Upload finished!")
}

// GetFileMetaHandler: 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	//fmeta :=meta.GetFileMeta(filehash)
	fmeta, err :=meta.GetFileMetaDB(filehash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fmeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form["filehash"][0]
	fmt.Printf("!!!!!!fsha1:%s", fsha1)
	fmeta := meta.GetFileMeta(fsha1)
	fmt.Printf("!!!!!!fsha1:%+v", fmeta)
	file, err := os.Open(fmeta.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	fmt.Printf("!!!!!!fsha1:%s", data)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//w.Header().Set("Content-Type", http.DetectContentType(data)) : 这种方式会导致问价格式错误
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fmeta.FileName+"\"")
	w.Write(data)
	//return
}

// 更新文件元信息，目前只支持修改文件名
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	opType := r.Form.Get("op")
	newFileName := r.Form.Get("fileName")
	if opType != "0"{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)
	fmt.Printf("curFileMeta: %+v", curFileMeta)
	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// FileDeleteHandler： 删除文件
func FileDeleteHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	// 在删除元信息之前将文件的存储位置取出来, 利用os删除该文件在磁盘上的物理数据
	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)
	// 删除元信息
	meta.RemoveFileMeta(fileSha1)
	// 返回200
	w.WriteHeader(http.StatusOK)
}