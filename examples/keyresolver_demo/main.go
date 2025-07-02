package main

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"log"

	"github.com/lucas-de-lima/signet-go/signet"
)

func main() {
	// Gerar dois pares de chaves (v1, v2)
	pubV1, privV1, _ := ed25519.GenerateKey(nil)
	pubV2, privV2, _ := ed25519.GenerateKey(nil)

	// Emissor 1: token com kid "v1"
	tokenV1, err := signet.NewPayload().WithSubject("user-abc").WithKeyID("v1").Sign(privV1)
	if err != nil {
		log.Fatalf("erro ao assinar token v1: %v", err)
	}
	// Emissor 2: token com kid "v2"
	tokenV2, err := signet.NewPayload().WithSubject("user-xyz").WithKeyID("v2").Sign(privV2)
	if err != nil {
		log.Fatalf("erro ao assinar token v2: %v", err)
	}

	// Validador: resolve a chave correta pelo kid
	keyMap := map[string]ed25519.PublicKey{
		"v1": pubV1,
		"v2": pubV2,
	}
	keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
		key, ok := keyMap[kid]
		if !ok {
			return nil, fmt.Errorf("kid desconhecido: %s", kid)
		}
		return key, nil
	}

	// Validar token v1
	payload1, err := signet.Parse(context.Background(), tokenV1, keyResolver)
	if err != nil {
		log.Fatalf("falha ao validar token v1: %v", err)
	}
	fmt.Printf("Token v1 validado! Subject: %s, kid: %s\n", payload1.Sub, payload1.Kid)

	// Validar token v2
	payload2, err := signet.Parse(context.Background(), tokenV2, keyResolver)
	if err != nil {
		log.Fatalf("falha ao validar token v2: %v", err)
	}
	fmt.Printf("Token v2 validado! Subject: %s, kid: %s\n", payload2.Sub, payload2.Kid)
}
