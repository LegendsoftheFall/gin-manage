package user

import (
	"manage/controller"
	"manage/logic/user"
	"manage/model"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// LikeHandler 点赞请求的处理函数
func LikeHandler(c *gin.Context) {
	p := new(model.ParamLike)

	//获取文章id和当前用户id参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("TagDetailHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}

	if err := user.LikeArticle(p.ArticleID, p.UserID); err != nil {
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, nil)
}
