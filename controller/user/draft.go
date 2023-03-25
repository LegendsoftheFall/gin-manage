package user

import (
	"manage/controller"
	"manage/logic/user"
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
		controller.ValidateError(c, err)
		return
	}
	userID, _ := controller.GetCurrentUserID(c)
	p.AuthorID = userID
	id, err := user.CreateDraft(p)
	if err != nil {
		zap.L().Error("logic.SaveDraft(p) failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
		"id": strconv.FormatInt(id, 10),
	})
}

// SaveDraftHandler 处理保存草稿请求的处理函数
func SaveDraftHandler(c *gin.Context) {
	p := new(model.Draft)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdateDraft with invalid param", zap.Error(err))
		controller.ValidateError(c, err)
		return
	}
	userID, _ := controller.GetCurrentUserID(c)
	p.AuthorID = userID
	err := user.SaveDraft(p)
	if err != nil {
		zap.L().Error("logic.UpdateDraft(p) failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, nil)
}

func DeleteDraftHandler(c *gin.Context) {
	draftID := c.Query("id")
	userID, _ := controller.GetCurrentUserID(c)
	if err := user.DeleteDraft(draftID, strconv.FormatInt(userID, 10)); err != nil {
		zap.L().Error("logic.DeleteDraft failed",
			zap.String("draftID", draftID),
			zap.Int64("userID", userID),
			zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, nil)
}

func DeleteAllDraftHandler(c *gin.Context) {
	userID, _ := controller.GetCurrentUserID(c)
	if err := user.DeleteAllDraft(strconv.FormatInt(userID, 10)); err != nil {
		zap.L().Error("logic.DeleteDraft failed",
			zap.Int64("userID", userID),
			zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, nil)
}

func DraftsHandler(c *gin.Context) {
	userID, err := controller.GetCurrentUserID(c)
	if err != nil {
		controller.ResponseError(c, controller.CodeInvalidAuth)
		return
	}
	draftInfoList, err := user.GetDraftInfoByUserID(strconv.FormatInt(userID, 10))
	if err != nil {
		zap.L().Error("logic.GetDraftInfoByUserID failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
		"list": draftInfoList,
	})
}

func DraftHandler(c *gin.Context) {
	draftID := c.Param("id")
	draftDetail, err := user.GetDraftDetailByID(draftID)
	if err != nil {
		if err != nil {
			zap.L().Error("logic.GetDraftDetailByID failed", zap.Error(err))
			controller.ResponseError(c, controller.CodeServerBusy)
			return
		}
	}
	controller.ResponseSuccess(c, gin.H{
		"draft": draftDetail,
	})
}
