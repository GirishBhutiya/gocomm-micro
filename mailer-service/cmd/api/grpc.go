package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"mailer-service/proto"
	"net"

	"google.golang.org/grpc"
)

type MailerServer struct {
	proto.UnimplementedMailerServiceServer
	Config
}

func (l *MailerServer) SendRegisterEmailVerification(ctx context.Context, req *proto.RegisterEmailRequest) (*proto.RegisterEmailResponse, error) {
	link := req.GetLink()
	email := req.GetEmail()

	var data struct {
		Link template.HTML
	}
	data.Link = template.HTML(link)

	//send mail with verify email link
	go l.Config.SendMail(l.Config.EnvVars.FromEmail, email, l.Config.EnvVars.VerifyEmailSubject, l.Config.EnvVars.VerifyEmailTemplate, data)

	res := &proto.RegisterEmailResponse{
		IsError: false,
		Message: "",
	}
	return res, nil
}
func (l *MailerServer) SendForgotPasswordLinkEmail(ctx context.Context, req *proto.ForgotPasswordRequestEmail) (*proto.ForgotPasswordResponseEmail, error) {
	link := req.GetLink()
	email := req.GetEmail()

	var data struct {
		Link template.HTML
	}
	data.Link = template.HTML(link)

	//send mail with verify email link
	go l.Config.SendMail(l.Config.EnvVars.FromEmail, email, l.Config.EnvVars.ForgotPasswordEmailSubject, l.Config.EnvVars.ForgotPasswordEmailTemplate, data)

	res := &proto.ForgotPasswordResponseEmail{
		IsError: false,
		Message: "",
	}
	return res, nil
}

func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("failed to listen for grpc: %v", err)
	}

	s := grpc.NewServer()

	proto.RegisterMailerServiceServer(s, &MailerServer{Config: *app})
	log.Printf("gRPC server started on port %s", gRpcPort)
	log.Println(fmt.Printf("add1: %s and add2: %s", lis.Addr(), lis.Addr().String()))

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to listen for grpc: %v", err)
	}

}
