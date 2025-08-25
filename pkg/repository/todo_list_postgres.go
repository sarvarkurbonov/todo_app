package repository

import "database/sql"

type TodoListRepository struct {
	db *sql.DB
}

func NewTodoListPostgres(db *sql.DB) *TodoListRepository {
	return &TodoListRepository{
		db: db,
	}
}

func (r *TodoListRepository) CreateList(userId int, list TodoList) (int, error) {
	return 0, nil
}
