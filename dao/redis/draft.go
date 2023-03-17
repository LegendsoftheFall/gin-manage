package redis

import (
	"time"

	"go.uber.org/zap"
)

const ExpireTime = 720 * time.Hour

func CreateDraft(did, uid string, data []byte) (err error) {
	pipeline := rdb.TxPipeline()
	pipeline.Set(getRedisKey(KeyDraft+did), data, ExpireTime)
	pipeline.SAdd(getRedisKey(KeyUserDraft+uid), did)
	_, err = pipeline.Exec()
	return
}

func SaveDraft(did string, data []byte) (err error) {
	err = rdb.Set(getRedisKey(KeyDraft+did), data, ExpireTime).Err()
	return
}

func DeleteDraft(did, uid string) (err error) {
	pipeline := rdb.TxPipeline()
	pipeline.SRem(getRedisKey(KeyUserDraft+uid), did)
	pipeline.Del(getRedisKey(KeyDraft + did))
	_, err = pipeline.Exec()
	return
}

func DeleteAllDraft(uid string) (err error) {
	ids, err := rdb.SMembers(getRedisKey(KeyUserDraft + uid)).Result()
	if err != nil {
		zap.L().Error("rdb.SMembers failed",
			zap.String("key", getRedisKey(KeyUserDraft+uid)),
			zap.Error(err))
	}
	if len(ids) == 0 {
		return nil
	}
	pipeline := rdb.TxPipeline()
	for _, id := range ids {
		pipeline.Del(getRedisKey(KeyDraft + id))
	}
	pipeline.Del(getRedisKey(KeyUserDraft + uid))
	_, err = pipeline.Exec()
	return
}

func GetDraftInfoByUserID(uid string) (dataList [][]byte, err error) {
	ids, err := rdb.SMembers(getRedisKey(KeyUserDraft + uid)).Result()
	if err != nil {
		zap.L().Error("rdb.SMembers failed",
			zap.String("key", getRedisKey(KeyUserDraft+uid)),
			zap.Error(err))
	}
	if len(ids) == 0 {
		return nil, err
	}
	dataList = make([][]byte, 0, len(ids))

	for i, id := range ids {
		data, ok := rdb.Get(getRedisKey(KeyDraft + ids[i])).Bytes()
		if ok != nil {
			zap.L().Error("rdb.Get failed",
				zap.String("key", getRedisKey(KeyDraft+id)),
				zap.Error(err))
			continue
		}
		dataList = append(dataList, data)
	}

	return
}

// GetDraftDetailByID 根据draftID获取draftDetail
func GetDraftDetailByID(did string) (data string, err error) {
	data, err = rdb.Get(getRedisKey(KeyDraft + did)).Result()
	return
}
