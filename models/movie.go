package  models

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Overview    string `json:"overview"`
	PosterPath  string `json:"poster_path"`
	GenreIDs    []int  `json:"genre_ids"`
	GenreNames  []string `json:"genre_names"`
	
}
type MovieListResponse struct {
	Results []Movie `json:"results"`
	Page    int     `json:"page"`
	TotalPages int  `json:"total_pages"`
}

