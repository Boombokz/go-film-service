package repositories

import (
	"context"
	"filmservice/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WatchListRepository struct {
	db *pgxpool.Pool
}

func NewWatchListRepository(conn *pgxpool.Pool) *WatchListRepository {
	return &WatchListRepository{db: conn}
}

func (r *WatchListRepository) GetAll(c context.Context) (
	[]models.Movie,
	error,
) {
	sql := `SELECT
    m.id,
    m.title,
    m.description,
    m.release_year,
    m.director,
    m.rating,
    m.is_watched,
    m.trailer_url,
    m.poster_url,
    g.id,
    g.title
FROM watch_list wl
JOIN movies m ON m.id = wl.movie_id
JOIN movies_genres mg ON mg.movie_id = m.id
JOIN genres g ON g.id = mg.genre_id
ORDER BY wl.added_at ASC;
	`

	rows, err := r.db.Query(
		c,
		sql,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := make(
		[]*models.Movie,
		0,
	)
	moviesMap := make(map[int]*models.Movie)

	for rows.Next() {
		var m models.Movie
		var g models.Genre

		err := rows.Scan(
			&m.Id,
			&m.Title,
			&m.Description,
			&m.ReleaseYear,
			&m.Director,
			&m.Rating,
			&m.IsWatched,
			&m.TrailerUrl,
			&m.PosterUrl,
			&g.Id,
			&g.Title,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := moviesMap[m.Id]; !exists {
			moviesMap[m.Id] = &m
			movies = append(
				movies,
				&m,
			)
		}

		moviesMap[m.Id].Genres = append(
			moviesMap[m.Id].Genres,
			g,
		)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	concreteMovies := make(
		[]models.Movie,
		0,
		len(movies),
	)
	for _, v := range movies {
		concreteMovies = append(
			concreteMovies,
			*v,
		)
	}

	return concreteMovies, nil
}

func (r *WatchListRepository) Exists(
	c context.Context,
	id int,
) (
	bool,
	error,
) {
	var exists bool
	err := r.db.QueryRow(
		c,
		`SELECT EXISTS (SELECT 1 FROM watch_list WHERE movie_id = $1)`,
		id,
	).Scan(&exists)
	return exists, err
}

func (r *WatchListRepository) Add(
	c context.Context,
	id int,
) error {
	_, err := r.db.Exec(
		c,
		`INSERT INTO watch_list (movie_id) VALUES ($1)`,
		id,
	)
	return err
}

func (r *WatchListRepository) Delete(
	c context.Context,
	id int,
) error {
	_, err := r.db.Exec(
		c,
		`DELETE FROM watch_list WHERE movie_id = $1`,
		id,
	)
	return err
}
