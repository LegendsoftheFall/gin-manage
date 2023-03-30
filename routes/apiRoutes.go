package routes

import (
	"manage/controller"
	"manage/controller/user"
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
		apiRoutes.GET("/userInfo", user.InfoOfUserHandler)
		apiRoutes.GET("/user/profile", user.ProfileOfUserHandler)
		apiRoutes.PATCH("/user/profile/update", user.UpdateProfileHandler)
	}
	{
		apiRoutes.GET("/selectTags", user.SelectTagsHandler)
	}
	{
		apiRoutes.POST("/upload", controller.UpLoadHandler)
		apiRoutes.POST("/createArticle", user.CreateArticleHandler)
		apiRoutes.PATCH("/editArticle", user.EditArticleHandler)
		apiRoutes.DELETE("/deleteArticle/:id", user.DeleteArticleHandler)
	}
	{
		apiRoutes.POST("/createDraft", user.CreateDraftHandler)
		apiRoutes.POST("/saveDraft", user.SaveDraftHandler)
		apiRoutes.POST("/deleteDraft", user.DeleteDraftHandler)
		apiRoutes.POST("/deleteAllDraft", user.DeleteAllDraftHandler)
		apiRoutes.GET("/drafts", user.DraftsHandler)
		apiRoutes.GET("/draft/:id", user.DraftHandler)
	}
	{
		apiRoutes.POST("/like", user.LikeHandler)
	}
	{
		apiRoutes.POST("/collect", user.CollectHandler)
		apiRoutes.GET("/bookmark", user.BookMarkHandler)
	}
	{
		apiRoutes.POST("tag/follow/do", user.TagDoFollowHandler)
		apiRoutes.POST("tag/follow/undo", user.TagUnDoFollowHandler)
		apiRoutes.POST("user/follow/do", user.DoFollowHandler)
		apiRoutes.POST("user/follow/undo", user.UnDoFollowHandler)
	}
	{
		apiRoutes.POST("/comment/create", user.CreateCommentHandler)
		apiRoutes.POST("/comment/delete", user.DeleteCommentHandler)
	}
}
