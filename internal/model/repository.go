package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

type AdDB struct {
	ID          int
	CreatedAt   time.Time
	Title       string
	Description string
	Price       int
	PhotoLinks  []string
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

// Get reads a record by id
func (r *Repository) Get(ctx context.Context, id int) (AdDB, error) {
	var ad AdDB
	row := r.db.QueryRowContext(ctx,
		"SELECT id, created_at, title, description, price, photo_links FROM public.ad WHERE id = $1", id)
	err := row.Scan(&ad.ID, &ad.CreatedAt, &ad.Title, &ad.Description, &ad.Price, pq.Array(&ad.PhotoLinks))
	if err != nil {
		return AdDB{}, err
	}
	return ad, err
}

// Create saves a new ad record
func (r *Repository) Create(ctx context.Context, ad AdDB) (int , error) {
	var id int
	err := r.db.QueryRowContext(ctx, "INSERT INTO public.ad (id, created_at, title, description, price, photo_links ) VALUES (DEFAULT, CURRENT_TIMESTAMP, $1, $2, $3, $4 ) RETURNING id",
		ad.Title, ad.Description, ad.Price, pq.Array(ad.PhotoLinks)).Scan(&id)
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return 0, err
	}

	return id, nil
}

// Count returns the number of records.
func (r *Repository) Count(ctx context.Context) (int, error) {
	var count int
	row := r.db.QueryRowContext(ctx, "select COUNT(*) from public.ad")
	err := row.Err()
	if err != nil{
		log.Printf("Unable to execute the query. %v", err)
		return 0, err
	}

	err = row.Scan(&count)
	if err != nil{
		log.Printf("Unable to scan the row. %v", err)
		return 0, err
	}

	return count, err
}

// List retrieves the ad records.
func (r *Repository) List(ctx context.Context, offset, limit int, sort string, desc bool) ([]AdDB, error) {
	var ads []AdDB
	var s string
	switch sort {
	case "price":
		s = "ORDER BY price"
	case "created":
		fallthrough
	default:
		s = "ORDER BY created_at"
	}

	if desc {
		s += " DESC"
	}

	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf("SELECT id, created_at, title, description, price, photo_links FROM public.ad %s LIMIT $1 OFFSET $2", s),
		limit, offset)
	if err != nil {
		log.Printf("Unable to execute the query. %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ad AdDB
		err = rows.Scan(&ad.ID, &ad.CreatedAt, &ad.Title, &ad.Description, &ad.Price, pq.Array(&ad.PhotoLinks))
		if err != nil {
			log.Printf("Unable to scan the row. %v", err)
			return nil, err
		}

		ads = append(ads, ad)
	}

	return ads, err
}
