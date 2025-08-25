package repository

import (
	todo "Todo_rest_api"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type TodoListRepository struct {
	db *sql.DB
}

func NewTodoListPostgres(db *sql.DB) *TodoListRepository {
	return &TodoListRepository{
		db: db,
	}
}

func (r *TodoListRepository) CreateList(userId int, list todo.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf(
		"INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id",
		TodoListTable,
	)

	if err := tx.QueryRow(createListQuery, list.Title, list.Description).Scan(&id); err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}
	createUserList := fmt.Sprintf("Insert into %s (user_id, list_id) VALUES ($1, $2)", UsersListTable)
	if _, err := tx.Exec(createUserList, userId, id); err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *TodoListRepository) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList
	query := fmt.Sprintf(`
    SELECT tl.id, tl.title, tl.description
    FROM %s tl
    INNER JOIN %s ul ON tl.id = ul.list_id
    WHERE ul.user_id = $1
`, TodoListTable, UsersListTable)

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	// Iterate over the rows
	for rows.Next() {
		var list todo.TodoList
		if err := rows.Scan(
			&list.Id,
			&list.Title,
			&list.Description,
		); err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}

	// Check for errors from iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return lists, nil
}

func (r *TodoListRepository) GetById(userId, listId int) (todo.TodoList, error) {
	var list todo.TodoList

	query := fmt.Sprintf(`
        SELECT tl.id, tl.title, tl.description
        FROM %s tl
        INNER JOIN %s ul ON tl.id = ul.list_id
        WHERE tl.id = $1 AND ul.user_id = $2
    `, TodoListTable, UsersListTable)

	err := r.db.QueryRow(query, listId, userId).
		Scan(&list.Id, &list.Title, &list.Description)
	if err != nil {
		return list, err
	}

	return list, nil
}

func (r *TodoListRepository) DeleteList(userId, listId int) error {
	query := fmt.Sprintf(`
        DELETE FROM %s tl
        USING %s ul
        WHERE tl.id = ul.list_id
          AND ul.user_id = $1
          AND ul.list_id = $2
    `, TodoListTable, UsersListTable)

	_, err := r.db.Exec(query, userId, listId)
	return err
}

func (r *TodoListRepository) UpdateList(userId, listId int, input todo.UpdateList) error {
	var setValues []string
	var args []interface{}
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`
        UPDATE %s tl
        SET %s
        FROM %s ul
        WHERE tl.id = ul.list_id
          AND ul.user_id = $%d
          AND ul.list_id = $%d
    `, TodoListTable, setQuery, UsersListTable, argId, argId+1)

	args = append(args, userId, listId)
	logrus.Debugf("update query: %s", query)
	logrus.Debugf("update args: %v", args)

	_, err := r.db.Exec(query, args...)
	return err
}
