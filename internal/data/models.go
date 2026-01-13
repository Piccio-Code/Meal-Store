package data

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Models struct {
	Stores StoreModel
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Stores: StoreModel{DB: db},
	}
}
