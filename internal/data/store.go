package data

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Store struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	UserID    string    `json:"-"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type StoreInput struct {
	Name string `json:"name" validate:"required,gte=3,lte=15"`
}

type StoreModel struct {
	DB *pgxpool.Pool
}

func (m StoreModel) Insert(ctx context.Context, newStore StoreInput, userId string) (storeId int, err error) {
	stmt := `INSERT INTO stores(name, user_id) 
			 VALUES ($1, $2)
			 RETURNING id`

	args := []interface{}{newStore.Name, userId}

	err = m.DB.QueryRow(ctx, stmt, args...).Scan(&storeId)

	if err != nil {
		return 0, err
	}

	return storeId, nil
}

func (m StoreModel) Get(ctx context.Context, storeId int, userId string) (store Store, err error) {

	stmt := `SELECT id, name, user_id, created_at
			 FROM stores
			 WHERE id = $1 AND user_id = $2`

	args := []interface{}{storeId, userId}

	err = m.DB.QueryRow(ctx, stmt, args...).Scan(&store.ID, &store.Name, &store.UserID, &store.CreatedAt)

	if err != nil {
		return Store{}, err
	}

	return store, nil
}

func (m StoreModel) List(ctx context.Context, userId string) (stores []Store, err error) {

	tx, err := m.DB.Begin(ctx)

	if err != nil {
		return nil, err
	}

	stmt := `SELECT id, name, user_id, created_at
			 FROM stores
			 WHERE user_id = $1`

	rows, err := tx.Query(ctx, stmt, userId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var store Store

		err := rows.Scan(&store.ID, &store.Name, &store.UserID, &store.CreatedAt)

		if err != nil {
			return nil, err
		}

		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stores, nil
}

func (m StoreModel) Update(ctx context.Context, newStore StoreInput, storeId int, userId string) error {
	stmt := `UPDATE stores
			 SET name = $1, created_at = $2
			 WHERE id = $3 AND user_id = $4
			`

	args := []interface{}{newStore.Name, time.Now(), storeId, userId}

	_, err := m.DB.Exec(ctx, stmt, args...)

	return err
}

func (m StoreModel) Delete(ctx context.Context, storeId int, userId string) error {
	stmt := `DELETE FROM stores
			 WHERE id = $1 AND user_id = $2`

	args := []interface{}{storeId, userId}

	_, err := m.DB.Exec(ctx, stmt, args...)

	return err
}
