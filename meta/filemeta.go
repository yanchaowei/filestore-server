package meta

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

// 根据hash值获取元信息
func GetFileMeta(fileSha1 string) FielMeta {
	return fileMetas[fileSha1]
}

// RemoveFileMeta： 根据fileSha1删除文件元信息
// 需要考虑线程安全的问题
func RemoveFileMeta(fileSha1 string)  {
	delete(fileMetas, fileSha1)
}
