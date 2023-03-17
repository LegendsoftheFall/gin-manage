package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"manage/model"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

var (
	ErrInvalidID = errors.New("id不存在")
)

func GetTrendingTags() (tagList []*model.Tag, err error) {
	sqlStr := `select tag_id,article_number,tag_name,image from tag order by article_number desc limit 6`
	if err = db.Select(&tagList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no tag in db")
			err = nil
		}
	}
	return
}

func SelectTags() (tagList []*model.Tag, err error) {
	sqlStr := `select tag_id,article_number,tag_name,image from tag`
	if err = db.Select(&tagList, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no tag in db")
			err = nil
		}
	}
	return
}

func GetTagDetailByID(id int) (detail *model.TagDetail, err error) {
	detail = new(model.TagDetail)
	sqlStr := `select tag_id,article_number,follower_number,tag_name,image,introduction from tag where tag_id = ?`
	if err = db.Get(detail, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrInvalidID
		}
	}
	return
}

func GetTagNameByArticleID(id int64) (tagName []string, err error) {
	sqlStr := `select tag_name from tag where tag_id in (
    select tag_id from article_tag where article_id = ?
)`
	err = db.Select(&tagName, sqlStr, id)
	return
}

func GetArticleIDByTagID(id int) (articleID []int64, err error) {
	sqlStr := `select article_id from article_tag where tag_id = ?`
	err = db.Select(&articleID, sqlStr, id)
	return
}

func GetTagIDByTagName(name string) (id int, err error) {
	sqlStr := `select id from tag where tag_name = ?`
	err = db.Get(&id, sqlStr, name)
	return
}

func GetTagNameByTagID(id int) (name string, err error) {
	sqlStr := `select tag_name from tag where id = ?`
	err = db.Get(&name, sqlStr, id)
	return
}

func GetAllTags(p *model.ParamTagList) (tagList []*model.Tag, total int, err error) {
	sqlStr1 := `select tag_id,article_number,tag_name,image from tag limit ?,?`
	if err = db.Select(&tagList, sqlStr1, (p.Page-1)*p.Size, p.Size); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no tag in db")
			err = nil
		}
	}
	sqlStr2 := `select count(*) from tag`
	err = db.Get(&total, sqlStr2)
	return
}

func GetFollowingTags(p *model.ParamTagList) (tagList []*model.Tag, total int, err error) {
	tagNumList := make([]int, 0)
	sqlStr1 := `select follow_tag_id from follow_tag where user_id = ?`
	err = db.Select(&tagNumList, sqlStr1, p.CurrentUserID)
	if err != nil {
		return nil, 0, err
	}
	total = len(tagNumList)
	if total == 0 {
		return nil, 0, err
	}
	sqlStr2 := `select tag_id,article_number,tag_name,image from tag where tag_id in (?) limit ?,?`
	query, args, err := sqlx.In(sqlStr2, tagNumList, (p.Page-1)*p.Size, p.Size)
	if err != nil {
		return nil, 0, err
	}
	query = db.Rebind(query)
	err = db.Select(&tagList, query, args...)
	return
}

func DeleteTagByID(id int64) (err error) {
	tx, err := db.Beginx() // 开启事务
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

	var tags []int

	// 1. 根据文章id找到对应的tag数组
	sqlStr1 := `select tag_id from article_tag where article_id = ?`
	if err = tx.Select(&tags, sqlStr1, id); err != nil {
		return err
	}

	// 2. 根据文章id删除tag
	sqlStr2 := `delete from article_tag where  article_id = ?`
	rs, err := tx.Exec(sqlStr2, id)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n < 1 {
		return errors.New("DeleteTagByID exec sqlStr2 failed")
	}

	//3.更新标签表文章数量
	sqlStr3 := `update tag set article_number = article_number - 1 where tag_id in (?)`
	query, args, err := sqlx.In(sqlStr3, tags)
	query = db.Rebind(query)
	rs, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if int(n) != len(tags) {
		return errors.New("DeleteTagByID exec sqlStr3 failed")
	}
	return
}

func GetTagIDsByArticleID(id int64) (tagIDs []int, err error) {
	sqlStr := `select tag_id from article_tag where article_id = ?`
	err = db.Select(&tagIDs, sqlStr, id)
	return
}

func GetSearchTags(p *model.ParamSearch) (tagList []*model.Tag, total int, err error) {
	sqlStr1 := `select tag_id,article_number,tag_name,image from tag where tag_name like concat('%',?,'%') order by follower_number desc limit ?,?`
	if err = db.Select(&tagList, sqlStr1, p.Key, (p.Page-1)*p.Size, p.Size); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no tag in db")
			err = nil
		}
	}
	total = len(tagList)
	return
}
