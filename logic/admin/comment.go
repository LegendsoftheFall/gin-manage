package admin

import (
	"manage/dao/mysql/admin"
	"manage/model"
	"time"

	"go.uber.org/zap"
)

func GetAllComment(p *model.ParamAdminComment) (commentList []*model.Comment, total int, err error) {
	commentList, total, err = admin.GetAllComment(p.Page, p.Size, p.Order, time.Now().String())
	if err != nil {
		zap.L().Error("admin.GetAllComment failed", zap.Error(err))
		return
	}
	for _, comment := range commentList {
		comment.UserName, err = admin.GetUserByID(comment.UserID)
		if err != nil {
			zap.L().Error("admin.GetUserByID failed", zap.Int64("userID", comment.UserID), zap.Error(err))
			continue
		}
	}
	return
}
