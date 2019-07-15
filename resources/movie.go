package resources

// Movie represents a Movie entity
type Movie struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Score    string   `json:"imdb_score"`
	Genre    []string `json:"genre"`
	Director string   `json:"director"`
	//Popularity float32  `json:"99popularity"`
}

// IsValid tells if a movie is valid based on its properties
func (m *Movie) IsValid() bool {
	if m.Name == "" {
		return false
	}
	return true
}

// IsEmpty checks if the movie instance is empty ie no values set
func (m *Movie) IsEmpty() bool {
	if m.ID == 0 {
		return true
	}
	return false
}

type Review struct {
	UserID  int
	MovieID int
	Rating  int
}

// MovieService represents Movie DB operations
type MovieService interface {
	CreateMovie(name, director string, genre []string, score string) (int, error)
	SaveMovie(ID int, name, director string, genre []string, score string) error
	GetMovie(id int) (*Movie, error)
	SearchMovie(str string) ([]*Movie, error)
}
