package main

import (
	"spotsync/internal/config"
	"spotsync/internal/database"
	"spotsync/internal/server"
)

func main() {
	cnfg := config.LoadEnv()

	db := database.ConnectDB(cnfg)

	server.Start(db, cnfg)
}