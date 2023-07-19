package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/rs/cors"
	"uread/router"
)

func main() {
	conn, err := pgx.Connect(context.Background(), "postgresql://uread:uread@localhost:5432/uread")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	mux := router.SetupRouter(conn)

	handler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:8081"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"*"},
	}).Handler(mux)

	fmt.Printf("Connecting on port 8080")
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to server: %v\n", err)
		os.Exit(1)
	}
}
