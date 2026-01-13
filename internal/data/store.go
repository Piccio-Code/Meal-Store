package data

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Store struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	UserID    int       `json:"-"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type StoreInput struct {
	Name string `json:"name" validate:"required,gte=3,lte=15"`
}

type ModelStore struct {
	DB *pgxpool.Pool
}

func (m ModelStore) Insert(ctx context.Context, newStore StoreInput, userId string) (storeId int, err error) {
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

func (m ModelStore) Get(ctx context.Context, storeId int, userId string) (store Store, err error) {

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

func (m ModelStore) Update(ctx context.Context, newStore StoreInput, storeId int, userId string) error {
	stmt := `UPDATE stores
			 SET name= $1
			 WHERE id = $2 AND user_id = $3
			`

	args := []interface{}{newStore.Name, storeId, userId}

	_, err := m.DB.Exec(ctx, stmt, args)

	return err
}

func (m ModelStore) Delete(ctx context.Context, storeId int, userId string) error {
	stmt := `DELETE FROM stores
			 WHERE id = $1 AND user_id = $2`

	args := []interface{}{storeId, userId}

	_, err := m.DB.Exec(ctx, stmt, args)

	return err
}
