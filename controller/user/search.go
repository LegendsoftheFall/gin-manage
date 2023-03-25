package user

import (
	"manage/controller"
	"manage/logic/user"
	"manage/model"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SearchArticleHandler 搜索文章的处理函数
func SearchArticleHandler(c *gin.Context) {
	p := &model.ParamSearch{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
		Category:      model.SearchTop,
		Key:           "",
	}

	//获取参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("SearchArticleHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}

	// 根据分页数据和排序参数查询文章列表
	articleList, total, err := user.GetSearchArticleList(p)
	if err != nil {
		zap.L().Error("logic.GetArticleList failed", zap.Error(err))
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

// SearchTagHandler 搜索标签的处理函数
func SearchTagHandler(c *gin.Context) {
	p := &model.ParamSearch{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
		Category:      model.SearchTag,
		Key:           "",
	}

	//获取参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("SearchTagHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}

	data, total, err := user.GetSearchTags(p)
	if err != nil {
		zap.L().Error("logic.GetTrendingTags failed", zap.Error(err))
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

// SearchUserHandler 搜索用户的处理函数
func SearchUserHandler(c *gin.Context) {
	p := &model.ParamSearch{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
		Category:      model.SearchUser,
		Key:           "",
	}

	//获取参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("SearchUserHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}

	data, total, err := user.GetSearchUsers(p)
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
