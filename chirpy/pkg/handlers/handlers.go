package handlers

import (
	"gitlab.com/daneofmanythings/chirpy/pkg/config"
)

var Repo *Repository

type Repository struct {
	App *config.Config
}

func NewRepo(c *config.Config) *Repository {
	return &Repository{
		App: c,
	}
}

func LinkRepository(r *Repository) {
	Repo = r
}
