package service

import (
	todo "Todo_rest_api"
	"Todo_rest_api/pkg/repository"
)

type TodoListServer struct {
	repo repository.TodoList
}

func NewTodoList(repo repository.TodoList) *TodoListServer {
	return &TodoListServer{repo: repo}
}

func (s *TodoListServer) CreateList(userId int, list todo.TodoList) (int, error) {
	return s.repo.CreateList(userId, list)
}

func (s *TodoListServer) GetAll(userId int) ([]todo.TodoList, error) {
	return s.repo.GetAll(userId)
}
func (s *TodoListServer) GetById(userId, listId int) (todo.TodoList, error) {
	return s.repo.GetById(userId, listId)
}
func (s *TodoListServer) DeleteList(userId, listId int) error {
	return s.repo.DeleteList(userId, listId)
}
func (s *TodoListServer) UpdateList(userId, listId int, input todo.UpdateList) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateList(userId, listId, input)
}
