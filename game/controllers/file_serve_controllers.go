package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

var FileServerCtl = &FileServerControllor{}

type FileServerControllor struct {
}

const (
	debug      = false
	accessRoot = "/www/wwwlogs/"
	testRoot   = "./"
)

func getRoot() string {
	if debug {
		return testRoot
	}
	return accessRoot
}

func getLastFile() os.FileInfo {
	dirs, err := ioutil.ReadDir(getRoot())
	if err != nil {
		fmt.Printf("没有这个目录！")
		return nil
	}
	var lastFile os.FileInfo
	for _, dir := range dirs {
		if !dir.IsDir() && dir.Name()[len(dir.Name())-3:] == ".gz" {
			if lastFile == nil || lastFile.ModTime().Before(dir.ModTime()) {
				dir1 := dir
				lastFile = dir1
			}
		}

	}
	return lastFile
}

func giveDownload(c *gin.Context, rootDir, filename string) {
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename)) //fmt.Sprintf("attachment; filename=%s", filename)对下载的文件重命名
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(rootDir + filename)
}

func (*FileServerControllor) GetLastAccess(c *gin.Context) {
	file := getLastFile()
	if file != nil {
		giveDownload(c, getRoot(), file.Name())
		//c.String(http.StatusOK, "最新文件：%s， 时间%s\n", file.Name(), file.ModTime())
	} else {
		c.String(http.StatusOK, "暂时没有最新文件！")
	}
}

func (*FileServerControllor) GetNowAccess(c *gin.Context) {
	giveDownload(c, getRoot(), "poke.ccc.log")
}
