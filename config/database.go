package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"os"
	"strconv"
	"time"
)

func NewDbPool(logger Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_SCHEMA"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Error("failed to connect DB", zap.Error(err))
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logger.Error("failed ping connect DB", zap.Error(err))
		return nil, err
	}

	maxLifeTime, _ := time.ParseDuration(os.Getenv("DATABASE_CONN_MAX_LIFETIME"))
	openConns, _ := strconv.Atoi(os.Getenv("DATABASE_MAX_OPEN_CONN"))
	idleConns, _ := strconv.Atoi(os.Getenv("DATABASE_MAX_IDLE_CONN"))

	db.SetConnMaxLifetime(maxLifeTime)
	db.SetMaxOpenConns(openConns)
	db.SetMaxIdleConns(idleConns)

	return db, nil
}
