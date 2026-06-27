package main

import (
	"booking-app/internal/config"
	"booking-app/internal/database"
	"booking-app/internal/server"
)

func main() {
	env := config.LoadEnv()
	
	db := database.ConnectDB(env)
	
	server.Start(db, env)
}	