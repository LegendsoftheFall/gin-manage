package controller

import (
	"manage/logic"
	"manage/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateDraftHandler 处理创建草稿请求的处理函数
func CreateDraftHandler(c *gin.Context) {
	p := new(model.Draft)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SaveDraft with invalid param", zap.Error(err))
		validateError(c, err)
		return
	}
	userID, _ := getCurrentUserID(c)
	p.AuthorID = userID
	id, err := logic.CreateDraft(p)
	if err != nil {
		zap.L().Error("logic.SaveDraft(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"id": strconv.FormatInt(id, 10),
	})
}

// SaveDraftHandler 处理保存草稿请求的处理函数
func SaveDraftHandler(c *gin.Context) {
	p := new(model.Draft)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdateDraft with invalid param", zap.Error(err))
		validateError(c, err)
		return
	}
	userID, _ := getCurrentUserID(c)
	p.AuthorID = userID
	err := logic.SaveDraft(p)
	if err != nil {
		zap.L().Error("logic.UpdateDraft(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

func DeleteDraftHandler(c *gin.Context) {
	draftID := c.Query("id")
	userID, _ := getCurrentUserID(c)
	if err := logic.DeleteDraft(draftID, strconv.FormatInt(userID, 10)); err != nil {
		zap.L().Error("logic.DeleteDraft failed",
			zap.String("draftID", draftID),
			zap.Int64("userID", userID),
			zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

func DeleteAllDraftHandler(c *gin.Context) {
	userID, _ := getCurrentUserID(c)
	if err := logic.DeleteAllDraft(strconv.FormatInt(userID, 10)); err != nil {
		zap.L().Error("logic.DeleteDraft failed",
			zap.Int64("userID", userID),
			zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

func DraftsHandler(c *gin.Context) {
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeInvalidAuth)
		return
	}
	draftInfoList, err := logic.GetDraftInfoByUserID(strconv.FormatInt(userID, 10))
	if err != nil {
		zap.L().Error("logic.GetDraftInfoByUserID failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, gin.H{
		"list": draftInfoList,
	})
}

func DraftHandler(c *gin.Context) {
	draftID := c.Param("id")
	draftDetail, err := logic.GetDraftDetailByID(draftID)
	if err != nil {
		if err != nil {
			zap.L().Error("logic.GetDraftDetailByID failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
	}
	ResponseSuccess(c, gin.H{
		"draft": draftDetail,
	})
}
