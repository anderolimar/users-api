package users

type UserDB struct {
	Name     string    `bson:"name"`
	Age      string    `bson:"age"`
	Email    string    `bson:"email"`
	Password string    `bson:"password"`
	Address  AddressDB `bson:"address"`
}

type AddressDB struct {
	Street  string `bson:"street"`
	Number  string `bson:"number"`
	ZIP     string `bson:"zip"`
	City    string `bson:"city"`
	State   string `bson:"state"`
	Country string `bson:"country"`
}

type UserRepository interface {
}

type userRepository struct {
}

func (repo *userRepository) InsertUser(user UserDB) {

}
