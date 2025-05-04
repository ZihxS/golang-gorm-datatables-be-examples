package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	datatables "github.com/ZihxS/golang-gorm-datatables"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type User struct {
	ID   int
	Name string
	Age  int
}

func main() {
	dsn := "..." // [👈🏼 ADJUST]
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/regex", func(w http.ResponseWriter, r *http.Request) {
		req, err := datatables.ParseRequest(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error processing request: %v", err), http.StatusInternalServerError)
			return
		}

		tx := db.Model(&User{})
		response, err := datatables.New(tx).Req(*req).WithNumber().Make()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error processing datatables: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	})

	log.Fatal(http.ListenAndServe(":6969", r))
}
