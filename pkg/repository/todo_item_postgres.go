package repository

import (
	todo "Todo_rest_api"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const cacheTTL = time.Minute

type TodoItemRepository struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewTodoItemPostgres(db *sql.DB, rdb *redis.Client) *TodoItemRepository {
	return &TodoItemRepository{
		db:  db,
		rdb: rdb,
	}
}

func (r *TodoItemRepository) CreateItem(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf(
		"INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id",
		TodoItemTable,
	)

	if err := tx.QueryRow(createItemQuery, item.Title, item.Description).Scan(&itemId); err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}
	createUserItem := fmt.Sprintf("Insert into %s (list_id, item_id) VALUES ($1, $2)", ListItemsTable)
	if _, err := tx.Exec(createUserItem, listId, itemId); err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return itemId, nil
}

func (r *TodoItemRepository) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := fmt.Sprintf(`
    SELECT ti.id, ti.title, ti.description,ti.done
    FROM %s ti
    INNER JOIN %s li ON ti.id = li.item_id inner join %s ul on li.list_id = ul.list_id
    WHERE li.list_id = $1 and ul.user_id = $2
    `, TodoItemTable, ListItemsTable, UsersListTable)
	rows, err := r.db.Query(query, listId, userId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var item todo.TodoItem
		if err := rows.Scan(&item.Id, &item.Title, &item.Description, &item.Done); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *TodoItemRepository) GetById(ctx context.Context, userId, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem

	key := fmt.Sprintf("u:%d:item:%d", userId, itemId)

	// 1) try cache
	if r.rdb != nil {
		if b, err := r.rdb.Get(ctx, key).Bytes(); err == nil {
			if err := json.Unmarshal(b, &item); err == nil {
				logrus.Info("item found in cache")
				return item, nil

			}
		}
	}

	query := fmt.Sprintf(`
        SELECT ti.id, ti.title, ti.description, ti.done
        FROM %s ti
        INNER JOIN %s li ON ti.id = li.item_id
        INNER JOIN %s ul ON li.list_id = ul.list_id
        WHERE ti.id = $1 AND ul.user_id = $2
    `, TodoItemTable, ListItemsTable, UsersListTable)

	err := r.db.QueryRow(query, itemId, userId).
		Scan(&item.Id, &item.Title, &item.Description, &item.Done)

	if err != nil {
		return item, err
	}

	if r.rdb != nil {
		if b, err := json.Marshal(item); err == nil {
			_ = r.rdb.Set(ctx, key, b, cacheTTL).Err()
		}
	}

	return item, nil

}

func (r *TodoItemRepository) Delete(ctx context.Context, userId, itemId int) error {
	query := fmt.Sprintf(`
        DELETE FROM %s ti
        USING %s li, %s ul
        WHERE ti.id = li.item_id
          AND li.list_id = ul.list_id
          AND ti.id = $1
          AND ul.user_id = $2
    `, TodoItemTable, ListItemsTable, UsersListTable)

	res, err := r.db.Exec(query, itemId, userId)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no item found with id=%d for user=%d", itemId, userId)
	}

	if r.rdb != nil {
		key := fmt.Sprintf("u:%d:item:%d", userId, itemId)
		_ = r.rdb.Del(ctx, key).Err()
	}

	return nil
}

func (r *TodoItemRepository) UpdateItem(ctx context.Context, userId int, itemId int, input todo.UpdateItem) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done = $%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`
        UPDATE %s ti
        SET %s
        FROM %s li, %s ul
        WHERE ti.id = li.item_id
          AND li.list_id = ul.list_id
          AND ti.id = $%d
          AND ul.user_id = $%d
    `, TodoItemTable, setQuery, ListItemsTable, UsersListTable, argId, argId+1)

	args = append(args, itemId, userId)

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("no item found with id=%d for user=%d", itemId, userId)
	}

	if r.rdb != nil {
		key := fmt.Sprintf("u:%d:item:%d", userId, itemId)
		_ = r.rdb.Del(ctx, key).Err()
	}

	return nil
}
