package users

type User struct {
	Name     string  `json:"name"`
	Age      string  `json:"age"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Address  Address `json:"address"`
}

type Address struct {
	Street  string `json:"street"`
	Number  string `json:"number"`
	ZIP     string `json:"zip"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}
