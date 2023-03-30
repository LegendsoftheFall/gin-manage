package model

const (
	OrderTime    = "time"
	OrderScore   = "score"
	ItemComment  = 1
	ItemArticle  = 2
	SearchTop    = "top"
	SearchLatest = "latest"
	SearchTag    = "tag"
	SearchUser   = "user"
)

// ParamSignUp 定义注册请求参数的结构体
type ParamSignUp struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	//RePassword string `json:"rePassword" binding:"required,eqfield=Password"`
}

// ParamLogin 定义注册请求参数的结构体
type ParamLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamToken 定义刷新token请求参数的结构体
type ParamToken struct {
	AToken string `json:"aToken" binding:"required"`
	RToken string `json:"rToken" binding:"required"`
}

// ParamPage 定义分页请求参数的结构体
type ParamPage struct {
	Page int64 `json:"page" form:"page"`
	Size int64 `json:"size" form:"size"`
}

// ParamTagInfo 定义分页请求参数的结构体
type ParamTagInfo struct {
	TagID  int   `json:"tid" form:"tid" binding:"required"`
	UserID int64 `json:"uid" form:"uid"`
}

// ParamArticleList 定义文章列表请求参数的结构体
type ParamArticleList struct {
	Page          int64  `json:"page" form:"page"`
	Size          int64  `json:"size" form:"size"`
	CurrentUserID int64  `json:"uid" form:"uid"`
	Order         string `json:"order" form:"order"`
}

// ParamTagList 定义标签列表请求参数的结构体
type ParamTagList struct {
	Page          int64 `json:"page" form:"page"`
	Size          int64 `json:"size" form:"size"`
	CurrentUserID int64 `json:"uid" form:"uid"`
}

// ParamUserList 定义用户列表请求参数的结构体
type ParamUserList struct {
	Page          int64 `json:"page" form:"page"`
	Size          int64 `json:"size" form:"size"`
	CurrentUserID int64 `json:"uid" form:"uid"`
}

// ParamLike 定义点赞请求参数的结构体
type ParamLike struct {
	ArticleID string `json:"aid" form:"aid" binding:"required"`
	UserID    string `json:"uid" form:"uid" binding:"required"`
}

// ParamCollect 定义收藏请求参数的结构体
type ParamCollect struct {
	ArticleID string `json:"aid" form:"aid" binding:"required"`
	UserID    string `json:"uid" form:"uid" binding:"required"`
}

// ParamFollowTag 定义关注标签请求参数的结构体
type ParamFollowTag struct {
	TagID  int   `json:"tid" form:"tid" binding:"required"`
	UserID int64 `json:"uid" form:"uid" binding:"required"`
}

// ParamFollowUser 定义关注用户请求参数的结构体
type ParamFollowUser struct {
	FollowUserID int64 `json:"fid" form:"fid" binding:"required"`
	UserID       int64 `json:"uid" form:"uid" binding:"required"`
}

// ParamUserProfile 定义关注用户请求参数的结构体
type ParamUserProfile struct {
	UserID        int64 `json:"uid" form:"uid" binding:"required"`
	CurrentUserID int64 `json:"cid" form:"cid"`
}

// ParamComment 定义评论请求参数的结构体
type ParamComment struct {
	ItemType    int    `json:"itemType" binding:"required"`
	UserID      int64  `json:"userID,string" binding:"required"`
	ItemID      int64  `json:"itemID,string"  binding:"required"`
	CommentID   int64  `json:"commentID,string"`
	ToCommentID int64  `json:"toCommentID,string"`
	Content     string `json:"content" binding:"required"`
	Picture     string `json:"picture"`
}

// ParamDeleteComment 定义删除评论请求参数的结构体
type ParamDeleteComment struct {
	ItemType  int   `json:"itemType" binding:"required"`
	UserID    int64 `json:"userID,string" binding:"required"`
	CommentID int64 `json:"commentID,string" binding:"required"`
}

type ParamCommentList struct {
	Order         string `json:"order" form:"order"`
	EndTime       string `json:"endTime" form:"end"`
	Page          int64  `json:"page" form:"page"`
	Size          int64  `json:"size" form:"size"`
	CurrentUserID int64  `json:"uid" form:"uid"`
	ItemID        int64  `json:"itemID" form:"item" binding:"required"`
}

type ParamAdminComment struct {
	Order   string `json:"order" form:"order"`
	EndTime string `json:"endTime" form:"end"`
	Page    int64  `json:"page" form:"page"`
	Size    int64  `json:"size" form:"size"`
	ItemID  int64  `json:"itemID,string" form:"item"`
}

type ParamSearch struct {
	Page          int64  `json:"page" form:"page"`
	Size          int64  `json:"size" form:"size"`
	CurrentUserID int64  `json:"uid" form:"uid"`
	Category      string `json:"category" form:"category"`
	Key           string `json:"key" form:"key"`
}

type ParamDeleteTag struct {
	ID     int64  `json:"id,string" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}

type ParamDeleteArticle struct {
	ID     int64  `json:"id,string" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}

type ParamAdminDeleteComment struct {
	ItemType int    `json:"itemType" binding:"required"`
	ID       int64  `json:"id,string" binding:"required"`
	Secret   string `json:"secret" binding:"required"`
}

type ParamCreateTag struct {
	ID           int    `json:"id,string" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Image        string `json:"image"`
	Introduction string `json:"introduction"`
}

type ParamSetStatus struct {
	ID     int64 `json:"id,string" form:"id"`
	Status int   `json:"status" form:"status"`
}
