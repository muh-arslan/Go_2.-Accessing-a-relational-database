package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
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

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	album, err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", album)

	// albumID, err := addAlbum(Album{
	// 	Title:  "Rich Dad, Poor Dad",
	// 	Artist: "Robert T. Kiyosaki",
	// 	Price:  8.99,
	// })
	// albumID, err := addAlbum(Album{
	// 	Title:  "Start With Why",
	// 	Artist: "Simon Sinek",
	// 	Price:  18.00,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Id of added album: %v\n", albumID)

	allAlbums, err := showAllAlbums()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", allAlbums)
}

func albumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

func albumByID(id int64) (Album, error) {
	var album Album

	row := db.QueryRow("SELECT * FROM album WHERE id=?", id)

	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return album, fmt.Errorf("albumByID %d: no such album", id)
		}
		return album, fmt.Errorf("albumByID %d: %v", id, err)
	}
	return album, nil
}

func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}

func showAllAlbums() ([]Album, error) {
	var albums []Album
	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("showAllAlbums: %v - No Records Found", err)
	}
	defer rows.Close()

	for rows.Next() {
		var albs Album
		if err := rows.Scan(&albs.ID, &albs.Title, &albs.Artist, &albs.Price); err != nil {
			return nil, fmt.Errorf("showAllAlbums: %v", err)
		}
		albums = append(albums, albs)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("showAllAlbums: %v", err)
	}
	return albums, nil
}
