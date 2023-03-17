package controller

import (
	"fmt"
	"manage/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UpLoadHandler 上传图片请求的处理函数
func UpLoadHandler(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")

	if err != nil {
		zap.L().Error("UpLoad with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	fileSize := fileHeader.Size

	url, err := logic.UpLoad(file, fileSize)
	if err != nil {
		zap.L().Error("logic.UpLoad failed", zap.String("fileName", fileHeader.Filename), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	fmt.Println(url)

	ResponseSuccess(c, gin.H{
		"url": url,
	})
}
