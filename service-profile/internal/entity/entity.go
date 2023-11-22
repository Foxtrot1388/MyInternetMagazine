package entity

type User struct {
	Id         int
	Login      string
	Firstname  string
	Secondname string
	Lastname   string
	Email      string
}

type UserDB struct {
	Id         int
	Pass       []byte
	Login      string
	Firstname  string
	Secondname string
	Lastname   string
	Email      string
}

type NewUser struct {
	Id         int
	Pass       []byte
	Login      string
	Firstname  string
	Secondname string
	Lastname   string
	Email      string
}

type LoginUser struct {
	Token string
}
