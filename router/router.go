package router

import (
	"gorm.io/gorm"
	"net/http"
	"uread/controller"
)

var mainRoute = "/api"

func SetupRouter(db *gorm.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(mainRoute+"/books", controller.GetAllBooks(db))
	mux.HandleFunc(mainRoute+"/books/{id}", controller.GetBook(db))
	mux.HandleFunc(mainRoute+"/upload", controller.UploadBook(db))
	return mux
}
