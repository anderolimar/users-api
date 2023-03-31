package users

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	/*
		Method to create user

		Parameters

		user: The user data to create new user.
	*/
	CreateUser(user User) (string, error)
	/*
		Method to get user

		Parameters

		userID: User ID to find user data.
	*/
	GetUser(userID string) (*User, error)
	/*
		Method to update user

		Parameters

		userID: User ID to find user data.
		user: User data to update user.
	*/
	UpdateUser(userID string, user User) error
	DeleteUser(userID string) error
}

type userServiceError struct {
	code string
}

func (e *userServiceError) Error() string {
	return e.code
}

const CREATE_USER_FAILED string = "CREATE_USER_FAILED"
const GET_USER_FAILED string = "GET_USER_FAILED"
const USER_ID_INVALID string = "USER_ID_INVALID"
const USER_EXISTS string = "USER_EXISTS"
const USER_NOT_EXISTS string = "USER_EXISTS"
const UPDATE_USER_FAILED string = "UPDATE_USER_FAILED"
const DELETE_USER_FAILED string = "DELETE_USER_FAILED"

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) CreateUser(user User) (string, error) {
	projection := Projection{{Key: "_id", Value: 1}}
	existingUser, err := svc.repo.FindUserByEmail(user.Email, projection)
	if err != nil {
		fmt.Println(fmt.Errorf("Error on FindUserByEmail : %v", err))
		return "", &userServiceError{code: CREATE_USER_FAILED}
	}

	if existingUser != nil {
		fmt.Println(fmt.Errorf("User already exists"))
		return "", &userServiceError{code: USER_EXISTS}
	}

	user.Password = svc.hashPassword(user.Password)

	insertID, err := svc.repo.InsertUser(user)
	return insertID, nil
}

func (svc *userService) GetUser(userID string) (*User, error) {
	projection := Projection{{Key: "password", Value: 0}}
	user, err := svc.repo.FindUserByID(userID, projection)
	if err != nil {
		if err.Error() == INVALID_OBJECT_ID {
			fmt.Println(fmt.Errorf("Invalid user id : %v", err))
			return nil, &userServiceError{code: USER_ID_INVALID}
		}
		fmt.Println(fmt.Errorf("Error on FindUserByID : %v", err))
		return nil, &userServiceError{code: GET_USER_FAILED}
	}

	if user == nil {
		fmt.Println(fmt.Errorf("User not exists"))
		return nil, &userServiceError{code: USER_NOT_EXISTS}
	}

	return user, nil
}

func (svc *userService) UpdateUser(userID string, user User) error {
	if err := svc.repo.UpdateUser(userID, user); err != nil {
		if err.Error() == INVALID_OBJECT_ID {
			fmt.Println(fmt.Errorf("Invalid user id : %v", err))
			return &userServiceError{code: USER_ID_INVALID}
		}
		fmt.Println(fmt.Errorf("Error on UpdateUser : %v", err))
		return &userServiceError{code: UPDATE_USER_FAILED}
	}
	return nil
}

func (svc *userService) DeleteUser(userID string) error {
	if err := svc.repo.DeleteUser(userID); err != nil {
		fmt.Println(fmt.Errorf("Error on DeleteUser : %v", err))
		return &userServiceError{code: DELETE_USER_FAILED}
	}
	return nil
}

func (svc *userService) hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes)
}
