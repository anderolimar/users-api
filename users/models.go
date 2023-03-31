package users

type UserID struct {
	ID string
}

type User struct {
	ID       string  `json:"id" bson:"_id,omitempty"`
	Name     string  `json:"name"`
	Age      string  `json:"age"`
	Email    string  `json:"email"`
	Password string  `json:"password,omitempty"`
	Address  Address `json:"address"`
}

func (u *User) projection() Projection {
	m := make(Projection, 0)
	for k := range UserAccess {
		v := UserAccess[k](u)
		if v != "" {
			m = append(m, ProjectionsFields{Key: k, Value: v})
		}
	}
	return m
}

var UserAccess = map[string]UserGetter{
	"name":            func(v *User) string { return v.Name },
	"age":             func(v *User) string { return v.Age },
	"email":           func(v *User) string { return v.Email },
	"password":        func(v *User) string { return v.Password },
	"address.street":  func(v *User) string { return v.Address.Street },
	"address.number":  func(v *User) string { return v.Address.Number },
	"address.zip":     func(v *User) string { return v.Address.ZIP },
	"address.city":    func(v *User) string { return v.Address.City },
	"address.state":   func(v *User) string { return v.Address.State },
	"address.country": func(v *User) string { return v.Address.Country },
}

type UserGetter func(v *User) string

type Address struct {
	Street  string `json:"street"`
	Number  string `json:"number"`
	ZIP     string `json:"zip"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}
