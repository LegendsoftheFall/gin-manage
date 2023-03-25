package user

import (
	"manage/dao/mysql"
	"manage/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

func GetCollectedArticleByIDs(ids []string) (articleList []*model.ArticleInfo, err error) {
	articleList = make([]*model.ArticleInfo, 0, len(ids))
	sqlStr := `select article_id,title,subtitle,content,image,author_id,view_count,likes,comments,create_time
from article where article_id in (?) order by FIND_IN_SET(article_id, ?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	//fmt.Println(query, args)
	if err != nil {
		return
	}
	query = mysql.DB.Rebind(query)
	err = mysql.DB.Select(&articleList, query, args...)
	return
}
