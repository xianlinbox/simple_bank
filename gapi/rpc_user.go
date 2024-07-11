package gapi

import (
	"context"

	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

func (server *GapiServer) UpdateUser(ctx context.Context, req *proto_code.UpdateUserRequest) (*proto_code.UpdateUserResponse, error) {
	authPayload, err :=server.authorizaion(ctx)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "authentication failed: %v", err)
	}
	arg:=db.UpdateUserParams{
		Username: authPayload.Username,
		FullName: pgtype.Text{
			String: req.GetFullName(),
			Valid: req.GetFullName() != "",
		},
		Email: pgtype.Text{
			String: req.GetEmail(),
			Valid: req.GetEmail() != "",
		},
	}

	password := req.GetPassword(); 
	if password != "" {
		password, err:=util.EncryptPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to encypt the password: %v", err)
		}
		arg.Password = pgtype.Text{
			String: password,
			Valid: true,
		}
		arg.PasswordExpiredAt = pgtype.Timestamptz{
			Time: time.Now().Add(time.Hour * 24 * 90),
			Valid: true,
		}
	}

	updatedUser,err:=server.store.UpdateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "didn't find the user: %v", err)
	}
	return &proto_code.UpdateUserResponse{
		User: convertUserToProto(updatedUser),
	}, nil
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
	metadata := extractMetadata(c)
	_, err = server.store.AddSession(c, db.AddSessionParams{
		ID:	uuid.New(),
		Username: user.Username,
		RefreshToken: newToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIP,
		ExpiredAt:    time.Now().Add(time.Hour * 24),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save session: %v", err)
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
		PasswordExpiredAt: timestamppb.New(user.PasswordExpiredAt),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}