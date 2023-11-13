package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	pb "github.com/GirishBhutiya/gocomm-micro/backend/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) RequestCome(w http.ResponseWriter, r *http.Request) {
	log.Println("Request Come")
	http.ServeFile(w, r, "./cmd/web/templates/test.page.gohtml")

	// dir, file := filepath.Split("./templates/test.html")

	// log.Println("Directory is:", dir)
	// log.Println("file is:", file)

	// fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	// //http.ServeFile(w, r, "./templates/test.html")
	// t, err := template.ParseFiles("./cmd/web/test.html")
	// if err != nil {
	// 	log.Panic(err)
	// }
	// t.Execute(w, nil)
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))

	if err != nil {

		return err
	}

	_, err = app.Client.Do(request)

	if err != nil {
		return err
	}

	return nil

}
func (app *Config) AuthenticateViagRPC(w http.ResponseWriter, r *http.Request) {
	log.Println("AuthenticateViagRPC from backend")
	var requestPayload AuthPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	conn, err := grpc.Dial("localhost:8081/authenticate", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()
	_, err = c.UserAuthenticate(ctx, &pb.UserRequest{
		User: &pb.User{
			Email:    requestPayload.Email,
			Password: requestPayload.Password,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Logged"

	app.writeJSON(w, http.StatusAccepted, payload)
}
