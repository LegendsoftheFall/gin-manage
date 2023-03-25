package user

import "manage/dao/mysql"

func AddLikeNumByArticleID(id int64) (err error) {
	sqlStr := `update article set likes=likes+1 where article_id = ?`
	_, err = mysql.DB.Exec(sqlStr, id)
	return
}

func RemoveLikeNumByArticleID(id int64) (err error) {
	sqlStr := `update article set likes=likes-1 where article_id = ?`
	_, err = mysql.DB.Exec(sqlStr, id)
	return
}
