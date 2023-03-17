package logic

import (
	"manage/dao/mysql"
	"manage/model"
)

func FollowTag(p *model.ParamFollowTag) (err error) {
	return mysql.FollowTag(p)
}

func FollowTagCancel(p *model.ParamFollowTag) (err error) {
	return mysql.FollowTagCancel(p)
}

func FollowUser(p *model.ParamFollowUser) (err error) {
	return mysql.FollowUser(p)
}

func FollowUserCancel(p *model.ParamFollowUser) (err error) {
	return mysql.FollowUserCancel(p)
}
