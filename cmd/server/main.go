package main

import (
	"drupal-scheduler/api"
	"drupal-scheduler/service"
	"drupal-scheduler/store"
	"log"
	"net/http"
	"os"
)

func main() {
	path := os.Getenv("DB_PATH")
	if path == "" {
		path = "scheduler.db"
	}
	db, e := store.Open(path)
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	svc := service.New(db)
	defer svc.Close()
	log.Fatal(http.ListenAndServe(":8080", api.New(svc).Handler()))
}
