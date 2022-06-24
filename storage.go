package main

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	Create(album) album
	Read() []album
	ReadOne(string) (album, error)
	Update(string, album) (album, error)
	Delete(string) error
}

type MemeoryStorage struct {
	albums []album
}

func (s MemeoryStorage) Create(am album) album {
	s.albums = append(s.albums, am)
	return am
}
func (s MemeoryStorage) ReadOne(id string) (album, error) {
	for _, v := range s.albums {
		if v.ID == id {
			return v, nil
		}
	}
	return album{}, errors.New("not_found")
}
func (s MemeoryStorage) Read() []album {
	return s.albums
}
func (s MemeoryStorage) Update(id string, newAlbum album) (album, error) {
	for i, _ := range s.albums {
		if s.albums[i].ID == id {
			// c.BindJSON(&albums[i])
			s.albums[i] = newAlbum
			// c.IndentedJSON(http.StatusNoContent, albums[i])
			return s.albums[i], nil
		}
	}
	return album{}, errors.New("not_found")
}
func (s MemeoryStorage) Delete(id string) error {
	for i, v := range s.albums {
		if v.ID == id {
			s.albums = append(s.albums[:i], s.albums[i+1:]...)
			return nil
		}
	}
	return errors.New("not_found")

}
func NewMemoryStorage() MemeoryStorage {
	var albums = []album{
		{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
	return MemeoryStorage{albums: albums}
}

type PostgresStorage struct {
	db *sql.DB
}

func (p PostgresStorage) CreateSchema() error {
	_, err := p.db.Exec("Create TABLE if not EXISTS albums (ID char(16) primary key,Title char(128), Artist char(128), Price decimal)")
	return err
}

func NewPostgresStorage() PostgresStorage {
	connStr := "user=user dbname=db password=pass sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	storage := PostgresStorage{db: db}
	err = storage.CreateSchema()
	if err != nil {
		log.Fatal(err)
	}
	return storage
}

func (p PostgresStorage) Create(am album) album {
	p.db.QueryRow("insert into albums(ID,Title,Artist,Price) values($1,$2,$3,$4)", am.ID, am.Title, am.Artist, am.Price)
	return am
}
func (p PostgresStorage) ReadOne(id string) (album, error) {
	var album album
	p.db.QueryRow("select * from albums where id = $1", id).Scan(&album)
	return album, nil
}
func (p PostgresStorage) Update(id string, am album) (album, error) {
	err := p.db.QueryRow("update albums set Title=$1 and Artist=$2 and Price=$3 where id=$4 ", am.Title, am.Artist, am.Price, am.ID)
	if err != nil {
		return album{}, errors.New("not_found")
	}
	return am, nil
}
func (p PostgresStorage) Delete(id string) error {
	err := p.db.QueryRow("delete from albums where id=$1", id)
	if err != nil {
		return errors.New("not_found")
	}
	return nil
}
func (p PostgresStorage) Read() []album {
	// var album album
	// err := p.db.QueryRow("select * from albums where id = $1", id).Scan(&album)
	// if err != nil {
	// 	return []album{}, errors.New("not_found")
	// }
	// return album, nil
	var a album
	return []album{a}
}

func NewStorage() Storage {
	return NewPostgresStorage()
}

// var albums = []album{
// 	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
// 	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
// 	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// }
