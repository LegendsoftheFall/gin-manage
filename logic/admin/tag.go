package admin

import (
	"manage/dao/mysql/admin"
	"manage/model"

	"go.uber.org/zap"
)

func CreateTagForAdmin(p *model.ParamCreateTag) (err error) {
	//判断标签是否存在
	if err = admin.CheckTagExist(p.Name); err != nil {
		return err
	}
	if err = admin.IsLTTotal(p.ID); err != nil {
		return err
	}
	return admin.CreateTag(p)
}

func EditTag(p *model.ParamCreateTag) (err error) {
	//判断标签是否存在
	if err = admin.CheckTagExist(p.Name); err != nil {
		return err
	}
	return admin.UpdateTag(p)
}

func DeleteTag(id int64) (err error) {
	return admin.DeleteTag(id)
}

func GetAllTags(p *model.ParamTagList) (tagList []*model.TagDetail, total int, err error) {
	tagList, total, err = admin.GetAllTags(p)
	if err != nil {
		zap.L().Error("mysql.GetAllTags failed",
			zap.Int64("userID", p.CurrentUserID),
			zap.Error(err))
		return nil, 0, err
	}
	return
}

func SelectTags() (tagList []*model.Tag, err error) {
	return admin.SelectTags()
}
