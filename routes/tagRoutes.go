package routes

import (
	"manage/controller"

	"github.com/gin-gonic/gin"
)

func tagRoutesInit(r *gin.Engine) {
	tagRoutes := r.Group("/n")
	{
		tagRoutes.GET("/:id", controller.TagDetailHandler)
	}
}
