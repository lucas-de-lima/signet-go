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
	tokenBytes, _ := signet.NewPayload().WithSubject("user-1").WithKeyID("v1").Sign(priv)
	// KeyResolverFunc que retorna a chave correta para qualquer kid
	keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
		return pub, nil
	}
	interceptor := GRPCAuthInterceptor(keyResolver)
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

func TestGRPCAuthInterceptor_Fails(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
		return pub, nil
	}

	// Gerar tokens para cenários
	iat := time.Now().Unix() - 100
	exp := iat + 1 // expira há ~99 segundos
	expiredToken, _ := signet.NewPayload().WithIssuedAt(iat).WithExpiration(exp).WithKeyID("v1").Sign(priv)
	wrongRoleToken, _ := signet.NewPayload().WithRole("user").WithKeyID("v1").Sign(priv)

	testCases := []struct {
		name         string
		ctx          context.Context
		options      []signet.ValidationOption
		expectedCode codes.Code
	}{
		{"Token ausente", metadata.NewIncomingContext(context.Background(), metadata.Pairs()), nil, codes.Unauthenticated},
		{"Token corrompido", metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization-bin", "corrompido")), nil, codes.Unauthenticated},
		{"Token expirado", metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization-bin", string(expiredToken))), nil, codes.PermissionDenied},
		{"Role incorreto", metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization-bin", string(wrongRoleToken))), []signet.ValidationOption{signet.RequireRole("admin")}, codes.PermissionDenied},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			interceptor := GRPCAuthInterceptor(keyResolver, tc.options...)
			_, err := interceptor(tc.ctx, nil, nil, handlerFake)
			if status.Code(err) != tc.expectedCode {
				t.Errorf("esperava código %v, mas obteve %v", tc.expectedCode, status.Code(err))
			}
		})
	}
}
