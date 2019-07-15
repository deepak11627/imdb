package db

import (
	"database/sql"
	"errors"

	"github.com/deepak11627/imdb/logger"
	"github.com/deepak11627/imdb/resources"
)

type DatabaseService interface {
	CreateMovie(name, director string, genre []string, score string) (int, error)
	SaveMovie(ID int, name, director string, genre []string, score string) error
	GetMovie(id int) (*resources.Movie, error)
	SearchMovie(str string) ([]*resources.Movie, error)
}

type Database struct {
	db     *sql.DB
	logger logger.Logger
}

func NewDB(db *sql.DB, l logger.Logger) *Database {
	return &Database{
		db:     db,
		logger: l,
	}
}

func (d *Database) CreateMovie(name, director string, genre []string, score string) (int, error) {
	stmt, err := d.db.Prepare("INSERT INTO `movie` (`name`, `director_name`, `imdb_score`) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, director, score)
	if err != nil {
		return 0, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	if len(genre) > 0 {
		genreIDs, err := d.createGenre(genre)
		if err != nil {
			d.logger.Debug("error while creating genre", err)
			return int(lastInsertedID), err
		}
		d.mapGenre(int(lastInsertedID), genreIDs)
	}

	return int(lastInsertedID), nil
}

func (d *Database) SaveMovie(id int, name, director string, genre []string, score string) error {
	stmt, err := d.db.Prepare("UPDATE movie SET name=?, director_name=?, imdb_score=? WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(name, director, score, id); err != nil {
		return err
	}

	if len(genre) > 0 {
		genreIDs, err := d.createGenre(genre)
		if err != nil {
			return err
		}
		d.mapGenre(id, genreIDs)
	}
	return nil
}

func (d *Database) GetMovie(id int) (*resources.Movie, error) {
	row := d.db.QueryRow("SELECT `id`, `name`, `director_name`, `imdb_score` FROM `movie` WHERE `id` = ?", id)
	if row != nil {
		var m resources.Movie
		var score, director sql.NullString
		if err := row.Scan(&m.ID, &m.Name, &director, &score); err != nil {
			switch err {
			case sql.ErrNoRows:
				return &m, nil
			}
			return nil, err
		}
		if score.Valid {
			m.Score = score.String
		}
		if director.Valid {
			m.Director = director.String
		}
		genre, err := d.getGenre(m.ID)
		if err != nil {
			d.logger.Warn("failed to get genre for movie", "error", err)
		} else {
			m.Genre = genre
		}
		return &m, nil
	}
	return nil, errors.New("not found")
}

// SearchMovie search movies by string , matches movie name and genre
func (d *Database) SearchMovie(str string) ([]*resources.Movie, error) {
	rows, err := d.db.Query("SELECT m.name, m.director, m.imdb_score, GROUP_CONTACT(g.name, ',') as Genre FROM `movie` m LEFT JOIN `genre` g ON m.id = g.movie_id WHERE `m.name` LIKE '%?%' OR g.name LIKE '%?%'", str, str)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var movies []*resources.Movie

	for rows.Next() {
		var m resources.Movie
		var score, director, genre sql.NullString
		if err := rows.Scan(&m.ID, &m.Name, &director, &score, &genre); err != nil {
			return nil, err
		}
		if score.Valid {
			m.Score = score.String
		}
		if director.Valid {
			m.Director = director.String
		}
		movies = append(movies, &m)
	}

	return movies, nil
}

func (d *Database) createGenre(genre []string) ([]int, error) {
	var genreIDs []int
	for _, g := range genre {
		stmt, err := d.db.Prepare("INSERT IGNORE INTO `genre` (`name`) VALUES (?)")
		if err != nil {
			return []int{}, err
		}
		defer stmt.Close()

		result, err := stmt.Exec(g)
		if err != nil {
			return []int{}, err
		}

		lastInsertedID, err := result.LastInsertId()
		if err != nil {
			return []int{}, err
		}
		genreIDs = append(genreIDs, int(lastInsertedID))
	}
	return genreIDs, nil
}

func (d *Database) mapGenre(movieID int, genre []int) error {
	err := d.deleteGenreMapping(movieID)
	if err != nil {
		return err
	}

	params := []interface{}{}
	query := "INSERT IGNORE INTO `movie_genre` (`movie_id`, `genre_id`) VALUES "

	for _, v := range genre {
		query += "(?,?),"
		params = append(params, movieID, v)
	}
	query = query[0 : len(query)-1]

	stmt, err := d.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(params...)

	// if no err, then err will be nil
	return err
}

func (d *Database) deleteGenreMapping(movieID int) error {
	stmt, err := d.db.Prepare("DELETE FROM `movie_genre` WHERE `movie_id` = ?")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(movieID)

	// if no err, then err will be nil
	return err
}

func (d *Database) getGenre(movieID int) ([]string, error) {
	rows, err := d.db.Query(" SELECT g.name FROM `movie_genre` mg LEFT JOIN `genre` g ON mg.genre_id =g.id WHERE mg.movie_id = ?", movieID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var genre []string

	for rows.Next() {
		var g sql.NullString
		if err := rows.Scan(&g); err != nil {
			return nil, err
		}
		if g.Valid {
			genre = append(genre, g.String)
		}
	}

	return genre, nil
}
