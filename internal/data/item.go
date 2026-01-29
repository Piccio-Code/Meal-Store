package data

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type ItemModel struct {
	DB *pgxpool.Pool
}

type Item struct {
	Id              *int       `json:"id,omitempty"`
	Name            *string    `json:"name,omitempty" validate:"required"`
	CurrentCapacity *int       `json:"current_capacity,omitempty" validate:"gte=1"`
	StoreId         int        `json:"-"`
	Version         *string    `json:"version"`
	CreatedAt       *time.Time `json:"created_at"`
	ModifiedAt      *time.Time `json:"modified_at"`
}

type UpdateItem struct {
	Id              *int    `json:"id,omitempty" validate:"required"`
	Name            *string `json:"name,omitempty"`
	CurrentCapacity *int    `json:"current_capacity,omitempty" validate:"gte=1"`
	Version         *string `json:"version" validate:"required"`
}

func (m ItemModel) List(ctx context.Context, storeId int, onlyWarnings bool) (items []Item, err error) {
	tx, err := m.DB.Begin(ctx)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	stmt := `
			SELECT id, name, current_capacity, store_id, version, created_at, modified_at
			FROM items
			WHERE store_id = $1 AND (current_capacity <= 1 OR NOT $2)
	`

	rows, err := tx.Query(ctx, stmt, storeId, onlyWarnings)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var item Item

		err := rows.Scan(&item.Id, &item.Name, &item.CurrentCapacity, &item.StoreId, &item.Version, &item.CreatedAt, &item.ModifiedAt)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, tx.Commit(ctx)
}

func (m ItemModel) Insert(ctx context.Context, newItem *Item) error {
	tx, err := m.DB.Begin(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	stmt := `
			INSERT INTO items(name, current_capacity, store_id)
			VALUES ($1, $2, $3)
			RETURNING id, version, created_at
	`

	args := []interface{}{*newItem.Name, *newItem.CurrentCapacity, newItem.StoreId}

	err = tx.QueryRow(ctx, stmt, args...).Scan(&newItem.Id, &newItem.Version, &newItem.CreatedAt)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (m ItemModel) InsertList(ctx context.Context, newItemList []*Item, storeId int) error {
	tx, err := m.DB.Begin(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	stmt := `
			INSERT INTO items(name, current_capacity, store_id)
			VALUES ($1, $2, $3)
			RETURNING id, version, created_at
	`

	for _, newItem := range newItemList {
		args := []interface{}{*newItem.Name, *newItem.CurrentCapacity, storeId}

		err = tx.QueryRow(ctx, stmt, args...).Scan(&newItem.Id, &newItem.Version, &newItem.CreatedAt)

		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (m ItemModel) Get(ctx context.Context, itemId, storeId int) (item Item, err error) {
	tx, err := m.DB.Begin(ctx)

	if err != nil {
		return Item{}, err
	}

	defer tx.Rollback(ctx)

	stmt := `
			SELECT id, name, current_capacity, store_id, version, created_at, modified_at
			FROM items
			WHERE id = $1 AND store_id = $2
	`

	err = tx.QueryRow(ctx, stmt, itemId, storeId).Scan(&item.Id, &item.Name, &item.CurrentCapacity, &item.StoreId, &item.Version, &item.CreatedAt, &item.ModifiedAt)

	if err != nil {
		return Item{}, err
	}

	return item, tx.Commit(ctx)
}

func (m ItemModel) GetId(ctx context.Context, itemName string, storeId int) (itemId int, err error) {
	tx, err := m.DB.Begin(ctx)

	if err != nil {
		return 0, err
	}

	defer tx.Rollback(ctx)

	stmt := `
			SELECT id
			FROM items
			WHERE name = $1 AND store_id = $2
	`

	err = tx.QueryRow(ctx, stmt, itemName, storeId).Scan(&itemId)

	if err != nil {
		return 0, err
	}

	return itemId, tx.Commit(ctx)
}

func (m ItemModel) Update(ctx context.Context, item *Item) error {
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	stmt := `
			UPDATE items
			SET name = $1, current_capacity = $2, modified_at = now(), version = uuid_generate_v4()
			WHERE id = $3 AND version = $4 AND store_id = $5
			RETURNING version, modified_at, created_at
`

	args := []interface{}{*item.Name, *item.CurrentCapacity, *item.Id, *item.Version, item.StoreId}

	err = tx.QueryRow(ctx, stmt, args...).Scan(&item.Version, &item.ModifiedAt, &item.CreatedAt)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (m ItemModel) Delete(ctx context.Context, itemId, storeId int) error {
	stmt := `
			DELETE FROM items
			WHERE id = $1 AND store_id = $2
	`

	result, err := m.DB.Exec(ctx, stmt, itemId, storeId)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("0 effected rows")
	}

	return err
}
