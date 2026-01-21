package data

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Store struct {
	ID         *int       `json:"id,omitempty"`
	Name       *string    `json:"name,omitempty" validate:"required,gte=3"`
	UserID     *string    `json:"-"`
	Version    *string    `json:"version,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	ModifiedAt *time.Time `json:"updated_at,omitempty"`
}

type StoreModel struct {
	DB *pgxpool.Pool
}

func (m StoreModel) Insert(ctx context.Context, newStore *Store, userId string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO stores(name, user_id) 
			 VALUES ($1, $2)
			 RETURNING id, version, created_at`

	args := []interface{}{newStore.Name, userId}

	return m.DB.QueryRow(ctx, stmt, args...).Scan(&newStore.ID, &newStore.Version, &newStore.CreatedAt)
}

func (m StoreModel) Get(ctx context.Context, storeId int, userId string) (store Store, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `SELECT id, name, user_id, created_at, version, modified_at
			 FROM stores
			 WHERE id = $1 AND user_id = $2`

	args := []interface{}{storeId, userId}

	err = m.DB.QueryRow(ctx, stmt, args...).Scan(&store.ID, &store.Name, &store.UserID, &store.CreatedAt, &store.Version, &store.ModifiedAt)

	if err != nil {
		return Store{}, err
	}

	return store, nil
}

func (m StoreModel) List(ctx context.Context, userId string) (stores []Store, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := m.DB.Begin(ctx)

	if err != nil {
		return nil, err
	}

	stmt := `SELECT id, name, user_id, created_at, version, modified_at
			 FROM stores
			 WHERE user_id = $1`

	rows, err := tx.Query(ctx, stmt, userId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var store Store

		err := rows.Scan(&store.ID, &store.Name, &store.UserID, &store.CreatedAt, &store.Version, &store.ModifiedAt)

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

func (m StoreModel) Update(ctx context.Context, newStore *Store, userId string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `UPDATE stores
			 SET name = $1, version = uuid_generate_v4(), modified_at = NOW()
			 WHERE id = $2 AND user_id = $3 AND version = $4
			 RETURNING version, modified_at`

	args := []interface{}{newStore.Name, newStore.ID, userId, newStore.Version}

	return m.DB.QueryRow(ctx, stmt, args...).Scan(&newStore.Version, &newStore.ModifiedAt)
}

func (m StoreModel) Delete(ctx context.Context, storeId int, userId string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM stores
			 WHERE id = $1 AND user_id = $2`

	args := []interface{}{storeId, userId}

	result, err := m.DB.Exec(ctx, stmt, args...)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("0 effected rows")
	}

	return err
}
