package model

import "time"

// Draft 保存与更新草稿
type Draft struct {
	ID         int64     `json:"id,string"`
	AuthorID   int64     `json:"authorID,string"`
	Tags       []int     `json:"tags,omitempty"`
	Title      string    `json:"title"`
	SubTitle   string    `json:"subtitle"`
	Content    string    `json:"content"`
	Html       string    `json:"html"`
	MarkDown   string    `json:"markdown"`
	Image      string    `json:"image"`
	CreateTime time.Time `json:"createTime"`
}

type DraftInfo struct {
	ID         int64     `json:"id,string"`
	AuthorID   int64     `json:"authorID,string"`
	Title      string    `json:"title"`
	SubTitle   string    `json:"subtitle"`
	Content    string    `json:"content"`
	Image      string    `json:"image"`
	Format     string    `json:"format"`
	CreateTime time.Time `json:"createTime" db:"create_time"`
}

type DraftDetail struct {
	Tags       []int        `json:"tags"`
	ID         int64        `json:"id,string"`
	AuthorID   int64        `json:"authorID,string"`
	Title      string       `json:"title"`
	SubTitle   string       `json:"subtitle"`
	Content    string       `json:"content"`
	MarkDown   string       `json:"markdown"`
	Image      string       `json:"image"`
	Format     string       `json:"format"`
	TagSimple  []*TagSimple `json:"tagSimple"`
	CreateTime time.Time    `json:"createTime"`
}
