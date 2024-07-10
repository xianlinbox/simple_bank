package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	USER_AGENT_KEY = "user-agent"
	CLIENT_IP_KEY = "grpcgateway-client-ip"
)

type Metadata struct {
	UserAgent string
	ClientIP string
}

func extractMetadata(ctx context.Context) (*Metadata){
	mtdt := &Metadata{}
	if md,ok := metadata.FromIncomingContext(ctx); ok{
		log.Printf("metadata found, %v", md)
		mtdt.UserAgent = md.Get(USER_AGENT_KEY)[0]
	}
	if p,ok := peer.FromContext(ctx); ok{
		log.Printf("peer found, %+v", p)
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}