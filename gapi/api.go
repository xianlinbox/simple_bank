package gapi

import (
	security "github.com/xianlinbox/simple_bank/api/security"
	db "github.com/xianlinbox/simple_bank/db/sqlc"
	"github.com/xianlinbox/simple_bank/proto_code"
)

type GapiServer struct {
	store db.Store
	tokenMaker security.Maker
	proto_code.UnimplementedUsersServiceServer
}

func NewServer(store db.Store, tokenMaker security.Maker) *GapiServer {	
	server := &GapiServer{
		store: store,
		tokenMaker: tokenMaker,
	}
	return server
}

func (server *GapiServer) Start(address string) error {
	return nil
}