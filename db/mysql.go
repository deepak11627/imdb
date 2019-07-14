package db

import (
	"database/sql"

	"github.com/deepak11627/imdb/resources"
)

type DatabaseService interface {
	CreateMovie(name, director string, genre []string, score float32) (int, error)
	SaveMovie(ID int, name, director string, genre []string, score float32) (int, error)
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

func (d *Database) CreateMovie(name, director string, genre []string, score float32) (int, error) {
	stmt, err := d.db.Prepare("INSERT INTO `movie` (`name`, `director`, `score`) VALUES (?, ?, ?)")
	if err != nil {
		return nil, err
	}

	result, err := stmt.Exec(name, director, score)
	if err != nil {
		return nil, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return lastInsertedID,nil
}

func (d *Database) SaveMovie(id int, name, director string, genre []string, score float32) (int, error) {
	stmt, err := d.db.Prepare("UPDATE `movie` SET `name` = ?, ")
	if err != nil {
		return nil, err
	}

	result, err := stmt.Exec(name, director, score)
	if err != nil {
		return nil, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return lastInsertedID,nil
}

func (d *Database) GetMovie(id int) (*resources.Movie, error) {
	stmt, err := d.db.QueryRow("SELECT * FROM `movie` WHERE `id` = ?")
	if err != nil {
		return nil, err
	}

	result, err := stmt.Exec(name, director, score)
	if err != nil {
		return nil, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	return lastInsertedID,nil
	return &resources.Movie{ID: id}, nil
}
