package redis

// redis key
const (
	KeyPre            = "manage:"
	KeyArticleTime    = "article:time"     //zSet 文章创建时间
	KeyArticleScore   = "article:score"    // zSet 文章分数
	KeyArticleLike    = "article:like:"    // set 文章和点赞用户, 参数是articleID
	KeyArticleCollect = "article:collect:" // set 文章和收藏用户, 参数是articleID
	KeyUserCollect    = "user:collect:"    // zSet 用户和收藏文章, 参数是userID 分数为时间
	KeyUserDraft      = "user:draft:"      // set 用户和草稿, 参数是userID
	KeyDraft          = "draft:"           // hash 草稿和草稿内容, 参数是draftID
	KeyTag            = "tag:"             // set 标签和文章, 参数是tagID
)

func getRedisKey(key string) string {
	return KeyPre + key
}
