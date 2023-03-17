package redis

import (
	"manage/model"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

func CreateArticle(p *model.Article) (err error) {
	pipeline := rdb.TxPipeline()
	// 创建文章时间
	pipeline.ZAdd(getRedisKey(KeyArticleTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: p.ID,
	})
	// 文章分数
	pipeline.ZAdd(getRedisKey(KeyArticleScore), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: p.ID,
	})
	// 标签
	for _, tag := range p.Tags {
		pipeline.SAdd(getRedisKey(KeyTag+strconv.Itoa(tag)), p.ID)
	}
	_, err = pipeline.Exec()
	return
}

func DeleteArticle(id string, tIDs []int) (err error) {
	pipeline := rdb.TxPipeline()
	pipeline.ZRem(getRedisKey(KeyArticleTime), id)
	pipeline.ZRem(getRedisKey(KeyArticleScore), id)
	for _, tid := range tIDs {
		pipeline.SRem(getRedisKey(KeyTag+strconv.Itoa(tid)), id)
	}
	_, err = pipeline.Exec()
	return
}

func EditArticle(id string, oldIDs, newIDs []int) (err error) {
	pipeline := rdb.TxPipeline()
	for _, oldID := range oldIDs {
		pipeline.SRem(getRedisKey(KeyTag+strconv.Itoa(oldID)), id)
	}
	for _, newID := range newIDs {
		pipeline.SAdd(getRedisKey(KeyTag+strconv.Itoa(newID)), id)
	}
	_, err = pipeline.Exec()
	return
}

func getIDsFromKey(key string, page, size int64) ([]string, error) {
	//确认索引的起点和终点
	start := (page - 1) * size
	stop := start + size - 1
	//按分数从大到小查询指定元素
	return rdb.ZRevRange(key, start, stop).Result()
}

func GetArticleIDInOrder(p *model.ParamArticleList) ([]string, error) {
	OrderKey := getRedisKey(KeyArticleScore)
	if p.Order == model.OrderTime {
		OrderKey = getRedisKey(KeyArticleTime)
	}
	return getIDsFromKey(OrderKey, p.Page, p.Size)
}

func GetTagArticleIDInOrder(id int, p *model.ParamArticleList) ([]string, error) {
	//使用ZINTERStore 把标签的文章set与文章分数的ZSet生成一个新的ZSet,对新的ZSet取数据
	tKey := getRedisKey(KeyTag + strconv.Itoa(id)) // tag: ID
	OrderKey := getRedisKey(KeyArticleScore)
	if p.Order == model.OrderTime {
		OrderKey = getRedisKey(KeyArticleTime)
	}
	//利用缓存key减少ZINTERStore执行的次数
	newKey := OrderKey + strconv.Itoa(id) // article:post ID or article:time ID
	if rdb.Exists(newKey).Val() < 1 {
		//不存在需要计算key
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(newKey, redis.ZStore{
			Aggregate: "max",
		}, tKey, OrderKey)
		pipeline.Expire(newKey, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	//存在根据key查询ids
	return getIDsFromKey(newKey, p.Page, p.Size)
}

func GetTagArticleNum(id int) (int64, error) {
	return rdb.SCard(getRedisKey(KeyTag + strconv.Itoa(id))).Result()
}

func GetArticleNum() (int64, error) {
	return rdb.ZCard(getRedisKey(KeyArticleScore)).Result()
}
