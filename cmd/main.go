package main

import (
	"fmt"
	"os"

	"github.com/pratikjagrut/myreads-backend/pkg/config"
)

func main() {
	db_host := os.Getenv("DB_HOST")
	db_port := os.Getenv("DB_PORT")
	db_user := os.Getenv("DB_USER")
	db_pass := os.Getenv("DB_PASS")
	db_name := os.Getenv("DB_NAME")

	db, err := config.NewConnection(db_host, db_port, db_name, db_user, db_pass)
	if err != nil {
		panic(err)
	}
	_ = db
	fmt.Println("connection done!")
}
