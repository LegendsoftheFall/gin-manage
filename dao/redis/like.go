package redis

const ScorePreLike = 432 * 2 // 每个赞的分数

func IsUserLike(aid, uid string) (isLiked bool, err error) {
	isLiked, err = rdb.SIsMember(getRedisKey(KeyArticleLike+aid), uid).Result()
	return
}

func CreateArticleLike(aid, uid string) (err error) {
	pipeline := rdb.TxPipeline()
	pipeline.SAdd(getRedisKey(KeyArticleLike+aid), uid)
	pipeline.ZIncrBy(getRedisKey(KeyArticleScore), ScorePreLike, aid)
	_, err = pipeline.Exec()
	return
}

func RemoveArticleLike(aid, uid string) (err error) {
	pipeline := rdb.TxPipeline()
	pipeline.SRem(getRedisKey(KeyArticleLike+aid), uid)
	pipeline.ZIncrBy(getRedisKey(KeyArticleScore), -ScorePreLike, aid)
	_, err = pipeline.Exec()
	return
}
