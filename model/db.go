package model

import (
	"database/sql"
	"log"
	"log/slog"
	"time"

	"github.com/rousseau-romain/round-timing/config"

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
		slog.Error("db connection error", "error", err.Error())
		log.Fatal(err)
	}

	// Connection pool configuration
	db.SetMaxOpenConns(25)               // Max open connections to database
	db.SetMaxIdleConns(5)                // Max idle connections in pool
	db.SetConnMaxLifetime(5 * time.Minute) // Max connection reuse time

	// Verify connection works (sql.Open doesn't actually connect)
	if err := db.Ping(); err != nil {
		slog.Error("db ping error", "error", err.Error())
		log.Fatal(err)
	}

	return db
}

var db = ConnectDb()

// DB is exported for use by subpackages
var DB = db

// LogDBStats logs current database connection pool statistics
func LogDBStats() {
	stats := db.Stats()
	slog.Info("db pool stats",
		"open", stats.OpenConnections,
		"in_use", stats.InUse,
		"idle", stats.Idle,
		"wait_count", stats.WaitCount,
		"wait_duration", stats.WaitDuration,
		"max_idle_closed", stats.MaxIdleClosed,
		"max_lifetime_closed", stats.MaxLifetimeClosed,
	)
}

// StartDBStatsLogger starts a background goroutine that logs DB stats periodically
func StartDBStatsLogger(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			LogDBStats()
		}
	}()
}
