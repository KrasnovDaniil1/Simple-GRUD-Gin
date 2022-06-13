package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// type error struct {
// 	Error string `json:"error"`
// }

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}
func postAlbum(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "bad request"})
	}
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	for _, v := range albums {
		if v.ID == id {
			c.IndentedJSON(http.StatusOK, v)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
func deleteAlbumByID(c *gin.Context) {
	id := c.Param("id")
	for i, v := range albums {
		if v.ID == id {
			albums = append(albums[:i], albums[i+1:]...)
			c.IndentedJSON(http.StatusNoContent, v)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}
func updateAlbumByID(c *gin.Context) {
	id := c.Param("id")
	for i, v := range albums {
		if v.ID == id {
			c.BindJSON(&albums[i])
			albums[i] = v
			c.IndentedJSON(http.StatusNoContent, albums[i])
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
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
