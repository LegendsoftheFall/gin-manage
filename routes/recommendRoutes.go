package routes

import (
	"manage/controller/user"

	"github.com/gin-gonic/gin"
)

func recommendRoutesInit(r *gin.Engine) {
	recommendRoutes := r.Group("/recommend")
	{
		recommendRoutes.GET("/list", user.ArticleListHandler)
	}
}
