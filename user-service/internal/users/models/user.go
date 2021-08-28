package models

type UserId int64

type User struct {
	Id        UserId
	Username  string `validate:"required,max=256"`
	Firstname string `validate:"required,max=50"`
	Lastname  string `validate:"required,max=50"`
	Email     string `validate:"required,max=50,email"`
	Phone     string `validate:"required,len=12"`
	OwnerId   string
}
