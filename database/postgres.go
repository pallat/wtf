package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Config(dbUrl string) *pgxpool.Config {
	dbConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		log.Panic("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = int32(maxOpenConns)
	dbConfig.MaxConnLifetime = connMaxLifetime
	dbConfig.MaxConnIdleTime = connMaxIdleTime
	dbConfig.HealthCheckPeriod = healthCheckPeriod

	return dbConfig
}

// NewPostgresDB new postgress instance
// Example:
// db := database.NewPostgresDB(os.Getenv("PGURL"))
// defer db.Close()
// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// defer cancel()
//
// rows, err := db.Query(ctx, "select product_name,retail_price from products")
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for rows.Next() {
//		var pname string
//		var rprice float64
//		if err := rows.Scan(&pname, &rprice); err != nil {
//			log.Println(err)
//		}
//		fmt.Println(pname, rprice)
//	}
func NewPostgresDB(dbUrl string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	connPool, err := pgxpool.NewWithConfig(ctx, Config(dbUrl))
	if err != nil {
		log.Panic("error while creating connection to the database!!", err)
	}

	connection, err := connPool.Acquire(ctx)
	if err != nil {
		log.Panic("error while acquiring connection from the database pool!!", err)
	}
	defer connection.Release()

	err = connection.Ping(ctx)
	if err != nil {
		log.Panic("could not ping database", err)
	}

	return connPool
}
