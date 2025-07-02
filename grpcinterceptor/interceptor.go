package grpcinterceptor

import (
	"context"
	"crypto/ed25519"
	"log"

	"github.com/lucas-de-lima/signet-go/signet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCAuthInterceptor retorna um UnaryServerInterceptor que valida tokens Signet e injeta o payload no contexto.
func GRPCAuthInterceptor(publicKey ed25519.PublicKey, options ...signet.ValidationOption) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Extrai o token do metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata ausente no contexto gRPC")
		}
		tokens := md["authorization-bin"]
		if len(tokens) == 0 || len(tokens[0]) == 0 {
			return nil, status.Error(codes.Unauthenticated, "token Signet ausente no header authorization-bin")
		}
		tokenBytes := []byte(tokens[0])

		// Valida o token usando signet.Parse
		payload, err := signet.Parse(tokenBytes, publicKey, options...)
		if err != nil {
			// Mapeia erros sentinela para status gRPC apropriados
			switch err {
			case signet.ErrInvalidSignature, signet.ErrInvalidPayload:
				return nil, status.Error(codes.Unauthenticated, "token inválido ou corrompido")
			case signet.ErrTokenExpired, signet.ErrAudienceMismatch, signet.ErrMissingRequiredRole, signet.ErrTokenRevoked:
				return nil, status.Error(codes.PermissionDenied, "token não autorizado: "+err.Error())
			default:
				// Logar o erro inesperado no servidor para observabilidade.
				log.Printf("ERRO: erro de autenticação inesperado no interceptor Signet: %v", err)
				return nil, status.Error(codes.Unauthenticated, "falha de autenticação interna")
			}
		}

		// Injeta o payload validado no contexto
		ctx = signet.InjectPayloadIntoContext(ctx, payload)
		return handler(ctx, req)
	}
}
