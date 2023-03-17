package logic

import (
	"manage/dao/mysql"
	"manage/dao/redis"
	"strconv"

	"go.uber.org/zap"
)

func IsUserLiked(aid, uid string) (isLiked bool, err error) {
	return redis.IsUserLike(aid, uid)
}

func LikeArticle(aid, uid string) (err error) {
	// 1. 判断当前用户是否点过赞
	isLiked, err := redis.IsUserLike(aid, uid)
	if err != nil {
		zap.L().Debug("redis.IsUserLike",
			zap.String("articleID", aid),
			zap.String("userID", uid))
		return
	}
	articleID, _ := strconv.ParseInt(aid, 10, 64)
	if !isLiked {
		// 2. 若未点赞,可以点赞
		//2.1 数据库点赞数+1
		if err = mysql.AddLikeNumByArticleID(articleID); err != nil {
			zap.L().Error("mysql.AddLikeNumByArticleID",
				zap.Int64("articleID", articleID),
				zap.Error(err))
			return
		}
		//2.2 保存用户到redis set集合
		if err = redis.CreateArticleLike(aid, uid); err != nil {
			zap.L().Error("redis.CreateArticleLike",
				zap.String("articleID", aid),
				zap.String("userID", uid),
				zap.Error(err))
		}
	} else {
		// 3. 若已点赞,取消点赞
		//3.1 数据库点赞数-1
		if err = mysql.RemoveLikeNumByArticleID(articleID); err != nil {
			zap.L().Error("mysql.RemoveLikeNumByArticleID",
				zap.Int64("articleID", articleID),
				zap.Error(err))
			return
		}
		//3.2 将用户到从edis set集合移除
		if err = redis.RemoveArticleLike(aid, uid); err != nil {
			zap.L().Error("redis.RemoveArticleLike",
				zap.String("articleID", aid),
				zap.String("userID", uid),
				zap.Error(err))
		}
	}

	return
}
