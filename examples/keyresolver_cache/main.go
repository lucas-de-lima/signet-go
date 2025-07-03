package main

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"log"
	"time"

	"github.com/lucas-de-lima/signet-go/signet"
)

// imports locais para garantir visibilidade dos construtores
// (necessário para execução via go run .)
// _ "./key_provider.go"
// _ "./caching_resolver.go"

func main() {
	// Simula geração de chaves e registro no provider
	provider := NewSlowKeyProvider()
	kid := "v1"
	pub, priv, _ := ed25519.GenerateKey(nil)
	provider.RegisterKey(kid, pub)

	// Cria o resolver com cache (TTL 5s para demo)
	resolver := NewCachingKeyResolver(provider, 5*time.Second)

	// Gera um token com o kid
	token, err := signet.NewPayload().WithSubject("user-abc").WithKeyID(kid).Sign(priv)
	if err != nil {
		log.Fatalf("erro ao assinar token: %v", err)
	}

	ctx := context.Background()

	// Primeira validação: cache miss (lento)
	start := time.Now()
	_, err = signet.Parse(ctx, token, resolver.Resolve)
	missDuration := time.Since(start)
	if err != nil {
		log.Fatalf("falha na validação (miss): %v", err)
	}
	fmt.Printf("Validação (cache miss): %v\n", missDuration)

	// Segunda validação: cache hit (rápido)
	start = time.Now()
	_, err = signet.Parse(ctx, token, resolver.Resolve)
	hitDuration := time.Since(start)
	if err != nil {
		log.Fatalf("falha na validação (hit): %v", err)
	}
	fmt.Printf("Validação (cache hit): %v\n", hitDuration)

	fmt.Println("Demonstração concluída. Veja o README.md para detalhes.")
}
