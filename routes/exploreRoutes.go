package routes

import (
	"manage/controller"

	"github.com/gin-gonic/gin"
)

func exploreRoutesInit(r *gin.Engine) {
	exploreRoutes := r.Group("/explore")
	{
		exploreRoutes.GET("/tags", controller.TagsHandler)
		exploreRoutes.GET("/followingTags", controller.FollowingTagsHandler)
		exploreRoutes.GET("/followingUsers", controller.FollowingUsersHandler)
	}
}
