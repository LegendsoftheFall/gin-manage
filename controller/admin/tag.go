package admin

import (
	"manage/controller"
	"manage/logic/admin"
	"manage/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var MySecret = "202010214303"

// CreateTagHandler 新建标签请求的处理函数
func CreateTagHandler(c *gin.Context) {
	// 参数校验
	p := new(model.ParamCreateTag)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreateTag with invalid param", zap.Error(err))
		controller.ValidateError(c, err)
		return
	}

	//创建标签
	if err := admin.CreateTagForAdmin(p); err != nil {
		zap.L().Error("logic.CreateTag(p) failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeTagExist)
		return
	}
	//返回响应
	controller.ResponseSuccess(c, nil)

}

// EditTagHandler 编辑标签请求的处理函数
func EditTagHandler(c *gin.Context) {
	p := new(model.ParamCreateTag)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("EditTag with invalid param", zap.Error(err))
		controller.ValidateError(c, err)
		return
	}
	if err := admin.EditTag(p); err != nil {
		zap.L().Error("logic.EditArticle(p) failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeTagExist)
		return
	}

	controller.ResponseSuccess(c, nil)
}

// DeleteTagHandler 删除标签请求的处理函数
func DeleteTagHandler(c *gin.Context) {
	p := new(model.ParamDeleteTag)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("DeleteTag with invalid param", zap.Error(err))
		controller.ValidateError(c, err)
		return
	}
	if p.Secret != MySecret {
		controller.ResponseError(c, controller.CodeInvalidAuth)
		return
	}

	if err := admin.DeleteTag(p.ID); err != nil {
		zap.L().Error("logic.DeleteTag(id) failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, nil)
}

// TagsForAdminHandler 获取全部标签请求的处理函数
func TagsForAdminHandler(c *gin.Context) {
	p := &model.ParamTagList{
		Page:          1,
		Size:          10,
		CurrentUserID: 0,
	}
	//获取分页参数
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("TagsHandler with invalid param", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvalidParam)
		return
	}

	data, total, err := admin.GetAllTags(p)
	if err != nil {
		zap.L().Error("logic.GetTrendingTags failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, gin.H{
		"list":  data,
		"total": total,
		"page":  p.Page,
		"size":  p.Size,
	})
}

// SelectTagsForAdminHandler 获取选择标签请求的处理函数
func SelectTagsForAdminHandler(c *gin.Context) {
	data, err := admin.SelectTags()
	if err != nil {
		zap.L().Error("logic.GetTrendingTags failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeServerBusy)
		return
	}
	controller.ResponseSuccess(c, data)
}
