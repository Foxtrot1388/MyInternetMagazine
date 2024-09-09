package model

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ElementOfList struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
