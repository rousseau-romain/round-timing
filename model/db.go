package model

import (
	"database/sql"
	"log"
	"round-timing/config"

	"github.com/go-sql-driver/mysql"
)

func ConnectDb() *sql.DB {
	cfg := mysql.Config{
		User:                 config.DB_USER,
		Passwd:               config.DB_PASSWORD,
		Net:                  "tcp",
		Addr:                 config.DB_HOST + ":" + config.DB_PORT,
		DBName:               config.DB_NAME,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Println("err", err)

		log.Fatal(err)
	}
	return db
}

var db = ConnectDb()
