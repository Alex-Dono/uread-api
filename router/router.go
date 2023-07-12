package router

import (
	"github.com/jackc/pgx/v4"
	"net/http"
	"uread/controller"
)

var mainRoute = "/api"

func SetupRouter(conn *pgx.Conn) {
	http.HandleFunc(mainRoute+"/books", controller.GetBook(conn))
	http.HandleFunc(mainRoute+"/books/{id}", controller.GetBook(conn))
	http.HandleFunc(mainRoute+"/upload", controller.UploadBook(conn))
}
