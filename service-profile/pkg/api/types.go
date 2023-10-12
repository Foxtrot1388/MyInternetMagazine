package api

type GetUserRequest struct {
	Id         int
	Login      string
	Firstname  string
	Secondname string
	Lastname   string
	Email      string
}

func (GetUserRequest) TableName() string {
	return "users"
}

type CreateUserRequest struct {
	Id         int
	Pass       string
	Login      string
	Firstname  string
	Secondname string
	Lastname   string
	Email      string
}

func (CreateUserRequest) TableName() string {
	return "users"
}

type DeleteUserRequest struct {
	Id int
}

func (DeleteUserRequest) TableName() string {
	return "users"
}

type LoginUserRequest struct {
	Id         int
	Firstname  string
	Secondname string
	Lastname   string
	Email      string
}

func (LoginUserRequest) TableName() string {
	return "users"
}
