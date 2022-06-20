package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type HttpError struct {
	Error string `json:"error"`
}
type Storage interface {
	Create(album) album
	Read() []album
	ReadOne(string) (album, error)
	Update(string, album) (album, error)
	Delete(string) error
}

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type MemeoryStorage struct {
	albums []album
}

type PostgresStorage struct {
	db *sql.DB
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

func NewPostgresStorage() PostgresStorage {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return PostgresStorage{db: db}
}

func NewStorage() Storage {
	return NewMemoryStorage()
}

// var albums = []album{
// 	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
// 	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
// 	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// }

var storage = NewMemoryStorage()

func getAlbums(c *gin.Context) {
	// storage := NewMemeoryStorage()
	c.IndentedJSON(http.StatusOK, storage.Read())
}
func postAlbum(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
	}
	// storage := NewMemeoryStorage()
	storage.Create(newAlbum)
	// albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	// storage := NewMemeoryStorage()
	album, err := storage.ReadOne(id)
	if err != nil {
		c.IndentedJSON(http.StatusOK, HttpError{"not_found"})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}
func deleteAlbumByID(c *gin.Context) {
	id := c.Param("id")
	err := storage.Delete(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	c.IndentedJSON(http.StatusNoContent, album{})

}
func updateAlbumByID(c *gin.Context) {
	id := c.Param("id")
	var newAlbum album
	c.BindJSON(&newAlbum)
	album, err := storage.Update(id, newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}
	c.IndentedJSON(http.StatusNoContent, album)
}
func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbum)
	router.GET("/albums/:id", getAlbumByID)
	router.DELETE("/albums/:id", deleteAlbumByID)
	router.PUT("/albums/:id", updateAlbumByID)

	router.Run("localhost:8080")
}
