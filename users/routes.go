// Users module containing Controllers, Services e Repositories.
// Module responsable for users bussiness
package users

import (
	"userapi/config"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Method to add routes in api (gin.RouterGroup), using config (config.Config)
// and client (mongo.Client)
func AddRoutes(api *gin.RouterGroup, config config.Config, client *mongo.Client) {
	var userRepository UserRepository = NewUserRepository(client, config.Database)
	var userService UserService = NewUserService(userRepository)
	var userController UserController = NewUserController(userService)

	api.GET("/users/:id", userController.GetUser)
	api.POST("/users", userController.CreateUser)
	api.PUT("/users/:id", userController.UpdateUser)
	api.DELETE("/users/:id", userController.DeleteUser)
}
