package mysql

import (
	"manage/model"
	"testing"
)

func TestGetSearchUsers(t *testing.T) {
	p := &model.ParamSearch{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
		Category:      model.SearchUser,
		Key:           "测试",
	}
	userList, total, err := GetSearchUsers(p)
	if err != nil {
		t.Fatalf("GetSearchUsers failed , err=%v\n", err)
	}
	t.Logf("success - %v - %v", userList, total)
}
