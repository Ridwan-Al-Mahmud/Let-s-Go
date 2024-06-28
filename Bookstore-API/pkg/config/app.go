package config

import(
  "fmt"
  "github.com/jinzhu/gorm"
  _ "github.com/mattn/go-sqlite3"
)

//github.com/jinzhu/gorm/dialects/mysql

var (
  db * gorm.DB
)

func Connect (){
  d ,err := gorm.Open("sqlite3", "github.com/Ridwan-Al-Mahmud/Bookstore-API/pkg/database/bookstore.db")
  if err != nil{
    fmt.Println("Error: ", err)
  }
  db = d
}

func GetDB() *gorm.DB {
  return db
}