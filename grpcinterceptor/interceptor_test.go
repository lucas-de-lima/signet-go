package grpcinterceptor

import (
	"context"
	"crypto/ed25519"
	"testing"
	"time"

	signetv1 "github.com/lucas-de-lima/signet-go/proto/v1"
	"github.com/lucas-de-lima/signet-go/signet"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// handlerFake retorna o payload do contexto, se presente
func handlerFake(ctx context.Context, req interface{}) (interface{}, error) {
	payload, ok := signet.PayloadFromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "payload não encontrado no contexto")
	}
	return payload, nil
}

func TestGRPCAuthInterceptor_Success(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	tokenBytes, _ := signet.NewPayload().WithSubject("user-1").Sign(priv)
	interceptor := GRPCAuthInterceptor(pub)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization-bin", string(tokenBytes)))
	resp, err := interceptor(ctx, nil, nil, handlerFake)
	if err != nil {
		t.Fatalf("esperava sucesso, obteve erro: %v", err)
	}
	payload, ok := resp.(*signetv1.SignetPayload)
	if !ok || payload.Sub != "user-1" {
		t.Errorf("payload não injetado corretamente no contexto")
	}
}

func TestGRPCAuthInterceptor_TokenAusente(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(nil)
	interceptor := GRPCAuthInterceptor(pub)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs())
	_, err := interceptor(ctx, nil, nil, handlerFake)
	if status.Code(err) != codes.Unauthenticated {
		t.Errorf("esperava Unauthenticated para token ausente, obteve: %v", err)
	}
}

func TestGRPCAuthInterceptor_TokenInvalido(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(nil)
	interceptor := GRPCAuthInterceptor(pub)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization-bin", "corrompido"))
	_, err := interceptor(ctx, nil, nil, handlerFake)
	if status.Code(err) != codes.Unauthenticated {
		t.Errorf("esperava Unauthenticated para token inválido, obteve: %v", err)
	}
}

func TestGRPCAuthInterceptor_TokenExpirado(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	// Definir iat e exp no passado, com exp > iat, mas exp < agora
	iat := time.Now().Unix() - 100
	exp := iat + 1 // expira há ~99 segundos
	tokenBytes, _ := signet.NewPayload().WithIssuedAt(iat).WithExpiration(exp).Sign(priv)
	interceptor := GRPCAuthInterceptor(pub)
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization-bin", string(tokenBytes)))
	_, err := interceptor(ctx, nil, nil, handlerFake)
	if status.Code(err) != codes.PermissionDenied {
		t.Errorf("esperava PermissionDenied para token expirado, obteve: %v", err)
	}
}

func TestGRPCAuthInterceptor_PermissaoNegada(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	tokenBytes, _ := signet.NewPayload().WithRole("user").Sign(priv)
	interceptor := GRPCAuthInterceptor(pub, signet.RequireRole("admin"))
	ctx := metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization-bin", string(tokenBytes)))
	_, err := interceptor(ctx, nil, nil, handlerFake)
	if status.Code(err) != codes.PermissionDenied {
		t.Errorf("esperava PermissionDenied para role ausente, obteve: %v", err)
	}
}
