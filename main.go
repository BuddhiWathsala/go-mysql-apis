package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	http.HandleFunc("/", productResource)
	http.ListenAndServe(":3001", nil)
}

func productResource(w http.ResponseWriter, r *http.Request) {
	userName := os.Getenv("MYSQL_GENERIC_USERNAME")
	password := os.Getenv("MYSQL_GENERIC_PASSWORD")
	dbName := os.Getenv("MYSQL_AD_DB_NAME")
	dbURL := userName + ":" + password + "@/" + dbName + "?charset=utf8&parseTime=True&loc=Local"

	db, _ := gorm.Open("mysql", dbURL)
	defer db.Close()

	var product Product
	product = readProduct(db)
	fmt.Println(product.Code)

	js, err := json.Marshal(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func migrate(db *gorm.DB) {
	db.AutoMigrate(&Product{})
	db.Create(&Product{Code: "L1212", Price: 1000})
}

func readProduct(db *gorm.DB) Product {
	var product Product
	db.First(&product, 1) // find product with id 1
	fmt.Println(product.Price)
	db.First(&product, "code = ?", "L1212") // find product with code l1212

	return product
}
