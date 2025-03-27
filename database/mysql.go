package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQL(dbUrl string) *sql.DB {
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Panic("error while creating connection to the database!!", err)
	}

	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pingErr := db.PingContext(ctx)
	if pingErr != nil {
		log.Panic("could not ping database", pingErr)
	}
	return db
}

func IsMysqlReady() bool {
	// TODO: implement check if mySQL is ready
	return true
}
