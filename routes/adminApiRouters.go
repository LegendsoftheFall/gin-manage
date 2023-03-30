package routes

import (
	"manage/controller/admin"
	"manage/controller/user"
	"manage/middleware"

	"github.com/gin-gonic/gin"
)

func adminApiRoutesInit(r *gin.Engine) {
	adminApiRoutes := r.Group("/admin/api", middleware.JWTAuthMiddleware())
	{
		adminApiRoutes.GET("/info", admin.InfoOfAdminHandler)
		adminApiRoutes.POST("/logout", admin.LogoutForAdminHandler)
	}
	{
		adminApiRoutes.POST("/tag/create", admin.CreateTagHandler)
		adminApiRoutes.PATCH("/tag/update", admin.EditTagHandler)
		adminApiRoutes.POST("/tag/delete", admin.DeleteTagHandler)
		adminApiRoutes.GET("/tags", admin.TagsForAdminHandler)
		adminApiRoutes.GET("/tag/select", admin.SelectTagsForAdminHandler)
		adminApiRoutes.GET("/tag/:id", user.TagDetailHandler)
	}
	{
		adminApiRoutes.GET("/article/:id", user.ArticleHandler)
		adminApiRoutes.POST("/article/delete", admin.DeleteArticleForAdminHandler)
	}
	{
		adminApiRoutes.GET("/comment/all", admin.GetAllCommentHandler)
		adminApiRoutes.GET("/comment/info/:id", admin.GetCommentInfoHandler)
		adminApiRoutes.GET("/comment/search", admin.GetCommentByItemIDHandler)
		adminApiRoutes.POST("/comment/set", admin.SetCommentStatusHandler)
		adminApiRoutes.POST("/comment/delete", admin.DeleteCommentForAdminHandler)
	}
}
