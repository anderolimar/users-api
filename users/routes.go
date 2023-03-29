package users

import (
	"userapi/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddRoutes(api *gin.RouterGroup, config config.Config, client *mongo.Client) {
	var userService UserService = NewUserService()
	var userController UserController = NewUserController(userService)

	api.GET("/users/:id", userController.GetUser)
	api.POST("/users", userController.CreateUser)
	api.PUT("/users", userController.UpdateUser)
	api.DELETE("/users", userController.DeleteUser)
}
