package routes

import (
	"manage/controller/admin"

	"github.com/gin-gonic/gin"
)

func adminRoutesInit(r *gin.Engine) {
	adminRoutes := r.Group("/admin")
	{
		adminRoutes.POST("/login", admin.LoginForAdminHandler)
		adminRoutes.POST("/signup", admin.SignUpForAdminHandler)
	}
}
