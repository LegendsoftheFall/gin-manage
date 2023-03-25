package user

import (
	"manage/controller"
	"manage/logic/user"
	"manage/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CollectHandler 收藏请求的处理函数
func CollectHandler(c *gin.Context) {
	p := new(model.ParamCollect)
	//获取文章id和当前用户id参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("TagDetailHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}
	if p.UserID == "" {
		return
	}
	if err := user.CollectArticle(p.ArticleID, p.UserID); err != nil {
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, nil)
}

// BookMarkHandler 获取收藏文章列表请求的处理函数
func BookMarkHandler(c *gin.Context) {
	p := &model.ParamPage{
		Page: 1,
		Size: 10,
	}
	//获取分页参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("TagDetailHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}

	userID, err := controller.GetCurrentUserID(c)
	if err != nil {
		controller.ResponseError(c, controller.CodeInvalidAuth)
		return
	}
	uid := strconv.FormatInt(userID, 10)

	// 根据用户id和分页参数去redis获取收藏文章ids
	articleIDs, total, err := user.GetArticleIDsByUserID(uid, p)
	if err != nil {
		zap.L().Error("logic.GetArticleIDsByUserID", zap.String("userID", uid), zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	if len(articleIDs) == 0 {
		controller.ResponseSuccess(c, nil)
		return
	}
	// 根据文章ids去mysql 获取文章info
	articleList, err := user.GetCollectedArticle(articleIDs, uid)
	if err != nil {
		zap.L().Error("logic.GetTagDetail(id) failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
		"list":  articleList,
		"total": total,
		"page":  p.Page,
		"size":  p.Size,
	})
}
