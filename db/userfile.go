package db

import (
	mydb "filestore-server/db/mysql"
	"fmt"
	"time"
)

// UserFile:用户文件表列表
type UserFile struct {
	UserName string
	FileHash string
	FileName string
	FileSize int64
	UploadAt string
	LastUpdated string
}

func OnUserFileUploadFinished(username, filehash, filename string, fileSize int64) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user_file (`user_name`, `file_sha1`, `file_name`, `file_size`," +
		"upload_at) values(?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Printf(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, fileSize, time.Now())
	if err != nil {
		return false
	}
	return true
}


// QueryUserFileInfo: 批量获取用户文件信息
func QueryUserFileInfo(username string, limit int) ([]UserFile, error) {
	stmt, err := mydb.DBConn().Prepare("select file_sha1, file_size, file_name, upload_at, last_update " +
		"from tbl_user_file where user_name=? limit ?")
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err = rows.Scan(&ufile.FileHash, &ufile.FileSize, &ufile.FileName, &ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}
	return userFiles, nil
}