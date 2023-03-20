package controller

import (
	"manage/logic"
	"manage/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HotTagsHandler 获取热门标签请求的处理函数
func HotTagsHandler(c *gin.Context) {
	p := &model.ParamTagList{
		Page:          1,
		Size:          6,
		CurrentUserID: 0,
	}
	//获取分页参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("TagsHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, total, err := logic.GetHotTags(p)
	if err != nil {
		zap.L().Error("logic.GetTrendingTags failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"list":  data,
		"total": total,
		"page":  p.Page,
		"size":  p.Size,
	})
}
