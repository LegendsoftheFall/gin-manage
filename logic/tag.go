package logic

import (
	"fmt"
	"manage/dao/mysql"
	"manage/dao/redis"
	"manage/hooks"
	"manage/model"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func GetTrendingTags() (tagList []*model.Tag, err error) {
	return mysql.GetTrendingTags()
}

func SelectTags() (tagList []*model.Tag, err error) {
	return mysql.SelectTags()
}

func GetAllTags(p *model.ParamTagList) (tagList []*model.Tag, total int, err error) {
	tagList, total, err = mysql.GetAllTags(p)
	if err != nil {
		zap.L().Error("mysql.GetAllTags failed",
			zap.Int64("userID", p.CurrentUserID),
			zap.Error(err))
		return nil, 0, err
	}
	for _, tag := range tagList {
		tag.IsFollow, err = mysql.IsFollowTag(p.CurrentUserID, tag.ID)
		if err != nil {
			zap.L().Error("mysql.IsFollowTag failed",
				zap.Int64("userID", p.CurrentUserID),
				zap.Int("tagID", tag.ID),
				zap.Error(err))
			continue
		}
	}
	return
}

func GetSearchTags(p *model.ParamSearch) (tagList []*model.Tag, total int, err error) {
	tagList, total, err = mysql.GetSearchTags(p)
	if err != nil {
		zap.L().Error("mysql.GetAllTags failed",
			zap.Int64("userID", p.CurrentUserID),
			zap.Error(err))
		return nil, 0, err
	}
	for _, tag := range tagList {
		tag.IsFollow, err = mysql.IsFollowTag(p.CurrentUserID, tag.ID)
		if err != nil {
			zap.L().Error("mysql.IsFollowTag failed",
				zap.Int64("userID", p.CurrentUserID),
				zap.Int("tagID", tag.ID),
				zap.Error(err))
			continue
		}
	}
	return
}

func GetFollowingTags(p *model.ParamTagList) (tagList []*model.Tag, total int, err error) {
	tagList, total, err = mysql.GetFollowingTags(p)
	if err != nil {
		zap.L().Error("mysql.GetFollowingTags failed",
			zap.Int64("userID", p.CurrentUserID),
			zap.Error(err))
		return nil, 0, err
	}
	for _, tag := range tagList {
		tag.IsFollow = true
	}
	return
}

func GetTagInfo(p *model.ParamTagInfo) (tagDetail *model.TagDetail, err error) {
	// 获取tag信息
	tagDetail, err = mysql.GetTagDetailByID(p.TagID)
	if err != nil {
		zap.L().Error("mysql.GetTagDetailByID failed",
			zap.Int("tagID", p.TagID),
			zap.Error(err))
		return nil, err
	}
	tagDetail.IsFollow, err = mysql.IsFollowTag(p.UserID, p.TagID)
	fmt.Println(tagDetail.IsFollow)
	return
}

func GetTagDetail(id int, p *model.ParamArticleList) (ArticleList []*model.ApiArticleInfo, total int, err error) {
	// 根据tagID获取文章id列表
	//articleIDs, err := mysql.GetArticleIDByTagID(id)
	//if err != nil {
	//	zap.L().Error("mysql.GetArticleIDByTagID", zap.Error(err))
	//	return nil, 0, err
	//}

	articleIDs, err := redis.GetTagArticleIDInOrder(id, p)
	if err != nil {
		zap.L().Error("redis.GetTagArticleIDInOrder", zap.Error(err))
		return nil, 0, err
	}

	total64, err := redis.GetTagArticleNum(id)
	total = int(total64)
	// 待优化 ↑↑↑
	// 根据文章id列表获取文章info列表
	articleInfoList, err := mysql.GetArticleListByIDs(articleIDs)
	if err != nil {
		zap.L().Error("mysql.GetArticleListByIDs", zap.Error(err))
		return nil, 0, err
	}

	ArticleList = make([]*model.ApiArticleInfo, 0, len(articleInfoList))

	// 遍历文章列表,根据取得的author_id获取用户信息,根据文章id获取tagName
	for _, article := range articleInfoList {
		// 判断用户是否为登录用户 是则获取点赞信息
		var isLiked bool
		var isCollected bool
		if p.CurrentUserID == 0 {
			// 不是登录用户
			isLiked = false
			isCollected = false
		} else {
			isLiked, err = IsUserLiked(strconv.FormatInt(article.ID, 10),
				strconv.FormatInt(p.CurrentUserID, 10))
			if err != nil {
				zap.L().Error("IsUserLiked",
					zap.String("article_id", strconv.FormatInt(article.ID, 10)),
					zap.String("userID", strconv.FormatInt(p.CurrentUserID, 10)),
					zap.Error(err))
				continue
			}
			isCollected, err = IsUserCollected(strconv.FormatInt(article.ID, 10),
				strconv.FormatInt(p.CurrentUserID, 10))
			if err != nil {
				zap.L().Error("IsUserCollected",
					zap.String("article_id", strconv.FormatInt(article.ID, 10)),
					zap.String("userID", strconv.FormatInt(p.CurrentUserID, 10)),
					zap.Error(err))
				continue
			}

		}
		// 根据取得的author_id获取用户信息
		author, err := mysql.GetUserByID(article.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetArticleListByIDs",
				zap.Int64("author_id", article.AuthorID), zap.Error(err))
			continue
		}

		// 根据文章id获取tagName
		tagName, err := mysql.GetTagNameByArticleID(article.ID)
		if err != nil {
			zap.L().Error("mysql.GetTagNameByArticleID",
				zap.Int64("article_id", article.ID), zap.Error(err))
			continue
		}
		// 根据tagName获取id
		tags := make([]*model.TagSimple, 0, len(tagName))
		for _, name := range tagName {
			id, err := mysql.GetTagIDByTagName(name)
			if err != nil {
				zap.L().Error("mysql.GetTagIDByTagName",
					zap.String("tag_name", name), zap.Error(err))
				continue
			}
			tag := &model.TagSimple{
				ID:   id,
				Name: name,
			}
			tags = append(tags, tag)
		}
		article.Tags = tags
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
