package database

import (
	"database/sql"
	"log"
	"time"
)

const (
	connMaxLifetime   = time.Hour
	connMaxIdleTime   = time.Minute * 30
	maxIdleConns      = 10
	maxOpenConns      = 10
	healthCheckPeriod = time.Minute

	mongoMaxPoolSize        = 100
	mongoConnectTimeout     = 5 * time.Second
	mongoSocketTimeout      = time.Second
	mongoTransactionTimeout = time.Second
)

func NewSQLite() *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	return db
}
