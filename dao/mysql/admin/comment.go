package admin

import (
	"database/sql"
	"errors"
	"manage/dao/mysql"
	"manage/model"
)

func GetUserByID(id int64) (username string, err error) {
	sqlStr := `select username from user where user_id = ?`
	if err = mysql.DB.Get(&username, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("用户不存在")
		}
		return
	}
	return
}

func GetAllComment(page, size int64, order, endTime string) (commentList []*model.Comment, total int, err error) {
	orderKey := "create_time"
	if order == model.OrderScore {
		orderKey = "likes"
	}
	sqlStr1 := `select status,item_type,comment_id,user_id,item_id,comment_content from comments
        where create_time < ? order by ? desc limit ?,?`
	if err = mysql.DB.Select(&commentList, sqlStr1, endTime, orderKey, (page-1)*size, size); err != nil {
		return
	}
	sqlStr2 := `select count(1) from comments`
	err = mysql.DB.Get(&total, sqlStr2)
	return
}
