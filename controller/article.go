package controller

import (
	"manage/dao/mysql"
	"manage/logic"
	"manage/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateArticleHandler 新建文章请求的处理函数
func CreateArticleHandler(c *gin.Context) {
	// 参数校验
	p := new(model.Article)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreateArticle with invalid param", zap.Error(err))
		validateError(c, err)
		return
	}

	userID, _ := getCurrentUserID(c)
	p.AuthorID = userID
	//获取草稿ID
	draftID := p.DraftID
	//创建文章
	id, err := logic.CreateArticle(p, strconv.FormatInt(userID, 10), draftID)
	if err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//返回响应
	ResponseSuccess(c, gin.H{
		"id": strconv.FormatInt(id, 10),
	})

}

// EditArticleHandler 编辑文章请求的处理函数
func EditArticleHandler(c *gin.Context) {
	p := new(model.Article)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("EditArticle with invalid param", zap.Error(err))
		validateError(c, err)
		return
	}

	// 判断是否为作者
	userID, _ := getCurrentUserID(c)
	if userID != p.AuthorID {
		ResponseError(c, CodeInvalidAuth)
		return
	}

	err := logic.EditArticle(p)
	if err != nil {
		zap.L().Error("logic.EditArticle(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, gin.H{
		"id": strconv.FormatInt(p.ID, 10),
	})
}

// DeleteArticleHandler 删除文章请求的处理函数
func DeleteArticleHandler(c *gin.Context) {
	articleID := c.Param("id")
	id, ok := strconv.ParseInt(articleID, 10, 64)
	if ok != nil {
		zap.L().Error("articleID is invalid", zap.Int64("articleID", id), zap.Error(ok))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 判断是否为作者
	userID, _ := getCurrentUserID(c)
	uid, _ := mysql.GetUserIDByArticleID(id)
	if userID != uid {
		ResponseError(c, CodeInvalidAuth)
		return
	}

	if err := logic.DeleteArticle(id); err != nil {
		zap.L().Error("logic.DeleteArticle(id) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"id": uid,
	})
}

// ArticleHandler 文章详情请求的处理函数
func ArticleHandler(c *gin.Context) {
	// 获取文章ID
	articleID := c.Param("id")
	id, ok := strconv.ParseInt(articleID, 10, 64)
	if ok != nil {
		zap.L().Error("articleID is invalid", zap.Int64("articleID", id), zap.Error(ok))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 是否被当前用户点赞
	uid := c.Query("uid")
	var isLiked bool
	var isCollected bool
	currentUserID, ok := strconv.ParseInt(uid, 10, 64)
	if ok != nil {
		isLiked = false
		isCollected = false
	} else {
		isLiked, ok = logic.IsUserLiked(articleID, strconv.FormatInt(currentUserID, 10))
		if ok != nil {
			zap.L().Error("logic.IsUserLiked",
				zap.String("articleID", articleID),
				zap.Int64("currentID", currentUserID),
				zap.Error(ok))
		}
		isCollected, ok = logic.IsUserCollected(articleID, strconv.FormatInt(currentUserID, 10))
		if ok != nil {
			zap.L().Error("logic.IsUserCollected",
				zap.String("articleID", articleID),
				zap.Int64("currentID", currentUserID),
				zap.Error(ok))
		}
	}
	// 根据文章ID获取用户信息
	userInfo, err := logic.GetUserInfoByArticleID(id, currentUserID)
	if err != nil {
		zap.L().Error("logic.GetUserInfoByArticleID is failed", zap.Int64("userID", id), zap.Error(err))
		ResponseError(c, CodeUserNotExist)
		return
	}
	// 根据文章ID获取文章详情
	articleDetail, tags, err := logic.GetArticleDetailByID(id)
	if err != nil {
		zap.L().Error("logic.GetArticleDetailByID is failed", zap.Int64("articleID", id), zap.Error(ok))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 更新文章阅读数
	if err = logic.UpdateViewCount(id); err != nil {
		zap.L().Error("logic.UpdateViewCount is failed", zap.Int64("articleID", id), zap.Error(err))
	}
	ResponseSuccess(c, gin.H{
		"userInfo":    userInfo,
		"article":     articleDetail,
		"tags":        tags,
		"isLiked":     isLiked,
		"isCollected": isCollected,
	})
}

// ArticleListHandler 所有文章请求的处理函数
func ArticleListHandler(c *gin.Context) {
	p := &model.ParamArticleList{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
		Order:         model.OrderScore,
	}
	//获取分页参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("ArticleListHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 根据分页数据和排序参数查询文章列表
	articleList, total, err := logic.GetArticleList(p)
	if err != nil {
		zap.L().Error("logic.GetArticleList failed", zap.Error(err))
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
