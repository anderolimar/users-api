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
		input            string
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
			input: `{
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
			input:            ``,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"message":"Invalid User Data","code":"INVALID_USER_DATA"}`,
		},
		{
			name: "user already exists",
			setupMock: func(service *MockUserService) {
				service.
					EXPECT().
					CreateUser(gomock.Any()).
					Return("", &userServiceError{code: USER_EXISTS})
			},
			input:            `{"email": "test@test.com"}`,
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
			input:            `{"email": "test@test.com"}`,
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

			req, err := http.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(tc.input))
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
		input            string
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
			input:            userID,
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
			input:            `gdfhdhgh`,
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
			input:            userID,
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
			input:            userID,
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

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%s", tc.input), nil)
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
		input            string
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
			input: `{
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
			  }`,
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message":"User Updated","code":"USER_UPDATED"}`,
		},
		// {
		// 	name:             "invalid user data",
		// 	setupMock:        func(service *MockUserService) {},
		// 	input:            ``,
		// 	expectedStatus:   http.StatusBadRequest,
		// 	expectedResponse: `{"message":"Invalid User Data","code":"USER_UPDATED"}`,
		// },
		// {
		// 	name: "user already exists",
		// 	setupMock: func(service *MockUserService) {
		// 		service.
		// 			EXPECT().
		// 			CreateUser(gomock.Any()).
		// 			Return("", &userServiceError{code: USER_EXISTS})
		// 	},
		// 	input:            `{"email": "test@test.com"}`,
		// 	expectedStatus:   http.StatusBadRequest,
		// 	expectedResponse: `{"message":"User Already Exists","code":"USER_ALREADY_EXISTS"}`,
		// },
		// {
		// 	name: "user create failed",
		// 	setupMock: func(service *MockUserService) {
		// 		service.
		// 			EXPECT().
		// 			CreateUser(gomock.Any()).
		// 			Return("", &userServiceError{code: CREATE_USER_FAILED})
		// 	},
		// 	input:            `{"email": "test@test.com"}`,
		// 	expectedStatus:   http.StatusBadGateway,
		// 	expectedResponse: `{"message":"User Create Failed","code":"USER_CREATE_FAILED"}`,
		// },
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

			req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/users/%s", userID), strings.NewReader(tc.input))
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
