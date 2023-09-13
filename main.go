package main

import (
	"assignment-2/config"
	"assignment-2/handler"
	"assignment-2/models"
	"encoding/gob"
	"fmt"
	"net/http"
)

func main() {
	gob.Register(models.User{})

	config.StartDB()

	db := config.GetDB()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	http.HandleFunc("/", handler.RouteIndex)
	http.HandleFunc("/register", handler.RouteRegister)
	http.HandleFunc("/logout", handler.RouteLogout)

	fmt.Println("server started at localhost:5555")
	http.ListenAndServe(":5555", nil)
}
