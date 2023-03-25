package routes

import (
	"manage/controller/user"

	"github.com/gin-gonic/gin"
)

func exploreRoutesInit(r *gin.Engine) {
	exploreRoutes := r.Group("/explore")
	{
		exploreRoutes.GET("/hotTags", user.HotTagsHandler)
		exploreRoutes.GET("/tags", user.TagsHandler)
		exploreRoutes.GET("/followingTags", user.FollowingTagsHandler)
		exploreRoutes.GET("/followingUsers", user.FollowingUsersHandler)
	}
}
