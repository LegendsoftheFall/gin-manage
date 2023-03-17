package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RESTFUL API

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

//ResponseError 返回错误响应
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

//ResponseAuthError 返回带有权限错误的401的错误响应
func ResponseAuthError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusUnauthorized, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

//// ResponseErrorWithData 返回错误响应
//func ResponseErrorWithData(c *gin.Context, code ResCode, data interface{}) {
//	c.JSON(http.StatusOK, &ResponseData{
//		Code: code,
//		Msg:  code.Msg(),
//		Data: data,
//	})
//}

//ResponseErrorWithMsg 返回带有错误信息的错误响应
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

//ResponseSuccess 返回成功响应
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}
