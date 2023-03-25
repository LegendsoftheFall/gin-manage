package admin

import (
	"errors"
	"fmt"
	"manage/controller"
	user2 "manage/dao/mysql/user"
	"manage/logic/admin"
	"manage/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SignUpForAdminHandler 处理管理员注册请求的函数
func SignUpForAdminHandler(c *gin.Context) {
	//获取请求参数校验
	p := new(model.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		controller.ValidateError(c, err)
		return
	}
	//业务处理
	if err := admin.SignUpForAdmin(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, user2.ErrUserExist) {
			controller.ResponseError(c, controller.CodeUserExist)
		}
		return
	}
	controller.ResponseSuccess(c, nil)
}

// LoginForAdminHandler  处理登录请求的函数
func LoginForAdminHandler(c *gin.Context) {
	p := new(model.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invalid param", zap.Error(err))
		controller.ValidateError(c, err)
		return
	}

	//业务处理
	a, err := admin.LoginForAdmin(p)
	if err != nil {
		zap.L().Error("logic.LoginForAdmin failed", zap.String("email", p.Email), zap.Error(err))
		if errors.Is(err, user2.ErrUserNotExist) {
			controller.ResponseError(c, controller.CodeUserNotExist)
			return
		}
		controller.ResponseError(c, controller.CodeInvalidPassword)
		return
	}

	controller.ResponseSuccess(c, gin.H{

		"aToken":  a.AccessToken,
		"rToken":  a.RefreshToken,
		"adminID": strconv.FormatInt(a.AdminID, 10),
	})
}

// LogoutForAdminHandler  处理登录请求的函数
func LogoutForAdminHandler(c *gin.Context) {
	controller.ResponseSuccess(c, nil)
}

// InfoOfAdminHandler 处理登录后获取用户信息请求的函数
func InfoOfAdminHandler(c *gin.Context) {
	id, err := controller.GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("getCurrentID failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeNeedLogin)
		return
	}
	a, err := admin.GetAdminByID(id)
	if err != nil {
		zap.L().Error("GetAdminInfoByID failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
		"id":     fmt.Sprintf("%d", a.AdminID),
		"name":   a.AdminName,
		"email":  a.Email,
		"avatar": a.Avatar,
	})
}
