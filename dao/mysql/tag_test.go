package mysql

import (
	"manage/model"
	"manage/setting"
	"testing"
)

func init() {
	dbCfg := &setting.MySQLConfig{
		Host:        "8.134.222.37",
		User:        "root",
		Password:    "Tang@1112",
		DbName:      "manage",
		Port:        3306,
		MaxConn:     10,
		MaxIdleConn: 10,
	}
	if err := Init(dbCfg); err != nil {
		panic(err)
	}
}

func TestGetTagNameByArticleID(t *testing.T) {
	var id int64 = 2074895122960384
	tagName, err := GetTagNameByArticleID(id)
	if err != nil {
		t.Fatalf("GetTagNameByArticleID failed , err=%v\n", err)
	}
	t.Logf("success - %v", tagName)
}

func TestGetArticleIDByTagID(t *testing.T) {
	id := 1
	articleID, err := GetArticleIDByTagID(id)
	if err != nil {
		t.Fatalf("GetArticleIDByTagID failed , err=%v\n", err)
	}
	t.Logf("success - %v", articleID)
}

func TestGetTagIDByTagName(t *testing.T) {
	name := "Golang"
	id, err := GetTagIDByTagName(name)
	if err != nil {
		t.Fatalf("GetTagIDByTagName failed , err=%v\n", err)
	}
	t.Logf("success - %v", id)
}

func TestGetSearchTags(t *testing.T) {
	p := &model.ParamSearch{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
		Category:      model.SearchTag,
		Key:           "go",
	}
	tagList, total, err := GetSearchTags(p)
	if err != nil {
		t.Fatalf("GetSearchTags failed , err=%v\n", err)
	}
	t.Logf("success - %v - %v", tagList, total)
}
