package todo

import "fmt"

type TodoList struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" binding:"required" db:"title"`
	Description string `json:"description" db:"description"`
}

type UsersList struct {
	Id     int
	UserId int
	ListId int
}

type TodoItem struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Done        bool   `json:"done" db:"done"`
}
type ListItem struct {
	Id     int
	ListId int
	ItemId int
}

type UpdateList struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}
type UpdateItem struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (i *UpdateList) Validate() error {
	if i.Title == nil && i.Description == nil {
		return fmt.Errorf("nothing to update")
	}
	return nil
}
func (i *UpdateItem) Validate() error {
	if i.Title == nil && i.Description == nil {
		return fmt.Errorf("nothing to update")
	}
	return nil
}
