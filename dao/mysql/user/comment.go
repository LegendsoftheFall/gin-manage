package user

import (
	"errors"
	"fmt"
	"manage/dao/mysql"
	"manage/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreateRootComment(p *model.ParamComment) (err error) {
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

	// 将根评论信息插入评论表
	sqlStr1 := `insert into comments(comment_id, user_id, item_id, item_type, status, comment_picture, comment_content)
values (?,?,?,?,?,?,?)`
	rs, err := tx.Exec(sqlStr1, p.CommentID, p.UserID, p.ItemID, p.ItemType, 1, p.Picture, p.Content)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("CreateRootComment exec sqlStr1 failed")
	}

	// 将评论ID插入路径表，祖先为0，后代为评论ID，深度为0
	sqlStr2 := `insert into tree_path(ancestor, descendant, distance) values (?,?,?)`
	rs, err = tx.Exec(sqlStr2, p.CommentID, p.CommentID, 0)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("CreateRootComment exec sqlStr2 failed")
	}

	return
}

func CreateReplyComment(p *model.ParamComment) (err error) {
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

	// 将评论信息插入评论表
	sqlStr1 := `insert into comments(comment_id, user_id, item_id, item_type, status, comment_picture, comment_content)
values (?,?,?,?,?,?,?)`
	rs, err := tx.Exec(sqlStr1, p.CommentID, p.UserID, p.ItemID, p.ItemType, 1, p.Picture, p.Content)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("CreateReplyComment exec sqlStr1 failed")
	}

	// 查询回复的评论对象的所有祖先，将后代改为本评论ID，distance+1
	newPath := make([]*model.TreePath, 0)
	sqlStr2 := `select ancestor, descendant ,distance  from tree_path where descendant = ?`
	if err = tx.Select(&newPath, sqlStr2, p.ToCommentID); err != nil {
		return errors.New("CreateReplyComment exec sqlStr2 failed")
	}

	length := len(newPath)
	// 拼接字符串，批量插入
	if length > 0 {
		// 存放(?,?)的slice
		valueString := make([]string, 0, length)
		// 存放values的slice
		valueArgs := make([]interface{}, 0, length*3)
		for _, np := range newPath {
			valueString = append(valueString, "(?,?,?)")
			valueArgs = append(valueArgs, np.Ancestor, p.CommentID, np.Distance+1)
		}
		// 拼接sql语句
		sqlStr3 := fmt.Sprintf("insert into"+" tree_path(ancestor,descendant,distance) "+"values "+"%s",
			strings.Join(valueString, ","))
		// 插入操作
		rs, err = tx.Exec(sqlStr3, valueArgs...)
		if err != nil {
			return err
		}
		n, err = rs.RowsAffected()
		if err != nil {
			return err
		}
		if int(n) != length {
			return errors.New("CreateReplyComment exec sqlStr3 failed")
		}
	}

	// 最后插入本评论
	sqlStr4 := `insert into tree_path(ancestor, descendant, distance) values (?,?,?)`
	rs, err = tx.Exec(sqlStr4, p.CommentID, p.CommentID, 0)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("CreateReplyComment exec sqlStr4 failed")
	}
	return
}

func DeleteRootComment(commentID int64) (err error) {
	tx, err := mysql.DB.Beginx()
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return err
	}

	defer func() {
		if pc := recover(); pc != nil {
			err = tx.Rollback()
			panic(pc)
		} else if err != nil {
			fmt.Println("rollback")
			err = tx.Rollback()
		} else {
			err = tx.Commit()
			fmt.Println("commit")
		}
	}()

	// 操作tree_path
	//查询除了自身的所有后代
	descendants := make([]int64, 0)
	sqlStr1 := `select descendant from tree_path where ancestor = ? and distance>0`
	if err = tx.Select(&descendants, sqlStr1, commentID); err != nil {
		return errors.New("DeleteRootComment exec sqlStr1 failed")
	}
	fmt.Println(descendants)

	// 若存在后代
	if len(descendants) > 0 {
		// 删除所有以根后代为祖先的所有节点
		sqlStr2 := `delete from tree_path where ancestor in (?)`
		query, args, err1 := sqlx.In(sqlStr2, descendants)
		query = mysql.DB.Rebind(query)
		rs, err1 := tx.Exec(query, args...)
		if err1 != nil {
			err = err1
			return
		}
		n, err1 := rs.RowsAffected()
		if err1 != nil {
			err = err1
			return
		}
		if n < 0 {
			return errors.New("DeleteRootComment exec sqlStr2 failed")
		}
	}

	// 删除以自身节点为祖先的节点
	sqlStr3 := `delete from tree_path where ancestor = ?`
	rs, err := tx.Exec(sqlStr3, commentID)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if int(n) != len(descendants)+1 {
		return errors.New("DeleteRootComment exec sqlStr3 failed")
	}

	// 操作comment表
	// 删除评论
	descendants = append(descendants, commentID)
	sqlStr4 := `delete from comments where comment_id in (?)`
	query, args, err := sqlx.In(sqlStr4, descendants)
	query = mysql.DB.Rebind(query)
	rs, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if int(n) != len(descendants) {
		return errors.New("DeleteRootComment exec sqlStr4 failed")
	}

	return
}

