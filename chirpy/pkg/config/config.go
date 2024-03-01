package config

import (
	"gitlab.com/daneofmanythings/chirpy/internal/database"
)

type Config struct {
	fileserverHits int
	DB             *database.DB
	JWTSECRET      string
	APIKEYPOLKA    string
}

func (c *Config) ResetHits() {
	c.fileserverHits = 0
}

func (c *Config) HitRegistered() {
	c.fileserverHits++
}

func (c *Config) GetHits() int {
	return c.fileserverHits
}
