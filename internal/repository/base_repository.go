package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type BaseRepository struct {
	db *sql.DB
}

func (r *BaseRepository) Close() error {
	return r.db.Close()
}
