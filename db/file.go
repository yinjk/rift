package db

import (
	_ "crypto/md5"
	"database/sql"
	_ "encoding/hex"
	_ "io"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"file-uploader/util"
)

const DataSource = "root:123456@tcp(localhost)/file?charset=utf8"
const DriverName = "mysql"

type FileInfo struct {
	Id         int
	Name       string
	Path       string
	Url        string
	FileType   string
	Size       int64
	Md5        string
	Ip         string
	CreateTime time.Time
}

var db, _ = sql.Open(DriverName, DataSource)

/**
 * 将文件信息保存进数据库
 */
func SaveFile(fileInfo FileInfo) (string, error) {
	//db, _ = sql.Open(DriverName, DataSource)
	//插入数据
	stmt, err := db.Prepare("INSERT INTO file_info(name, path, url, file_type, size, md5, ip) VALUES(?, ?, ?, ?, ?, ?, ?)")
	util.CheckErr(err)

	_, err = stmt.Exec(fileInfo.Name, fileInfo.Path, fileInfo.Url, fileInfo.FileType, fileInfo.Size, fileInfo.Md5, fileInfo.Ip)
	if err != nil {
		os.Remove(fileInfo.Path) //如果插入数据库失败，应当删除服务器真实路径下对应的文件，避免造成脏数据
		return "", err
	}
	return fileInfo.Url, nil
}

/**
 * 通过md5校验码检查该文件是否存在
 */
func ContainsFile(md5s string) (bool, *FileInfo) {
	//db, err = sql.Open(DriverName, DataSource)
	rows, err := db.Query("SELECT id, name, path, url, file_type as fileType, size, create_time as createTime FROM file_info where md5 = ?", md5s)
	util.CheckErr(err)
	if rows.Next() {
		var fileInfo = FileInfo{}
		rows.Scan(&fileInfo.Id, &fileInfo.Name, &fileInfo.Path, &fileInfo.Url, &fileInfo.FileType, &fileInfo.Size, &fileInfo.CreateTime)
		return true, &fileInfo
	}
	return false, nil
}

/**
 * 查询所有文件列表
 */
func FileList() ([]FileInfo) {
	return nil
}

/**
 * 根据id删除文件
 */
func DeleteFile(id int) bool {
	return false
}

//func db(file *os.File, filePath string, ip string) {
//    //获取文件md5
//    md5s  := util.Md5File(file)
//    //获取文件信息
//    fileInfo, _ := file.Stat()
//    fileSize := fileInfo.Size()
//    fileName := fileInfo.Name()
//    fileType := path.Ext(fileName)
//    db, err := sql.Open("mysql", "root:123456@tcp(193.112.112.136)/file?charset=utf8")
//    //插入数据
//    stmt, err := db.Prepare("INSERT INTO file_info(name, path, file_type, size, ")
//    util.CheckErr(err)
//
//    res, err := stmt.Exec(fileName, filePath, fileType, fileSize, md5s, ip)
//    util.CheckErr(err)
//
//    id, err := res.LastInsertId()
//    util.CheckErr(err)
//
//    fmt.Println(id)
//
//    //查询数据
//    rows, err := db.Query("SELECT id, name, path as filePath, file_type as fileType, size, create_time FROM file_info")
//    util.CheckErr(err)
//
//    for rows.Next() {
//        var id, size int
//        var name, filePath, fileType, createTime string
//        err = rows.Scan(&id, &name, &filePath, &fileType, &size, &createTime)
//        util.CheckErr(err)
//        fmt.Println(id, name, filePath, fileType, size, createTime)
//    }
//
//    //删除数据
//    //stmt, err = db.Prepare("delete from file_info where id =?")
//    //CheckErr(err)
//    //
//    //res, err = stmt.Exec(id)
//    //CheckErr(err)
//    //
//    //affect, err = res.RowsAffected()
//    //CheckErr(err)
//    //
//    //fmt.Println(affect)
//
//    defer db.Close()
//
//}
