package users

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestServiceCreateUser(t *testing.T) {

	const userID string = "64260e1da4c0c814bda5734a"

	tests := []struct {
		name             string
		setupMock        func(service *MockUserRepository)
		inputParam       User
		expectedResponse string
		expectedError    error
	}{
		{
			name: "create user success",
			setupMock: func(repository *MockUserRepository) {
				repository.
					EXPECT().
					FindUserByEmail(gomock.Any(), gomock.Any()).
					Return(nil, nil)
				repository.
					EXPECT().
					InsertUser(gomock.Any()).
					Return(userID, nil)
			},
			inputParam: User{
				Address: Address{
					City:    "SP",
					Country: "BR",
					Number:  "111",
					State:   "SP",
					Street:  "Rua hum",
					ZIP:     "12345-678",
				},
				Age:      "33",
				Email:    "test@test.com",
				Name:     "Test",
				Password: "12345",
			},
			expectedResponse: userID,
			expectedError:    nil,
		},
		{
			name: "find user failed",
			setupMock: func(repository *MockUserRepository) {
				repository.
					EXPECT().
					FindUserByEmail(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf(INVALID_OBJECT_ID))
			},
			inputParam: User{
				Address: Address{
					City:    "SP",
					Country: "BR",
					Number:  "111",
					State:   "SP",
					Street:  "Rua hum",
					ZIP:     "12345-678",
				},
				Age:      "33",
				Email:    "test@test.com",
				Name:     "Test",
				Password: "12345",
			},
			expectedResponse: "",
			expectedError:    &userServiceError{code: CREATE_USER_FAILED},
		},
		{
			name: "user already exists",
			setupMock: func(repository *MockUserRepository) {
				repository.
					EXPECT().
					FindUserByEmail(gomock.Any(), gomock.Any()).
					Return(&User{}, nil)
			},
			inputParam: User{
				Address: Address{
					City:    "SP",
					Country: "BR",
					Number:  "111",
					State:   "SP",
					Street:  "Rua hum",
					ZIP:     "12345-678",
				},
				Age:      "33",
				Email:    "test@test.com",
				Name:     "Test",
				Password: "12345",
			},
			expectedResponse: "",
			expectedError:    &userServiceError{code: USER_EXISTS},
		},
		{
			name: "insert user failed",
			setupMock: func(repository *MockUserRepository) {
				repository.
					EXPECT().
					FindUserByEmail(gomock.Any(), gomock.Any()).
					Return(nil, nil)
				repository.
					EXPECT().
					InsertUser(gomock.Any()).
					Return("", errors.New("Any Error"))
			},
			inputParam: User{
				Address: Address{
					City:    "SP",
					Country: "BR",
					Number:  "111",
					State:   "SP",
					Street:  "Rua hum",
					ZIP:     "12345-678",
				},
				Age:      "33",
				Email:    "test@test.com",
				Name:     "Test",
				Password: "12345",
			},
			expectedResponse: "",
			expectedError:    &userServiceError{code: CREATE_USER_FAILED},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tu *testing.T) {
			ctrl := gomock.NewController(tu)
			repo := NewMockUserRepository(ctrl)
			tc.setupMock(repo)

			service := NewUserService(repo)

			result, err := service.CreateUser(tc.inputParam)

			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("Expecting error %d , but returns %d", tc.expectedError, err)
			}

			if result != tc.expectedResponse {
				t.Errorf("Expecting body %s , but returns %s", tc.expectedResponse, result)
			}
		})
	}
}

