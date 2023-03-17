package routes

import (
	"manage/controller"

	"github.com/gin-gonic/gin"
)

func defaultRoutesInit(r *gin.Engine) {
	defaultRoutes := r.Group("/")
	{
		defaultRoutes.POST("/signup", controller.SignUpHandler)
		defaultRoutes.POST("/login", controller.LoginHandler)
		defaultRoutes.POST("/refresh", controller.RefreshHandler)
	}
	{
		defaultRoutes.GET("/trendingTags", controller.TrendingTagHandler)
		defaultRoutes.GET("/tag", controller.TagInfoHandler)
	}
	{
		defaultRoutes.GET("/user/:id", controller.UserHomeHandler)
		defaultRoutes.GET("/user/profile", controller.ProfileHandler)
		defaultRoutes.GET("/user/following", controller.FollowingUsersHandler)
		defaultRoutes.GET("/user/followers", controller.FollowerUsersHandler)
	}
	{
		defaultRoutes.GET("/article/:id", controller.ArticleHandler)
	}
	{
		defaultRoutes.GET("/comment/list", controller.CommentListHandler)
	}
	{
		defaultRoutes.GET("/search/article", controller.SearchArticleHandler)
		defaultRoutes.GET("/search/tag", controller.SearchTagHandler)
		defaultRoutes.GET("/search/user", controller.SearchUserHandler)
	}
}