func IsLeafComment(commentID int64) (isLeaf bool, err error) {
	var leaf int
	sqlStr1 := `select count(1) from tree_path where ancestor = ? and distance > 0`
	if err = mysql.DB.Get(&leaf, sqlStr1, commentID); err != nil {
		return
	}
	if leaf == 0 {
		isLeaf = true
	} else {
		isLeaf = false
	}
	return
}

func HaveReplies(commentID int64) (have bool, err error) {
	// 查询是否存在多个子评论
	var reply int
	sqlStr := `select count(1) from tree_path where ancestor = ? and distance = 1`
	if err = mysql.DB.Get(&reply, sqlStr, commentID); err != nil {
		return
	}
	if reply > 1 {
		have = true
	} else {
		have = false
	}
	return
}

func GetAncestorReply(commentID int64) (replies []int64, err error) {
	sqlStr := `select ancestor from tree_path where descendant = ? and distance>0 order by distance`
	err = mysql.DB.Select(&replies, sqlStr, commentID)
	return
}

func IsReplyDeleted(commentID int64) (isDeleted bool, err error) {
	var status int
	sqlStr := `select status from comments where comment_id = ?`
	if err = mysql.DB.Get(&status, sqlStr, commentID); err != nil {
		return
	}
	if status == 0 {
		isDeleted = true
	} else {
		isDeleted = false
	}
	return
}

func DeleteReplyComment(commentID int64) (err error) {
	tx, err := mysql.DB.Beginx()
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return err
	}

	defer func() {
		if pc := recover(); pc != nil {
			err = tx.Rollback()
			panic(pc)
		} else if err != nil {
			fmt.Println("rollback")
			err = tx.Rollback()
		} else {
			err = tx.Commit()
			fmt.Println("commit")
		}
	}()

	// 删除待删除回复 从后代依次删除
	sqlStr1 := `delete from tree_path where descendant = ?`
	rs, err := tx.Exec(sqlStr1, commentID)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n < 0 {
		return errors.New("DeleteReplyComment exec sqlStr1 failed")
	}

	sqlStr2 := `delete from comments where comment_id = ?`
	rs, err = tx.Exec(sqlStr2, commentID)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("DeleteReplyComment exec sqlStr1 failed")
	}

	return
}

func SetReplyComment(commentID int64) (err error) {
	sqlStr := `update comments set status = 0 where comment_id = ?`
	_, err = mysql.DB.Exec(sqlStr, commentID)
	return
}

func GetRootComment(itemID, page, size int64, order, endTime string) (rootComment []*model.RootComment, total int, err error) {
	// 获取根评论信息
	if order == model.OrderScore {
		sqlStr1 := `select comment_id, user_id, item_id, item_type, status,likes, comment_picture, comment_content, create_time from comments
where item_id = ? and item_type = 2 and create_time < ? and status = 1 order by likes desc limit ?,?`
		if err = mysql.DB.Select(&rootComment, sqlStr1, itemID, endTime, (page-1)*size, size); err != nil {
			return
		}
	} else {
		sqlStr1 := `select comment_id, user_id, item_id, item_type, status,likes, comment_picture, comment_content, create_time from comments
where item_id = ? and item_type = 2 and create_time < ? and status = 1 order by create_time desc limit ?,?`
		if err = mysql.DB.Select(&rootComment, sqlStr1, itemID, endTime, (page-1)*size, size); err != nil {
			return
		}
	}
	// 获取根评论总数
	sqlStr2 := `select count(1) from comments where item_id = ? and status = 1 and item_type = 2 and create_time < ?`
	err = mysql.DB.Get(&total, sqlStr2, itemID, endTime)
	return
}

func GetReplyComment(rootCommentID int64) (replyComment []*model.ReplyComment, err error) {
	// 查询出所有回复评论
	sqlStr1 := `select comment_id, user_id, item_id, item_type, status, likes, comment_picture, comment_content, create_time  from comments
where comment_id in (
    select descendant from tree_path where ancestor = ? and distance > 0
) and status = 1 order by create_time`
	if err = mysql.DB.Select(&replyComment, sqlStr1, rootCommentID); err != nil {
		return
	}

	//
	for _, reply := range replyComment {
		reply.ReplyCommentID = rootCommentID
		// 查询回复ID(深度为1)
		sqlStr2 := `select ancestor from tree_path where descendant = ? and distance = 1`
		if err = mysql.DB.Get(&reply.ToReplyID, sqlStr2, reply.ReplyID); err != nil {
			continue
		}
		// 若是子评论
		if reply.ToReplyID == rootCommentID {
			reply.ToReplyID = 0
			reply.ToUserID = 0
			reply.Level = 1
		} else {
			reply.Level = 2
		}

		sqlStr3 := `select user_id from comments where comment_id = ? and status = 1`
		if err = mysql.DB.Get(&reply.ToUserID, sqlStr3, reply.ReplyID); err != nil {
			continue
		}
	}
	return
}
