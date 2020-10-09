package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/yinjk/go-utils/pkg/net/common"
	"github.com/yinjk/go-utils/pkg/net/http"
	arrays "github.com/yinjk/go-utils/pkg/utils/collection/list"

	"rift/util"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Err  string `json:"err"`
}
type Config struct {
	DataPath   string
	RootPath   string
	Domain     string
	Port       string
	Token      string
	StaticAuth bool
}

var config Config

func init() {
	if util.ExistsFile("conf/app.toml") {
		print("conf/app.toml")
	} else {
		if !util.ExistsPath("conf") {
			if err := os.MkdirAll("conf", 0777); err != nil {
				panic(err)
			}
		}
		f, err := os.Create("conf/app.toml")
		if err != nil {
			panic(err)
		}
		token := util.GenerateToken()
		ip, err := util.GetClientIp()
		if err != nil {
			ip = "localhost"
		}
		defaultConf :=
			`rootPath="root"
dataPath="data"
domain="http://%s:8088/"
port=":8088"
token="%s"
staticAuth=false
`
		_, err = f.WriteString(fmt.Sprintf(defaultConf, ip, token))
		if err != nil {
			panic(err)
		}
	}
}

func _init(config Config) {
	if !util.ExistsPath(config.DataPath) {
		if err := os.MkdirAll(config.DataPath, 0777); err != nil {
			panic(err)
		}
	}
	if !util.ExistsPath(config.RootPath) {
		if err := os.MkdirAll(config.RootPath, 0777); err != nil {
			panic(err)
		}
	}
}

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o rift .
func main() {
	if _, err := toml.DecodeFile("conf/app.toml", &config); err != nil {
		panic(err)
	}
	_init(config)
	engine := http.NewEngine(http.Config{
		Mode: "debug",
		Port: config.Port,
	})
	engine.Use(authorization)
	engine.POST("/upload", upload)
	engine.GET("/list", list)
	engine.POST("/mkdir", mkdir)
	engine.DELETE("/files", deleteFiles)
	engine.POST("/move", move)
	//engine.GET("/download", download)
	engine.Static("/static", config.RootPath)
	engine.ListenAndStartUp()
}

func authorization(ctx *gin.Context) {
	token := ctx.GetHeader("token")
	fmt.Println(ctx.Request.RequestURI)
	if !config.StaticAuth && strings.HasPrefix(ctx.Request.RequestURI, "/static") {
		return
	}
	if !config.StaticAuth && strings.HasPrefix(ctx.Request.RequestURI, "/upload") {
		token = ctx.PostForm("token")
	}
	if config.Token != "" && config.Token != token {
		ctx.JSON(401, common.NewFailResult(401, "Token authentication failed"))
		ctx.Abort()
		return
	}
}

func upload(ctx *gin.Context) {
	//t
	//fmt.Println(token)
	dir := ctx.PostForm("dir")
	files, err := ctx.FormFile("file_data")
	if err != nil {
		ctx.JSON(400, common.NewFailResult(400, err.Error()))
		return
	}
	if fileUrl, err := saveFileDirect(files, dir); err != nil {
		ctx.JSON(500, common.NewFailResult(500, err.Error()))
		return
	} else {
		ctx.JSON(200, common.NewSuccessResult(fileUrl, "upload success"))
	}
}

type FileInfo struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	SimpleName string `json:"simpleName"`
	Size       string `json:"size"`
	Time       string `json:"time"`
	IsDir      bool   `json:"isDir"`
	Url        string `json:"url"`
}

func list(ctx *gin.Context) {
	//token := ctx.PostForm("token")
	//fmt.Println(token)
	dir, _ := ctx.GetQuery("dir")
	if dir == "/" {
		dir = ""
	}
	absolutePath := getAbsolutePath(dir)
	files, err := ioutil.ReadDir(absolutePath)
	if err != nil {
		ctx.JSON(500, common.NewFailResult(500, err.Error()))
		return
	}
	res := make([]FileInfo, 0)
	for i, f := range files {
		res = append(res, FileInfo{
			ID:    i,
			Name:  f.Name(),
			Size:  util.FormatSize(f.Size()),
			Time:  util.FormatTime(f.ModTime()),
			IsDir: f.IsDir(),
			Url:   getFileUrl(dir, f.Name()),
		})
	}
	arrays.StreamOfSlice(res).Sorted(func(o1, o2 interface{}) bool {
		f1 := o1.(FileInfo)
		f2 := o2.(FileInfo)
		return (f1.IsDir == f2.IsDir && f1.Name < f2.Name) || (f1.IsDir && !f2.IsDir)
	}).Unmarshal(&res)
	ctx.JSON(200, common.NewSuccessResult(res))
}

