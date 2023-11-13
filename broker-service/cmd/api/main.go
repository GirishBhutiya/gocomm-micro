package main

import (
	"broker-service/token"
	"broker-service/util"
	"fmt"
	"log"
	"net/http"
)

const (
	webPort      = "80"
	authGrpcPort = "50001"
)

type Server struct {
	config     util.Config
	tokenMaker token.Maker
}

func main() {

	config, err := util.LoadConfig("./app")
	if err != nil {
		log.Fatal("can not load config:", err)
	}

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("can not create toekn maker: %v", err)
	}

	app := Server{
		config:     config,
		tokenMaker: tokenMaker,
	}

	log.Printf("Starting Broker service on port %s\n", webPort)

	//define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	//start the server
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
