package models

type Movie struct {
	Id          int
	Title       string
	Description string
	ReleaseYear int
	Director    string
	Rating      int
	IsWatched   bool
	TrailerUrl  string
	Genres      []Genre
	PosterUrl   string
}

type MovieFilters struct {
	SearchTerm string
	GenreId    string
	IsWatched  string
	Sort       string
}
