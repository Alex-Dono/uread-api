package router

import (
	"github.com/jackc/pgx/v4"
	"net/http"
	"uread/controller"
)

var mainRoute = "/api"

func SetupRouter(conn *pgx.Conn) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("api/books", controller.GetAllBooks(conn))
	mux.HandleFunc(mainRoute+"/books/{id}", controller.GetBook(conn))
	mux.HandleFunc(mainRoute+"/upload", controller.UploadBook(conn))
	return mux
}
