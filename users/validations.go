package users

type ValidationResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (r ValidationResponse) Error() string {
	return r.Code
}

var EMAIL_REQUIRED *ValidationResponse = &ValidationResponse{
	Code:    "EMAIL_REQUIRED",
	Message: "Email Required",
}

func ValidateUser(user User) *ValidationResponse {
	if user.Email == "" {
		return EMAIL_REQUIRED
	}
	return nil
}
