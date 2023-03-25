package user

import (
	"errors"
	"fmt"
	"manage/controller"
	user2 "manage/dao/mysql/user"
	"manage/logic/user"
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
		controller.ValidateError(c, err)
		return
	}
	//业务处理
	if err := user.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, user2.ErrUserExist) {
			controller.ResponseError(c, controller.CodeUserExist)
		}
		return
	}
	controller.ResponseSuccess(c, nil)
}

// LoginHandler  处理登录请求的函数
func LoginHandler(c *gin.Context) {
	p := new(model.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		controller.ValidateError(c, err)
		return
	}

	//业务处理
	u, err := user.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("email", p.Email), zap.Error(err))
		if errors.Is(err, user2.ErrUserNotExist) {
			controller.ResponseError(c, controller.CodeUserNotExist)
			return
		}
		controller.ResponseError(c, controller.CodeInvalidPassword)
		return
	}

	controller.ResponseSuccess(c, gin.H{
		//"user_id":   fmt.Sprintf("%d", user.UserID),
		//"user_name": user.Username,
		//"email":     user.Email,
		"aToken": u.AccessToken,
		"rToken": u.RefreshToken,
		"userID": strconv.FormatInt(u.UserID, 10),
	})
}

// RefreshHandler 处理过期token请求的函数
func RefreshHandler(c *gin.Context) {
	//获取token参数
	p := new(model.ParamToken)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Refresh with invalid param", zap.Error(err))
		controller.ValidateError(c, err)
		return
	}
	newAToken, newRToken, err := jwt.RefreshToken(p.AToken, p.RToken)
	if err != nil {
		if err == jwt.ErrNeedLogin {
			controller.ResponseError(c, controller.CodeInvalidToken)
		}
		return
	}
	controller.ResponseSuccess(c, gin.H{
		"newAToken": newAToken,
		"newRToken": newRToken,
	})
}

// InfoOfUserHandler 处理登录后获取用户信息请求的函数
func InfoOfUserHandler(c *gin.Context) {
	id, err := controller.GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("getCurrentUserID failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeNeedLogin)
		return
	}
	user, err := user.GetUserByID(id)
	if err != nil {
		zap.L().Error("GetUserInfoByID failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
		"id":     fmt.Sprintf("%d", user.UserID),
		"name":   user.Username,
		"email":  user.Email,
		"avatar": user.Avatar,
	})
}

// HomeHandler 处理用户主页请求的处理函数
func HomeHandler(c *gin.Context) {
	p := &model.ParamArticleList{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
	}
	// 获取分页参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("UserDoFollowHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}
	userID := c.Param("id")
	id, ok := strconv.ParseInt(userID, 10, 64)
	if ok != nil {
		zap.L().Error("userID is invalid", zap.Int64("userID", id), zap.Error(ok))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}
	userInfo, err := user.GetUserInfoByUserID(id, p.CurrentUserID)
	if err != nil {
		zap.L().Error("logic.GetUserInfoByID is failed", zap.Int64("userID", id), zap.Error(err))
		controller.ResponseError(c, controller.CodeUserNotExist)
		return
	}
	articleList, total, err := user.GetArticleByID(id, p.Page, p.Size)
	if err != nil {
		zap.L().Error("logic.GetArticleByID is failed", zap.Int64("userID", id), zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
		"userInfo":    userInfo,
		"articleInfo": articleList,
		"total":       total,
		"page":        p.Page,
		"size":        p.Size,
	})
}

// ArticleOfUserHandler 处理用户所有文章请求的处理函数
func ArticleOfUserHandler(c *gin.Context) {
	p := &model.ParamArticleList{
		CurrentUserID: 0,
	}
	// 获取分页参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("UserDoFollowHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}
	articleList, total, err := user.GetAllArticleByID(p.CurrentUserID)
	if err != nil {
		zap.L().Error("logic.GetArticleByID is failed", zap.Int64("userID", p.CurrentUserID), zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
		"articleInfo": articleList,
		"total":       total,
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
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}
	data, total, err := user.GetFollowingUsers(p)
	if err != nil {
		zap.L().Error("logic.GetFollowUsers failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
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
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}
	data, total, err := user.GetFollowerUsers(p)
	if err != nil {
		zap.L().Error("logic.GetFollowUsers failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
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
		controller.ValidateError(c, err)
		return
	}
	userID, err := controller.GetCurrentUserID(c)
	if err != nil {
		controller.ResponseError(c, controller.CodeInvalidAuth)
		return
	}
	if err = user.UpdateUserProfile(userID, p); err != nil {
		zap.L().Error("logic.UpdateUserProfile failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, nil)
}

// ProfileOfUserHandler 获取后台用户资料请求的处理函数
func ProfileOfUserHandler(c *gin.Context) {
	userID, err := controller.GetCurrentUserID(c)
	if err != nil {
		controller.ResponseError(c, controller.CodeInvalidAuth)
		return
	}
	profile, err := user.GetUserProfile(userID)
	if err != nil {
		zap.L().Error("logic.GetUserProfile failed",
			zap.Int64("userID", userID),
			zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
		"profile": profile,
	})
}

// ProfileHandler 获取用户资料请求的处理函数
func ProfileHandler(c *gin.Context) {
	p := new(model.ParamUserProfile)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("ProfileHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}
	profile, err := user.GetProfile(p)
	if err != nil {
		zap.L().Error("logic.GetProfile failed",
			zap.Int64("currentUserID", p.CurrentUserID),
			zap.Int64("userID", p.UserID),
			zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
		"profile": profile,
	})
}
