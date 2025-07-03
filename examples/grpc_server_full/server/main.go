package main

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/lucas-de-lima/signet-go/grpcinterceptor"
	"github.com/lucas-de-lima/signet-go/signet"
	"google.golang.org/grpc"
)

// Simples serviço gRPC de perfil
// (em produção, use proto gerado)
type ProfileService struct{}

func (s *ProfileService) GetProfile(ctx context.Context, req *Empty) (*Profile, error) {
	payload, ok := signet.PayloadFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("payload não encontrado no contexto")
	}
	return &Profile{
		Subject:   payload.Sub,
		Audience:  payload.Aud,
		Roles:     payload.Roles,
		ExpiresAt: payload.Exp,
	}, nil
}

type Empty struct{}
type Profile struct {
	Subject   string
	Audience  string
	Roles     []string
	ExpiresAt int64
}

type DemoMetrics struct{}

func (m *DemoMetrics) IncrementTokenValidation(ctx context.Context, success bool, reason string) {
	log.Printf("[metrics] success=%v, reason=%s", success, reason)
}

func main() {
	// Setup de chaves e provider
	provider := NewSlowKeyProvider()
	kid := "v1"
	pub, _, _ := ed25519.GenerateKey(nil)
	provider.RegisterKey(kid, pub)
	resolver := NewCachingKeyResolver(provider, 5*time.Minute)

	// Interceptor com KeyResolver, métricas e validação de audiência
	interceptor := grpcinterceptor.GRPCAuthInterceptor(
		resolver.Resolve,
		signet.WithMetricsRecorder(&DemoMetrics{}),
		signet.WithAudience("my-protected-service"),
	)

	server := grpc.NewServer(grpc.UnaryInterceptor(interceptor))
	// Em produção, use pb.RegisterProfileServiceServer(server, &ProfileService{})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("erro ao escutar: %v", err)
	}
	log.Println("Servidor gRPC ouvindo em :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("erro ao servir: %v", err)
	}
}
