package mysql

import (
	"fmt"
	"testing"
)

func TestDeleteRootComment(t *testing.T) {
	var commentID int64 = 15004403886985216
	if err := DeleteRootComment(commentID); err != nil {
		t.Fatalf("DeleteRootComment failed , err=%v\n", err)
	}
	t.Log("success")
}

func TestIsLeafComment(t *testing.T) {
	var commentID int64 = 15289892149923840
	isLeaf, err := IsLeafComment(commentID)
	if err != nil {
		t.Fatalf("IsLeafComment failed , err=%v\n", err)
	}
	t.Logf("success -- %v", isLeaf)
}

func TestHaveReplies(t *testing.T) {
	var commentID int64 = 15289835799449600
	have, err := HaveReplies(commentID)
	if err != nil {
		t.Fatalf("HaveReplies failed , err=%v\n", err)
	}
	t.Logf("success -- %v", have)
}

func TestGetAncestorReply(t *testing.T) {
	var commentID int64 = 15290072077176832
	replies, err := GetAncestorReply(commentID)
	if err != nil {
		t.Fatalf("GetAncestorReply failed , err=%v\n", err)
	}
	t.Logf("success -- %v", replies)
}

func TestGetRootComment(t *testing.T) {
	var (
		itemID  int64 = 14683242602958848
		page    int64 = 1
		size    int64 = 5
		order         = "score"
		endTime       = "2022-11-20T16:25:57Z"
		//endTime = time.Now().String()
	)
	fmt.Println(endTime)
	rootComment, total, err := GetRootComment(itemID, page, size, order, endTime)
	if err != nil {
		t.Fatalf("GetRootComment failed , err=%v\n", err)
	}
	t.Logf("success comment: %v, total: %d ", rootComment, total)
}

func TestGetReplyComment(t *testing.T) {
	var rootCommentID int64 = 16451466701049856
	replyComment, err := GetReplyComment(rootCommentID)
	if err != nil {
		t.Fatalf("GetReplyComment failed , err=%v\n", err)
	}
	t.Logf("success comment: %v ", replyComment)
}
