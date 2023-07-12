package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4"
	"uread/router"
)

func main() {
	conn, err := pgx.Connect(context.Background(), "postgresql://uread:uread@localhost:5432/uread")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	router.SetupRouter(conn)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to server: %v\n", err)
		os.Exit(1)
	}
}
