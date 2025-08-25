package service

import (
	todo "Todo_rest_api"
	"Todo_rest_api/pkg/repository"
	"context"
)

type Authorization interface {
	CreateUser(user todo.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	CreateList(userId int, list todo.TodoList) (int, error)
	GetAll(userId int) ([]todo.TodoList, error)
	GetById(userId, listId int) (todo.TodoList, error)
	DeleteList(userId, listId int) error
	UpdateList(userId int, listId int, input todo.UpdateList) error
}
type TodoItem interface {
	CreateItem(userId, listId int, item todo.TodoItem) (int, error)
	GetAll(userId, listId int) ([]todo.TodoItem, error)
	GetById(ctx context.Context, userId, itemId int) (todo.TodoItem, error)
	Delete(ctx context.Context, userId, itemId int) error
	UpdateItem(ctx context.Context, userId int, itemId int, input todo.UpdateItem) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoList(repos.TodoList),
		TodoItem:      NewTodoItem(repos.TodoItem, repos.TodoList),
	}
}
