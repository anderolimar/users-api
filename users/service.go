package users

type UserService interface {
	CreateUser(user User) (User, error)
	GetUser(userID int) (User, error)
	UpdateUser(user User) error
	DeleteUser(userID int) error
}

type userService struct {
}

func NewUserService() UserService {
	return &userService{}
}

func (service *userService) CreateUser(user User) (User, error) {
	return User{}, nil
}

func (service *userService) GetUser(userID int) (User, error) {
	return User{}, nil
}

func (service *userService) UpdateUser(user User) error {
	return nil
}

func (service *userService) DeleteUser(userID int) error {
	return nil
}
