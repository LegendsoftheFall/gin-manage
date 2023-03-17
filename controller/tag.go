package controller

import (
	"manage/logic"
	"manage/model"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// TrendingTagHandler 趋势标签请求的处理函数
func TrendingTagHandler(c *gin.Context) {
	//查询排名前六的标签并以列表方式返回
	data, err := logic.GetTrendingTags()
	if err != nil {
		zap.L().Error("logic.GetTrendingTags failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// SelectTagsHandler 获取选择标签请求的处理函数
func SelectTagsHandler(c *gin.Context) {
	data, err := logic.SelectTags()
	if err != nil {
		zap.L().Error("logic.GetTrendingTags failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// TagInfoHandler 标签信息请求的处理函数
func TagInfoHandler(c *gin.Context) {
	p := new(model.ParamTagInfo)
	//获取标签id
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("TagInfoHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	tagDetail, err := logic.GetTagInfo(p)
	if err != nil {
		zap.L().Error("logic.GetTagInfo(id) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"tag": tagDetail,
	})
}

// TagsHandler 获取全部标签请求的处理函数
func TagsHandler(c *gin.Context) {
	p := &model.ParamTagList{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
	}
	//获取分页参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("TagsHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, total, err := logic.GetAllTags(p)
	if err != nil {
		zap.L().Error("logic.GetTrendingTags failed", zap.Error(err))
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

// FollowingTagsHandler 获取已关注标签请求的处理函数
func FollowingTagsHandler(c *gin.Context) {
	p := &model.ParamTagList{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
	}
	//获取分页参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("FollowingTagsHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	data, total, err := logic.GetFollowingTags(p)
	if err != nil {
		zap.L().Error("logic.GetFollowingTags failed", zap.Error(err))
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

// TagDetailHandler 标签文章请求的处理函数
func TagDetailHandler(c *gin.Context) {
	p := &model.ParamArticleList{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
		Order:         model.OrderScore,
	}
	//获取分页参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("TagDetailHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//获取标签id
	tagID := c.Param("id")
	id, err := strconv.Atoi(tagID)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 根据标签ID和分页数据查询文章列表
	articleList, total, err := logic.GetTagDetail(id, p)
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
