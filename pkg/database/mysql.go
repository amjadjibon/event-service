package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"event-service/conf"
)

const (
	driverName = "mysql"
)

func GetDB(cfg conf.Config) *sql.DB {
	db, err := sql.Open(driverName, cfg.DbUrl)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxIdleTime(time.Minute * 3)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// ping to check if database is alive
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
