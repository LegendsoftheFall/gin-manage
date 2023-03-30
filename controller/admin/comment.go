package admin

import (
	"manage/controller"
	"manage/logic/admin"
	"manage/model"
	"strconv"
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

// GetCommentByItemIDHandler 后台获取资源ID的评论请求
func GetCommentByItemIDHandler(c *gin.Context) {
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

	commentList, total, err := admin.GetCommentByItemID(p)
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

// GetCommentInfoHandler 后台获取评论信息请求
func GetCommentInfoHandler(c *gin.Context) {
	commentID := c.Param("id")
	id, ok := strconv.ParseInt(commentID, 10, 64)
	if ok != nil {
		zap.L().Error("commentID is invalid", zap.Int64("commentID", id), zap.Error(ok))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}

	commentInfo, err := admin.GetCommentInfo(id)
	if err != nil {
		zap.L().Error("admin.GetCommentInfo failed",
			zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}

	controller.ResponseSuccess(c, commentInfo)
}

func SetCommentStatusHandler(c *gin.Context) {
	p := &model.ParamSetStatus{
		ID:     0,
		Status: 0,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("SetCommentStatusHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}

	err := admin.SetCommentStatus(p.ID, p.Status)
	if err != nil {
		zap.L().Error("admin.SetCommentStatus failed",
			zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}

	controller.ResponseSuccess(c, nil)
}

// DeleteCommentForAdminHandler 后台删除文章请求的处理函数
func DeleteCommentForAdminHandler(c *gin.Context) {
	p := new(model.ParamAdminDeleteComment)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("DeleteArticle with invalid param", zap.Error(err))
		controller.ValidateError(c, err)
		return
	}
	if p.Secret != MySecret {
		controller.ResponseError(c, controller.CodeInvalidAuth)
		return
	}
	if err := admin.DeleteCommentForAdmin(p); err != nil {
		zap.L().Error("admin.DeleteCommentForAdmin(p) failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, nil)
}
