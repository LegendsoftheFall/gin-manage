package routes

import (
	"github.com/gin-gonic/gin"
	"manage/controller"
)

func recommendRoutesInit(r *gin.Engine) {
	recommendRoutes := r.Group("/recommend")
	{
		recommendRoutes.GET("/list", controller.ArticleListHandler)
	}
}
