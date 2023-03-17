package routes

import (
	"manage/controller"
	"manage/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserIDHandler(c *gin.Context) {
	userID, _ := c.Get("userID")
	c.JSON(http.StatusOK, gin.H{
		"code": 2000,
		"msg":  "success",
		"data": gin.H{"id": userID},
	})
}

func apiRoutesInit(r *gin.Engine) {
	apiRoutes := r.Group("/api", middleware.JWTAuthMiddleware())
	{
		apiRoutes.GET("/user_id", UserIDHandler)
		apiRoutes.GET("/userInfo", controller.UserInfoHandler)
		apiRoutes.GET("/user/profile", controller.UserProfileHandler)
		apiRoutes.PATCH("/user/profile/update", controller.UpdateProfileHandler)
	}
	{
		apiRoutes.GET("/selectTags", controller.SelectTagsHandler)
	}
	{
		apiRoutes.POST("/upload", controller.UpLoadHandler)
		apiRoutes.POST("/createArticle", controller.CreateArticleHandler)
		apiRoutes.PATCH("/editArticle", controller.EditArticleHandler)
		apiRoutes.DELETE("/deleteArticle/:id", controller.DeleteArticleHandler)
	}
	{
		apiRoutes.POST("/createDraft", controller.CreateDraftHandler)
		apiRoutes.POST("/saveDraft", controller.SaveDraftHandler)
		apiRoutes.POST("/deleteDraft", controller.DeleteDraftHandler)
		apiRoutes.POST("/deleteAllDraft", controller.DeleteAllDraftHandler)
		apiRoutes.GET("/drafts", controller.DraftsHandler)
		apiRoutes.GET("/draft/:id", controller.DraftHandler)
	}
	{
		apiRoutes.POST("/like", controller.LikeHandler)
	}
	{
		apiRoutes.POST("/collect", controller.CollectHandler)
		apiRoutes.GET("/bookmark", controller.BookMarkHandler)
	}
	{
		apiRoutes.POST("tag/follow/do", controller.TagDoFollowHandler)
		apiRoutes.POST("tag/follow/undo", controller.TagUnDoFollowHandler)
		apiRoutes.POST("user/follow/do", controller.UserDoFollowHandler)
		apiRoutes.POST("user/follow/undo", controller.UserUnDoFollowHandler)
	}
	{
		apiRoutes.POST("/comment/create", controller.CreateCommentHandler)
		apiRoutes.POST("/comment/delete", controller.DeleteCommentHandler)
	}
}
