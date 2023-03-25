package user

import (
	"manage/controller"
	"manage/logic/user"
	"manage/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// TagDoFollowHandler 发起关注标签请求的处理函数
func TagDoFollowHandler(c *gin.Context) {
	p := new(model.ParamFollowTag)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("TagDoFollowHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}
	userID, _ := controller.GetCurrentUserID(c)
	if userID != p.UserID {
		controller.ResponseError(c, controller.CodeInvalidAuth)
		return
	}
	if err := user.FollowTag(p); err != nil {
		zap.L().Error("logic.FollowTag failed",
			zap.Int64("userID", p.UserID),
			zap.Int("tagID", p.TagID),
			zap.Error(err))
	}
	controller.ResponseSuccess(c, nil)
}

// TagUnDoFollowHandler 发起取消关注标签请求的处理函数
func TagUnDoFollowHandler(c *gin.Context) {
	p := new(model.ParamFollowTag)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("TagUnDoFollowHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}
	userID, _ := controller.GetCurrentUserID(c)
	if userID != p.UserID {
		controller.ResponseError(c, controller.CodeInvalidAuth)
		return
	}
	if err := user.FollowTagCancel(p); err != nil {
		zap.L().Error("logic.FollowTagCancel failed",
			zap.Int64("userID", p.UserID),
			zap.Int("tagID", p.TagID),
			zap.Error(err))
	}
	controller.ResponseSuccess(c, nil)
}

// UserDoFollowHandler 发起关注用户请求的处理函数
func UserDoFollowHandler(c *gin.Context) {
	p := new(model.ParamFollowUser)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("UserDoFollowHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}
	userID, _ := controller.GetCurrentUserID(c)
	if userID != p.UserID {
		controller.ResponseError(c, controller.CodeInvalidAuth)
		return
	}
	if err := user.FollowUser(p); err != nil {
		zap.L().Error("logic.FollowUser failed",
			zap.Int64("userID", p.UserID),
			zap.Int64("follow_user_ID", p.FollowUserID),
			zap.Error(err))
	}
	controller.ResponseSuccess(c, nil)
}

// UserUnDoFollowHandler 发起取消关注用户请求的处理函数
func UserUnDoFollowHandler(c *gin.Context) {
	p := new(model.ParamFollowUser)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("UserDoFollowHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}
	userID, _ := controller.GetCurrentUserID(c)
	if userID != p.UserID {
		controller.ResponseError(c, controller.CodeInvalidAuth)
		return
	}
	if err := user.FollowUserCancel(p); err != nil {
		zap.L().Error("logic.FollowUserCancel failed",
			zap.Int64("userID", p.UserID),
			zap.Int64("follow_user_ID", p.FollowUserID),
			zap.Error(err))
	}
	controller.ResponseSuccess(c, nil)
}