//
//func download(ctx *gin.Context) {
//	response, err := nethttp.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
//	if err != nil || response.StatusCode != nethttp.StatusOK {
//		ctx.Status(nethttp.StatusServiceUnavailable)
//		return
//	}
//
//	reader := response.Body
//	contentLength := response.ContentLength
//	contentType := response.Header.Get("Content-Type")
//
//	extraHeaders := map[string]string{
//		"Content-Disposition": `attachment; filename="gopher.png"`,
//	}
//
//	file, err := os.Open("/data")
//	if err != nil {
//		panic(err)
//	}
//	defer file.Close()
//	bytes, err := ioutil.ReadAll(file)
//	if err != nil {
//		panic(err)
//	}
//	fileName := file.Name()
//	ctx.Header("Content-Type", "application/octet-stream")
//	ctx.Header("Content-Disposition", "attachment; filename="+fileName)
//	_, err = ctx.Writer.Write(bytes)
//	if err != nil {
//		panic(err)
//	}
//	ctx.DataFromReader(nethttp.StatusOK, contentLength, contentType, reader, extraHeaders)
//
//}

func mkdir(ctx *gin.Context) {
	dir := ctx.PostForm("dir")
	absolutePath := getAbsolutePath(dir)
	if err := os.MkdirAll(absolutePath, 0777); err != nil {
		ctx.JSON(500, common.NewFailResult(500, err.Error()))
		return
	}
	ctx.JSON(200, common.NewSuccessResult("success"))
	return
}

func deleteFiles(ctx *gin.Context) {
	var req map[string][]string
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, common.NewFailResult(400, err.Error()))
		return
	}
	files := req["files"]
	if files == nil || len(files) == 0 {
		ctx.JSON(400, common.NewFailResult(400, "request file list is null"))
		return
	}
	for _, filePath := range files {
		if err := os.RemoveAll(getAbsolutePath(filePath)); err != nil {
			ctx.JSON(500, common.NewFailResult(500, err.Error()))
			return
		}
	}
	ctx.JSON(200, common.NewSuccessResult("success", "success"))
}

func move(ctx *gin.Context) {
	var req map[string][]string
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(400, common.NewFailResult(400, err.Error()))
		return
	}
	files := req["files"]
	newPath := req["path"][0]
	if files == nil || len(files) == 0 {
		ctx.JSON(400, common.NewFailResult(400, "request file list is null"))
		return
	}
	for _, filePath := range files {
		base := path.Base(filePath)
		if err := os.Rename(getAbsolutePath(filePath), path.Join(getAbsolutePath(newPath), "/"+base)); err != nil {
			ctx.JSON(500, common.NewFailResult(500, err.Error()))
			return
		}
	}
	ctx.JSON(200, common.NewSuccessResult("success", "success"))
}

func saveFileDirect(fileHeader *multipart.FileHeader, dir string) (source string, err error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	fileName := fileHeader.Filename
	path.Join()
	fileDir := path.Join(config.RootPath, dir)
	fileAbsolutePath := path.Join(config.RootPath, dir, fileName)
	fileRelativePath := path.Join(dir, fileName)
	if !util.ExistsFile(fileDir) {
		if err := os.MkdirAll(fileAbsolutePath, 0777); err != nil {
			return "", err
		}
	}
	f, err := os.OpenFile(fileAbsolutePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	//从头开始读文件，保存文件到服务器
	if _, err = file.Seek(0, 0); err != nil { //offset偏移位置，whence为0时表示从文件开始偏移，为1时表示从当前位置偏移，为2时表示从文件结尾偏移
		return "", err
	}
	if _, err = io.Copy(f, file); err != nil {
		return "", err
	}
	return path.Join(config.Domain, "/static", "/"+fileRelativePath), nil
}

func saveFile(fileHeader *multipart.FileHeader) (source string, err error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	fMd5 := util.Md5File(file)
	fmt.Println(fMd5)
	fileName := fileHeader.Filename
	fileType := path.Ext(fileName)
	now := time.Now()
	timePath := now.Format("2006/01/02/")
	fileDirPath := config.RootPath + "/" + timePath
	realFileName := strconv.Itoa(now.Nanosecond()) + fileType
	filePath := fileDirPath + realFileName
	if !util.ExistsFile(fileDirPath) {
		if err := os.MkdirAll(fileDirPath, 0777); err != nil {
			return "", err
		}
	}
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	//从头开始读文件，保存文件到服务器
	if _, err = file.Seek(0, 0); err != nil { //offset偏移位置，whence为0时表示从文件开始偏移，为1时表示从当前位置偏移，为2时表示从文件结尾偏移
		return "", err
	}
	if _, err = io.Copy(f, file); err != nil {
		return "", err
	}
	return config.Domain + timePath + realFileName, nil
}
func getAbsolutePath(path string) string {
	if strings.HasPrefix(path, "/") {
		return config.RootPath + path
	}
	return config.RootPath + "/" + path
}

func getFileUrl(dir, fileName string) string {
	return path.Join(config.Domain, "static", dir, fileName)
}
