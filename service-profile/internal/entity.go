package entity

type User struct {
	Id         int
	Login      string
	Firstname  string
	Secondname string
	Lastname   string
	Email      string
}

type NewUser struct {
	Id         int
	Pass       string
	Login      string
	Firstname  string
	Secondname string
	Lastname   string
	Email      string
}

type LoginUser struct {
	Id         int
	Firstname  string
	Secondname string
	Lastname   string
	Email      string
}
