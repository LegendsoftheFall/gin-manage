package model

import "time"

// Article 文章详情
type Article struct {
	Likes      int       `json:"likes" db:"likes"`
	Comments   int       `json:"comments" db:"comments"`
	ViewCount  int       `json:"viewCount" db:"view_count"`
	ID         int64     `json:"id,string" db:"article_id"`
	AuthorID   int64     `json:"authorID,string" db:"author_id"`
	Tags       []int     `json:"tags,omitempty" binding:"required"`
	DraftID    string    `json:"did"`
	Title      string    `json:"title" db:"title" binding:"required"`
	SubTitle   string    `json:"subtitle" db:"subtitle"`
	Content    string    `json:"content" db:"content"`
	Html       string    `json:"html" db:"html" binding:"required"`
	MarkDown   string    `json:"markdown" db:"markdown" binding:"required"`
	Image      string    `json:"image" db:"image"`
	Source     string    `json:"source" db:"source"`
	Format     string    `json:"format"`
	CreateTime time.Time `json:"createTime" db:"create_time"`
}

type ArticleInfo struct {
	Likes       int          `json:"likes" db:"likes"`
	Comments    int          `json:"comments" db:"comments"`
	ViewCount   int          `json:"viewCount" db:"view_count"`
	ID          int64        `json:"id,string" db:"article_id"`
	AuthorID    int64        `json:"authorID,string" db:"author_id"`
	Title       string       `json:"title" db:"title"`
	SubTitle    string       `json:"subtitle" db:"subtitle"`
	Content     string       `json:"content" db:"content"`
	Image       string       `json:"image" db:"image"`
	Format      string       `json:"format"`
	IsLiked     bool         `json:"isLiked"`
	IsCollected bool         `json:"isCollected"`
	Tags        []*TagSimple `json:"tags"`
	CreateTime  time.Time    `json:"createTime" db:"create_time"`
}

// ApiArticleInfo 获取文章信息
type ApiArticleInfo struct {
	AuthorName   string `json:"authorName"`
	Avatar       string `json:"avatar"`
	*ArticleInfo `json:"articleInfo"`
}
