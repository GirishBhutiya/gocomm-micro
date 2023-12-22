package main

import (
	"authentication/data"
	"authentication/proto"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LoginServer struct {
	proto.UnimplementedLoginServiceServer
	Models data.Models
}

func (l *LoginServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {

	input := req.GetLoginReq()

	var u data.User
	user, err := u.GetByEmail(input.GetUsername())
	if err != nil {
		log.Println(err)
		res := &proto.LoginResponse{
			Authenticated: false,
			User:          nil,
		}
		return res, err
	}

	match, err := user.PasswordMatches(input.GetPassword())
	if err != nil {
		log.Println(err)
		res := &proto.LoginResponse{
			Authenticated: false,
			User:          nil,
		}
		return res, err
	}

	if !match {
		res := &proto.LoginResponse{
			Authenticated: false,
			User:          nil,
		}
		return res, errors.New("password incorrect")
	}
	//check is email verified or not.
	if user.Active == 0 {
		res := &proto.LoginResponse{
			Authenticated: false,
			User:          nil,
		}
		return res, errors.New("your email is not verified, please check your mailbox including spam folder")
	}
	// return response
	usr := &proto.User{
		ID:                int64(user.ID),
		Email:             user.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		RollId:            int64(user.RollId),
		Roll:              user.Roll,
		Password:          user.Password,
		Active:            int32(user.Active),
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt),
		UpdatedAt:         timestamppb.New(user.UpdatedAt),
	}
	res := &proto.LoginResponse{
		User:          usr,
		Authenticated: true,
	}
	return res, nil

}

func (l *LoginServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {

	input := req.GetUser()

	u := data.User{
		Email:     input.GetEmail(),
		FirstName: input.GetFirstName(),
		LastName:  input.GetLastName(),
		Password:  input.GetPassword(),
		RollId:    int(input.GetRollId()),
		//Active:    int(input.GetActive()),
	}
	//log.Println("register input password:", input.GetPassword())
	userId, err := u.Insert(u)
	if err != nil {
		log.Println(err)
		res := &proto.RegisterResponse{
			IsError: true,
			Error:   err.Error(),
			User:    nil,
		}
		return res, err
	}
	//send verify email
	conn, err := grpc.Dial(fmt.Sprintf("mailer-service:%s", mailerGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println(err)
		res := &proto.RegisterResponse{
			IsError: true,
			Error:   err.Error(),
			User:    nil,
		}
		return res, err
	}

	defer conn.Close()

	c := proto.NewMailerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.SendRegisterEmailVerification(ctx, &proto.RegisterEmailRequest{
		ID:    int64(userId),
		Email: u.Email,
	})

	if err != nil {
		log.Println(err)
		res := &proto.RegisterResponse{
			IsError: true,
			Error:   err.Error(),
			User:    nil,
		}
		return res, err
	}

	userN, err := u.GetOne(userId)
	if err != nil {
		log.Println(err)
		res := &proto.RegisterResponse{
			IsError: true,
			Error:   err.Error(),
			User:    nil,
		}
		return res, err
	}
	// return response
	usr := &proto.User{
		ID:        int64(userN.ID),
		Email:     userN.Email,
		FirstName: userN.FirstName,
		LastName:  userN.LastName,
		RollId:    int64(userN.RollId),
		Password:  userN.Password,
		Active:    int32(userN.Active),
		CreatedAt: timestamppb.New(userN.CreatedAt),
		UpdatedAt: timestamppb.New(userN.UpdatedAt),
	}
	res := &proto.RegisterResponse{
		User:    usr,
		IsError: false,
		Error:   "",
	}
	return res, nil

}
func (l *LoginServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {

	input := req.GetUser()

	u := data.User{
		ID:        int(input.GetID()),
		Email:     input.GetEmail(),
		FirstName: input.GetFirstName(),
		LastName:  input.GetLastName(),
		RollId:    int(input.GetRollId()),
		Active:    int(input.GetActive()),
	}

	err := u.Update()
	if err != nil {
		log.Println(err)
		res := &proto.UpdateUserResponse{
			IsError: true,
			Error:   err.Error(),
			User:    nil,
		}
		return res, err
	}
	input.UpdatedAt = timestamppb.Now()

	res := &proto.UpdateUserResponse{
		User:    input,
		IsError: false,
		Error:   "",
	}
	return res, nil

}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("failed to listen for grpc: %v", err)
	}

	s := grpc.NewServer()

	proto.RegisterLoginServiceServer(s, &LoginServer{Models: app.Models})
	log.Printf("gRPC server started on port %s", gRpcPort)
	log.Println(fmt.Printf("add1: %s and add2: %s", lis.Addr(), lis.Addr().String()))

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen for grpc: %v", err)
	}

}
