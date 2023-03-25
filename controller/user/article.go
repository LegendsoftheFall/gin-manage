package user

import (
	"manage/controller"
	user2 "manage/dao/mysql/user"
	"manage/logic/user"
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
		controller.ValidateError(c, err)
		return
	}

	userID, _ := controller.GetCurrentUserID(c)
	p.AuthorID = userID
	//获取草稿ID
	draftID := p.DraftID
	//创建文章
	id, err := user.CreateArticle(p, strconv.FormatInt(userID, 10), draftID)
	if err != nil {
		zap.L().Error("logic.CreateArticle(p) failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	//返回响应
	controller.ResponseSuccess(c, gin.H{
		"id": strconv.FormatInt(id, 10),
	})

}

// EditArticleHandler 编辑文章请求的处理函数
func EditArticleHandler(c *gin.Context) {
	p := new(model.Article)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("EditArticle with invalid param", zap.Error(err))
		controller.ValidateError(c, err)
		return
	}

	// 判断是否为作者
	userID, _ := controller.GetCurrentUserID(c)
	if userID != p.AuthorID {
		controller.ResponseError(c, controller.CodeInvalidAuth)
		return
	}

	err := user.EditArticle(p)
	if err != nil {
		zap.L().Error("logic.EditArticle(p) failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}

	controller.ResponseSuccess(c, gin.H{
		"id": strconv.FormatInt(p.ID, 10),
	})
}

// DeleteArticleHandler 删除文章请求的处理函数
func DeleteArticleHandler(c *gin.Context) {
	articleID := c.Param("id")
	id, ok := strconv.ParseInt(articleID, 10, 64)
	if ok != nil {
		zap.L().Error("articleID is invalid", zap.Int64("articleID", id), zap.Error(ok))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}

	// 判断是否为作者
	userID, _ := controller.GetCurrentUserID(c)
	uid, _ := user2.GetUserIDByArticleID(id)
	if userID != uid {
		controller.ResponseError(c, controller.CodeInvalidAuth)
		return
	}

	if err := user.DeleteArticle(id); err != nil {
		zap.L().Error("logic.DeleteArticle(id) failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
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
		controller.ResponseError(c, controller.CodeInvalidParam)
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
		isLiked, ok = user.IsUserLiked(articleID, strconv.FormatInt(currentUserID, 10))
		if ok != nil {
			zap.L().Error("logic.IsUserLiked",
				zap.String("articleID", articleID),
				zap.Int64("currentID", currentUserID),
				zap.Error(ok))
		}
		isCollected, ok = user.IsUserCollected(articleID, strconv.FormatInt(currentUserID, 10))
		if ok != nil {
			zap.L().Error("logic.IsUserCollected",
				zap.String("articleID", articleID),
				zap.Int64("currentID", currentUserID),
				zap.Error(ok))
		}
	}
	// 根据文章ID获取用户信息
	userInfo, err := user.GetUserInfoByArticleID(id, currentUserID)
	if err != nil {
		zap.L().Error("logic.GetUserInfoByArticleID is failed", zap.Int64("userID", id), zap.Error(err))
		controller.ResponseError(c, controller.CodeUserNotExist)
		return
	}
	// 根据文章ID获取文章详情
	articleDetail, tags, err := user.GetArticleDetailByID(id)
	if err != nil {
		zap.L().Error("logic.GetArticleDetailByID is failed", zap.Int64("articleID", id), zap.Error(ok))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	// 更新文章阅读数
	if err = user.UpdateViewCount(id); err != nil {
		zap.L().Error("logic.UpdateViewCount is failed", zap.Int64("articleID", id), zap.Error(err))
	}
	controller.ResponseSuccess(c, gin.H{
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
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}

	// 根据分页数据和排序参数查询文章列表
	articleList, total, err := user.GetArticleList(p)
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
