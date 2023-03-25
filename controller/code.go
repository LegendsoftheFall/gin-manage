package controller

//定义状态码并创建结构体返回对应的信息

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeInvalidToken
	CodeNeedLogin
	CodeInvalidAuth

	CodeTagExist
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "邮箱或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeInvalidToken: "无效的Token,请重新登录",
	CodeNeedLogin:    "需要登录",
	CodeInvalidAuth:  "没有权限",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
