package data

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/url"
	"time"
)

type EatenItemsModel struct {
	DB *pgxpool.Pool
}

type EatenItem struct {
	Id        int       `json:"id,omitempty" `
	Quantity  int       `json:"quantity,omitempty" validate:"required,gte=1"`
	EatenDate time.Time `json:"eaten_date"`
	ItemId    int       `json:"item_id,omitempty"`
}

type EatenItemResponse struct {
	Id        int       `json:"id,omitempty" `
	Quantity  int       `json:"quantity,omitempty" validate:"required,gte=1"`
	EatenDate time.Time `json:"eaten_date"`
	ItemName  string    `json:"item_name,omitempty"`
}

type EatenItemFilters struct {
	Span      string `validate:"oneof=week month year all"`
	AfterDate time.Time
}

func NewEatenItemFilters(query url.Values) (filters EatenItemFilters, err error) {

	filters.Span = query.Get("time_span")

	if filters.Span == "" {
		filters.Span = "week"
	}

	v := validator.New()
	err = v.Struct(filters)

	if err != nil {
		return EatenItemFilters{}, err
	}

	switch filters.Span {
	case "week":
		filters.AfterDate = time.Now().Add(-time.Hour * 24 * 7)
		break
	case "month":
		filters.AfterDate = time.Now().Add(-time.Hour * 24 * 7 * 31)
		break
	}

	return filters, nil
}

func (e EatenItemsModel) Create(ctx context.Context, item EatenItem) error {
	stmt := `
			INSERT INTO eatenitems(quantity, item_id) 
			VALUES ($1, $2)
`

	args := []interface{}{item.Quantity, item.ItemId}

	_, err := e.DB.Exec(ctx, stmt, args)

	if err != nil {
		return err
	}

	return nil
}

func (e EatenItemsModel) CreateList(ctx context.Context, items []*EatenItem) error {

	tx, err := e.DB.Begin(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	for _, item := range items {
		stmt := `
			INSERT INTO eatenitems(quantity, item_id) 
			VALUES ($1, $2)
	`

		args := []interface{}{item.Quantity, item.ItemId}

		_, err := tx.Exec(ctx, stmt, args...)

		if err != nil {
			return err
		}

	}

	return tx.Commit(ctx)
}

func (e EatenItemsModel) Get(ctx context.Context, itemId int, filters EatenItemFilters) (eatenItems []EatenItemResponse, err error) {

	stmt := `
		SELECT e.id, quantity, eaten_date, i.name
		FROM eatenitems e 
			JOIN public.items i on i.id = e.item_id
		WHERE item_id = $1 AND eaten_date >= $2
	`

	rows, err := e.DB.Query(ctx, stmt, itemId, filters.AfterDate)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var eatenItem EatenItemResponse

		err := rows.Scan(&eatenItem.Id, &eatenItem.Quantity, &eatenItem.EatenDate, &eatenItem.ItemName)

		if err != nil {
			return nil, err
		}

		eatenItems = append(eatenItems, eatenItem)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return eatenItems, nil
}
