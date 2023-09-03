package service

import (
	"bytes"
	"chat-server/models/common/response"
	"chat-server/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"path"
	"strings"
	"time"
)

func UploadImage(c *gin.Context) {
	w := c.Writer
	req := c.Request

	// 接收文件
	_, head, err := req.FormFile("file")

	if err != nil {
		response.FailWithDetailed(w, err.Error(), c)
		return
	}
	// 判断文件是否为jpg/png格式
	ext := strings.ToLower(path.Ext(head.Filename))

	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		response.FailWithMessage("只支持jpg/jpeg/png图片上传", c)
		return
	}

	// 打开上传文件
	fileHandle, fe := head.Open()
	if err != nil {
		response.FailWithDetailed(w, fe.Error(), c)
		return
	}
	defer fileHandle.Close()

	// 获取上传文件字节流
	fileByte, be := io.ReadAll(fileHandle)

	if be != nil {
		response.FailWithDetailed(w, be.Error(), c)
		return
	}

	client, err := utils.InitOss()

	list, _ := client.DescribeRegions()

	for _, region := range list.Regions {
		fmt.Println(region)
	}

	if err != nil {
		response.FailWithDetailed(w, err.Error(), c)
		return
	}

	bucket, e := client.Bucket("spence-oss-bucket")

	if e != nil {
		response.FailWithDetailed(w, e.Error(), c)
	}

	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), ext)
	yunFileTmpPath := "uploads" + "/" + fileName

	e = bucket.PutObject(yunFileTmpPath, bytes.NewReader([]byte(fileByte)))

	if e != nil {
		response.FailWithDetailed(w, e.Error(), c)
	}

	/*// 设置文件名
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), ext)

	// 创建文件并且存在某个目录下
	dstFile, e := os.Create("./upload/" + fileName)

	if e != nil {
		response.FailWithDetailed(w, e.Error(), c)
		return
	}

	_, copyErr := io.Copy(dstFile, file)
	if copyErr != nil {
		response.FailWithDetailed(w, e.Error(), c)
		return
	}

	// 返回文件路径
	url := "/upload/" + fileName
	response.OkWithData(url, "上传成功", c)
	defer dstFile.Close()*/
}
