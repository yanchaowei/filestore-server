package meta

import (
	mydb "filestore-server/db"
)

// FileMeta： 文件元信息结构
type FielMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FielMeta

func init() {
	fileMetas = make(map[string]FielMeta)
}

// UpdateFileMeta: 新增/更新 Filemeta 元信息
func UpdateFileMeta(fmeta FielMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// UpdateFileMetaDB: 更新源文件信息到MySQL中
func UpdateFileMetaDB(fmeta FielMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// 根据hash值获取元信息
func GetFileMeta(fileSha1 string) FielMeta {
	return fileMetas[fileSha1]
}

// RemoveFileMeta： 根据fileSha1删除文件元信息
// 需要考虑线程安全的问题
func RemoveFileMeta(fileSha1 string)  {
	delete(fileMetas, fileSha1)
}

// GetFileMetaDB: 从MySQL获取文件元信息
func GetFileMetaDB(filesha1 string) (FielMeta, error) {
	tfile, err := mydb.GetFileMeta(filesha1)
	if err != nil {
		return FielMeta{}, nil
	}
	fmeta := FielMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		Location: tfile.FileAddr.String,
		FileSize: tfile.FileSize.Int64,
	}
	return fmeta, nil
}
