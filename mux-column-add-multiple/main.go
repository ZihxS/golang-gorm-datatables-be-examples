package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	datatables "github.com/ZihxS/golang-gorm-datatables" // [👈🏼 FOCUS HERE]
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
	dsn := "..." // [👈🏼 ADJUST HERE]
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

	r.HandleFunc("/column-add-multiple", func(w http.ResponseWriter, r *http.Request) {
		req, err := datatables.ParseRequest(r) // [👈🏼 FOCUS HERE]
		if err != nil {
			http.Error(w, fmt.Sprintf("Error processing request: %v", err), http.StatusInternalServerError)
			return
		}

		// [👇🏼 FOCUS HERE]
		generateInitialsName := func(fullName string) string {
			var initials []string

			parts := strings.Fields(fullName)
			for _, part := range parts {
				if len(part) > 0 {
					initials = append(initials, strings.ToUpper(part[:1]))
				}
			}

			return strings.Join(initials, "")
		}
		// [👆🏻 FOCUS HERE]

		// [👇🏼 FOCUS HERE]
		generateInitialsNameColumn := datatables.Column{Data: "initials_name", RenderFunc: func(m map[string]any) any {
			return generateInitialsName(m["name"].(string))
		}}
		// [👆🏻 FOCUS HERE]

		// [👇🏼 FOCUS HERE]
		actionColumn := datatables.Column{Data: "action", RenderFunc: func(m map[string]any) any {
			return `<a href="https://alwaysngoding.com" class="btn btn-primary btn-sm">EXAMPLE ACTION</a>`
		}}
		// [👆🏻 FOCUS HERE]

		tx := db.Model(&User{})
		response, err := datatables.New(tx).Req(*req).AddColumns(actionColumn, generateInitialsNameColumn).Make() // [👈🏼 FOCUS HERE]
		if err != nil {
			http.Error(w, fmt.Sprintf("Error processing datatables: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	})

	log.Fatal(http.ListenAndServe(":6969", r))
}
