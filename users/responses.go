package users

type UserResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type UserListResponse struct {
	Users []User `json:"users"`
}

var INVALID_USER UserResponse = UserResponse{
	Message: "Invalid User",
	Code:    "INVALID_USER",
}

var INVALID_USER_ID UserResponse = UserResponse{
	Message: "Invalid User ID",
	Code:    "INVALID_USER_ID",
}

var USER_NOT_FOUND UserResponse = UserResponse{
	Message: "User Not Found",
	Code:    "USER_NOT_FOUND",
}

var USER_CREATED UserResponse = UserResponse{
	Message: "User Created",
	Code:    "USER_CREATED",
}

var USER_CREATE_FAILED UserResponse = UserResponse{
	Message: "User Create Failed",
	Code:    "USER_CREATE_FAILED",
}

var USER_UPDATED UserResponse = UserResponse{
	Message: "User Updated",
	Code:    "USER_UPDATED",
}

var USER_UPDATE_FAILED UserResponse = UserResponse{
	Message: "User Failed",
	Code:    "USER_UPDATE_FAILED",
}

var USER_DELETED UserResponse = UserResponse{
	Message: "User Deleted",
	Code:    "USER_DELETED",
}

var USER_DELETE_FAILED UserResponse = UserResponse{
	Message: "User Delete Failed",
	Code:    "USER_DELETE_FAILED",
}
