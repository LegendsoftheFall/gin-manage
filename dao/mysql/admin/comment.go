package admin

import (
	"database/sql"
	"errors"
	"fmt"
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

func GetAllComment(page, size int64, endTime string) (commentList []*model.Comment, total int, err error) {
	sqlStr1 := `select status,item_type,comment_id,user_id,item_id,comment_content,create_time from comments
        where create_time < ? order by create_time desc limit ?,?`
	if err = mysql.DB.Select(&commentList, sqlStr1, endTime, (page-1)*size, size); err != nil {
		return
	}
	sqlStr2 := `select count(1) from comments`
	err = mysql.DB.Get(&total, sqlStr2)
	return
}

func GetCommentByItemID(id, page, size int64, endTime string) (commentList []*model.Comment, total int, err error) {
	sqlStr1 := `select status,item_type,comment_id,user_id,item_id,comment_content,create_time from comments
        where item_id = ? and create_time < ? order by create_time desc limit ?,?`
	fmt.Println(id)
	if err = mysql.DB.Select(&commentList, sqlStr1, id, endTime, (page-1)*size, size); err != nil {
		return
	}
	sqlStr2 := `select count(1) from comments where item_id = ?`
	err = mysql.DB.Get(&total, sqlStr2, id)
	return
}

func CheckSuperiorComment(commentID int64) (isHave bool) {
	var count int
	sqlStr2 := `select count(ancestor) from tree_path where descendant = ? and distance = 1`
	if err := mysql.DB.Get(&count, sqlStr2, commentID); err != nil {
		if err == sql.ErrNoRows {
			isHave = false
			return
		}
	}
	if count == 1 {
		isHave = true
	}
	return
}

func GetSuperiorCommentInfo(commentID int64) (superiorCommentInfo model.Comment, err error) {
	var superiorCommentID int64
	sqlStr2 := `select ancestor from tree_path where descendant = ? and distance = 1`
	if err = mysql.DB.Get(&superiorCommentID, sqlStr2, commentID); err != nil {
		return
	}
	sqlStr3 := `select comment_id, user_id, item_id, item_type, status, comment_content, create_time from comments
where comment_id = ?`
	if err = mysql.DB.Get(&superiorCommentInfo, sqlStr3, superiorCommentID); err != nil {
		return
	}
	return
}

func GetCommentInfo(commentID int64) (commentInfo model.Comment, err error) {
	sqlStr1 := `select comment_id, user_id, item_id, item_type, status, comment_content, create_time from comments
where comment_id = ?`
	fmt.Println("commentID", commentID)
	if err = mysql.DB.Get(&commentInfo, sqlStr1, commentID); err != nil {
		return
	}
	return
}

func AllowStatus(commentID int64) (err error) {
	sqlStr := `update comments set status = 1 where comment_id = ?`
	_, err = mysql.DB.Exec(sqlStr, commentID)
	return
}

func BanStatus(commentID int64) (err error) {
	sqlStr := `update comments set status = 0 where comment_id = ?`
	_, err = mysql.DB.Exec(sqlStr, commentID)
	return
}
