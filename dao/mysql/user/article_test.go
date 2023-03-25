package user

import (
	"manage/model"
	"testing"
)

func TestGetArticleListByIDs(t *testing.T) {
	ids := make([]string, 0)
	ids = append(ids, "2074895122960384")
	ids = append(ids, "2324781982552064")

	articleList, err := GetArticleListByIDs(ids)
	if err != nil {
		t.Fatalf("GetArticleListByIDs failed , err=%v\n", err)
	}
	t.Logf("success - %v", articleList)
}

func TestGetArticleListByID(t *testing.T) {
	var id int64 = 1006167354511360
	articleList, err := GetArticleListByID(id, 1, 10)
	if err != nil {
		t.Fatalf("GetArticleListByID failed , err=%v\n", err)
	}
	t.Logf("success - %v", articleList)
}

func TestGetUserIDByArticleID(t *testing.T) {
	var id int64 = 2074895122960384
	userID, err := GetUserIDByArticleID(id)
	if err != nil {
		t.Fatalf("GetUserIDByArticleID failed , err=%v\n", err)
	}
	t.Logf("success - %v", userID)
}

func TestGetArticleNumByID(t *testing.T) {
	var id int64 = 1006167354511360
	num, err := GetArticleNumByID(id)
	if err != nil {
		t.Fatalf("GetArticleNumByID failed , err=%v\n", err)
	}
	t.Logf("success - %v", num)
}

func TestGetSearchArticleList(t *testing.T) {
	p := &model.ParamSearch{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
		Category:      model.SearchTop,
		Key:           "e",
	}
	articleList, total, err := GetSearchArticleListByTime(p)
	if err != nil {
		t.Fatalf("GetSearchArticleList failed , err=%v\n", err)
	}
	t.Logf("success - %v - %v", articleList, total)
}
