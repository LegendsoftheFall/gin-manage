package redis

import (
	"manage/model"
	"time"

	"github.com/go-redis/redis"
)

func IsUserCollected(aid, uid string) (isCollected bool, err error) {
	isCollected, err = rdb.SIsMember(getRedisKey(KeyArticleCollect+aid), uid).Result()
	return
}

func CreateArticleCollected(aid, uid string) (err error) {
	pipeline := rdb.TxPipeline()
	pipeline.SAdd(getRedisKey(KeyArticleCollect+aid), uid)
	// 将用户点赞的文章按时间排序
	pipeline.ZAdd(getRedisKey(KeyUserCollect+uid), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: aid,
	})
	_, err = pipeline.Exec()
	return
}

func RemoveArticleCollected(aid, uid string) (err error) {
	pipeline := rdb.TxPipeline()
	pipeline.SRem(getRedisKey(KeyArticleCollect+aid), uid)
	pipeline.ZRem(getRedisKey(KeyUserCollect+uid), aid)
	_, err = pipeline.Exec()
	return
}

func GetArticleIDsByUserID(uid string, p *model.ParamPage) (ids []string, total int64, err error) {
	key := getRedisKey(KeyUserCollect + uid)
	start := (p.Page - 1) * p.Size
	stop := start + p.Size - 1
	total = rdb.ZCard(key).Val()
	ids, err = rdb.ZRevRange(key, start, stop).Result()
	return
}

func DeleteArticleCollectAndLike(aid string) (err error) {
	ids := rdb.SMembers(getRedisKey(KeyArticleCollect + aid)).Val()
	pipeline := rdb.TxPipeline()
	// 查询出收藏文章的所有用户
	//ids := pipeline.SMembers(getRedisKey(KeyArticleCollect + aid)).Val()
	// 遍历所有用户,移除zSet中的文章id
	for _, id := range ids {
		pipeline.ZRem(getRedisKey(KeyUserCollect+id), aid)
	}
	// 删除set中的文章id
	pipeline.Del(getRedisKey(KeyArticleCollect + aid))
	pipeline.Del(getRedisKey(KeyArticleLike + aid))
	_, err = pipeline.Exec()
	return
}
