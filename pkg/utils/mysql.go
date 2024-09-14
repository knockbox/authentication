package utils

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLConnection func for creating a mysql connection.
func MySQLConnection() (*sqlx.DB, error) {
	maxConnections, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConnections, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConnections, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	db, err := sqlx.Connect("mysql", getDatbaseConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database, %w", err)
	}

	// Set database options from environment
	db.SetMaxOpenConns(maxConnections)
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetConnMaxLifetime(time.Duration(maxLifetimeConnections))

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database, %w", err)
	}

	return db, nil
}

// getDatbaseConnectionString constructs the default connection string from environment variables.
func getDatbaseConnectionString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SCHEMA"),
	)
}
