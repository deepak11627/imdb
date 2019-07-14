package resources

// Movie represents a Movie entity
type Movie struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Score    float32  `json:"imdb_score"`
	Genre    []string `json:"genre"`
	Director string   `json:"director"`
	//Popularity float32  `json:"99popularity"`
}

type Review struct {
	UserID  int
	MovieID int
	Rating  int
}

// MovieService represents Movie DB operations
type MovieService interface {
	CreateMovie(name, director string, genre []string, score float32) (int, error)
	SaveMovie(ID int, name, director string, genre []string, score float32) (int, error)
	GetMovie(id int) (*Movie, error)
}
