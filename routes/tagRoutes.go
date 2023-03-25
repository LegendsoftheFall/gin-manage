package routes

import (
	"manage/controller/user"

	"github.com/gin-gonic/gin"
)

func tagRoutesInit(r *gin.Engine) {
	tagRoutes := r.Group("/n")
	{
		tagRoutes.GET("/:id", user.TagDetailHandler)
	}
}
