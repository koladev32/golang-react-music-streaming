package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

// Song represents the song model
type Song struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Author    string `json:"author"`
	Thumbnail string `json:"thumbnail"`
}

func initDB() {
	// Initialize SQLite connection
	db, err = gorm.Open(sqlite.Open("songs.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	db.AutoMigrate(&Song{})
}

func getSongs(w http.ResponseWriter, r *http.Request) {
	var songs []Song
	db.Find(&songs)
	json.NewEncoder(w).Encode(songs)
}

func getSong(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	var song Song
	db.First(&song, id)
	if song.ID == 0 {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(song)
}

func main() {
	initDB()

	// Initialize router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/api/songs", getSongs).Methods("GET")
	r.HandleFunc("/api/songs/{id}", getSong).Methods("GET")

	// Start server
	log.Println("Server is running on port 8005")
	log.Fatal(http.ListenAndServe(":8005", r))
}
