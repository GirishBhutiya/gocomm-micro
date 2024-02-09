package main

import (
	"fmt"
	"log"
	"mailer-service/util"
	"net/http"
)

const (
	webPort  = "80"
	gRpcPort = "50002"
)

var counts int64

type Config struct {
	EnvVars util.Config
}

func main() {
	config, err := util.LoadConfig("./app")
	if err != nil {
		log.Fatal("can not load config:", err)
	}

	log.Println("Starting authentication service")

	app := Config{
		EnvVars: config,
	}

	log.Println("Authentication service Started")
	go app.gRPCListen()
	//go app.gRPCListen()

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
}
