package admin

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"manage/dao/mysql"
	"manage/model"
)

const secret = "TangXiongSheng"

var (
	ErrAdminExist      = errors.New("用户已存在")
	ErrAdminNotExist   = errors.New("用户不存在")
	ErrInvalidPassword = errors.New("密码错误")
)

// encryptPassword md5加密
func encryptPassword(originPasswd string) string {
	h := md5.New()
	h.Write([]byte(secret)) //加盐
	return hex.EncodeToString(h.Sum([]byte(originPasswd)))
}

func CheckAdminExist(email string) (err error) {
	sqlStr := `select count(admin_id) from admin where email = ?`
	var count int
	if err = mysql.DB.Get(&count, sqlStr, email); err != nil {
		return
	}
	if count > 0 {
		return ErrAdminExist
	}
	return
}

func InsertAdmin(admin *model.Admin) (err error) {
	//密码加密
	admin.Password = encryptPassword(admin.Password)
	//插入数据库
	sqlStr := `insert into admin(admin_id,email,admin_name,password) values (?,?,?,?)`
	_, err = mysql.DB.Exec(sqlStr, admin.AdminID, admin.Email, admin.AdminName, admin.Password)
	return
}

func Login(a *model.Admin) (err error) {
	originPassword := a.Password
	sqlStr := `select admin_id,admin_name,email,password from admin where email = ?`
	if err = mysql.DB.Get(a, sqlStr, a.Email); err != nil {
		if err == sql.ErrNoRows {
			return ErrAdminNotExist
		}
		return
	}
	//判断密码是否正确
	password := encryptPassword(originPassword)
	if password != a.Password {
		return ErrInvalidPassword
	}

	return
}

func GetAdminByID(id int64) (a *model.Admin, err error) {
	a = new(model.Admin)
	sqlStr := `select admin_id,admin_name,email,avatar from admin where admin_id = ?`
	if err = mysql.DB.Get(a, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrAdminNotExist
		}
		return
	}
	return
}
