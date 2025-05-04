package model

import "github.com/tomato3713/storyline/server/graph/model"

type User struct {
	ID   string
	Name string
}

func (u *User) ToGQLModel() *model.User {
	return &model.User{
		ID:   u.ID,
		Name: u.Name,
	}
}
