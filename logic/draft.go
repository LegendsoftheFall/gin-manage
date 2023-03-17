package logic

import (
	"encoding/json"
	"fmt"
	"manage/dao/mysql"
	"manage/dao/redis"
	"manage/hooks"
	"manage/model"
	"manage/pkg/snowflake"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func CreateDraft(p *model.Draft) (did int64, err error) {
	p.ID = snowflake.GenID()
	p.Content = hooks.TrimHtml(p.Html)
	p.CreateTime = time.Now()
	did = p.ID
	data, err := json.Marshal(p)
	fmt.Println(data)
	if err != nil {
		return
	}
	err = redis.CreateDraft(strconv.FormatInt(p.ID, 10),
		strconv.FormatInt(p.AuthorID, 10),
		data)
	return
}

func SaveDraft(p *model.Draft) (err error) {
	p.CreateTime = time.Now()
	p.Content = hooks.TrimHtml(p.Html)
	data, err := json.Marshal(p)
	if err != nil {
		return
	}
	return redis.SaveDraft(strconv.FormatInt(p.ID, 10), data)
}

func DeleteDraft(did, uid string) (err error) {
	return redis.DeleteDraft(did, uid)
}

func DeleteAllDraft(uid string) (err error) {
	return redis.DeleteAllDraft(uid)
}

func GetDraftInfoByUserID(uid string) (draftInfoList []*model.DraftInfo, err error) {
	dataList, err := redis.GetDraftInfoByUserID(uid)
	if err != nil {
		zap.L().Error("redis.GetDraftInfoByUserID failed",
			zap.String("userID", uid),
			zap.Error(err))
		return
	}
	draftInfoList = make([]*model.DraftInfo, 0, len(dataList))
	for _, data := range dataList {
		draftInfo := new(model.DraftInfo)
		ok := json.Unmarshal(data, draftInfo)
		if ok != nil {
			zap.L().Error("json.Unmarshal failed",
				zap.ByteString("[]byte", data),
				zap.Error(ok))
			continue
		}
		draftInfoList = append(draftInfoList, draftInfo)
	}
	for _, draft := range draftInfoList {
		draft.Format = hooks.DateDay(draft.CreateTime)
		content := []rune(draft.Content)
		if len(content) > 270 {
			draft.Content = string(content[:270])
		}
	}
	return
}

func GetDraftDetailByID(did string) (draftDetail *model.DraftDetail, err error) {
	data, err := redis.GetDraftDetailByID(did)
	if err != nil {
		zap.L().Error("redis.GetDraftDetailByID failed",
			zap.String("draftID", did),
			zap.Error(err))
		return
	}
	b := hooks.String2Bytes(data)
	err = json.Unmarshal(b, &draftDetail)
	if err != nil {
		zap.L().Error("json.Unmarshal failed",
			zap.ByteString("[]byte", b),
			zap.Error(err))
	}
	tags := make([]*model.TagSimple, 0, len(draftDetail.Tags))
	for _, t := range draftDetail.Tags {
		name, ok := mysql.GetTagNameByTagID(t)
		if ok != nil {
			zap.L().Error("mysql.GetTagNameByTagID failed",
				zap.Int("tagID", t),
				zap.Error(err))
			continue
		}
		tag := &model.TagSimple{
			ID:   t,
			Name: name,
		}
		tags = append(tags, tag)
	}
	draftDetail.TagSimple = tags
	draftDetail.Format = hooks.DateDay(draftDetail.CreateTime)
	return
}
