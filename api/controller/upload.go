package controller

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"zendea/util"
	"zendea/util/log"
	"zendea/util/uploader"
)

type UploadController struct {
	BaseController
}

const uploadMaxBytes int64 = 1024 * 1024 * 3 // 1M

// Upload upload file
func (c *UploadController) Upload(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("image")
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	defer file.Close()
	if header.Size > uploadMaxBytes {
		c.Fail(ctx, util.NewErrorMsg("图片不能超过3M"))
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}

	log.Info("上传文件：%s, size: %s", header.Filename, header.Size)

	url, err := uploader.PutImage(fileBytes)
	if err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}
	data := make(map[string]string)
	data["url"] = url
	c.Success(ctx, data)
}

// UploadFromEditor upload file from editor
func (c *UploadController) UploadFromEditor(ctx *gin.Context) {
	errFiles := make([]string, 0)
	succMap := make(map[string]string)

	user := c.GetCurrentUser(ctx)
	if user == nil {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"code":    1,
			"message": "请先登录",
			"data": gin.H{
				"errFiles": errFiles,
				"succMap":  succMap,
			},
		})
		return
	}

	mForm, _ := ctx.MultipartForm()
	files := mForm.File["file[]"]
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			log.Error(err.Error())
			errFiles = append(errFiles, file.Filename)
			continue
		}
		fileBytes, err := ioutil.ReadAll(f)
		if err != nil {
			log.Error(err.Error())
			errFiles = append(errFiles, file.Filename)
			continue
		}
		url, err := uploader.PutImage(fileBytes)
		if err != nil {
			log.Error(err.Error())
			errFiles = append(errFiles, file.Filename)
			continue
		}

		succMap[file.Filename] = url
	}

	c.Success(ctx, gin.H{
		"errFiles": errFiles,
		"succMap":  succMap,
	})

}

// UploadFromURL fetch file by URL
func (c *UploadController) UploadFromURL(ctx *gin.Context) {
	user := c.GetCurrentUser(ctx)
	if user == nil {
		c.Fail(ctx, util.ErrorNotLogin)
		return
	}

	data := make(map[string]string)
	if err := ctx.ShouldBindJSON(&data); err != nil {
		c.Fail(ctx, util.FromError(err))
		return
	}

	url := data["url"]
	output, err := uploader.CopyImage(url)
	if err != nil {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"code":    0,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}
	c.Success(ctx, gin.H{
		"originalURL": url,
		"url":         output,
	})
}
