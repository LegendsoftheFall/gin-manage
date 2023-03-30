package admin

import (
	"fmt"
	"manage/dao/mysql/admin"
	"manage/dao/mysql/user"
	"manage/hooks"
	"manage/model"
	"time"

	"go.uber.org/zap"
)

func GetAllComment(p *model.ParamAdminComment) (commentList []*model.Comment, total int, err error) {
	commentList, total, err = admin.GetAllComment(p.Page, p.Size, time.Now().String())
	if err != nil {
		zap.L().Error("admin.GetAllComment failed", zap.Error(err))
		return
	}
	for _, comment := range commentList {
		if comment.ItemType == 2 {
			comment.Type = "文章评论"
		} else if comment.ItemType == 1 {
			comment.Type = "回复评论"
		}
		comment.Format = hooks.DateDay(comment.CreateTime)
		comment.UserName, err = admin.GetUserByID(comment.UserID)
		if err != nil {
			zap.L().Error("admin.GetUserByID failed", zap.Int64("userID", comment.UserID), zap.Error(err))
			continue
		}
	}
	return
}

func GetCommentByItemID(p *model.ParamAdminComment) (commentList []*model.Comment, total int, err error) {
	commentList, total, err = admin.GetCommentByItemID(p.ItemID, p.Page, p.Size, time.Now().String())
	if err != nil {
		zap.L().Error("admin.GetAllComment failed", zap.Error(err))
		return
	}
	for _, comment := range commentList {
		if comment.ItemType == 2 {
			comment.Type = "文章评论"
		} else if comment.ItemType == 1 {
			comment.Type = "回复评论"
		}
		comment.Format = hooks.DateDay(comment.CreateTime)
		comment.UserName, err = admin.GetUserByID(comment.UserID)
		if err != nil {
			zap.L().Error("admin.GetUserByID failed", zap.Int64("userID", comment.UserID), zap.Error(err))
			continue
		}
	}
	return
}

func GetCommentInfo(commentID int64) (commentInfo model.ApiCommentInfo, err error) {
	comment, err := admin.GetCommentInfo(commentID)
	if err != nil {
		zap.L().Error("admin.GetCommentInfo failed", zap.Error(err))
		return
	}
	isHave := admin.CheckSuperiorComment(commentID)
	if isHave {
		superiorComment, err1 := admin.GetSuperiorCommentInfo(commentID)
		if err1 != nil {
			zap.L().Error("admin.GetSuperiorCommentInfo failed", zap.Error(err1))
			return
		}
		if superiorComment.UserID != 0 {
			superiorComment.Format = hooks.DateDay(superiorComment.CreateTime)
			superiorComment.UserName, err = admin.GetUserByID(superiorComment.UserID)
			if err != nil {
				zap.L().Error("admin.GetUserByID failed", zap.Int64("userID", superiorComment.UserID), zap.Error(err))
			}
		}
		commentInfo.SuperiorCommentInfo = &superiorComment
	}
	comment.Format = hooks.DateDay(comment.CreateTime)
	commentInfo.CommentInfo = &comment
	commentInfo.HaveSuperior = isHave
	return
}

func SetCommentStatus(commentID int64, status int) (err error) {
	if status == 1 {
		err = admin.BanStatus(commentID)
		if err != nil {
			zap.L().Error("admin.BanStatus failed", zap.Error(err))
			return
		}
	} else if status == 0 {
		err = admin.AllowStatus(commentID)
		if err != nil {
			zap.L().Error("admin.BanStatus failed", zap.Error(err))
			return
		}
	}
	return
}

func DeleteCommentForAdmin(p *model.ParamAdminDeleteComment) (err error) {
	willDelete := make([]int64, 0)
	switch p.ItemType {
	case model.ItemArticle:
		err = user.DeleteRootComment(p.ID)
	case model.ItemComment:
		// 查询是否有回复
		isLeaf, err1 := user.IsLeafComment(p.ID)
		if err1 != nil {
			err = err1
			return
		}
		// 有回复，将评论状态设为0
		if !isLeaf {
			err = user.SetReplyComment(p.ID)
		} else {
			// 将本身加入待删除队列
			willDelete = append(willDelete, p.ID)
			// 若没有回复，查询所有祖先
			ancestors, err2 := user.GetAncestorReply(p.ID)
			if err2 != nil {
				err = err2
				return
			}
			fmt.Println(ancestors)
			// 若祖先回复满足 1.已删除 2.不存在多个回复 加入待删除队列
			// 若祖先中不满足以上条件，停止遍历
			for _, c := range ancestors {
				//判断是否删除
				isDeleted, err3 := user.IsReplyDeleted(c)
				if err3 != nil {
					err = err3
					break
				}
				//若未删除,直接结束循环
				if !isDeleted {
					break
				} else {
					// 若删除，查询是否有多个回复
					haveReplies, err4 := user.HaveReplies(c)
					if err4 != nil {
						err = err4
						return
					}
					if !haveReplies {
						// 若只有一个回复，且这个回复将要被删除
						willDelete = append(willDelete, c)
					} else {
						// 有多个回复，不能影响到其他为被删除的回复,结束循环
						break
					}
				}
			}
			// 删除操作
			for _, commentID := range willDelete {
				if err = user.DeleteReplyComment(commentID); err != nil {
					return
				}
			}
		}
	default:
		return
	}
	return
}
