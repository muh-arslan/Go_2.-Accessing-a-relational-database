package dataaccess

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID     int64   `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float32 `json:"price"`
}

func DBConnection() {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		// User:   "root",
		// Passwd: "127586",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "recordings",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		fmt.Println("Error!")
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println("Ping Error!")
		log.Fatal(pingErr)
	}

	fmt.Println("Connected")

	// ShowAllProcesses()
}

// func ShowAllProcesses() {
// 	fetchAlbumByArtist()
// 	fetchAlbumByID()
// 	addAndFetchAlbum()
// 	fetchAllAlbums()
// 	deleteAndFetchAlbumByID()
// 	fetchAllAlbums()
// 	deleteAndFetchAllAlbums()
// 	fetchAllAlbums()
// }

// func fetchAlbumByArtist() {
// 	albums, err := albumsByArtist("John Coltrane")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Albums found: %v\n", albums)
// }

// func fetchAlbumByID() {
// 	album, err := albumByID(2)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Album found: %v\n", album)

// }

// func addAndFetchAlbum() {
// 	albumID, err := addAlbum(Album{
// 		Title:  "Rich Dad, Poor Dad",
// 		Artist: "Robert T. Kiyosaki",
// 		Price:  8.99,
// 	})
// 	// albumID, err := addAlbum(Album{
// 	// 	Title:  "Start With Why",
// 	// 	Artist: "Simon Sinek",
// 	// 	Price:  18.00,
// 	// })
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Id of added album: %v\n", albumID)
// }

// func fetchAllAlbums() {
// 	allAlbums, err := showAllAlbums()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Albums found: %v\n", allAlbums)
// }

// func deleteAndFetchAlbumByID() {
// 	id := int64(4)
// 	res, err := deleteAlbumByID(id)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if !res {
// 		log.Fatalf("Record with id : %v not found\n", id)
// 	}
// 	fmt.Printf("Record with id : %v deleted\n", id)
// }

// func deleteAndFetchAllAlbums() {
// 	result, err := deleteAllAlbums()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if !result {
// 		log.Fatal("Not Record found to Delete")
// 	}
// 	fmt.Println("All Record are deleted")
// }

func AlbumsByArtist(c *gin.Context) {
	name := c.Param("name")
	name = strings.Replace(name, "+", " ", 5)

	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}
		albums = append(albums, alb)
	}

	if albums == nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": "no record found"})
		return
	}
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func AlbumByID(c *gin.Context) {
	id := c.Param("id")
	var album Album

	row := db.QueryRow("SELECT * FROM album WHERE id=?", id)

	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
		return
	}
	c.IndentedJSON(http.StatusOK, album)
}

func AddAlbum(c *gin.Context) {
	var alb Album

	if err := c.BindJSON(&alb); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid syntax"})
		return
	}

	if (alb.Title == "") || (alb.Artist == "") || (alb.Price == 0) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "pass all the parameters"})
		return
	}

	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	id, er := result.LastInsertId()

	if er != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	alb.ID = id
	c.IndentedJSON(http.StatusCreated, alb)
}

func ShowAllAlbums(c *gin.Context) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No Record Found"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var albs Album
		if err := rows.Scan(&albs.ID, &albs.Title, &albs.Artist, &albs.Price); err != nil {
			c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
			return
		}
		albums = append(albums, albs)
	}
	if err := rows.Err(); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	if albums == nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "No Record Found"})
		return
	}
	c.IndentedJSON(http.StatusOK, albums)
}

func UpdateAlbumByID(c *gin.Context) {
	var alb Album

	if err := c.BindJSON(&alb); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invlalid request"})
		return
	}

	id := int64(alb.ID)
	if id == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id is required"})
		return
	}

	query := "UPDATE album SET "
	if alb.Title != "" {
		query += fmt.Sprintf("title = '%s', ", alb.Title)
	}
	if alb.Artist != "" {
		query += fmt.Sprintf("artist = '%s', ", alb.Artist)
	}
	if alb.Price != 0 {
		query += fmt.Sprintf("price = %f, ", alb.Price)
	}
	query = strings.TrimSuffix(query, ", ")
	query += " WHERE id = ?"

	result, err := db.Exec(query, id)
	if err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	rowsAffected, er := result.RowsAffected()
	if er != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": er.Error()})
		return
	}

	if rowsAffected == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "no album found with given id"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "album updated successfully"})
}

func DeleteAllAlbums(c *gin.Context) {
	result, err := db.Exec("DELETE FROM album")
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	affectedRows, err := result.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusNotModified, gin.H{"message": err.Error()})
		return
	}
	resp := fmt.Sprintf("%v rows deleted", affectedRows)
	c.IndentedJSON(http.StatusOK, gin.H{"message": resp})
}

func DeleteAlbumByID(c *gin.Context) {
	id := c.Param("id")
	result, err := db.Exec("DELETE FROM album WHERE id=?", id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	affectedRow, err := result.RowsAffected()
	if err != nil {
		c.IndentedJSON(http.StatusNotModified, gin.H{"message": err.Error()})
		return
	}
	if affectedRow <= 0 {
		c.IndentedJSON(http.StatusNotModified, gin.H{"message": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "success"})
}
