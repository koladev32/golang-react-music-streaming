package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

// Song represents the song model in the existing music_song table
type Song struct {
	ID        uint   `gorm:"column:id"`
	Name      string `gorm:"column:name"`
	File      string `gorm:"column:file"`
	Author    string `gorm:"column:author"`
	Thumbnail string `gorm:"column:thumbnail"`
}

// TableName overrides the table name used by Gorm
func (Song) TableName() string {
	return "music_song"
}

func initDB() {
	// Initialize SQLite connection
	db, err = gorm.Open(sqlite.Open("../backend/db.sqlite3"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

// function to handle the streaming of the file via HTTP range requests.
// First, we retrieve song data from the db,
func streamHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	var song Song
	if err := db.First(&song, id).Error; err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	fileURL := "https://res.cloudinary.com/kolawole31/video/upload/v1722638406/uyqfcqlect4y80lvjylq.mp3" // Construct the full file URL

	resp, err := http.Get(fileURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "File not found on the server", http.StatusNotFound)
		return
	}
	defer resp.Body.Close()

	fileSize := resp.ContentLength

	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		http.ServeFile(w, r, fileURL)
		return
	}

	// Extract the byte range from the request
	bytesRange := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
	start, err := strconv.ParseInt(bytesRange[0], 10, 64)
	if err != nil {
		http.Error(w, "Invalid range", http.StatusBadRequest)
		return
	}

	var end int64
	if len(bytesRange) > 1 && bytesRange[1] != "" {
		end, err = strconv.ParseInt(bytesRange[1], 10, 64)
		if err != nil {
			http.Error(w, "Invalid range", http.StatusBadRequest)
			return
		}
	} else {
		end = fileSize - 1
	}

	if start > end || end >= fileSize {
		http.Error(w, "Invalid range", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// Set headers for partial content
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
	w.Header().Set("Content-Type", "audio/mpeg")
	w.WriteHeader(http.StatusPartialContent)

	// Skip the bytes until the start position
	io.CopyN(io.Discard, resp.Body, start)

	// Read and write the specified byte range
	buffer := make([]byte, 1024) // 1KB buffer size
	bytesToRead := end - start + 1
	for bytesToRead > 0 {
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		if n == 0 {
			break
		}
		if int64(n) > bytesToRead {
			n = int(bytesToRead)
		}
		if _, err := w.Write(buffer[:n]); err != nil {
			http.Error(w, "Error writing response", http.StatusInternalServerError)
			return
		}
		bytesToRead -= int64(n)
	}
}

func main() {
	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/songs/listen/{id}", streamHandler).Methods("GET")

	log.Println("Server is running on port 8005")
	log.Fatal(http.ListenAndServe(":8005", r))
}
