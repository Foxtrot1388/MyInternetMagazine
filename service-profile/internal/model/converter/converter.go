package converter

import (
	"v1/internal/entity"
	"v1/internal/model"
)

func GetUser(e *entity.User) *model.User {

	return &model.User{
		Id:         e.Id,
		Login:      e.Login,
		Firstname:  e.Firstname,
		Secondname: e.Secondname,
		Lastname:   e.Lastname,
		Email:      e.Email,
	}

}
