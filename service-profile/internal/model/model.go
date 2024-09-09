package model

type User struct {
	Id         int    `json:"id"`
	Login      string `json:"login"`
	Firstname  string `json:"firstname"`
	Secondname string `json:"secondname"`
	Lastname   string `json:"lastname"`
	Email      string `json:"email"`
}

type UserDB struct {
	Id         int    `json:"id"`
	Pass       []byte `json:"pass"`
	Login      string `json:"login"`
	Firstname  string `json:"firstname"`
	Secondname string `json:"secondname"`
	Lastname   string `json:"lastname"`
	Email      string `json:"email"`
}

type NewUser struct {
	Id         int    `json:"id"`
	Pass       []byte `json:"pass"`
	Login      string `json:"login"`
	Firstname  string `json:"firstname"`
	Secondname string `json:"secondname"`
	Lastname   string `json:"lastname"`
	Email      string `json:"email"`
}

type LoginUser struct {
	Token string `json:"token"`
}
