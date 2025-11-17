package main

import (
	"log"
	"os"

	"github.com/mukailasam/akoko/internal/api"
	"github.com/mukailasam/akoko/internal/api/handler"
	"github.com/mukailasam/akoko/internal/api/routes"
	"github.com/mukailasam/akoko/internal/database"
	"github.com/mukailasam/akoko/internal/database/repository"
	"github.com/mukailasam/akoko/pkg/auth"
	"github.com/mukailasam/akoko/pkg/logger"
	_ "github.com/mukailasam/ayika"
	"github.com/mukailasam/igo"
)

func main() {
	logFile, err := os.OpenFile("logfile.json", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalln(err)
	}
	logger := logger.NewAkokoLogger(logFile)

	db, err := database.DBConnect("./storage.db")
	if err != nil {
		panic(err)
	}
	repo := repository.NewRepository(db)

	apiService := api.NewAPIService(repo, logger)
	hdler := handler.NewHandler(apiService, logger)

	app := igo.NewIgo()
	app.Use(igo.Logger)
	//app.Use(igo.Recovery)
	app.Use(auth.Auth)

	r := routes.NewRouter(app, hdler)
	r.RegisterRoutes()

	app.Run(":9000")
}
