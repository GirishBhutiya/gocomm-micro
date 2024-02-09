package main

import (
	"authentication/cmd/db"
	"authentication/data/models"
	"authentication/data/proto"
	"authentication/internal/util"
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var usr models.User
var us models.UserService
var config util.ConfigVars

func setup() models.UserService {
	var err error
	config, err = util.LoadConfig("./../../")
	if err != nil {
		log.Fatal("can not load config:", err)
	}

	db := db.ConnectToTestDB()

	return &models.UserServiceStruct{
		DB: db,
	}

}

func init() {
	us = setup()
	us.ResetUsers()
	usr = CreateRandomUser()

}

func TestLoginServer_Register(t *testing.T) {

	client := getServer(t)

	tests := []struct {
		name    string
		args    models.User
		want    models.User
		isErr   bool
		wantErr string
	}{
		{name: "register with correct email",
			args:    usr,
			want:    usr,
			isErr:   false,
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := us.ConvertToProtoUser(tt.args)

			res, err := client.Register(context.Background(), &proto.RegisterRequest{
				User: &user,
			})
			if err != nil {
				t.Fatalf("grpc.Register() %v", err)
			}
			if res.IsError {
				t.Fatalf("grpc.Register %v", err)
			}
			if !CompareUser(&usr, res.User) {
				t.Fatalf("added user and returned user are not same.")
			}
			if tt.name == "register with correct email" {
				usr.ID = int(res.User.ID)
			}
		})
	}

}
func TestLoginServer_VerifyRegisteredEmail(t *testing.T) {

	client := getServer(t)

	tests := []struct {
		name     string
		link     string
		want     models.User
		verified bool
		isErr    bool
		wantErr  string
	}{
		{name: "verify correct user",
			link:     util.GetFullVerifyEmailLink(usr.Email, config.FrontEndDomain, config.HashSecretKeyVerifyEmail),
			want:     usr,
			verified: true,
			isErr:    false,
			wantErr:  "",
		},
		{name: "wrong hash",
			link:     util.GetFullVerifyEmailLink(usr.Email, config.FrontEndDomain, ""),
			want:     usr,
			verified: false,
			isErr:    false,
			wantErr:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := client.VerifyRegisteredEmail(context.Background(), &proto.VerifyRegisteredEmailRequest{
				Link: tt.link,
			})
			if err != nil {
				t.Fatalf("grpc.VerifyRegisteredEmail() %v", err)
			}

			if tt.verified != res.Verified {
				t.Fatalf("User must be verified but not.")
			}
		})
	}

}
func TestLoginServer_Login(t *testing.T) {

	client := getServer(t)

	tests := []struct {
		name          string
		username      string
		password      string
		want          models.User
		authenticated bool
		isErr         bool
		wantErr       string
	}{
		{name: "login with correct password and user",
			username:      usr.Email,
			password:      usr.Password,
			want:          usr,
			authenticated: true,
			isErr:         false,
			wantErr:       "",
		},
		{name: "login with incorrect password",
			username:      usr.Email,
			password:      "ascvf123",
			want:          usr,
			authenticated: false,
			isErr:         true,
			wantErr:       "",
		},
		{name: "login with incorrect user",
			username:      "abc@email.com",
			password:      "ascvf123",
			want:          usr,
			authenticated: false,
			isErr:         true,
			wantErr:       "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := client.Login(context.Background(), &proto.LoginRequest{
				LoginReq: &proto.LoginData{
					Username: tt.username,
					Password: tt.password,
				},
			})
			if res.Authenticated != tt.authenticated {
				t.Fatalf("gRPC.Login() Authenticated must be %v but we get %v", tt.authenticated, res.Authenticated)
			}
			if err != nil && !tt.isErr {
				t.Fatalf("grpc.Login() %v", err)
			}

		})
	}

}

// RandomEmail generates a random email
func CreateRandomUser() models.User {
	return models.User{
		Email:     util.RandomEmail(),
		FirstName: util.RandStringBytesMaskImprSrcUnsafe(6),
		LastName:  util.RandStringBytesMaskImprSrcUnsafe(6),
		Password:  fmt.Sprintf("%s%s", util.RandStringBytesMaskImprSrcUnsafe(7), "123!"),
		RollId:    1,
	}
}
func getServer(t *testing.T) *LoginServer {
	// Server Initialization
	lis := bufconn.Listen(1024 * 1024)

	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := LoginServer{
		Config:      config,
		userService: us,
	}

	proto.RegisterLoginServiceServer(srv, &svc)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("srv.Serve %v", err)
		}
	}()

	// Test
	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(dialer), grpc.WithInsecure())

	t.Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		t.Fatalf("grpc.DialContex() %v", err)
	}
	return &LoginServer{userService: us, Config: config}
}
func CompareUser(got *models.User, want *proto.User) bool {
	if got.Email != want.Email {
		return false
	}
	if got.FirstName != want.FirstName {
		return false
	}
	if got.LastName != want.LastName {
		return false
	}
	if got.RollId != int(want.RollId) {
		return false
	}
	return true
}
