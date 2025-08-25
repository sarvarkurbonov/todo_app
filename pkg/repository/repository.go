package repository

import (
	todo "Todo_rest_api"
	"context"
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GetUser(username string) (todo.User, error)
}
type TodoList interface {
	CreateList(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId, listId int) (todo.TodoList, error)
	DeleteList(userId, listId int) error
	UpdateList(userId int, listId int, input todo.UpdateList) error
}
type TodoItem interface {
	CreateItem(listId int, item todo.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo.TodoItem, error)
	GetById(ctx context.Context, userId, itemId int) (todo.TodoItem, error)
	Delete(ctx context.Context, userId, itemId int) error
	UpdateItem(ctx context.Context, userId int, itemId int, input todo.UpdateItem) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sql.DB, rdb *redis.Client) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db, rdb),
	}
}
