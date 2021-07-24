package apiserver

type getUserRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type deleteUserRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type createUserRequest struct {
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
}

type updateUserRequestBody struct {
	Username  string `json:"username" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
}

type updateUserRequestUri struct {
	Id int64 `uri:"id" binding:"required"`
}
