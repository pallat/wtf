package ticket

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct{}

func NewStorage(*sql.DB) *Storage {
	return &Storage{}
}
