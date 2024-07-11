package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/xianlinbox/simple_bank/api/security"
	"google.golang.org/grpc/metadata"
)

const (
	authorization_header="authorization"
)
func (server *GapiServer) authorizaion(ctx context.Context) (*security.Payload, error) {
	mtdt, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata not found")
	}
	auth_header := mtdt.Get(authorization_header)

	if len(auth_header) == 0 {
		return nil, fmt.Errorf("authorisation header not found")
	}
	token := auth_header[0]
	token_fields := strings.Fields(token)
	if len(token_fields) != 2 || strings.ToLower(token_fields[0]) != "bearer" {
		return nil, fmt.Errorf("authorisation header format must be Bearer {token}")
	}
	payload, err := server.tokenMaker.VerifyToken(token_fields[1])
	if err != nil {
		return nil, fmt.Errorf("authorisation token verification failed: %v", err)
	}
	return payload, nil
}