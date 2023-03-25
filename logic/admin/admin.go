package admin

import (
	"manage/dao/mysql/admin"
	"manage/model"
	"manage/pkg/jwt"
	"manage/pkg/snowflake"
)

func SignUpForAdmin(p *model.ParamSignUp) (err error) {
	//判断用户是否储存在
	if err = admin.CheckAdminExist(p.Email); err != nil {
		return err
	}
	//生成UID
	adminID := snowflake.GenID()
	a := &model.Admin{
		AdminID:   adminID,
		Email:     p.Email,
		AdminName: p.Username,
		Password:  p.Password,
	}
	//加密
	//保存到数据库
	return admin.InsertAdmin(a)
}

func LoginForAdmin(p *model.ParamLogin) (a *model.Admin, err error) {
	a = &model.Admin{
		Email:    p.Email,
		Password: p.Password,
	}
	if err = admin.Login(a); err != nil {
		return nil, err
	}

	aToken, rToken, err := jwt.GenToken(a.AdminID, a.Email)
	if err != nil {
		return
	}
	a.AccessToken = aToken
	a.RefreshToken = rToken
	return
}

func GetAdminByID(id int64) (a *model.Admin, err error) {
	return admin.GetAdminByID(id)
}
