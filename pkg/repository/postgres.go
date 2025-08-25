package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	UsersTable     = "users"
	TodoListTable  = "todo_list"
	TodoItemTable  = "todo_item"
	UsersListTable = "users_lists"
	ListItemsTable = "list_items"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
	SslMode  string
}

func NewPostgresDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DbName, cfg.SslMode))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	s := db.Stats()
	logrus.WithFields(logrus.Fields{
		"Open":  s.OpenConnections,
		"InUse": s.InUse,
		"Idle":  s.Idle,
	}).Info("DB pool configured")

	return db, nil
}
