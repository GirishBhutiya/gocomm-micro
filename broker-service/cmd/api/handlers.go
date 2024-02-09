package main

import (
	"broker-service/data"
	auth "broker-service/proto"
	"broker-service/util"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type ResetPasswordPayload struct {
	Email    string `json:"email"`
	Password string `json:"new_password"`
}
type UserResponse struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	FirstName         string    `json:"first_name,omitempty"`
	LastName          string    `json:"last_name,omitempty"`
	RollId            int       `json:"roll_id"`
	Active            int       `json:"active"`
	Roll              string    `json:"roll"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
type LoginResponse struct {
	UserResponse  UserResponse `json:"user"`
	AccessToken   string       `json:"access_token"`
	Authenticated bool         `json:"authenticated"`
	ErrorMessage  string       `json:"message"`
}

func (app *Server) Brocker(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "Hit the brocker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}
func (app *Server) Login(w http.ResponseWriter, r *http.Request) {
	var authPayload AuthPayload

	err := app.readJSON(w, r, &authPayload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	if !util.CheckPasswordValidity(authPayload.Password) {
		log.Println("password error")
		app.errorJSON(w, errors.New("password not valid. please set password more than 8 character without spaces"))
		return
	}
	conn, err := grpc.Dial(fmt.Sprintf("authentication-service:%s", authGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	defer conn.Close()
	c := auth.NewLoginServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.Login(ctx, &auth.LoginRequest{
		LoginReq: &auth.LoginData{
			Username: authPayload.Email,
			Password: authPayload.Password,
		},
	})

	if err != nil {
		log.Println(err)
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}
	userRes := CreateUserResponse(res.User)

	accessToken, err := app.tokenMaker.CreateToken(res.User.Email, app.config.AccessTokenDuration)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	response := LoginResponse{
		UserResponse:  userRes,
		AccessToken:   accessToken,
		Authenticated: true,
		ErrorMessage:  "",
	}
	app.writeJSON(w, http.StatusAccepted, response)

}
func (app *Server) Register(w http.ResponseWriter, r *http.Request) {
	var user data.User

	err := app.readJSON(w, r, &user)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	if util.CheckPasswordValidity(user.Password) {
		log.Println("password error")
		app.errorJSON(w, errors.New("password not valid. please set password more than 8 character without spaces"))
		return
	}
	//log.Println("Brocker 1 pass :", user.Password)
	conn, err := grpc.Dial(fmt.Sprintf("authentication-service:%s", authGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	defer conn.Close()

	c := auth.NewLoginServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.Register(ctx, &auth.RegisterRequest{
		User: &auth.User{
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Password:  user.Password,
			RollId:    int64(user.RollId),
			Active:    int32(user.Active),
		},
	})

	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted, res)

}
func (app *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var user data.User

	err := app.readJSON(w, r, &user)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	conn, err := grpc.Dial(fmt.Sprintf("authentication-service:%s", authGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	defer conn.Close()

	c := auth.NewLoginServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	res, err := c.UpdateUser(ctx, &auth.UpdateUserRequest{
		User: &auth.User{
			ID:        int64(user.ID),
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			RollId:    int64(user.RollId),
			Active:    int32(user.Active),
		},
	})
	log.Println(res)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted, res)

}
func (app *Server) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	theURL := r.RequestURI
	log.Println(theURL)

	conn, err := grpc.Dial(fmt.Sprintf("authentication-service:%s", authGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	defer conn.Close()

	c := auth.NewLoginServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	res, err := c.VerifyRegisteredEmail(ctx, &auth.VerifyRegisteredEmailRequest{
		Link: theURL,
	})
	log.Println(res)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	app.writeJSON(w, http.StatusAccepted, res)

}
func (app *Server) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	//log.Println("Email is:", email)

	conn, err := grpc.Dial(fmt.Sprintf("authentication-service:%s", authGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println("dial authentication service")
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	defer conn.Close()

	c := auth.NewLoginServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	res, err := c.ForgotPassword(ctx, &auth.ForgotPasswordRequest{
		Email: email,
	})
	log.Println(res)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	app.writeJSON(w, http.StatusAccepted, res)

}
func (app *Server) ResetPassword(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	//log.Println("Email is:", email)
	theURL := r.RequestURI
	//log.Println(theURL)

	var resetPasswordPayload ResetPasswordPayload

	err := app.readJSON(w, r, &resetPasswordPayload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	if util.CheckPasswordValidity(resetPasswordPayload.Password) {
		log.Println("password error")
		app.errorJSON(w, errors.New("password not valid. please set password more than 8 character without spaces"))
		return
	}

	conn, err := grpc.Dial(fmt.Sprintf("authentication-service:%s", authGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {

		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	defer conn.Close()

	c := auth.NewLoginServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	res, err := c.ResetPassword(ctx, &auth.ResetPasswordRequest{
		Email:    email,
		Password: resetPasswordPayload.Password,
		Link:     theURL,
	})
	log.Println(res)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	app.writeJSON(w, http.StatusAccepted, res)

}
