package main

import (
	"authentication/data/models"
	"authentication/data/proto"
	"authentication/internal/urlsigner"
	"authentication/internal/util"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type LoginServer struct {
	proto.UnimplementedLoginServiceServer
	//Models data.Models
	Config      util.ConfigVars
	userService models.UserService
}

func (l *LoginServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {

	input := req.GetLoginReq()
	//us := models.UserService{}
	//var u data.User
	user, err := l.userService.GetByEmail(input.GetUsername())
	if err != nil {
		log.Println(err)
		res := &proto.LoginResponse{
			Authenticated: false,
			User:          nil,
		}
		return res, err
	}

	match, err := l.userService.PasswordMatches(input.GetPassword(), user)
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

	u := models.User{
		Email:     input.GetEmail(),
		FirstName: input.GetFirstName(),
		LastName:  input.GetLastName(),
		Password:  input.GetPassword(),
		RollId:    int(input.GetRollId()),
		//Active:    int(input.GetActive()),
	}

	//log.Println("register input password:", input.GetPassword())
	userN, err := l.userService.Insert(u)
	if err != nil {
		log.Println(err)
		res := &proto.RegisterResponse{
			IsError: true,
			Error:   err.Error(),
			User:    nil,
		}
		return res, err
	}

	//create verify email link
	/* link := util.GetVeifyEmailLink(u.Email, l.Config.FrontEndDomain)

	sign := urlsigner.Signer{
		Secret: []byte(fmt.Sprintf("%s%s", l.Config.HashSecretKeyVerifyEmail, u.Email)),
	}
	signedLink := sign.GenerateTokenFromString(link) */
	signedLink := util.GetFullVerifyEmailLink(u.Email, l.Config.FrontEndDomain, l.Config.HashSecretKeyVerifyEmail)
	//send verify email
	go l.callRegisterEmailVerification(u.Email, signedLink)

	if err != nil {

		log.Println(err)
		res := &proto.RegisterResponse{
			IsError: true,
			Error:   err.Error(),
			User:    nil,
		}
		return res, err
	}

	/* userN, err := l.userService.GetOne(userId)
	if err != nil {
		log.Println(err)
		res := &proto.RegisterResponse{
			IsError: true,
			Error:   err.Error(),
			User:    nil,
		}
		return res, err
	} */
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

	u := &models.User{
		ID:        int(input.GetID()),
		Email:     input.GetEmail(),
		FirstName: input.GetFirstName(),
		LastName:  input.GetLastName(),
		RollId:    int(input.GetRollId()),
		Active:    int(input.GetActive()),
	}

	err := l.userService.Update(u)
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
func (l *LoginServer) VerifyRegisteredEmail(ctx context.Context, req *proto.VerifyRegisteredEmailRequest) (*proto.VerifyRegisteredEmailResponse, error) {

	link := req.GetLink()

	//mUrl := fmt.Sprintf("%s%s", l.Config.FrontEndDomain, link)

	parsedURL, err := url.Parse(link)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	params, err := url.ParseQuery(parsedURL.RawQuery)
	if err != nil {
		return nil, err
	}

	email := params.Get("email")

	sign := urlsigner.Signer{
		Secret: []byte(fmt.Sprintf("%s%s", l.Config.HashSecretKeyVerifyEmail, email)),
	}

	valid := sign.VerifyToken(link)

	res := &proto.VerifyRegisteredEmailResponse{
		Verified: valid,
	}
	if !valid {
		res.Message = "Email not verified."
	}

	u := &models.User{
		Email: email,
	}
	err = l.userService.ValidateEmail(u)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	res.Message = "Email sucessfully verified."

	return res, nil
}
func (l *LoginServer) ForgotPassword(ctx context.Context, req *proto.ForgotPasswordRequest) (*proto.ForgotPasswordResponse, error) {
	email := req.GetEmail()

	/*u := data.User{
		Email: email,
	}*/
	_, err := l.userService.GetByEmail(email)
	if err != nil {
		log.Println(err)
		res := &proto.ForgotPasswordResponse{
			IsError: true,
			Message: "No account found.",
		}
		return res, err
	}

	link := util.GetForgotPasswordLink(email, l.Config.FrontEndDomain)

	sign := urlsigner.Signer{
		Secret: []byte(fmt.Sprintf("%s%s", l.Config.HashSecretKeyForgotPassword, email)),
	}
	signedLink := sign.GenerateTokenFromString(link)

	//send verify email
	conn, err := grpc.Dial(fmt.Sprintf("mailer-service:%s", mailerGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {

		log.Println(err)
		res := &proto.ForgotPasswordResponse{
			IsError: true,
			Message: "Sending Error in mail",
		}
		return res, err
	}

	defer conn.Close()

	c := proto.NewMailerServiceClient(conn)

	_, err = c.SendForgotPasswordLinkEmail(ctx, &proto.ForgotPasswordRequestEmail{
		Email: email,
		Link:  signedLink,
	})

	if err != nil {

		log.Println(err)
		res := &proto.ForgotPasswordResponse{
			IsError: true,
			Message: "Sending Error in mail",
		}
		return res, err
	}
	res := &proto.ForgotPasswordResponse{
		IsError: false,
		Message: "Forgot Password Link sent via mail. Please check your mailbox including spam folder.",
	}
	return res, nil

}
func (l *LoginServer) ResetPassword(ctx context.Context, req *proto.ResetPasswordRequest) (*proto.ResetPasswordResponse, error) {
	email := req.GetEmail()
	theLink := req.GetLink()
	password := req.GetPassword()

	sign := urlsigner.Signer{
		Secret: []byte(fmt.Sprintf("%s%s", l.Config.HashSecretKeyForgotPassword, email)),
	}
	mUrl := fmt.Sprintf("%s%s", l.Config.FrontEndDomain, theLink)
	valid := sign.VerifyToken(mUrl)

	res := &proto.ResetPasswordResponse{}
	if !valid {
		res.Message = "Password Reset Link not veririfed"
		res.IsError = true
		return res, errors.New("password reset link not veririfed")
	}
	//make sure not expired
	expired := sign.Expired(mUrl, 60)
	if expired {
		res.Message = "Link Expired"
		res.IsError = true
		return res, errors.New("password reset link expired")
	}
	u := &models.User{
		Email:    email,
		Password: password,
	}
	//log.Println("Password in grpc:", password)
	err := l.userService.ResetPassword(password, u)
	if err != nil {
		log.Println(err)
		res := &proto.ResetPasswordResponse{
			IsError: true,
			Message: err.Error(),
		}
		return res, err
	}

	res = &proto.ResetPasswordResponse{
		IsError: false,
		Message: "Password changed",
	}
	return res, nil

}
func (l *LoginServer) callRegisterEmailVerification(userEmail, link string) {

	conn, err := grpc.Dial(fmt.Sprintf("mailer-service:%s", mailerGrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Println(err)
	}

	defer conn.Close()
	c := proto.NewMailerServiceClient(conn)

	_, err = c.SendRegisterEmailVerification(context.Background(), &proto.RegisterEmailRequest{
		Email: userEmail,
		Link:  link,
	})
	if err != nil {
		log.Println(err)
	}
}
func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("failed to listen for grpc: %v", err)
	}

	s := grpc.NewServer()

	proto.RegisterLoginServiceServer(s, &LoginServer{
		//Models: app.Models,
		userService: app.userService,
		Config:      app.EnvVars,
	})
	log.Printf("gRPC server started on port %s", gRpcPort)
	log.Println(fmt.Printf("add1: %s and add2: %s", lis.Addr(), lis.Addr().String()))

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen for grpc: %v", err)
	}

}
