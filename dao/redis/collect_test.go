package redis

import (
	"manage/model"
	"manage/setting"
	"testing"
)

func init() {
	rdbCfg := &setting.RedisConfig{
		Host:     "8.134.222.37",
		Password: "663117",
		DB:       0,
		Port:     6397,
		PoolSize: 100,
	}
	if err := Init(rdbCfg); err != nil {
		panic(err)
	}
}

func TestGetArticleIDsByUserID(t *testing.T) {
	uid := "1893129523302400"
	p := &model.ParamPage{
		Page: 1,
		Size: 10,
	}
	ids, err := GetArticleIDsByUserID(uid, p)
	if err != nil {
		t.Fatalf("GetArticleIDsByUserID failed , err=%v\n", err)
	}
	t.Logf("success - %v", ids)
}
