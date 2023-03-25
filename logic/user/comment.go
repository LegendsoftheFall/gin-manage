package user

import (
	"fmt"
	"manage/dao/mysql/user"
	"manage/hooks"
	"manage/model"
	"manage/pkg/snowflake"
	"time"

	"go.uber.org/zap"
)

func CreateComment(p *model.ParamComment) (commentID int64, err error) {
	p.CommentID = snowflake.GenID()
	commentID = p.CommentID
	switch p.ItemType {
	case model.ItemArticle:
		err = user.CreateRootComment(p)
	case model.ItemComment:
		err = user.CreateReplyComment(p)
	default:
		return
	}
	return
}

func DeleteComment(p *model.ParamDeleteComment) (err error) {
	willDelete := make([]int64, 0)
	switch p.ItemType {
	case model.ItemArticle:
		err = user.DeleteRootComment(p.CommentID)
	case model.ItemComment:
		// 查询是否有回复
		isLeaf, err1 := user.IsLeafComment(p.CommentID)
		if err1 != nil {
			err = err1
			return
		}
		// 有回复，将评论状态设为0
		if !isLeaf {
			err = user.SetReplyComment(p.CommentID)
		} else {
			// 将本身加入待删除队列
			willDelete = append(willDelete, p.CommentID)
			// 若没有回复，查询所有祖先
			ancestors, err2 := user.GetAncestorReply(p.CommentID)
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

func GetCommentList(p *model.ParamCommentList) (commentList []*model.ApiComment, total int, err error) {
	// 取得目标ID 根据分页参数获取根评论切片
	rootComment, commentTotal, err := user.GetRootComment(p.ItemID, p.Page, p.Size, p.Order, time.Now().String())
	if err != nil {
		zap.L().Error("mysql.GetRootComment failed", zap.Error(err))
		return
	}
	total = commentTotal
	// 根评论和回复评论的映射关系  根评论 -> (回复ID -> 回复)
	commentReplyMap := make(map[int64]map[int64]*model.ReplyComment, len(rootComment))
	// 根评论和回复的切片映射关系
	replySliceMap := make(map[int64][]*model.ReplyComment, len(rootComment))
	// 用户ID与用户信息的映射关系
	userInfoMap := make(map[int64]*model.UserInfo)

	for _, comment := range rootComment {
		// 判断map中是否存在对应userID
		_, ok := userInfoMap[comment.UserID]
		if !ok {
			userInfoMap[comment.UserID], err = user.GetUserInfoByUserID(comment.UserID)
			if err != nil {
				zap.L().Error("mysql.GetUserInfoByUserID",
					zap.Int64("userID", comment.UserID))
			}
		}
		// 获取根评论的所有回复
		replyComments, err1 := user.GetReplyComment(comment.CommentID)
		fmt.Println(replyComments)
		if err1 != nil {
			zap.L().Error("mysql.GetReplyComment",
				zap.Int64("commentID", comment.CommentID),
				zap.Error(err1))
			continue
		}
		// RootCommentInfo中的回复切片
		replySliceMap[comment.CommentID] = make([]*model.ReplyComment, 0)
		replySliceMap[comment.CommentID] = replyComments
		// 将获取的所有回复根据ID保存到map中
		replyMap := make(map[int64]*model.ReplyComment, len(replyComments))
		for _, reply := range replyComments {
			replyMap[reply.ReplyID] = reply
			// 回复用户
			_, ok = userInfoMap[reply.UserID]
			if !ok {
				userInfoMap[reply.UserID], err = user.GetUserInfoByUserID(reply.UserID)
				if err != nil {
					zap.L().Error("mysql.GetUserInfoByUserID",
						zap.Int64("userID", reply.UserID))
					continue
				}
			}
			// 被回复用户
			_, ok = userInfoMap[reply.ToUserID]
			if !ok {
				// 有回复对象
				if reply.ToUserID != 0 {
					userInfoMap[reply.ToUserID], err = user.GetUserInfoByUserID(reply.ToUserID)
					if err != nil {
						zap.L().Error("mysql.GetUserInfoByUserID",
							zap.Int64("userID", reply.ToUserID))
						continue
					}
				}
			}
		}
		commentReplyMap[comment.CommentID] = replyMap
	}

	commentList = make([]*model.ApiComment, 0, len(rootComment))
	for _, comment := range rootComment {
		list := &model.ApiComment{
			CommentID:   0,
			CommentInfo: nil,
			UserInfo:    nil,
			ReplyInfos:  []*model.ReplyCommentInfo{},
		}

		comment.Format = hooks.TimeSub(time.Now(), comment.CreatTime)
		// json:"comment_id"
		list.CommentID = comment.CommentID
		// json:"comment_info"
		commentInfo := &model.RootCommentInfo{
			ReplyCount:     0,
			RootComment:    comment,
			CommentReplies: []*model.ReplyComment{},
		}
		if replySliceMap[comment.CommentID] != nil {
			commentInfo.CommentReplies = replySliceMap[comment.CommentID]
		}
		commentInfo.ReplyCount = len(commentInfo.CommentReplies)
		list.CommentInfo = commentInfo
		// json:"user_info"
		list.UserInfo = userInfoMap[comment.UserID]
		// json:"reply_infos"
		replies := replySliceMap[comment.CommentID]
		for _, reply := range replies {
			reply.Format = hooks.TimeSub(time.Now(), reply.CreatTime)
			replyInfo := new(model.ReplyCommentInfo)
			replyInfo.ReplyID = reply.ReplyID
			replyInfo.ReplyInfo = commentReplyMap[comment.CommentID][reply.ReplyID]
			replyInfo.UserInfo = userInfoMap[reply.UserID]
			// 若回复评论和对应的用户ID是根评论，说明是子评论，将父回复和回复用户置为空
			if reply.ToReplyID == 0 && reply.ToUserID == 0 {
				replyInfo.ParentReply = nil
				replyInfo.ReplyUser = nil
			} else {
				replyInfo.ParentReply = commentReplyMap[comment.CommentID][reply.ToReplyID]
				replyInfo.ReplyUser = userInfoMap[reply.ToUserID]
			}
			list.ReplyInfos = append(list.ReplyInfos, replyInfo)
		}
		// 组合
		commentList = append(commentList, list)
	}
	return
}
