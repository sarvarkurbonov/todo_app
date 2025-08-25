package repository

import (
	todo "Todo_rest_api"
	"database/sql"
	"fmt"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) CreateUser(user todo.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO users (name, username, password) VALUES ($1, $2, $3) RETURNING id")
	if err := r.db.QueryRow(query, user.Name, user.Username, user.Password).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username string) (todo.User, error) {
	// get user from database
	var user todo.User
	query := fmt.Sprintf("SELECT id, name, username, password FROM %s WHERE username=$1", UsersTable)
	err := r.db.QueryRow(query, username).Scan(&user.Id, &user.Name, &user.Username, &user.Password)

	return user, err
}
