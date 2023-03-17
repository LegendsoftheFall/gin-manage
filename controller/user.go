package controller

import (
	"errors"
	"fmt"
	"manage/dao/mysql"
	"manage/logic"
	"manage/model"
	"manage/pkg/jwt"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	//获取请求参数校验
	p := new(model.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		validateError(c, err)
		return
	}
	//业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrUserExist) {
			ResponseError(c, CodeUserExist)
		}
		return
	}
	ResponseSuccess(c, nil)
}

// LoginHandler  处理登录请求的函数
func LoginHandler(c *gin.Context) {
	p := new(model.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		validateError(c, err)
		return
	}

	//业务处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("email", p.Email), zap.Error(err))
		if errors.Is(err, mysql.ErrUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	ResponseSuccess(c, gin.H{
		//"user_id":   fmt.Sprintf("%d", user.UserID),
		//"user_name": user.Username,
		//"email":     user.Email,
		"aToken": user.AccessToken,
		"rToken": user.RefreshToken,
		"userID": strconv.FormatInt(user.UserID, 10),
	})
}

// RefreshHandler 处理过期token请求的函数
func RefreshHandler(c *gin.Context) {
	//获取token参数
	p := new(model.ParamToken)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Refresh with invalid param", zap.Error(err))
		validateError(c, err)
		return
	}
	newAToken, newRToken, err := jwt.RefreshToken(p.AToken, p.RToken)
	if err != nil {
		if err == jwt.ErrNeedLogin {
			ResponseError(c, CodeInvalidToken)
		}
		return
	}
	ResponseSuccess(c, gin.H{
		"newAToken": newAToken,
		"newRToken": newRToken,
	})
}

// UserInfoHandler 处理登录后获取用户信息请求的函数
func UserInfoHandler(c *gin.Context) {
	id, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("getCurrentUserID failed", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	user, err := logic.GetUserByID(id)
	if err != nil {
		zap.L().Error("GetUserInfoByID failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"id":     fmt.Sprintf("%d", user.UserID),
		"name":   user.Username,
		"email":  user.Email,
		"avatar": user.Avatar,
	})
}

// UserHomeHandler 处理用户主页请求的处理函数
func UserHomeHandler(c *gin.Context) {
	p := &model.ParamArticleList{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
	}
	// 获取分页参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("UserDoFollowHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	userID := c.Param("id")
	id, ok := strconv.ParseInt(userID, 10, 64)
	if ok != nil {
		zap.L().Error("userID is invalid", zap.Int64("userID", id), zap.Error(ok))
		ResponseError(c, CodeInvalidParam)
		return
	}
	userInfo, err := logic.GetUserInfoByUserID(id, p.CurrentUserID)
	if err != nil {
		zap.L().Error("logic.GetUserInfoByID is failed", zap.Int64("userID", id), zap.Error(err))
		ResponseError(c, CodeUserNotExist)
		return
	}
	articleList, total, err := logic.GetArticleByID(id, p.Page, p.Size)
	if err != nil {
		zap.L().Error("logic.GetArticleByID is failed", zap.Int64("userID", id), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"userInfo":    userInfo,
		"articleInfo": articleList,
		"total":       total,
		"page":        p.Page,
		"size":        p.Size,
	})
}

// FollowingUsersHandler 获取关注者信息请求的处理函数
func FollowingUsersHandler(c *gin.Context) {
	p := &model.ParamUserList{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
	}
	//获取分页参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("FollowingUsersHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, total, err := logic.GetFollowingUsers(p)
	if err != nil {
		zap.L().Error("logic.GetFollowUsers failed", zap.Error(err))
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

// FollowerUsersHandler 获取追随者信息请求的处理函数
func FollowerUsersHandler(c *gin.Context) {
	p := &model.ParamUserList{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
	}
	//获取分页参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("FollowingUsersHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, total, err := logic.GetFollowerUsers(p)
	if err != nil {
		zap.L().Error("logic.GetFollowUsers failed", zap.Error(err))
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

// UpdateProfileHandler 更新用户资料请求的处理函数
func UpdateProfileHandler(c *gin.Context) {
	p := new(model.UserProfile)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdateProfile with invalid param", zap.Error(err))
		validateError(c, err)
		return
	}
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeInvalidAuth)
		return
	}
	if err = logic.UpdateUserProfile(userID, p); err != nil {
		zap.L().Error("logic.UpdateUserProfile failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// UserProfileHandler 获取后台用户资料请求的处理函数
func UserProfileHandler(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeInvalidAuth)
		return
	}
	profile, err := logic.GetUserProfile(userID)
	if err != nil {
		zap.L().Error("logic.GetUserProfile failed",
			zap.Int64("userID", userID),
			zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"profile": profile,
	})
}

// ProfileHandler 获取用户资料请求的处理函数
func ProfileHandler(c *gin.Context) {
	p := new(model.ParamUserProfile)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("ProfileHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	profile, err := logic.GetProfile(p)
	if err != nil {
		zap.L().Error("logic.GetProfile failed",
			zap.Int64("currentUserID", p.CurrentUserID),
			zap.Int64("userID", p.UserID),
			zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"profile": profile,
	})
}
