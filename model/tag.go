package model

type Tag struct {
	ID            int    `json:"id,string" db:"tag_id"`
	ArticleNumber int    `json:"num,string" db:"article_number"`
	Name          string `json:"name" db:"tag_name"`
	Image         string `json:"image" db:"image"`
	IsFollow      bool   `json:"isFollow"`
}

type TagSimple struct {
	ID   int    `json:"id,string" db:"tag_id"`
	Name string `json:"name" db:"tag_name"`
}

type TagDetail struct {
	ID             int    `json:"id,string" db:"tag_id"`
	ArticleNumber  int    `json:"articleNum" db:"article_number"`
	FollowerNumber int    `json:"followerNum" db:"follower_number"`
	Name           string `json:"name" db:"tag_name"`
	Image          string `json:"image" db:"image"`
	Introduction   string `json:"introduction,omitempty" db:"introduction"`
	IsFollow       bool   `json:"isFollow"`
}
