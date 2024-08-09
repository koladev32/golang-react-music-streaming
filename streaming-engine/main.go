package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

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

// Extracts the song ID from the request URL
func getSongID(r *http.Request) (int, error) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	return id, err
}

// Retrieves the song details from the database
func getSongFromDB(id int) (Song, error) {
	var song Song
	err := db.First(&song, id).Error
	return song, err
}

// Fetches the file from the URL
func fetchFile(fileURL string) (*http.Response, error) {
	fullURL := "http://localhost:8000/media/" + fileURL
	resp, err := http.Get(fullURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("file not found on the server")
	}
	return resp, nil
}

// Parses the Range header to get the start and end bytes
func parseRangeHeader(rangeHeader string, fileSize int64) (int64, int64, error) {
	bytesRange := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
	start, err := strconv.ParseInt(bytesRange[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	var end int64
	if len(bytesRange) > 1 && bytesRange[1] != "" {
		end, err = strconv.ParseInt(bytesRange[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}
	} else {
		end = fileSize - 1
	}

	if start > end || end >= fileSize {
		return 0, 0, fmt.Errorf("invalid range")
	}

	return start, end, nil
}

// Writes the partial content to the response
func writePartialContent(w http.ResponseWriter, start, end, fileSize int64, resp *http.Response) error {
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
	w.Header().Set("Content-Type", "audio/mpeg")
	w.WriteHeader(http.StatusPartialContent)

	// Create a channel for the buffered data and a wait group for synchronization
	dataChan := make(chan []byte)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
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
			dataChan <- buffer[:n]
			bytesToRead -= int64(n)
		}
		close(dataChan)
	}()

	go func() {
		defer wg.Wait()
		for chunk := range dataChan {
			if _, err := w.Write(chunk); err != nil {
				http.Error(w, "Error writing response", http.StatusInternalServerError)
				return
			}
		}
	}()

	// Skip the bytes until the start position
	io.CopyN(io.Discard, resp.Body, start)

	return nil
}

// Handles streaming of the file via HTTP range requests
func streamHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getSongID(r)
	if err != nil {
		http.Error(w, "Invalid song ID", http.StatusBadRequest)
		return
	}

	song, err := getSongFromDB(id)
	if err != nil {
		http.Error(w, "Song not found", http.StatusNotFound)
		return
	}

	resp, err := fetchFile(song.File)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer resp.Body.Close()

	fileSize := resp.ContentLength

	rangeHeader := r.Header.Get("Range")
	if rangeHeader == "" {
		http.ServeFile(w, r, song.File)
		return
	}

	start, end, err := parseRangeHeader(rangeHeader, fileSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := writePartialContent(w, start, end, fileSize, resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/songs/listen/{id}", streamHandler).Methods("GET")

	log.Println("Server is running on port 8005")
	log.Fatal(http.ListenAndServe(":8005", r))
}
