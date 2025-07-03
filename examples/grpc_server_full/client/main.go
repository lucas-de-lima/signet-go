package main

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"log"
	"time"

	"github.com/lucas-de-lima/signet-go/signet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Empty struct{}
type Profile struct {
	Subject   string
	Audience  string
	Roles     []string
	ExpiresAt int64
}

func main() {
	// Gera par de chaves e token válido
	kid := "v1"
	_, priv, _ := ed25519.GenerateKey(nil)
	// Para simular token inválido, altere o kid ou a audiência abaixo
	token, err := signet.NewPayload().
		WithSubject("user-xyz").
		WithAudience("my-protected-service").
		WithRole("admin").
		WithKeyID(kid).
		Sign(priv)
	if err != nil {
		log.Fatalf("erro ao assinar token: %v", err)
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("erro ao conectar: %v", err)
	}
	defer conn.Close()

	client := NewProfileClient(conn)
	ctx := context.Background()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization-bin", string(token))

	resp, err := client.GetProfile(ctx, &Empty{})
	if err != nil {
		log.Fatalf("erro na chamada GetProfile: %v", err)
	}
	fmt.Printf("Perfil recebido: subject=%s, audience=%s, roles=%v, expira=%v\n",
		resp.Subject, resp.Audience, resp.Roles, time.Unix(resp.ExpiresAt, 0))
}

// Mock do client gRPC (em produção, use pb.NewProfileServiceClient)
type ProfileClient struct {
	cc *grpc.ClientConn
}

func NewProfileClient(cc *grpc.ClientConn) *ProfileClient {
	return &ProfileClient{cc}
}

func (c *ProfileClient) GetProfile(ctx context.Context, in *Empty) (*Profile, error) {
	// Simula chamada gRPC: em produção, use o método gerado
	// Aqui, apenas para demonstração, retorna um mock
	return &Profile{
		Subject:   "user-xyz",
		Audience:  "my-protected-service",
		Roles:     []string{"admin"},
		ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
	}, nil
}
