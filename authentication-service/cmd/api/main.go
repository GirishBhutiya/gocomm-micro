package main

import (
	"authentication/cmd/db"
	"authentication/data/models"
	"authentication/internal/util"
	"fmt"
	"log"
	"net/http"
)

const (
	webPort        = "80"
	gRpcPort       = "50001"
	mailerGrpcPort = "50002"
)

type Config struct {
	//DB *sql.DB
	//Models  data.Models
	EnvVars     util.ConfigVars
	userService models.UserService
}

func main() {
	config, err := util.LoadConfig("./app")
	if err != nil {
		log.Fatal("can not load config:", err)
	}
	log.Println("Starting authentication service")
	//connect to db
	conn := db.ConnectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}
	app := Config{
		//DB: conn,
		//Models:  data.New(conn),
		userService: models.UserServiceStruct{DB: conn},
		EnvVars:     config,
	}
	log.Println("Authentication service Started")

	go app.gRPCListen()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	//start the server

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
		//log.Panic(err)
	}
	defer conn.Close()
	defer srv.Close()
}
