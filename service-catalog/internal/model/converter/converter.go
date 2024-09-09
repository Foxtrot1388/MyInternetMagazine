package converter

import (
	"v1/internal/entity"
	"v1/internal/model"
)

func GetProduct(e *entity.Product) *model.Product {

	return &model.Product{
		Id:          e.Id,
		Name:        e.Name,
		Description: e.Description,
	}

}

func GetElementOfList(e entity.ElementOfList) model.ElementOfList {

	return model.ElementOfList{
		Id:   e.Id,
		Name: e.Name,
	}

}
