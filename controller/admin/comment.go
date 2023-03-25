package admin

import (
	"manage/controller"
	"manage/logic/admin"
	"manage/model"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetAllCommentHandler 后台获取所有评论请求
func GetAllCommentHandler(c *gin.Context) {
	p := &model.ParamAdminComment{
		Order:   model.OrderTime,
		EndTime: time.Now().String(),
		Page:    1,
		Size:    10,
		ItemID:  0,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetAllCommentHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}

	commentList, total, err := admin.GetAllComment(p)
	if err != nil {
		zap.L().Error("logic.GetCommentList failed",
			zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
		"list":  commentList,
		"page":  p.Page,
		"size":  p.Size,
		"total": total,
	})
}
