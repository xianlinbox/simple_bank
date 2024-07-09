package gapi

import (
	"context"

	"time"

	db "github.com/xianlinbox/simple_bank/db/sqlc"
	"github.com/xianlinbox/simple_bank/proto_code"
	"github.com/xianlinbox/simple_bank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *GapiServer) CreateUser(ctx context.Context, req *proto_code.CreateUserRequest) (*proto_code.CreateUserResponse, error) {
	hashedPassword, err := util.EncryptPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed encrupt password: %v", err)
	}
	user, err := server.store.AddUser(ctx, db.AddUserParams{
		Username: req.Username,
		Password: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
	})

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed save user: %v", err)
	}

	response := proto_code.CreateUserResponse{
		User: convertUserToProto(user),
	}
	return &response, nil
}

func (server *GapiServer) Login( c context.Context, req *proto_code.LoginRequest) (*proto_code.LoginResponse, error) {
	user, err := server.store.GetUser(c, req.Username)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "User is not found: %v", err)
	}

	err = util.CheckPassword(user.Password, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Password is not correct: %v", err)
	}
	newToken,err := server.tokenMaker.GenerateToken(user.Username, time.Minute*15)
	if err != nil{
		return nil, status.Errorf(codes.Internal, "Failed to generate new token: %v", err)
	}

	response:= proto_code.LoginResponse{
		User: convertUserToProto(user),
		AccessToken: newToken,
	}
	
	return &response, nil
}

func convertUserToProto(user db.User) *proto_code.User {
	return &proto_code.User{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordExpiredAt: timestamppb.New(user.PasswordExpiredAt.Time),
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
	}
}