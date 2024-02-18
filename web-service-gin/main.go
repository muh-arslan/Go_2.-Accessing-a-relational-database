package main

import (
	dataaccess "example/data-access"

	"github.com/gin-gonic/gin"
)

func main() {
	dataaccess.DBConnection()

	router := gin.Default()

	router.GET("/albums", dataaccess.ShowAllAlbums)
	router.GET("/albums/:id", dataaccess.AlbumByID)
	router.GET("/albums/filterByArtist/:name", dataaccess.AlbumsByArtist)
	router.POST("/albums", dataaccess.AddAlbum)
	router.PATCH("/albums", dataaccess.UpdateAlbumByID)
	router.DELETE("/albums", dataaccess.DeleteAllAlbums)
	router.DELETE("/albums/:id", dataaccess.DeleteAlbumByID)

	router.Run("localhost:8080")
}
