Signet for Go (signet-go)
signet-go é a implementação de referência em Go para a especificação Signet, um padrão moderno para cargas de segurança de aplicação. Ele fornece uma alternativa segura, performática e nativa a tokens baseados em texto, otimizada para ecossistemas de alta performance como gRPC.

Porquê Signet?
Em um mundo de microsserviços, a segurança e a performance da comunicação são fundamentais. Enquanto o JWT foi uma ferramenta útil para a web baseada em JSON, ele introduz uma sobrecarga desnecessária e uma superfície de risco em ambientes binários.

O Signet foi projetado desde o primeiro dia para este novo paradigma, oferecendo:

Segurança por Design: Usa criptografia moderna (Ed25519) por padrão. A validação temporal e de integridade é obrigatória, não opcional.

Performance Inerente: Utiliza Protocol Buffers para uma serialização binária extremamente rápida e compacta, eliminando o overhead do JSON e Base64.

Clareza Operacional: Com uma API fluente, tratamento de erros robusto e ferramentas de observabilidade, o signet-go foi feito para ser usado com confiança em produção.

Pronto para o Mundo Real: Suporte nativo para rotação de chaves (KeyResolver), revogação de tokens (perfil STATEFUL) e integração plug-and-play com gRPC.

Instalação
go get github.com/signet/signet-go

Início Rápido (Quick Start)
Este exemplo leva você do zero a um token validado em menos de um minuto.

package main

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"log"

	"github.com/signet/signet-go/signet"
)

func main() {
	// 1. Gere um par de chaves Ed25519. Em produção, você as carregaria de
	// um local seguro (ex: KMS, Vault).
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatalf("Falha ao gerar chaves: %v", err)
	}

	// 2. Crie e assine um novo token Signet usando a API fluente.
	// O 'kid' (Key ID) ajuda o validador a encontrar a chave correta.
	tokenBytes, err := signet.NewPayload().
		WithSubject("user-12345").
		WithAudience("billing-service").
		WithRole("user").
		WithKeyID("v1").
		Sign(privateKey)
	if err != nil {
		log.Fatalf("Falha ao assinar o token: %v", err)
	}

	fmt.Printf("Token Signet gerado com sucesso! (%d bytes)\n", len(tokenBytes))

	// 3. Crie um KeyResolver. Esta função ensina o Parse a encontrar a chave pública
	// correta com base no 'kid' do token.
	keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
		if kid == "v1" {
			return publicKey, nil
		}
		return nil, fmt.Errorf("kid desconhecido: %s", kid)
	}

	// 4. Valide o token. O Parse verifica a assinatura, expiração e outros claims.
	payload, err := signet.Parse(context.Background(), tokenBytes, keyResolver)
	if err != nil {
		log.Fatalf("Falha ao validar o token: %v", err)
	}

	fmt.Printf("Token validado com sucesso para o sujeito: %s\n", payload.GetSub())
}

Funcionalidades Avançadas
O signet-go oferece um controle granular sobre a validação através de Opções Funcionais.

Validando Claims Específicos
Você pode exigir que um token tenha uma audiência e papéis específicos.

payload, err := signet.Parse(ctx, tokenBytes, keyResolver,
    signet.WithAudience("billing-service"), // Exige que 'aud' seja "billing-service"
    signet.RequireRole("user"),             // Exige que o papel "user" esteja presente
)

Protegendo um Servidor gRPC
Proteger todos os seus endpoints gRPC é tão simples quanto adicionar o interceptor na criação do servidor.

import "github.com/signet/signet-go/grpcinterceptor"

// ... seu keyResolver ...

server := grpc.NewServer(
    grpc.UnaryInterceptor(
        grpcinterceptor.GRPCAuthInterceptor(keyResolver,
            signet.WithAudience("billing-service"),
        ),
    ),
)

Revogação de Tokens (Perfil STATEFUL)
Para tokens que precisam ser revogados antes de expirarem, use o perfil STATEFUL.

// Emissor: Crie um token com um ID de sessão único.
sid := []byte("session-xyz-789")
tokenBytes, _ := signet.NewPayload().
    WithSessionID(sid).
    // ... outros claims ...
    Sign(privateKey)

// Validador: Forneça uma função que verifica se o SID está na sua blacklist.
revocationChecker := func(sidToCheck []byte) bool {
    // Lógica para verificar no Redis, banco de dados, etc.
    return myBlacklist.Contains(string(sidToCheck))
}

payload, err := signet.Parse(ctx, tokenBytes, keyResolver,
    signet.WithRevocationCheck(revocationChecker),
)

Próximos Passos
Para entender a filosofia e os princípios por trás do projeto, leia a Especificação Signet.

Para uma referência completa da API, visite a Documentação GoDoc.

Para exemplos de produção, incluindo KeyResolver com cache e integração com métricas, explore o diretório /examples.