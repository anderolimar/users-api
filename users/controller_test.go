package users

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestCreateUser(t *testing.T) {

	const userID string = "64260e1da4c0c814bda5734a"

	tests := []struct {
		name             string
		setupMock        func(service *MockUserService)
		inputBody        string
		expectedResponse string
		expectedStatus   int
	}{
		{
			name: "create user success",
			setupMock: func(service *MockUserService) {
				service.
					EXPECT().
					CreateUser(gomock.Any()).
					Return(userID, nil)
			},
			inputBody: `{
				"address": {
				  "city": "SP",
				  "country": "BR",
				  "number": "111",
				  "state": "SP",
				  "street": "Rua hum",
				  "zip": "12345-678"
				},
				"age": "33",
				"email": "test@test.com",
				"name": "Test",
				"password": "12345"
			  }`,
			expectedStatus:   http.StatusCreated,
			expectedResponse: fmt.Sprintf(`{"ID":"%s"}`, userID),
		},
		{
			name:             "invalid user data",
			setupMock:        func(service *MockUserService) {},
			inputBody:        ``,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid User Data","code":"INVALID_USER_DATA"}`,
		},
		{
			name:             "email required",
			setupMock:        func(service *MockUserService) {},
			inputBody:        `{}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"Email Required","code":"EMAIL_REQUIRED"}`,
		},
		{
			name: "user already exists",
			setupMock: func(service *MockUserService) {
				service.
					EXPECT().
					CreateUser(gomock.Any()).
					Return("", &userServiceError{code: USER_EXISTS})
			},
			inputBody:        `{"email": "test@test.com"}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"User Already Exists","code":"USER_ALREADY_EXISTS"}`,
		},
		{
			name: "user create failed",
			setupMock: func(service *MockUserService) {
				service.
					EXPECT().
					CreateUser(gomock.Any()).
					Return("", &userServiceError{code: CREATE_USER_FAILED})
			},
			inputBody:        `{"email": "test@test.com"}`,
			expectedStatus:   http.StatusBadGateway,
			expectedResponse: `{"message":"User Create Failed","code":"USER_CREATE_FAILED"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tu *testing.T) {
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(tu)
			svc := NewMockUserService(ctrl)
			tc.setupMock(svc)

			controller := NewUserController(svc)
			r := gin.Default()
			r.POST("/api/v1/users", controller.CreateUser)

			req, err := http.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(tc.inputBody))
			if err != nil {
				t.Errorf("Error in request : %v", err)
			}
			r.ServeHTTP(w, req)

			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("Expecting statusCode %d , but returns %d", tc.expectedStatus, w.Result().StatusCode)
			}

			if r := w.Body.String(); r != tc.expectedResponse {
				t.Errorf("Expecting body %s , but returns %s", tc.expectedResponse, r)
			}
		})
	}
}

func TestGetUser(t *testing.T) {

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
		setupMock        func(service *MockUserService)
		inputParam       string
		expectedResponse string
		expectedStatus   int
	}{
		{
			name: "get user success",
			setupMock: func(service *MockUserService) {
				service.
					EXPECT().
					GetUser(gomock.Any()).
					Return(&user, nil)
			},
			inputParam:       userID,
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"id":"","name":"Test","age":"33","email":"test@test.com","password":"12345","address":{"street":"Rua hum","number":"111","zip":"12345-678","city":"SP","state":"SP","country":"BR"}}`,
		},
		{
			name: "invalid user id",
			setupMock: func(service *MockUserService) {
				service.
					EXPECT().
					GetUser(gomock.Any()).
					Return(nil, &userServiceError{code: USER_ID_INVALID})
			},
			inputParam:       `gdfhdhgh`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid User ID","code":"INVALID_USER_ID"}`,
		},
		{
			name: "user not found",
			setupMock: func(service *MockUserService) {
				service.
					EXPECT().
					GetUser(gomock.Any()).
					Return(nil, &userServiceError{code: USER_NOT_EXISTS})
			},
			inputParam:       userID,
			expectedStatus:   http.StatusNotFound,
			expectedResponse: `{"message":"User Not Found","code":"USER_NOT_FOUND"}`,
		},
		{
			name: "user find failed",
			setupMock: func(service *MockUserService) {
				service.
					EXPECT().
					GetUser(gomock.Any()).
					Return(nil, &userServiceError{code: GET_USER_FAILED})
			},
			inputParam:       userID,
			expectedStatus:   http.StatusBadGateway,
			expectedResponse: `{"message":"User Find Failed","code":"USER_FIND_FAILED"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tu *testing.T) {
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(tu)
			svc := NewMockUserService(ctrl)
			tc.setupMock(svc)

			controller := NewUserController(svc)

			r := gin.Default()
			r.GET("/api/v1/users/:id", controller.GetUser)

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s", tc.inputParam), nil)
			if err != nil {
				t.Errorf("Error in request : %v", err)
			}

			r.ServeHTTP(w, req)

			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("Expecting statusCode %d , but returns %d", tc.expectedStatus, w.Result().StatusCode)
			}

			if r := w.Body.String(); r != tc.expectedResponse {
				t.Errorf("Expecting body %s , but returns %s", tc.expectedResponse, r)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {

	const userID string = "64260e1da4c0c814bda5734a"

	tests := []struct {
		name             string
		setupMock        func(service *MockUserService)
		inputBody        string
		inputParam       string
		expectedResponse string
		expectedStatus   int
	}{
		{
			name: "update user success",
			setupMock: func(service *MockUserService) {
				service.
					EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			inputBody: `{
				"age": "33",
				"email": "test@test.com",
				"name": "Test"
			  }`,
			inputParam:       userID,
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"User Updated","code":"USER_UPDATED"}`,
		},
		{
			name:             "invalid user data",
			setupMock:        func(service *MockUserService) {},
			inputBody:        ``,
			inputParam:       userID,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid User Data","code":"INVALID_USER_DATA"}`,
		},
		{
			name: "invalid user id",
			setupMock: func(service *MockUserService) {
				service.
					EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(&userServiceError{code: USER_ID_INVALID})
			},
			inputBody: `{
				"age": "33",
				"email": "test@test.com",
				"name": "Test"
			  }`,
			inputParam:       `gdfhdhgh`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid User ID","code":"INVALID_USER_ID"}`,
		},
		{
			name: "user update failed",
			setupMock: func(service *MockUserService) {
				service.
					EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Return(&userServiceError{code: UPDATE_USER_FAILED})
			},
			inputBody: `{
				"age": "33",
				"email": "test@test.com",
				"name": "Test"
			  }`,
			inputParam:       userID,
			expectedStatus:   http.StatusBadGateway,
			expectedResponse: `{"message":"User Update Failed","code":"USER_UPDATE_FAILED"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tu *testing.T) {
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(tu)
			svc := NewMockUserService(ctrl)
			tc.setupMock(svc)

			controller := NewUserController(svc)
			r := gin.Default()
			r.PUT("/api/v1/users/:id", controller.UpdateUser)

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", tc.inputParam), strings.NewReader(tc.inputBody))
			if err != nil {
				t.Errorf("Error in request : %v", err)
			}
			r.ServeHTTP(w, req)

			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("Expecting statusCode %d , but returns %d", tc.expectedStatus, w.Result().StatusCode)
			}

			if r := w.Body.String(); r != tc.expectedResponse {
				t.Errorf("Expecting body %s , but returns %s", tc.expectedResponse, r)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {

	const userID string = "64260e1da4c0c814bda5734a"

	tests := []struct {
		name             string
		setupMock        func(service *MockUserService)
		inputParam       string
		expectedResponse string
		expectedStatus   int
	}{
		{
			name: "update user success",
			setupMock: func(service *MockUserService) {
				service.
					EXPECT().
					DeleteUser(gomock.Any()).
					Return(nil)
			},
			inputParam:       userID,
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"User Deleted","code":"USER_DELETED"}`,
		},
		{
			name: "invalid user id",
			setupMock: func(service *MockUserService) {
				service.
					EXPECT().
					DeleteUser(gomock.Any()).
					Return(&userServiceError{code: USER_ID_INVALID})
			},
			inputParam:       `gdfhdhgh`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid User ID","code":"INVALID_USER_ID"}`,
		},
		{
			name: "user delete failed",
			setupMock: func(service *MockUserService) {
				service.
					EXPECT().
					DeleteUser(gomock.Any()).
					Return(&userServiceError{code: DELETE_USER_FAILED})
			},
			inputParam:       userID,
			expectedStatus:   http.StatusBadGateway,
			expectedResponse: `{"message":"User Delete Failed","code":"USER_DELETE_FAILED"}`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(tu *testing.T) {
			w := httptest.NewRecorder()
			ctrl := gomock.NewController(tu)
			svc := NewMockUserService(ctrl)
			tc.setupMock(svc)

			controller := NewUserController(svc)
			r := gin.Default()
			r.DELETE("/api/v1/users/:id", controller.DeleteUser)

			req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/users/%s", tc.inputParam), nil)
			if err != nil {
				t.Errorf("Error in request : %v", err)
			}
			r.ServeHTTP(w, req)

			if w.Result().StatusCode != tc.expectedStatus {
				t.Errorf("Expecting statusCode %d , but returns %d", tc.expectedStatus, w.Result().StatusCode)
			}

			if r := w.Body.String(); r != tc.expectedResponse {
				t.Errorf("Expecting body %s , but returns %s", tc.expectedResponse, r)
			}
		})
	}
}
