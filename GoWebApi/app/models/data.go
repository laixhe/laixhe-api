package models

type Data struct {
	User *User
}

func NewData() *Data {
	return &Data{
		User: &User{},
	}
}
