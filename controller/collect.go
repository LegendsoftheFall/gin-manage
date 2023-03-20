package controller

import (
	"manage/logic"
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
		ResponseError(c, CodeInvalidParam)
		return
	}
	if err := logic.CollectArticle(p.ArticleID, p.UserID); err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
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
		ResponseError(c, CodeInvalidParam)
		return
	}

	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeInvalidAuth)
		return
	}
	uid := strconv.FormatInt(userID, 10)

	// 根据用户id和分页参数去redis获取收藏文章ids
	articleIDs, total, err := logic.GetArticleIDsByUserID(uid, p)
	if err != nil {
		zap.L().Error("logic.GetArticleIDsByUserID", zap.String("userID", uid), zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	if len(articleIDs) == 0 {
		ResponseSuccess(c, nil)
		return
	}
	// 根据文章ids去mysql 获取文章info
	articleList, err := logic.GetCollectedArticle(articleIDs, uid)
	if err != nil {
		zap.L().Error("logic.GetTagDetail(id) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"list":  articleList,
		"total": total,
		"page":  p.Page,
		"size":  p.Size,
	})
}
