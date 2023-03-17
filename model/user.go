package model

import "time"

type User struct {
	UserID       int64  `db:"user_id"`
	Email        string `db:"email"`
	Username     string `db:"username"`
	Password     string `db:"password"`
	Avatar       string `db:"avatar,omitempty"`
	AccessToken  string
	RefreshToken string
}

type UserInfo struct {
	UserID       int64  `json:"userID,string" db:"user_id"`
	Follower     int    `json:"follower" db:"follower"`
	Following    int    `json:"following" db:"following"`
	UserName     string `json:"username" db:"username"`
	Avatar       string `json:"avatar" db:"avatar"`
	Email        string `json:"email" db:"email"`
	Introduction string `json:"introduction" db:"introduction"`
	HomePage     string `json:"homePage" db:"homepage"`
	Github       string `json:"github" db:"github"`
	Position     string `json:"position" db:"position"`
	IsFollow     bool   `json:"isFollow"`
}

type UserProfile struct {
	UserName     string `json:"username" db:"username"`
	Avatar       string `json:"avatar" db:"avatar"`
	Location     string `json:"location" db:"location"`
	Company      string `json:"company" db:"company"`
	Position     string `json:"position" db:"position"`
	Introduction string `json:"introduction" db:"introduction"`
	HomePage     string `json:"homePage" db:"homepage"`
	Github       string `json:"github" db:"github"`
}

type Profile struct {
	IsFollow     bool      `json:"isFollow"`
	UserID       int64     `json:"userID,string" db:"user_id"`
	Follower     int       `json:"follower" db:"follower"`
	Following    int       `json:"following" db:"following"`
	UserName     string    `json:"username" db:"username"`
	Avatar       string    `json:"avatar" db:"avatar"`
	Email        string    `json:"email" db:"email"`
	Introduction string    `json:"introduction" db:"introduction"`
	Location     string    `json:"location" db:"location"`
	Company      string    `json:"company" db:"company"`
	Position     string    `json:"position" db:"position"`
	HomePage     string    `json:"homePage" db:"homepage"`
	Github       string    `json:"github" db:"github"`
	Format       string    `json:"format"`
	CreateTime   time.Time `json:"createTime" db:"create_time"`
}
