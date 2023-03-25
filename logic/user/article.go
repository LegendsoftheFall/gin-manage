package user

import (
	"manage/dao/mysql/user"
	"manage/dao/redis"
	"manage/hooks"
	"manage/model"
	"manage/pkg/snowflake"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func CreateArticle(p *model.Article, uid, did string) (id int64, err error) {
	p.ID = snowflake.GenID()
	p.Content = hooks.TrimHtml(p.Html)
	id = p.ID
	if err = redis.CreateArticle(p); err != nil {
		return -1, err
	}
	if err = redis.DeleteDraft(did, uid); err != nil {
		return -1, err
	}
	err = user.CreateArticle(p)
	return
}

func EditArticle(p *model.Article) (err error) {
	p.Content = hooks.TrimHtml(p.Html)
	// 查询旧标签
	tagIDs, err := user.GetTagIDsByArticleID(p.ID)
	if err != nil {
		return
	}
	// 删除所有标签
	if err = user.DeleteTagByID(p.ID); err != nil {
		return
	}
	// 更新redis中的标签文章set
	if err = redis.EditArticle(strconv.FormatInt(p.ID, 10), tagIDs, p.Tags); err != nil {
		return
	}
	return user.UpdateArticle(p)
}

func DeleteArticle(id int64) (err error) {
	// 查询旧标签
	tagIDs, err := user.GetTagIDsByArticleID(id)
	if err != nil {
		return
	}
	// 删除所有标签
	if err = user.DeleteTagByID(id); err != nil {
		return
	}
	// 删除文章相关的点赞和收藏记录
	if err = redis.DeleteArticleCollectAndLike(strconv.FormatInt(id, 10)); err != nil {
		return
	}
	// 删除文章的分数和时间
	if err = redis.DeleteArticle(strconv.FormatInt(id, 10), tagIDs); err != nil {
		return
	}
	return user.DeleteArticle(id)
}

// GetArticleByID 根据用户ID获取用户主页文章列表
func GetArticleByID(id, page, size int64) (articleList []*model.ArticleInfo, total int, err error) {
	articleList, err = user.GetArticleListByID(id, page, size)
	if len(articleList) == 0 {
		return nil, 0, err
	}
	total, err = user.GetArticleNumByID(id)
	for _, article := range articleList {
		article.Format = hooks.TimeSub(time.Now(), article.CreateTime)
		// 裁剪字符串
		content := []rune(article.Content)
		if len(content) > 270 {
			article.Content = string(content[:270])
		}
	}
	return
}

// GetAllArticleByID 根据用户ID获取用户主页文章列表
func GetAllArticleByID(id int64) (articleList []*model.ArticleInfo, total int, err error) {
	articleList, err = user.GetAllArticleListByID(id)
	if len(articleList) == 0 {
		return nil, 0, err
	}
	total, err = user.GetArticleNumByID(id)
	for _, article := range articleList {
		article.Format = hooks.TimeSub(time.Now(), article.CreateTime)
		// 裁剪字符串
		content := []rune(article.Content)
		if len(content) > 270 {
			article.Content = string(content[:270])
		}
	}
	return
}

func GetUserIDByArticleID(id int64) (userID int64, err error) {
	return user.GetUserIDByArticleID(id)
}

func UpdateViewCount(id int64) (err error) {
	return user.UpdateViewCount(id)
}

func GetArticleDetailByID(id int64) (article model.Article, tags []*model.TagSimple, err error) {
	article, err = user.GetArticleDetailByID(id)
	article.Format = hooks.Date(article.CreateTime)

	// 标签信息
	// 根据文章id获取tagName
	tagName, err := user.GetTagNameByArticleID(id)
	if err != nil {
		zap.L().Error("mysql.GetTagNameByArticleID",
			zap.Int64("article_id", id), zap.Error(err))
	}
	// 根据tagName获取id
	tags = make([]*model.TagSimple, 0, len(tagName))
	for _, name := range tagName {
		id, err := user.GetTagIDByTagName(name)
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

	return
}

func GetArticleList(p *model.ParamArticleList) (ArticleList []*model.ApiArticleInfo, total int, err error) {
	// TODO : 优化按时间排序由于新插入数据导致分页数据出错的问题
	articleIDs, err := redis.GetArticleIDInOrder(p)
	if err != nil {
		zap.L().Error("redis.GetArticleIDInOrder", zap.Error(err))
		return nil, 0, err
	}

	total64, err := redis.GetArticleNum()
	total = int(total64)
	// 根据文章id列表获取文章info列表
	articleInfoList, err := user.GetArticleListByIDs(articleIDs)
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
		author, err := user.GetUserByID(article.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetArticleListByIDs",
				zap.Int64("author_id", article.AuthorID), zap.Error(err))
			continue
		}

		// 根据文章id获取tagName
		tagName, err := user.GetTagNameByArticleID(article.ID)
		if err != nil {
			zap.L().Error("mysql.GetTagNameByArticleID",
				zap.Int64("article_id", article.ID), zap.Error(err))
			continue
		}
		// 根据tagName获取id
		tags := make([]*model.TagSimple, 0, len(tagName))
		for _, name := range tagName {
			id, err := user.GetTagIDByTagName(name)
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

func GetSearchArticleList(p *model.ParamSearch) (ArticleList []*model.ApiArticleInfo, total int, err error) {
	// 根据参数获取文章info列表
	articleInfoList := make([]*model.ArticleInfo, 0)
	if p.Category == model.SearchTop {
		articleInfoList, total, err = user.GetSearchArticleListByView(p)
		if err != nil {
			zap.L().Error("mysql.GetSearchArticleListByView", zap.Error(err))
			return nil, 0, err
		}
	} else if p.Category == model.SearchLatest {
		articleInfoList, total, err = user.GetSearchArticleListByTime(p)
		if err != nil {
			zap.L().Error("mysql.GetSearchArticleListByTime", zap.Error(err))
			return nil, 0, err
		}
	}

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
		author, err := user.GetUserByID(article.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetArticleListByIDs",
				zap.Int64("author_id", article.AuthorID), zap.Error(err))
			continue
		}

		// 根据文章id获取tagName
		tagName, err := user.GetTagNameByArticleID(article.ID)
		if err != nil {
			zap.L().Error("mysql.GetTagNameByArticleID",
				zap.Int64("article_id", article.ID), zap.Error(err))
			continue
		}
		// 根据tagName获取id
		tags := make([]*model.TagSimple, 0, len(tagName))
		for _, name := range tagName {
			id, err := user.GetTagIDByTagName(name)
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
