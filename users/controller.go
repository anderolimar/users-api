package users

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service UserService
}

func NewUserController(service UserService) UserController {
	return UserController{
		service: service,
	}
}

func (ctr UserController) CreateUser(c *gin.Context) {
	var user User

	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.JSON(400, INVALID_USER)
		return
	}

	var createdUser User
	createdUser, err = ctr.service.CreateUser(user)
	if err != nil {
		c.JSON(502, USER_CREATE_FAILED)
		return
	}

	c.JSON(201, createdUser)
}

func (ctr UserController) GetUser(c *gin.Context) {
	var userIDStr string = c.Param("id")

	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		c.JSON(400, INVALID_USER_ID)
		return
	}

	user, err := ctr.service.GetUser(userID)
	if err != nil {
		c.JSON(404, USER_NOT_FOUND)
		return
	}

	c.JSON(200, user)
}

func (ctr UserController) UpdateUser(c *gin.Context) {
	var user User

	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.JSON(400, INVALID_USER)
		return
	}

	if err := ctr.service.UpdateUser(user); err != nil {
		c.JSON(502, USER_UPDATE_FAILED)
		return
	}

	c.JSON(200, USER_UPDATED)
}

func (ctr UserController) DeleteUser(c *gin.Context) {
	var userIDStr string = c.Param("id")

	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		c.JSON(400, INVALID_USER_ID)
		return
	}

	err = ctr.service.DeleteUser(userID)
	if err != nil {
		c.JSON(502, USER_DELETE_FAILED)
		return
	}

	c.JSON(200, USER_DELETED)
}
