package mysql

import (
	"errors"
	"fmt"
	"manage/model"
)

func FollowTag(p *model.ParamFollowTag) (err error) {
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

	sqlStr1 := `insert into follow_tag(user_id, follow_tag_id) values (?,?) `
	rs, err := tx.Exec(sqlStr1, p.UserID, p.TagID)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("FollowTag exec sqlStr1 failed")
	}

	sqlStr2 := `update tag set follower_number = follower_number + 1 where tag_id = ?`
	rs, err = tx.Exec(sqlStr2, p.TagID)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("FollowTag exec sqlStr2 failed")
	}

	return
}

func FollowTagCancel(p *model.ParamFollowTag) (err error) {
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

	sqlStr1 := `delete from follow_tag where user_id = ? and follow_tag_id = ?`
	rs, err := tx.Exec(sqlStr1, p.UserID, p.TagID)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("FollowTagCancel exec sqlStr1 failed")
	}

	sqlStr2 := `update tag set follower_number = follower_number - 1 where tag_id = ?`
	rs, err = tx.Exec(sqlStr2, p.TagID)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("FollowTagCancel exec sqlStr2 failed")
	}

	return
}

func IsFollowTag(uid int64, tid int) (isFollow bool, err error) {
	sqlStr := `select count(1) from follow_tag where user_id = ? and follow_tag_id = ?`
	var count int
	if err = db.Get(&count, sqlStr, uid, tid); err != nil {
		return
	}
	if count != 0 {
		isFollow = true
	} else {
		isFollow = false
	}
	return
}

func FollowUser(p *model.ParamFollowUser) (err error) {
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

	sqlStr1 := `insert into follow_user(user_id, follow_user_id) values (?,?) `
	rs, err := tx.Exec(sqlStr1, p.UserID, p.FollowUserID)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("FollowUser exec sqlStr1 failed")
	}

	sqlStr2 := `update user set follower = follower + 1 where user_id = ?`
	rs, err = tx.Exec(sqlStr2, p.FollowUserID)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("FollowUser exec sqlStr2 failed")
	}

	sqlStr3 := `update user set following = following + 1 where user_id = ?`
	rs, err = tx.Exec(sqlStr3, p.UserID)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("FollowUser exec sqlStr3 failed")
	}

	return
}

func FollowUserCancel(p *model.ParamFollowUser) (err error) {
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

	sqlStr1 := `delete from follow_user where user_id = ? and follow_user_id = ?`
	rs, err := tx.Exec(sqlStr1, p.UserID, p.FollowUserID)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("FollowUserCancel exec sqlStr1 failed")
	}

	sqlStr2 := `update user set follower = follower - 1 where user_id = ?`
	rs, err = tx.Exec(sqlStr2, p.FollowUserID)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("FollowUserCancel exec sqlStr2 failed")
	}

	sqlStr3 := `update user set following = following - 1 where user_id = ?`
	rs, err = tx.Exec(sqlStr3, p.UserID)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("FollowUserCancel exec sqlStr3 failed")
	}

	return
}

func IsFollowUser(uid, fid int64) (isFollow bool, err error) {
	sqlStr := `select count(1) from follow_user where user_id = ? and follow_user_id = ?`
	var count int
	if err = db.Get(&count, sqlStr, uid, fid); err != nil {
		return
	}
	if count != 0 {
		isFollow = true
	} else {
		isFollow = false
	}
	return
}
