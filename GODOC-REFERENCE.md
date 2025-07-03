# Referência Técnica GoDoc — signet-go

Este documento serve como guia de documentação GoDoc para todos os símbolos públicos da biblioteca signet-go, conforme checklist da História 5.2.

---

## Pacote signet

### Tipos Públicos

#### PayloadBuilder
Constrói um payload Signet de forma fluente e segura.

- Permite definir claims obrigatórios e opcionais (sub, aud, roles, custom_claims, sid, kid).
- Garante segurança por padrão: iat = agora, exp = agora + 15min.
- Use os métodos encadeáveis para construir o payload e finalize com Build() ou Sign().

**Exemplo:**
```go
payload, err := signet.NewPayload().
    WithSubject("user-123").
    WithAudience("api-backend").
    WithRole("admin").
    Build()
```

#### ValidationOption
Função que customiza o comportamento de validação do Parse.

- Usada para exigir audiência, papéis, revogação, métricas, etc.
- Componível: pode passar várias opções para Parse.

#### KeyResolverFunc
Função que resolve a chave pública correta para validação do token.

- Recebe o contexto e o kid extraído do token.
- Retorna a chave pública correspondente ou ErrUnknownKeyID.
- Permite integração com cache, JWKS, banco, etc.

**Exemplo:**
```go
keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
    if kid == "" {
        return myDefaultPublicKey, nil // fallback para tokens antigos
    }
    key, ok := keyMap[kid]
    if !ok {
        return nil, signet.ErrUnknownKeyID
    }
    return key, nil
}
```

#### MetricsRecorder
Interface para instrumentação de métricas de validação de tokens.

- Permite integração com Prometheus, OpenTelemetry, etc.
- Recebe o contexto, sucesso e razão padronizada.

**Exemplo:**
```go
type MeuRecorder struct{}
func (r *MeuRecorder) IncrementTokenValidation(ctx context.Context, success bool, reason string) {
    fmt.Printf("Token validado: sucesso=%v, motivo=%s\n", success, reason)
}
```

### Funções Públicas

#### NewPayload()
Cria um PayloadBuilder seguro (iat = agora, exp = agora + 15min).

#### Parse()
Deserializa e valida rigorosamente um SignetToken.

- Executa: deserialização, resolução de chave via KeyResolverFunc, verificação Ed25519, validação de exp/iat, claims e métricas.
- Retorna o payload validado ou erro sentinela. Use errors.Is para checagem.

**Exemplo:**
```go
payload, err := signet.Parse(ctx, tokenBytes, keyResolver, signet.WithAudience("api-backend"))
```

#### InjectPayloadIntoContext()
Injeta o payload validado no contexto para uso downstream (ex: gRPC).

#### PayloadFromContext()
Extrai o payload do contexto, se presente.

#### WithAudience()
Exige que o claim aud do payload seja igual ao fornecido.

#### RequireRole()
Exige que o payload contenha o papel fornecido.

#### RequireRoles()
Exige que o payload contenha todos os papéis fornecidos.

#### WithRevocationCheck()
Ativa validação STATEFUL, usando função checker para revogação.

#### WithMetricsRecorder()
Registra um implementador de MetricsRecorder para capturar métricas.

#### WithKeyID() (no PayloadBuilder)
Define o Key ID (kid) do payload para rotação/seleção de chave pública.

#### Sign() (no PayloadBuilder)
Serializa, assina e retorna o token final.

#### Build() (no PayloadBuilder)
Valida e retorna o payload pronto para uso.

### Variáveis e Constantes Públicas (Sentinel Errors e Razões de Falha)

- ErrTokenExpired: token expirado
- ErrTokenNotYetValid: iat no futuro
- ErrInvalidSignature: assinatura inválida
- ErrInvalidPayload: payload inválido
- ErrInvalidExpIat: exp <= iat
- ErrAudienceMismatch: audiência não corresponde
- ErrMissingRequiredRole: papel obrigatório ausente
- ErrTokenRevoked: token revogado
- ErrUnknownKeyID: kid não corresponde a nenhuma chave conhecida
- ReasonSuccess, ReasonInvalidSignature, ReasonTokenExpired, ReasonAudienceMismatch, ReasonInvalidPayload, ReasonTokenNotYetValid, ReasonMissingRequiredRole, ReasonTokenRevoked: razões padronizadas para métricas

---

## Pacote grpcinterceptor

### Funções Públicas

#### GRPCAuthInterceptor()
Retorna um grpc.UnaryServerInterceptor que protege endpoints gRPC validando tokens Signet.

- Extrai o token do header 'authorization-bin'.
- Valida usando signet.Parse com KeyResolverFunc e opções.
- Injeta o payload validado no contexto.
- Mapeia erros sentinela para status gRPC apropriados (Unauthenticated, PermissionDenied).

**Exemplo:**
```go
server := grpc.NewServer(
    grpc.UnaryInterceptor(
        grpcinterceptor.GRPCAuthInterceptor(keyResolver, signet.WithAudience("api-backend")),
    ),
)
```

---

Cada item acima pode ser expandido com exemplos mais avançados conforme necessidade do usuário ou publicação no pkg.go.dev.

> **Atenção:** Após validação deste esqueleto, cada seção será preenchida com a documentação GoDoc final, pronta para ser aplicada no código-fonte e publicada no pkg.go.dev. 