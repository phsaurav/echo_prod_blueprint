package server

import (
	"github.com/phsaurav/echo_prod_blueprint/internal/database"
)

type Store struct {
	db database.Service
}

func NewStore(db database.Service) Store {
	return Store{
		db: db,
	}
}

func (r *Store) DBHealth() map[string]string {
	return r.db.Health()
}
