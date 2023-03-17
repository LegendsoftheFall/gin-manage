package logic

import (
	"manage/dao/mysql"
	"manage/dao/redis"
	"manage/hooks"
	"manage/model"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func IsUserCollected(aid, uid string) (isLiked bool, err error) {
	return redis.IsUserCollected(aid, uid)
}

func CollectArticle(aid, uid string) (err error) {
	isCollected, err := redis.IsUserCollected(aid, uid)
	if err != nil {
		zap.L().Debug("redis.IsUserCollected",
			zap.String("articleID", aid),
			zap.String("userID", uid))
		return
	}
	if !isCollected {
		if err = redis.CreateArticleCollected(aid, uid); err != nil {
			zap.L().Error("redis.CreateArticleCollected",
				zap.String("articleID", aid),
				zap.String("userID", uid),
				zap.Error(err))
		}
	} else {
		if err = redis.RemoveArticleCollected(aid, uid); err != nil {
			zap.L().Error("redis.CreateArticleCollected",
				zap.String("articleID", aid),
				zap.String("userID", uid),
				zap.Error(err))
		}
	}
	return
}

func GetArticleIDsByUserID(uid string, p *model.ParamPage) (ids []string, total int64, err error) {
	return redis.GetArticleIDsByUserID(uid, p)
}

func GetCollectedArticle(ids []string, uid string) (ArticleList []*model.ApiArticleInfo, err error) {
	// 根据文章id列表获取文章info列表
	articleInfoList, err := mysql.GetCollectedArticleByIDs(ids)
	if err != nil {
		zap.L().Error("mysql.GetCollectedArticleByIDs", zap.Error(err))
		return nil, err
	}

	ArticleList = make([]*model.ApiArticleInfo, 0, len(articleInfoList))

	// 遍历文章列表,根据取得的author_id获取用户信息,根据文章id获取tagName
	for _, article := range articleInfoList {
		// 判断用户是否为登录用户 是则获取点赞信息
		var isLiked bool
		var isCollected bool
		isLiked, err = IsUserLiked(strconv.FormatInt(article.ID, 10), uid)
		if err != nil {
			zap.L().Error("IsUserLiked",
				zap.String("article_id", strconv.FormatInt(article.ID, 10)),
				zap.String("userID", uid),
				zap.Error(err))
			continue
		}
		isCollected, err = IsUserCollected(strconv.FormatInt(article.ID, 10), uid)
		if err != nil {
			zap.L().Error("IsUserCollected",
				zap.String("article_id", strconv.FormatInt(article.ID, 10)),
				zap.String("userID", uid),
				zap.Error(err))
			continue
		}

		// 根据取得的author_id获取用户信息
		author, err2 := mysql.GetUserByID(article.AuthorID)
		if err2 != nil {
			zap.L().Error("mysql.GetArticleListByIDs",
				zap.Int64("author_id", article.AuthorID), zap.Error(err2))
			continue
		}
		// 日期转换
		article.Format = hooks.TimeSub(time.Now(), article.CreateTime)
		// 裁剪字符串
		content := []rune(article.Content)
		if len(content) > 270 {
			article.Content = string(content[:270])
		}
		// 是否被当前用户点赞收藏
		article.IsLiked = isLiked
		article.IsCollected = isCollected
		//赋值拼接
		articleInfo := &model.ApiArticleInfo{
			AuthorName:  author.Username,
			Avatar:      author.Avatar,
			ArticleInfo: article,
		}
		ArticleList = append(ArticleList, articleInfo)
	}
	return
}
