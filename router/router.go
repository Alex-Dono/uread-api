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
	mux.HandleFunc(mainRoute+"/book", controller.GetBook(db)) // complete endpoint: /book/id-of-the-idea
	mux.HandleFunc(mainRoute+"/upload", controller.UploadBook(db))
	return mux
}
