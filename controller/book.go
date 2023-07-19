package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"uread/model"
)

func GetAllBooks(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Query the database for all books
		rows, err := conn.Query(context.Background(), "SELECT * FROM books")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Iterate through the result set and append each book to a slice
		var books []model.Book
		for rows.Next() {
			var book model.Book
			err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.FilePath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			books = append(books, book)
		}

		// Check for errors from iterating over rows.
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert the slice to JSON
		booksJson, err := json.Marshal(books)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the JSON to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(booksJson)
	}
}

func GetBook(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the book's ID from the URL parameters
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get the book from the database
		var book model.Book
		err = conn.QueryRow(context.Background(), "SELECT * FROM books WHERE id = $1", id).Scan(&book.ID, &book.Title, &book.Author, &book.FilePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the book file to the client
		http.ServeFile(w, r, book.FilePath)
	}
}

func UploadBook(conn *pgx.Conn) http.HandlerFunc {
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
		_, err = conn.Exec(context.Background(), "INSERT INTO books (id, title, author, file_path) VALUES ($1, $2, $3, $4)",
			uuid.New(),
			r.FormValue("title"), r.FormValue("author"), filepath.Join("./uploads", handler.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Respond with a success message
		fmt.Fprintln(w, "Successfully uploaded file")
	}
}
