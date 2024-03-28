package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/PriyanshuSharma23/greenlight/internal/validator"
	"github.com/lib/pq"
)

type Movie struct {
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Genres    []string  `json:"genres,omitempty"`
	ID        int       `json:"id"`
	Year      int       `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Version   int       `json:"version"`
}

func ValidateMovie(v *validator.Validator, m *Movie) {
	v.Check(validator.NotBlank(m.Title), "title", "must be provided")
	v.Check(validator.MaxChars(m.Title, 500), "title", "must not be more than 500 characters")

	v.Check(m.Year != 0, "year", "must be provided")
	v.Check(validator.Min(m.Year, 1888), "year", "must be greater or equal to than 1888")
	v.Check(validator.Max(m.Year, time.Now().Year()), "year", "must not be in future")

	v.Check(m.Runtime != 0, "runtime", "must be provided")
	v.Check(validator.Min(int(m.Runtime), 1), "runtime", "must be a positive value")

	v.Check(m.Genres != nil, "genres", "must be provided")
	v.Check(len(m.Genres) <= 5, "genres", "length must be less than or equal to 5")
	v.Check(validator.Unique(m.Genres), "genres", "values must be unique")
}

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(mov *Movie) error {
	stmt := `INSERT INTO movies (title, year, runtime, genres)
          VALUES ($1, $2, $3, $4)
          RETURNING id, created_at, version`

	args := []any{mov.Title, mov.Year, mov.Runtime, pq.Array(mov.Genres)}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	return m.DB.QueryRowContext(ctx, stmt, args...).Scan(&mov.ID, &mov.CreatedAt, &mov.Version)
}

func (m MovieModel) Get(id int) (*Movie, error) {
	if id < 0 {
		return nil, ErrNoRecordFound
	}

	var movie Movie

	stmt := `SELECT id, created_at, title, year, runtime, genres, version
           FROM movies
           WHERE id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecordFound
		}
		return nil, err
	}

	return &movie, nil
}

func (m MovieModel) Update(mov *Movie) error {
	stmt := `UPDATE movies 
           SET title=$1, year=$2, runtime=$3, genres=$4, version = version + 1
           WHERE id=$5 AND version=$6
           RETURNING version;`

	args := []any{mov.Title, mov.Year, mov.Runtime, pq.Array(mov.Genres), mov.ID, mov.Version}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, stmt, args...).Scan(&mov.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m MovieModel) Delete(id int) error {
	stmt := `DELETE FROM movies
           WHERE id=$1`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	r, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNoRecordFound
	}

	return nil
}

func (m MovieModel) GetAll(title string, genres []string, f Filters) ([]*Movie, Metadata, error) {
	stmt := fmt.Sprintf(`SELECT COUNT(*) OVER(),id, created_at, title, year, runtime, genres, version
           FROM movies
		   WHERE ((to_tsvector('simple', title) @@ plainto_tsquery('simple', $1)) OR $1 = '')
		   AND (genres @> $2 OR $2 = '{}')
		   ORDER BY %s %s, id ASC
		   LIMIT $3 OFFSET $4`, f.sortColumn(), f.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{
		title,
		pq.Array(genres),
		f.limit(),
		f.offset(),
	}

	rows, err := m.DB.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	var totalRecords int

	movies := make([]*Movie, 0)

	for rows.Next() {
		var m Movie

		err := rows.Scan(
			&totalRecords,
			&m.ID,
			&m.CreatedAt,
			&m.Title,
			&m.Year,
			&m.Runtime,
			pq.Array(&m.Genres),
			&m.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		movies = append(movies, &m)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, f.Page, f.PageSize)

	return movies, metadata, nil
}
