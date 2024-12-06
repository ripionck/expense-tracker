package routes

import (
	"expense-tracker/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func WishlistRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/wishlist", controllers.GetWishlist)
		v1.POST("/wishlist/add", controllers.CreateWishlist)
		v1.PUT("/wishlist/:id", controllers.UpdateWishlist)
		v1.DELETE("/wishlist/:id", controllers.DeleteWishlist)
	}
}
