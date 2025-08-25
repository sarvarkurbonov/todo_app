package service

import (
	todo "Todo_rest_api"
	"Todo_rest_api/pkg/repository"
	"context"
)

type TodoItemServer struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItem(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemServer {
	return &TodoItemServer{repo: repo, listRepo: listRepo}
}

func (s *TodoItemServer) CreateItem(userId, listId int, item todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}
	return s.repo.CreateItem(listId, item)
}

func (s *TodoItemServer) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(userId, listId)

}
func (s *TodoItemServer) GetById(ctx context.Context, userId, itemId int) (todo.TodoItem, error) {
	return s.repo.GetById(ctx, userId, itemId)
}
func (s *TodoItemServer) Delete(ctx context.Context, userId, itemId int) error {
	return s.repo.Delete(ctx, userId, itemId)
}
func (s *TodoItemServer) UpdateItem(ctx context.Context, userId int, itemId int, input todo.UpdateItem) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateItem(ctx, userId, itemId, input)
}
