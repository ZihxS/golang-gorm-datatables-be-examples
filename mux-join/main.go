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

	r.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		req, err := datatables.ParseRequest(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error processing request: %v", err), http.StatusInternalServerError)
			return
		}

		tx := db.
			Select("SELECT o.product_name, o.price, o.quantity, (o.quantity*o.price) total_price, u.name buyer, p.details buyer_profile_detail, t.tag_name").
			Table("orders o").
			Joins("INNER JOIN users u on u.id = o.user_id").
			Joins("INNER JOIN profiles p on p.user_id = u.id").
			Joins("LEFT JOIN user_tags ut on ut.user_id = u.id").
			Joins("LEFT JOIN tags t on t.id = ut.tag_id")

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
