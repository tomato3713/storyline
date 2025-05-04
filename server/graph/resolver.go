package graph

import "github.com/tomato3713/storyline/server/repository"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repository repository.Repository
}
