package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"manage/dao/mysql"
	"manage/model"

	"go.uber.org/zap"
)

var (
	ErrTagExist = errors.New("标签已存在")
)

func CheckTagExist(name string) (err error) {
	sqlStr := `select count(tag_id) from tag where tag_name = ?`
	var count int
	if err = mysql.DB.Get(&count, sqlStr, name); err != nil {
		return
	}
	if count > 0 {
		return ErrTagExist
	}
	return
}

func IsLTTotal(id int) (err error) {
	var total int
	sqlStr := `select count(*) from tag`
	if err = mysql.DB.Get(&total, sqlStr); err != nil {
		return err
	}
	if id != total+1 {
		return ErrTagExist
	}
	return
}

func CreateTag(p *model.ParamCreateTag) (err error) {
	tx, err := mysql.DB.Beginx() // 开启事务
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return err
	}
	defer func() {
		if pc := recover(); pc != nil {
			err = tx.Rollback()
			panic(pc) // re-throw panic after Rollback
		} else if err != nil {
			fmt.Println("rollback")
			err = tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
			fmt.Println("commit")
		}
	}()

	// 1.插入标签表
	sqlStr1 := `insert into tag (tag_id,tag_name,image,introduction) values (?,?,?,?)`
	rs, err := tx.Exec(sqlStr1, p.ID, p.Name, p.Image, p.Introduction)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("CreateTag exec sqlStr1 failed")
	}

	return
}

func UpdateTag(p *model.ParamCreateTag) (err error) {
	tx, err := mysql.DB.Beginx() // 开启事务
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return err
	}
	defer func() {
		if pc := recover(); pc != nil {
			err = tx.Rollback()
			panic(pc) // re-throw panic after Rollback
		} else if err != nil {
			fmt.Println("rollback")
			err = tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
			fmt.Println("commit")
		}
	}()

	// 1.更新标签
	sqlStr1 := `update tag set tag_name=?,image=?,introduction=? where tag_id = ?`
	rs, err := tx.Exec(sqlStr1, p.Name, p.Image, p.Introduction, p.ID)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n < 0 {
		return errors.New("UpdateArticle exec sqlStr1 failed")
	}

	return
}

func DeleteTag(id int64) (err error) {
	sqlStr := `delete from tag where tag_id = ?`
	_, err = mysql.DB.Exec(sqlStr, id)
	return
}

func GetAllTags(p *model.ParamTagList) (tagList []*model.TagDetail, total int, err error) {
	sqlStr1 := `select tag_id,article_number,follower_number,tag_name,image,introduction from tag limit ?,?`
	if err = mysql.DB.Select(&tagList, sqlStr1, (p.Page-1)*p.Size, p.Size); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no tag in DB")
			err = nil
		}
	}
	sqlStr2 := `select count(*) from tag`
	err = mysql.DB.Get(&total, sqlStr2)
	return
}

func SelectTags() (tagList []*model.Tag, err error) {
	sqlStr := `select tag_id,article_number,follower_number,tag_name,image from tag`
	if err = mysql.DB.Select(&tagList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no tag in DB")
			err = nil
		}
	}
	return
}
