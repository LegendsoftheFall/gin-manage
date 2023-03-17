package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"manage/model"

	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
)

const secret = "TangXiongSheng"

var (
	ErrUserExist       = errors.New("用户已存在")
	ErrUserNotExist    = errors.New("用户不存在")
	ErrInvalidPassword = errors.New("密码错误")
)

// encryptPassword md5加密
func encryptPassword(originPasswd string) string {
	h := md5.New()
	h.Write([]byte(secret)) //加盐
	return hex.EncodeToString(h.Sum([]byte(originPasswd)))
}

func CheckUserExist(email string) (err error) {
	sqlStr := `select count(user_id) from user where email = ?`
	var count int
	if err = db.Get(&count, sqlStr, email); err != nil {
		return
	}
	if count > 0 {
		return ErrUserExist
	}
	return
}

func InsertUser(user *model.User) (err error) {
	//密码加密
	user.Password = encryptPassword(user.Password)
	//插入数据库
	sqlStr := `insert into user(user_id,email,username,password) values (?,?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Email, user.Username, user.Password)
	return
}

func Login(user *model.User) (err error) {
	originPassword := user.Password
	sqlStr := `select user_id,username,email,password from user where email = ?`
	if err = db.Get(user, sqlStr, user.Email); err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotExist
		}
		return
	}
	//判断密码是否正确
	password := encryptPassword(originPassword)
	if password != user.Password {
		return ErrInvalidPassword
	}

	return
}

func GetUserByID(id int64) (user *model.User, err error) {
	user = new(model.User)
	sqlStr := `select user_id,username,email,avatar from user where user_id = ?`
	if err = db.Get(user, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrUserNotExist
		}
		return
	}
	return
}

func GetUserInfoByUserID(id int64) (userInfo *model.UserInfo, err error) {
	userInfo = new(model.UserInfo)
	sqlStr := `select user_id,username,email,avatar,introduction,follower,following,homepage,github,position from user where user_id = ?`
	if err = db.Get(userInfo, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrUserNotExist
		}
		return
	}
	return
}

func GetUserInfoByArticleID(id int64) (userInfo *model.UserInfo, err error) {
	userInfo = new(model.UserInfo)
	sqlStr := `select user_id,username,email,avatar,introduction,follower from user where user_id in (
    select author_id from article where article_id = ?
)`
	if err = db.Get(userInfo, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrUserNotExist
		}
		return
	}
	return
}

func GetFollowingUsers(p *model.ParamUserList) (userInfoList []*model.UserInfo, total int, err error) {
	userIDList := make([]int64, 0)
	sqlStr1 := `select follow_user_id from follow_user where user_id = ? order by create_time desc`
	err = db.Select(&userIDList, sqlStr1, p.CurrentUserID)
	if err != nil {
		return nil, 0, err
	}
	total = len(userIDList)
	if total == 0 {
		return nil, 0, err
	}

	sqlStr2 := `select user_id,follower,username,avatar,email,introduction from user where user_id in (?) limit ?,?`
	query, args, err := sqlx.In(sqlStr2, userIDList, (p.Page-1)*p.Size, p.Size)
	if err != nil {
		return nil, 0, err
	}
	query = db.Rebind(query)
	err = db.Select(&userInfoList, query, args...)
	return
}

func GetFollowerUsers(p *model.ParamUserList) (userInfoList []*model.UserInfo, total int, err error) {
	userIDList := make([]int64, 0)
	sqlStr1 := `select user_id from follow_user where follow_user_id = ? order by create_time desc `
	err = db.Select(&userIDList, sqlStr1, p.CurrentUserID)
	if err != nil {
		return nil, 0, err
	}
	total = len(userIDList)
	if total == 0 {
		return nil, 0, err
	}

	sqlStr2 := `select user_id,follower,username,avatar,email,introduction from user where user_id in (?) limit ?,?`
	query, args, err := sqlx.In(sqlStr2, userIDList, (p.Page-1)*p.Size, p.Size)
	if err != nil {
		return nil, 0, err
	}
	query = db.Rebind(query)
	err = db.Select(&userInfoList, query, args...)
	return
}

func GetUserProfile(userID int64) (profile model.UserProfile, err error) {
	sqlStr := `select username, avatar, location, company, position, homepage, github, introduction from user where user_id = ?`
	if err = db.Get(&profile, sqlStr, userID); err != nil {
		if err == sql.ErrNoRows {
			err = ErrUserNotExist
		}
		return
	}
	return
}

func UpdateUserProfile(userID int64, p *model.UserProfile) (err error) {
	sqlStr := `update user set username = ?, avatar = ?, location = ?, company = ?, position =?, introduction = ?,
                homepage = ?, github = ? where user_id = ?`
	_, err = db.Exec(sqlStr, p.UserName, p.Avatar, p.Location, p.Company, p.Position, p.Introduction, p.HomePage, p.Github, userID)
	return
}

func GetProfile(userID int64) (profile model.Profile, err error) {
	sqlStr := `select user_id, username, email, location, position, company, homepage, github, avatar, introduction, follower, following, create_time from user where user_id = ?`
	if err = db.Get(&profile, sqlStr, userID); err != nil {
		if err == sql.ErrNoRows {
			err = ErrUserNotExist
		}
		return
	}
	return
}

func GetSearchUsers(p *model.ParamSearch) (userInfoList []*model.UserInfo, total int, err error) {
	sqlStr1 := `select user_id,follower,username,avatar,email,introduction from user where username like concat('%',?,'%') order by follower desc limit ?,?`
	if err = db.Select(&userInfoList, sqlStr1, p.Key, (p.Page-1)*p.Size, p.Size); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no user in db")
			err = nil
		}
	}
	total = len(userInfoList)
	return
}
