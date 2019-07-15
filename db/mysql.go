package db

import (
	"database/sql"
	"errors"

	"github.com/deepak11627/imdb/resources"
)

type DatabaseService interface {
	CreateMovie(name, director string, genre []string, score string) (int, error)
	SaveMovie(ID int, name, director string, genre []string, score string) error
	GetMovie(id int) (*resources.Movie, error)
}

type Database struct {
	db *sql.DB
	//	logger interface{}
}

func NewDB(db *sql.DB) *Database {
	return &Database{
		db: db,
		//logger: logger,
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
	return nil
}

func (d *Database) GetMovie(id int) (*resources.Movie, error) {
	row := d.db.QueryRow("SELECT `id`, `name`, `director_name`, `imdb_score` FROM `movie` WHERE `id` = ?", id)
	if row != nil {
		m := new(resources.Movie)
		if err := row.Scan(m.ID, m.Name, m.Director, m.Score); err != nil {
			return nil, err
		}
		return m, nil
	}
	return nil, errors.New("not found")
}
