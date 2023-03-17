package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"manage/model"
	"strings"

	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
)

func CreateArticle(p *model.Article) (err error) {
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

	// 1.插入文章表
	sqlStr1 := `insert into article (article_id, title, subtitle, content, html, markdown, author_id,image,source)
values (?,?,?,?,?,?,?,?,?)`
	rs, err := tx.Exec(sqlStr1, p.ID, p.Title, p.SubTitle, p.Content, p.Html, p.MarkDown, p.AuthorID, p.Image, p.Source)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("CreateArticle exec sqlStr1 failed")
	}

	// 2.批量插入文章标签表
	tags := p.Tags
	length := len(tags)
	if length > 0 {
		// 存放(?,?)的slice
		valueString := make([]string, 0, length)
		// 存放values的slice
		valueArgs := make([]interface{}, 0, length*2)
		for _, t := range tags {
			valueString = append(valueString, "(?,?)")
			valueArgs = append(valueArgs, p.ID)
			valueArgs = append(valueArgs, t)
		}
		// 拼接sql语句
		sqlStr2 := fmt.Sprintf("insert into"+" article_tag(article_id, tag_id) "+"values "+"%s",
			strings.Join(valueString, ","))

		// 插入操作
		rs, err = tx.Exec(sqlStr2, valueArgs...)
		if err != nil {
			return err
		}
		n, err = rs.RowsAffected()
		if err != nil {
			return err
		}
		if int(n) != length {
			return errors.New("CreateArticle exec sqlStr2 failed")
		}
	}

	//3.更新标签表文章数量
	sqlStr3 := `update tag set article_number = article_number + 1 where tag_id in (?)`
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
	if int(n) != length {
		return errors.New("CreateArticle exec sqlStr3 failed")
	}

	return
}

func UpdateArticle(p *model.Article) (err error) {
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

	// 1.更新文章
	sqlStr1 := `update article set title=?,subtitle=?,content=?,image=?,markdown=?,html=? where article_id = ?`
	rs, err := tx.Exec(sqlStr1, p.Title, p.SubTitle, p.Content, p.Image, p.MarkDown, p.Html, p.ID)
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

	// 2.批量插入文章标签表
	tags := p.Tags
	length := len(tags)
	if length > 0 {
		// 存放(?,?)的slice
		valueString := make([]string, 0, length)
		// 存放values的slice
		valueArgs := make([]interface{}, 0, length*2)
		for _, t := range tags {
			valueString = append(valueString, "(?,?)")
			valueArgs = append(valueArgs, p.ID)
			valueArgs = append(valueArgs, t)
		}
		// 拼接sql语句
		sqlStr2 := fmt.Sprintf("insert into"+" article_tag(article_id, tag_id) "+"values "+"%s",
			strings.Join(valueString, ","))

		// 插入操作
		rs, err = tx.Exec(sqlStr2, valueArgs...)
		if err != nil {
			return err
		}
		n, err = rs.RowsAffected()
		if err != nil {
			return err
		}
		if int(n) != length {
			return errors.New("UpdateArticle exec sqlStr2 failed")
		}
	}

	//3.更新标签表文章数量
	sqlStr3 := `update tag set article_number = article_number + 1 where tag_id in (?)`
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
	if int(n) != length {
		return errors.New("UpdateArticle exec sqlStr3 failed")
	}

	return
}

func DeleteArticle(id int64) (err error) {
	sqlStr := `delete from article where article_id = ?`
	_, err = db.Exec(sqlStr, id)

	return
}

//func GetArticleListByIDs(ids []int64, page, size int64) (articleList []*model.ArticleInfo, err error) {
//	articleList = make([]*model.ArticleInfo, 0, len(ids))
//	sqlStr := `select article_id,title,subtitle,content,image,author_id,view_count,likes,comments,create_time
//from article where article_id in (?) limit ?,?`
//	query, args, err := sqlx.In(sqlStr, ids, (page-1)*size, size)
//	//fmt.Println(query, args)
//	if err != nil {
//		return
//	}
//	query = db.Rebind(query)
//	err = db.Select(&articleList, query, args...)
//	return
//}

func GetArticleListByIDs(ids []string) (articleList []*model.ArticleInfo, err error) {
	articleList = make([]*model.ArticleInfo, 0, len(ids))
	sqlStr := `select article_id,title,subtitle,content,image,author_id,view_count,likes,comments,create_time
from article where article_id in (?)
order by find_in_set(article_id,?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	//fmt.Println(query, args)
	if err != nil {
		return
	}
	query = db.Rebind(query)
	err = db.Select(&articleList, query, args...)
	return
}

func GetArticleListByID(id, page, size int64) (articleList []*model.ArticleInfo, err error) {
	sqlStr := `select article_id,title,subtitle,content,image,author_id,view_count,likes,comments,create_time
from article where author_id = ?  order by create_time desc limit ?,?`
	err = db.Select(&articleList, sqlStr, id, (page-1)*size, size)
	if err != nil {
		return
	}
	return
}

func GetArticleNumByID(id int64) (number int, err error) {
	sqlStr := `select count(article_id) from article where author_id = ?`
	err = db.Get(&number, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrInvalidID
		}
	}
	return
}

func GetArticleDetailByID(id int64) (article model.Article, err error) {
	sqlStr := `select article_id, title, subtitle, content, markdown, image, author_id, view_count, likes, comments, create_time
from article where article_id = ?`
	err = db.Get(&article, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrInvalidID
		}
		zap.L().Error("GetArticleDetailByID failed", zap.Error(err))
	}
	return
}

func GetUserIDByArticleID(id int64) (userID int64, err error) {
	sqlStr := `select author_id from article where article_id = ?`
	err = db.Get(&userID, sqlStr, id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrUserNotExist
		}
	}
	return
}

func UpdateViewCount(id int64) (err error) {
	sqlStr := `update article set view_count = view_count+1 where article_id = ?`
	_, err = db.Exec(sqlStr, id)
	return
}

func GetSearchArticleListByView(p *model.ParamSearch) (articleList []*model.ArticleInfo, total int, err error) {
	articleList = make([]*model.ArticleInfo, 0)
	sqlStr := `select article_id,title,subtitle,content,image,author_id,view_count,likes,comments,create_time
from article where title like concat('%',?,'%') order by view_count desc limit ?,?`
	if err = db.Select(&articleList, sqlStr, p.Key, (p.Page-1)*p.Size, p.Size); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no article in db")
			err = nil
		}
	}
	total = len(articleList)
	return
}

func GetSearchArticleListByTime(p *model.ParamSearch) (articleList []*model.ArticleInfo, total int, err error) {
	articleList = make([]*model.ArticleInfo, 0)
	sqlStr := `select article_id,title,subtitle,content,image,author_id,view_count,likes,comments,create_time
from article where title like concat('%',?,'%') order by create_time desc limit ?,?`
	if err = db.Select(&articleList, sqlStr, p.Key, (p.Page-1)*p.Size, p.Size); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no article in db")
			err = nil
		}
	}
	total = len(articleList)
	return
}
