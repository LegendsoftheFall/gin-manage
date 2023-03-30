package admin

import (
	"manage/controller"
	"manage/logic/user"
	"manage/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// DeleteArticleForAdminHandler 后台删除文章请求的处理函数
func DeleteArticleForAdminHandler(c *gin.Context) {
	p := new(model.ParamDeleteArticle)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("DeleteArticle with invalid param", zap.Error(err))
		controller.ValidateError(c, err)
		return
	}
	if p.Secret != MySecret {
		controller.ResponseError(c, controller.CodeInvalidAuth)
		return
	}
	if err := user.DeleteArticle(p.ID); err != nil {
		zap.L().Error("user.DeleteArticle(p.ID) failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, nil)
}
