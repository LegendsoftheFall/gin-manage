package logic

import (
	"manage/dao/mysql"
	"manage/hooks"
	"manage/model"
	"manage/pkg/jwt"
	"manage/pkg/snowflake"

	"go.uber.org/zap"
)

func SignUp(p *model.ParamSignUp) (err error) {
	//判断用户是否储存在
	if err = mysql.CheckUserExist(p.Email); err != nil {
		return err
	}
	//生成UID
	userID := snowflake.GenID()
	user := &model.User{
		UserID:   userID,
		Email:    p.Email,
		Username: p.Username,
		Password: p.Password,
	}
	//加密
	//保存到数据库
	return mysql.InsertUser(user)
}

func Login(p *model.ParamLogin) (user *model.User, err error) {
	user = &model.User{
		Email:    p.Email,
		Password: p.Password,
	}
	if err = mysql.Login(user); err != nil {
		return nil, err
	}

	aToken, rToken, err := jwt.GenToken(user.UserID, user.Email)
	if err != nil {
		return
	}
	user.AccessToken = aToken
	user.RefreshToken = rToken
	return
}

func GetUserByID(id int64) (user *model.User, err error) {
	return mysql.GetUserByID(id)
}

func GetUserInfoByUserID(id int64, cid int64) (userInfo *model.UserInfo, err error) {
	userInfo, err = mysql.GetUserInfoByUserID(id)
	if err != nil {
		zap.L().Error("mysql.GetUserInfoByUserID",
			zap.Int64("UserID", id),
			zap.Error(err))
		return nil, err
	}
	if cid == 0 {
		return
	} else {
		userInfo.IsFollow, err = mysql.IsFollowUser(cid, id)
		if err != nil {
			zap.L().Error("mysql.IsFollowUser",
				zap.Int64("UserID", cid),
				zap.Int64("followUserID", id),
				zap.Error(err))
		}
	}
	return
}

func GetUserInfoByArticleID(id, cid int64) (userInfo *model.UserInfo, err error) {
	userInfo, err = mysql.GetUserInfoByArticleID(id)
	if err != nil {
		zap.L().Error("mysql.GetUserInfoByArticleID",
			zap.Int64("articleID", id),
			zap.Error(err))
	}
	userInfo.IsFollow, err = mysql.IsFollowUser(cid, userInfo.UserID)
	return
}

func GetFollowingUsers(p *model.ParamUserList) (userInfoList []*model.UserInfo, total int, err error) {
	userInfoList, total, err = mysql.GetFollowingUsers(p)
	if err != nil {
		zap.L().Error(" mysql.GetFollowUsers failed",
			zap.Int64("userID", p.CurrentUserID),
			zap.Error(err))
		return nil, 0, err
	}
	for _, userInfo := range userInfoList {
		userInfo.IsFollow = true
	}
	return
}

func GetFollowerUsers(p *model.ParamUserList) (userInfoList []*model.UserInfo, total int, err error) {
	userInfoList, total, err = mysql.GetFollowerUsers(p)
	if err != nil {
		zap.L().Error(" mysql.GetFollowUsers failed",
			zap.Int64("userID", p.CurrentUserID),
			zap.Error(err))
		return nil, 0, err
	}
	for _, userInfo := range userInfoList {
		userInfo.IsFollow, err = mysql.IsFollowUser(p.CurrentUserID, userInfo.UserID)
		if err != nil {
			zap.L().Error(" mysql.IsFollowUser failed",
				zap.Int64("userID", p.CurrentUserID),
				zap.Int64("follow_userID", userInfo.UserID),
				zap.Error(err))
			continue
		}
	}
	return
}

func GetUserProfile(userID int64) (profile model.UserProfile, err error) {
	return mysql.GetUserProfile(userID)
}

func UpdateUserProfile(userID int64, p *model.UserProfile) error {
	return mysql.UpdateUserProfile(userID, p)
}

func GetProfile(p *model.ParamUserProfile) (profile model.Profile, err error) {
	profile, err = mysql.GetProfile(p.UserID)
	if err != nil {
		zap.L().Error("mysql.GetProfile failed",
			zap.Int64("currentUserID", p.CurrentUserID),
			zap.Int64("userID", p.UserID),
			zap.Error(err))
		return
	}
	profile.Format = hooks.Date(profile.CreateTime)
	if p.UserID == 0 {
		profile.IsFollow = false
	} else {
		profile.IsFollow, err = mysql.IsFollowUser(p.CurrentUserID, p.UserID)
		if err != nil {
			zap.L().Error("mysql.IsFollowUser failed",
				zap.Int64("currentUserID", p.CurrentUserID),
				zap.Int64("userID", p.UserID),
				zap.Error(err))
		}
	}
	return
}

func GetSearchUsers(p *model.ParamSearch) (userInfoList []*model.UserInfo, total int, err error) {
	userInfoList, total, err = mysql.GetSearchUsers(p)
	if err != nil {
		zap.L().Error(" mysql.GetFollowUsers failed",
			zap.Int64("userID", p.CurrentUserID),
			zap.Error(err))
		return nil, 0, err
	}
	for _, userInfo := range userInfoList {
		userInfo.IsFollow, err = mysql.IsFollowUser(p.CurrentUserID, userInfo.UserID)
		if err != nil {
			zap.L().Error(" mysql.IsFollowUser failed",
				zap.Int64("userID", p.CurrentUserID),
				zap.Int64("follow_userID", userInfo.UserID),
				zap.Error(err))
			continue
		}
	}
	return
}