func TestServiceGetUser(t *testing.T) {

	const userID string = "64260e1da4c0c814bda5734a"

	var user User = User{
		Address: Address{
			City:    "SP",
			Country: "BR",
			Number:  "111",
			State:   "SP",
			Street:  "Rua hum",
			ZIP:     "12345-678",
		},
		Age:      "33",
		Email:    "test@test.com",
		Name:     "Test",
		Password: "12345",
	}

	tests := []struct {
		name             string
		setupMock        func(service *MockUserRepository)
		inputParam       string
		expectedResponse *User
		expectedError    error
	}{
		{
			name: "create user success",
			setupMock: func(repository *MockUserRepository) {
				repository.
					EXPECT().
					FindUserByID(gomock.Any(), gomock.Any()).
					Return(&user, nil)
			},
			inputParam:       userID,
			expectedResponse: &user,
			expectedError:    nil,
		},
		{
			name: "invalid user id",
			setupMock: func(repository *MockUserRepository) {
				repository.
					EXPECT().
					FindUserByID(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf(INVALID_OBJECT_ID))
			},
			inputParam:       "dfsgrgerg",
			expectedResponse: nil,
			expectedError:    &userServiceError{code: USER_ID_INVALID},
		},
		{
			name: "get user fail",
			setupMock: func(repository *MockUserRepository) {
				repository.
					EXPECT().
					FindUserByID(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("Any Error"))
			},
			inputParam:       userID,
			expectedResponse: nil,
			expectedError:    &userServiceError{code: GET_USER_FAILED},
		},
		{
			name: "user not exists",
			setupMock: func(repository *MockUserRepository) {
				repository.
					EXPECT().
					FindUserByID(gomock.Any(), gomock.Any()).
					Return(nil, nil)
			},
			inputParam:       userID,
			expectedResponse: nil,
			expectedError:    &userServiceError{code: USER_NOT_EXISTS},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tu *testing.T) {
			ctrl := gomock.NewController(tu)
			repo := NewMockUserRepository(ctrl)
			tc.setupMock(repo)

			service := NewUserService(repo)

			result, err := service.GetUser(tc.inputParam)

			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("Expecting error %d , but returns %d", tc.expectedError, err)
			}

			if !reflect.DeepEqual(result, tc.expectedResponse) {
				t.Errorf("Expecting body %s , but returns %s", tc.expectedResponse, result)
			}
		})
	}
}

func TestServiceUpdateUser(t *testing.T) {

	const userID string = "64260e1da4c0c814bda5734a"

	var user User = User{
		Address: Address{
			City:    "SP",
			Country: "BR",
			Number:  "111",
			State:   "SP",
			Street:  "Rua hum",
			ZIP:     "12345-678",
		},
		Age:      "33",
		Email:    "test@test.com",
		Name:     "Test",
		Password: "12345",
	}

	type updateParams struct {
		UserID string
		User   User
	}

	tests := []struct {
		name          string
		setupMock     func(service *MockUserRepository)
		inputParam    updateParams
		expectedError error
	}{
		{
			name: "create user success",
			setupMock: func(repository *MockUserRepository) {
				repository.
					EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			inputParam:    updateParams{UserID: userID, User: user},
			expectedError: nil,
		},
		{
			name: "invalid user id",
			setupMock: func(repository *MockUserRepository) {
				repository.
					EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf(INVALID_OBJECT_ID))
			},
			inputParam:    updateParams{UserID: "any id invalid", User: user},
			expectedError: &userServiceError{code: USER_ID_INVALID},
		},
		{
			name: "update user fail",
			setupMock: func(repository *MockUserRepository) {
				repository.
					EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("Any error"))
			},
			inputParam:    updateParams{UserID: userID, User: user},
			expectedError: &userServiceError{code: UPDATE_USER_FAILED},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tu *testing.T) {
			ctrl := gomock.NewController(tu)
			repo := NewMockUserRepository(ctrl)
			tc.setupMock(repo)

			service := NewUserService(repo)

			err := service.UpdateUser(tc.inputParam.UserID, tc.inputParam.User)

			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("Expecting error %d , but returns %d", tc.expectedError, err)
			}
		})
	}
}

func TestServiceDeleteUser(t *testing.T) {

	const userID string = "64260e1da4c0c814bda5734a"

	tests := []struct {
		name          string
		setupMock     func(service *MockUserRepository)
		inputParam    string
		expectedError error
	}{
		{
			name: "delete user success",
			setupMock: func(repository *MockUserRepository) {
				repository.
					EXPECT().
					DeleteUser(gomock.Any()).
					Return(nil)
			},
			inputParam:    userID,
			expectedError: nil,
		},
		{
			name: "invalid user id",
			setupMock: func(repository *MockUserRepository) {
				repository.
					EXPECT().
					DeleteUser(gomock.Any()).
					Return(fmt.Errorf(INVALID_OBJECT_ID))
			},
			inputParam:    "any id invalid",
			expectedError: &userServiceError{code: USER_ID_INVALID},
		},
		{
			name: "delete user fail",
			setupMock: func(repository *MockUserRepository) {
				repository.
					EXPECT().
					DeleteUser(gomock.Any()).
					Return(fmt.Errorf("Any error"))
			},
			inputParam:    userID,
			expectedError: &userServiceError{code: DELETE_USER_FAILED},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tu *testing.T) {
			ctrl := gomock.NewController(tu)
			repo := NewMockUserRepository(ctrl)
			tc.setupMock(repo)

			service := NewUserService(repo)

			err := service.DeleteUser(tc.inputParam)

			if err != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("Expecting error %d , but returns %d", tc.expectedError, err)
			}
		})
	}
}
