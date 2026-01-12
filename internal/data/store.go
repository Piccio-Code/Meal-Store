package data

import (
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
