package main

import (
	"fmt"

	"github.com/taaag51/smart-pantry/backend-api/db"
	"github.com/taaag51/smart-pantry/backend-api/model"
)

func main() {
	dbConn := db.NewDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.User{}, &model.FoodItem{})
}
