package main

import (
	"task-5-vix-btpns-HaiqalRamanizarAlFajri/database"
	"task-5-vix-btpns-HaiqalRamanizarAlFajri/router"
)

func main() {
	db := database.SetupDB()

	r := router.SetupRoutes(db)
	r.Run()
}
