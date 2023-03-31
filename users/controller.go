package users

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// Controller containing all User request handlers
type UserController struct {
	service UserService
}

// Returns new UserController instance
func NewUserController(service UserService) UserController {
	return UserController{
		service: service,
	}
}

// CreateUser godoc
//
//	@Summary		Create new user
//	@Description	This endpoint creates a new user from user data in request body.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		User	true	"body"
//	@Success		201		{object}	UserID
//	@Failure		401
//	@Failure		400		{object}	UserResponse
//	@Failure		400		{object}	UserResponse
//	@Failure		502		{object}	UserResponse
//	@Router			/users [post]
func (ctr UserController) CreateUser(c *gin.Context) {
	var user User

	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.JSON(400, INVALID_USER_DATA)
		return
	}

	var userID string
	userID, err = ctr.service.CreateUser(user)

	if err != nil {
		if err.Error() == USER_EXISTS {
			c.JSON(400, USER_ALREADY_EXISTS)
			return
		}
		c.JSON(502, USER_CREATE_FAILED)
		return
	}

	c.JSON(201, UserID{ID: userID})
}

// GetUser godoc
//
//	@Summary		Return user data
//	@Description	This endpoint returns a user by user id.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"userID"
//	@Success		201		{object}	UserID
//	@Failure		401
//	@Failure		400		{object}	UserResponse
//	@Failure		400		{object}	UserResponse
//	@Failure		502		{object}	UserResponse
//	@Router			/users/{id} [get]
func (ctr UserController) GetUser(c *gin.Context) {
	var userID string = c.Param("id")

	user, err := ctr.service.GetUser(userID)
	if err != nil {
		if err.Error() == USER_ID_INVALID {
			c.JSON(400, INVALID_USER_ID)
			return
		}

		if err.Error() == USER_NOT_EXISTS {
			c.JSON(404, USER_NOT_FOUND)
			return
		}
		c.JSON(502, USER_FIND_FAILED)
		return
	}

	c.JSON(200, user)
}

// UpdateUser godoc
//
//	@Summary		Update user
//	@Description	This endpoint updates a user from user data in request body.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"userID"
//	@Param			request	body		User	true	"body"
//	@Success		200	{object}	UserResponse
//	@Failure		401
//	@Failure		400	{object}	UserResponse
//	@Failure		400	{object}	UserResponse
//	@Failure		502	{object}	UserResponse
//	@Router			/users/{id} [put]
func (ctr UserController) UpdateUser(c *gin.Context) {
	var user User
	var userID string = c.Param("id")

	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		if err.Error() == USER_ID_INVALID {
			c.JSON(400, INVALID_USER_ID)
			return
		}

		c.JSON(400, INVALID_USER_DATA)
		return
	}

	if err := ctr.service.UpdateUser(userID, user); err != nil {
		c.JSON(502, USER_UPDATE_FAILED)
		return
	}

	c.JSON(200, USER_UPDATED)
}

// DeleteUser godoc
//
//	@Summary		Delete user
//	@Description	This endpoint delete a user by user id.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"userID"
//	@Success		200	{object}	UserResponse
//	@Failure		401
//	@Failure		400	{object}	UserResponse
//	@Failure		502	{object}	UserResponse
//	@Router			/users/{id} [delete]
func (ctr UserController) DeleteUser(c *gin.Context) {
	var userID string = c.Param("id")

	if userID == "" {
		c.JSON(400, INVALID_USER_ID)
		return
	}

	err := ctr.service.DeleteUser(userID)
	if err != nil {
		c.JSON(502, USER_DELETE_FAILED)
		return
	}

	c.JSON(200, USER_DELETED)
}
