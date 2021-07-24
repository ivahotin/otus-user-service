package apiserver

type getUserByIdResponse struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
