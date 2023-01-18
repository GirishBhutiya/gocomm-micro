package main

import (
	"context"

	"github.com/GirishBhutiya/gocomm-micro/db/data"
	pb "github.com/GirishBhutiya/gocomm-micro/db/proto"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	Models data.Repository
}

func (u *UserServer) UserAuthenticate(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	input := req.GetUser()

	user := data.User{
		ID:        int(input.Id),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
		Active:    int(input.Active),
		CreatedAt: input.CreateAt.AsTime(),
		UpdatedAt: input.UpdatedAt.AsTime(),
	}

	//validate the user against the database
	userR, err := u.Models.GetByEmail(user.Email)
	if err != nil {
		res := &pb.UserResponse{Result: "invalid credentials"}
		return res, err
	}

	valid, err := u.Models.PasswordMatches(user.Password, *userR)
	if err != nil || !valid {
		//app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		res := &pb.UserResponse{Result: "invalid credentials"}
		return res, err
	}
	if valid {
		res := &pb.UserResponse{Result: "Correct Credentials, Logged In!!!"}
		return res, err
	} else {
		res := &pb.UserResponse{Result: "invalid credentials"}
		return res, err
	}
	
}
