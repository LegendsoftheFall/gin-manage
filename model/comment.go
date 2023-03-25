package model

import "time"

type TreePath struct {
	Ancestor   int64 `db:"ancestor"`
	Descendant int64 `db:"descendant"`
	Distance   int   `db:"distance"`
}

type RootComment struct {
	IsLike    bool      `json:"is_like"`
	IsClick   bool      `json:"is_click"`
	ItemType  int       `json:"item_type" db:"item_type"`
	Status    int       `json:"comment_status" db:"status"`
	Likes     int       `json:"likes" db:"likes"`
	CommentID int64     `json:"comment_id,string" db:"comment_id"`
	UserID    int64     `json:"user_id,string" db:"user_id"`
	ItemID    int64     `json:"item_id,string" db:"item_id"`
	Content   string    `json:"content" db:"comment_content"`
	Picture   string    `json:"picture" db:"comment_picture"`
	Format    string    `json:"format"`
	Binding   string    `json:"binding"`
	CreatTime time.Time `json:"create_time" db:"create_time"`
}

type ReplyComment struct {
	IsLike         bool      `json:"is_like"`
	IsClick        bool      `json:"is_click"`
	Level          int       `json:"level"`
	ItemType       int       `json:"item_type" db:"item_type"`
	Likes          int       `json:"likes" db:"likes"`
	Status         int       `json:"reply_status" db:"status"`
	ReplyCommentID int64     `json:"reply_comment_id"`
	ReplyID        int64     `json:"reply_id" db:"comment_id"`
	ToReplyID      int64     `json:"reply_to_reply_id" db:"ancestor"`
	UserID         int64     `json:"user_id" db:"user_id"`
	ToUserID       int64     `json:"reply_to_user_id" db:"user_id"`
	ItemID         int64     `json:"item_id" db:"item_id"`
	Content        string    `json:"reply_content" db:"comment_content"`
	Picture        string    `json:"reply_picture" db:"comment_picture"`
	Format         string    `json:"format"`
	Binding        string    `json:"binding"`
	CreatTime      time.Time `json:"create_time" db:"create_time"`
}

type RootCommentInfo struct {
	ReplyCount int `json:"reply_count"`
	*RootComment
	CommentReplies []*ReplyComment `json:"comment_replies"`
}

type ReplyCommentInfo struct {
	ReplyID     int64         `json:"reply_id"`
	ReplyInfo   *ReplyComment `json:"reply_info"`
	ParentReply *ReplyComment `json:"parent_reply"`
	UserInfo    *UserInfo     `json:"user_info"`
	ReplyUser   *UserInfo     `json:"reply_user"`
}

type ApiComment struct {
	CommentID   int64               `json:"comment_id"`
	CommentInfo *RootCommentInfo    `json:"comment_info"`
	UserInfo    *UserInfo           `json:"user_info"`
	ReplyInfos  []*ReplyCommentInfo `json:"reply_infos"`
}

type Comment struct {
	Status    int    `json:"status" db:"status"`
	ItemType  int    `json:"item_type" db:"item_type"`
	CommentID int64  `json:"comment_id" db:"comment_id"`
	UserID    int64  `json:"user_id" db:"user_id"`
	ItemID    int64  `json:"item_id" db:"item_id"`
	UserName  string `json:"user_name" db:"username"`
	Content   string `json:"comment_content" db:"comment_content"`
}
