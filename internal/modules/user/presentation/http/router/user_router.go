package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ndxbinh1922001/VNalo-be/internal/modules/user/presentation/http/handler"
)

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(router *gin.RouterGroup, userHandler *handler.UserHandler) {
	users := router.Group("/users")
	{
		users.POST("", userHandler.CreateUser)
		users.GET("", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)

		// Additional user actions
		users.POST("/:id/promote-vip", userHandler.PromoteToVIP)
		users.POST("/:id/demote-vip", userHandler.DemoteFromVIP)
		users.POST("/:id/change-password", userHandler.ChangePassword)
		users.POST("/:id/activate", userHandler.ActivateUser)
		users.POST("/:id/deactivate", userHandler.DeactivateUser)
	}
}

