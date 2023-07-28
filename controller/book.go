package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"uread/model"
)

func GetAllBooks(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var books []model.Book
		result := db.Find(&books)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		err := json.NewEncoder(w).Encode(books)
		if err != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetBook(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the book's ID from the URL parameters
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "ID is required", http.StatusBadRequest)
			return
		}

		// Get the book from the database
		var book model.Book
		err := db.
			Where("id = ?", id).
			First(&book).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Book not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Send the book file to the client
		http.ServeFile(w, r, book.FilePath)
	}
}

func UploadBook(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the multipart form in the request
		err := r.ParseMultipartForm(10 << 20) // limit your maxMultipartMemory
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Retrieve the file from form data
		file, handler, err := r.FormFile("file") // retrieve the file from form data
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Create the file in the server's file system
		dst, err := os.Create(filepath.Join("./uploads", handler.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create a new book in the database with the file's path
		book := model.Book{
			ID:       uuid.New(),
			Title:    r.FormValue("title"),
			Author:   r.FormValue("author"),
			FilePath: filepath.Join("./uploads", handler.Filename),
		}
		result := db.Create(&book)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with a success message
		fmt.Fprintln(w, "Successfully uploaded file")
	}
}
