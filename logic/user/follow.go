package user

import (
	"manage/dao/mysql/user"
	"manage/model"
)

func FollowTag(p *model.ParamFollowTag) (err error) {
	return user.FollowTag(p)
}

func FollowTagCancel(p *model.ParamFollowTag) (err error) {
	return user.FollowTagCancel(p)
}

func FollowUser(p *model.ParamFollowUser) (err error) {
	return user.FollowUser(p)
}

func FollowUserCancel(p *model.ParamFollowUser) (err error) {
	return user.FollowUserCancel(p)
}
