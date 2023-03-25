package routes

import (
	"manage/controller/user"

	"github.com/gin-gonic/gin"
)

func defaultRoutesInit(r *gin.Engine) {
	defaultRoutes := r.Group("/")
	{
		defaultRoutes.POST("/signup", user.SignUpHandler)
		defaultRoutes.POST("/login", user.LoginHandler)
		defaultRoutes.POST("/refresh", user.RefreshHandler)
	}
	{
		defaultRoutes.GET("/trendingTags", user.TrendingTagHandler)
		defaultRoutes.GET("/tag", user.TagInfoHandler)
	}
	{
		defaultRoutes.GET("/user/:id", user.HomeHandler)
		defaultRoutes.GET("/user/articles", user.ArticleOfUserHandler)
		defaultRoutes.GET("/user/profile", user.ProfileHandler)
		defaultRoutes.GET("/user/following", user.FollowingUsersHandler)
		defaultRoutes.GET("/user/followers", user.FollowerUsersHandler)
	}
	{
		defaultRoutes.GET("/article/:id", user.ArticleHandler)
	}
	{
		defaultRoutes.GET("/comment/list", user.CommentListHandler)
	}
	{
		defaultRoutes.GET("/search/article", user.SearchArticleHandler)
		defaultRoutes.GET("/search/tag", user.SearchTagHandler)
		defaultRoutes.GET("/search/user", user.SearchUserHandler)
	}
}
