"# signet-go" 

## Observabilidade e Métricas

A biblioteca signet-go permite instrumentação de métricas de validação de tokens via a interface `MetricsRecorder`, sem acoplamento a sistemas específicos (Prometheus, OpenTelemetry, etc).

### Como usar

Implemente a interface:

```go
import (
    "context"
    "crypto/ed25519"
    "github.com/lucas-de-lima/signet-go/signet"
)

type MeuRecorder struct{}
func (r *MeuRecorder) IncrementTokenValidation(ctx context.Context, success bool, failureReason string) {
    // Exemplo: enviar para Prometheus, logar, etc.
    fmt.Printf("Token validado: sucesso=%v, motivo=%s\n", success, failureReason)
}
```

### Rotação de Chaves e KeyResolverFunc

Para suportar múltiplos emissores e rotação de chaves, use o campo `kid` e um `KeyResolverFunc`:

```go
// Emissor: gera tokens com um kid específico
payload, _ := signet.NewPayload().WithSubject("user-1").WithKeyID("v1").Sign(privateKeyV1)

// Validador: resolve a chave correta pelo kid
keyMap := map[string]ed25519.PublicKey{
    "v1": publicKeyV1,
    "v2": publicKeyV2,
}
keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
    key, ok := keyMap[kid]
    if !ok {
        return nil, fmt.Errorf("kid desconhecido: %s", kid)
    }
    return key, nil
}

// Validação do token
payload, err := signet.Parse(context.Background(), tokenBytes, keyResolver, signet.WithMetricsRecorder(&MeuRecorder{}))
```

- O campo `kid` identifica a chave usada para assinar o token.
- O `KeyResolverFunc` permite rotação de chaves e múltiplos emissores sem downtime.
- `success` indica se a validação foi bem-sucedida.
- `failureReason` é uma string padronizada (ex: `signet.ReasonTokenExpired`, `signet.ReasonInvalidSignature`).

Assim, você pode criar dashboards, alertas e monitorar a saúde do seu sistema de autenticação em tempo real, com máxima segurança e flexibilidade. 
