package user

import (
	"manage/controller"
	"manage/logic/user"
	"manage/model"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateCommentHandler 创建评论请求的处理函数
func CreateCommentHandler(c *gin.Context) {
	p := new(model.ParamComment)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreateComment with invalid param", zap.Error(err))
		controller.ValidateError(c, err)
		return
	}
	userID, _ := controller.GetCurrentUserID(c)
	if p.UserID != userID {
		controller.ResponseAuthError(c, controller.CodeInvalidAuth)
		return
	}
	commentID, err := user.CreateComment(p)
	if err != nil {
		zap.L().Error("logic.CreateComment failed",
			zap.Int64("userID", p.UserID),
			zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
		"id": commentID,
	})
}

// DeleteCommentHandler 删除评论请求的处理函数
func DeleteCommentHandler(c *gin.Context) {
	p := new(model.ParamDeleteComment)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("DeleteComment with invalid param", zap.Error(err))
		controller.ValidateError(c, err)
		return
	}
	userID, _ := controller.GetCurrentUserID(c)
	if p.UserID != userID {
		controller.ResponseAuthError(c, controller.CodeInvalidAuth)
		return
	}
	if err := user.DeleteComment(p); err != nil {
		zap.L().Error("logic.DeleteComment failed",
			zap.Int64("userID", p.UserID),
			zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, nil)
}

// CommentListHandler 获取评论列表请求的处理函数
func CommentListHandler(c *gin.Context) {
	p := &model.ParamCommentList{
		Order:         model.OrderTime,
		EndTime:       time.Now().String(),
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("CommentListHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}
	commentList, total, err := user.GetCommentList(p)
	if err != nil {
		zap.L().Error("logic.GetCommentList failed",
			zap.Int64("itemID", p.ItemID),
			zap.Int64("userID", p.CurrentUserID),
			zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
		"list":  commentList,
		"total": total,
	})
}