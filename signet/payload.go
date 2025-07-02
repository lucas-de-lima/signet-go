package signet

import (
	"context"
	"crypto/ed25519"
	"errors"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/lucas-de-lima/signet-go/internal/core"
	signetv1 "github.com/lucas-de-lima/signet-go/proto/v1"
)

// Sentinel errors para validação e segurança
var (
	ErrTokenExpired        = errors.New("token expirado")
	ErrTokenNotYetValid    = errors.New("token com iat no futuro")
	ErrInvalidSignature    = errors.New("assinatura inválida")
	ErrInvalidPayload      = errors.New("payload inválido")
	ErrInvalidPrivateKey   = errors.New("chave privada inválida")
	ErrInvalidPublicKey    = errors.New("chave pública inválida")
	ErrInvalidExpIat       = errors.New("exp deve ser maior que iat")
	ErrAudienceMismatch    = errors.New("audiência do token não corresponde à esperada")
	ErrMissingRequiredRole = errors.New("payload não possui o(s) papel(is) requerido(s)")
	ErrTokenRevoked        = errors.New("token revogado (sid presente na lista de revogação)")
)

// Razões padronizadas para métricas de validação
const (
	ReasonSuccess             = "success"
	ReasonInvalidSignature    = "invalid_signature"
	ReasonTokenExpired        = "token_expired"
	ReasonAudienceMismatch    = "audience_mismatch"
	ReasonInvalidPayload      = "invalid_payload"
	ReasonTokenNotYetValid    = "token_not_yet_valid"
	ReasonMissingRequiredRole = "missing_required_role"
	ReasonTokenRevoked        = "token_revoked"
)

// PayloadBuilder implementa a API fluente para construção de payloads Signet.
// Segurança por padrão: iat = agora, exp = agora + 15min.
// Todos os métodos retornam o builder para encadeamento.
// O método Build() valida as regras de negócio.
type PayloadBuilder struct {
	payload *signetv1.SignetPayload
}

// NewPayload cria um builder com iat = agora e exp = agora + 15min.
// Este é o ponto de entrada recomendado para criar um novo payload seguro.
// O desenvolvedor pode sobrescrever esses valores usando WithIssuedAt/WithExpiration.
func NewPayload() *PayloadBuilder {
	now := time.Now().Unix()
	return &PayloadBuilder{
		payload: &signetv1.SignetPayload{
			Iat: now,
			Exp: now + 15*60, // 15 minutos
		},
	}
}

// WithSubject define o subject do payload (ex: ID do usuário).
func (b *PayloadBuilder) WithSubject(sub string) *PayloadBuilder {
	b.payload.Sub = sub
	return b
}

// WithAudience define o público do payload (ex: sistema alvo).
func (b *PayloadBuilder) WithAudience(aud string) *PayloadBuilder {
	b.payload.Aud = aud
	return b
}

// WithRole adiciona um papel (role) ao payload.
func (b *PayloadBuilder) WithRole(role string) *PayloadBuilder {
	b.payload.Roles = append(b.payload.Roles, role)
	return b
}

// WithCustomClaim adiciona um claim customizado (chave-valor) ao payload.
func (b *PayloadBuilder) WithCustomClaim(key, value string) *PayloadBuilder {
	if b.payload.CustomClaims == nil {
		b.payload.CustomClaims = make(map[string]string)
	}
	b.payload.CustomClaims[key] = value
	return b
}

// WithSessionID define o ID de sessão (sid) para perfis STATEFUL.
func (b *PayloadBuilder) WithSessionID(sid []byte) *PayloadBuilder {
	b.payload.Sid = sid
	return b
}

// WithIssuedAt sobrescreve o campo iat (issued at) do payload.
func (b *PayloadBuilder) WithIssuedAt(iat int64) *PayloadBuilder {
	b.payload.Iat = iat
	return b
}

// WithExpiration sobrescreve o campo exp (expiration) do payload.
func (b *PayloadBuilder) WithExpiration(exp int64) *PayloadBuilder {
	b.payload.Exp = exp
	return b
}

// WithKeyID define o Key ID (kid) do payload, usado para rotação e seleção de chave pública.
func (b *PayloadBuilder) WithKeyID(kid string) *PayloadBuilder {
	b.payload.Kid = kid
	return b
}

// Build valida as regras de negócio e retorna o payload pronto para uso.
// Valida se exp > iat e se ambos são positivos.
// Retorna erro se as regras forem violadas.
func (b *PayloadBuilder) Build() (*signetv1.SignetPayload, error) {
	// Timestamps não podem ser zero ou negativos
	if b.payload.Iat <= 0 || b.payload.Exp <= 0 {
		return nil, ErrInvalidPayload
	}
	// Validação de sanidade: exp deve ser maior que iat
	if b.payload.Exp <= b.payload.Iat {
		return nil, ErrInvalidExpIat
	}
	return b.payload, nil
}

// Sign serializa o payload, assina com a chave privada e retorna o token final.
// Segue rigorosamente o processo da especificação Signet v1.0:
// 1. Valida e serializa o payload
// 2. Assina os bytes do payload usando Ed25519
// 3. Monta o SignetToken
// 4. Serializa o token final
func (b *PayloadBuilder) Sign(privateKey ed25519.PrivateKey) ([]byte, error) {
	// 1. Construir e validar o payload
	payload, err := b.Build()
	if err != nil {
		return nil, fmt.Errorf("falha ao construir payload: %w", err)
	}
	// 2. Serializar o payload para bytes canônicos (protobuf)
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("falha ao serializar payload para protobuf: %w", err)
	}
	// 3. Assinar os bytes do payload usando Ed25519
	signature, err := core.Sign(privateKey, payloadBytes)
	if err != nil {
		return nil, fmt.Errorf("falha ao assinar payload no núcleo criptográfico: %w", err)
	}
	// 4. Construir o SignetToken final
	token := &signetv1.SignetToken{
		Payload:   payloadBytes,
		Signature: signature,
	}
	// 5. Serializar o token final
	tokenBytes, err := proto.Marshal(token)
	if err != nil {
		return nil, fmt.Errorf("falha ao serializar SignetToken para protobuf: %w", err)
	}
	return tokenBytes, nil
}

// MetricsRecorder define um contrato para registrar métricas de validação de tokens.
// As implementações podem usar este hook para integrar com sistemas como Prometheus, OpenTelemetry, etc.
type MetricsRecorder interface {
	// IncrementTokenValidation é chamada para cada operação de Parse.
	// O parâmetro 'success' é true se a validação foi bem-sucedida.
	// Em caso de falha, 'failureReason' contém uma string curta e padronizada (use as constantes Reason*).
	IncrementTokenValidation(ctx context.Context, success bool, failureReason string)
}

// ValidationOption define opções de validação para a função Parse.
// Permite customizar o comportamento de validação (ex: pular expiração para testes).
type ValidationOption func(*validationConfig)

type validationConfig struct {
	skipExpirationCheck bool
	skipIssuedAtCheck   bool
	expectedAudience    string
	requiredRoles       []string
	revocationChecker   func([]byte) bool
	metricsRecorder     MetricsRecorder
}

// WithSkipExpirationCheck permite pular a verificação de expiração (útil para testes).
func WithSkipExpirationCheck() ValidationOption {
	return func(c *validationConfig) {
		c.skipExpirationCheck = true
	}
}

// WithSkipIssuedAtCheck permite pular a verificação de iat (útil para testes).
func WithSkipIssuedAtCheck() ValidationOption {
	return func(c *validationConfig) {
		c.skipIssuedAtCheck = true
	}
}

// WithExpectedAudience define a audiência esperada para validação.
func WithExpectedAudience(audience string) ValidationOption {
	return func(c *validationConfig) {
		c.expectedAudience = audience
	}
}

// WithAudience exige que o claim aud do payload seja igual ao fornecido.
func WithAudience(audience string) ValidationOption {
	return func(c *validationConfig) {
		c.expectedAudience = audience
	}
}

// RequireRole exige que o payload contenha o papel fornecido.
// Pode ser chamado múltiplas vezes para exigir múltiplos papéis.
func RequireRole(role string) ValidationOption {
	return func(c *validationConfig) {
		c.requiredRoles = append(c.requiredRoles, role)
	}
}

// RequireRoles exige que o payload contenha todos os papéis fornecidos.
func RequireRoles(roles ...string) ValidationOption {
	return func(c *validationConfig) {
		c.requiredRoles = append(c.requiredRoles, roles...)
	}
}

// WithRevocationCheck ativa a validação STATEFUL, usando a função checker fornecida.
// checker deve retornar true se o sid estiver revogado.
func WithRevocationCheck(checker func(sid []byte) bool) ValidationOption {
	return func(c *validationConfig) {
		c.revocationChecker = checker
	}
}

// WithMetricsRecorder registra um implementador de MetricsRecorder para capturar métricas de validação.
func WithMetricsRecorder(recorder MetricsRecorder) ValidationOption {
	return func(c *validationConfig) {
		c.metricsRecorder = recorder
	}
}

// Parse deserializa e valida rigorosamente um SignetToken.
//
// Fluxo seguro e obrigatório:
// 1. Deserializa o token (SignetToken)
// 2. Deserializa o payload para extrair os claims e o campo 'kid' (Key ID)
// 3. Resolve a chave pública correta via KeyResolverFunc, usando o contexto e o kid
// 4. Verifica a assinatura ANTES de processar o payload (garante integridade)
// 5. Valida expiração, iat, audiência, papéis, revogação, etc.
// 6. Emite métricas de sucesso/falha, se configurado
// 7. Retorna o payload válido ou erro sentinela específico
//
// Este fluxo garante máxima segurança, rastreabilidade e suporte a rotação de chaves e múltiplos emissores.
func Parse(ctx context.Context, tokenBytes []byte, keyResolver KeyResolverFunc, options ...ValidationOption) (*signetv1.SignetPayload, error) {
	config := &validationConfig{}
	for _, option := range options {
		option(config)
	}
	recordMetricAndReturn := func(ctx context.Context, success bool, reason string, payload *signetv1.SignetPayload, err error) (*signetv1.SignetPayload, error) {
		if config.metricsRecorder != nil {
			config.metricsRecorder.IncrementTokenValidation(ctx, success, reason)
		}
		return payload, err
	}
	// 1. Deserializar o token
	var token signetv1.SignetToken
	if err := proto.Unmarshal(tokenBytes, &token); err != nil {
		return recordMetricAndReturn(ctx, false, ReasonInvalidPayload, nil, fmt.Errorf("falha ao deserializar SignetToken: %w", err))
	}
	if token.Payload == nil || token.Signature == nil {
		return recordMetricAndReturn(ctx, false, ReasonInvalidPayload, nil, ErrInvalidPayload)
	}
	// 2. Deserializar o payload para extrair o kid
	var payload signetv1.SignetPayload
	if err := proto.Unmarshal(token.Payload, &payload); err != nil {
		return recordMetricAndReturn(ctx, false, ReasonInvalidPayload, nil, fmt.Errorf("falha ao deserializar SignetPayload: %w", err))
	}
	// 3. Resolver a chave pública via keyResolver
	pubKey, err := keyResolver(ctx, payload.Kid)
	if err != nil {
		return recordMetricAndReturn(ctx, false, ReasonInvalidSignature, nil, fmt.Errorf("falha ao resolver chave pública para kid '%s': %w", payload.Kid, err))
	}
	// 4. Verificar a assinatura
	if err := core.Verify(pubKey, token.Payload, token.Signature); err != nil {
		if errors.Is(err, core.ErrVerificationFailed) {
			return recordMetricAndReturn(ctx, false, ReasonInvalidSignature, nil, fmt.Errorf("falha na verificação criptográfica do núcleo: %w", ErrInvalidSignature))
		}
		return recordMetricAndReturn(ctx, false, ReasonInvalidSignature, nil, fmt.Errorf("falha na verificação criptográfica do núcleo: %w", err))
	}
	// 5. Validações temporais (a menos que explicitamente puladas)
	now := time.Now().Unix()
	if !config.skipExpirationCheck {
		if payload.Exp <= now {
			return recordMetricAndReturn(ctx, false, ReasonTokenExpired, nil, ErrTokenExpired)
		}
	}
	if !config.skipIssuedAtCheck {
		if payload.Iat > now {
			return recordMetricAndReturn(ctx, false, ReasonTokenNotYetValid, nil, ErrTokenNotYetValid)
		}
	}
	if config.expectedAudience != "" && payload.Aud != config.expectedAudience {
		return recordMetricAndReturn(ctx, false, ReasonAudienceMismatch, nil, ErrAudienceMismatch)
	}
	if len(config.requiredRoles) > 0 {
		rolesMap := make(map[string]struct{}, len(payload.Roles))
		for _, r := range payload.Roles {
			rolesMap[r] = struct{}{}
		}
		for _, required := range config.requiredRoles {
			if _, ok := rolesMap[required]; !ok {
				return recordMetricAndReturn(ctx, false, ReasonMissingRequiredRole, nil, ErrMissingRequiredRole)
			}
		}
	}
	if config.revocationChecker != nil && len(payload.Sid) > 0 {
		if config.revocationChecker(payload.Sid) {
			return recordMetricAndReturn(ctx, false, ReasonTokenRevoked, nil, ErrTokenRevoked)
		}
	}
	return recordMetricAndReturn(ctx, true, ReasonSuccess, &payload, nil)
}

// contextKey é uma chave privada para evitar colisão no contexto
// Não exportada para garantir isolamento
var contextKey = struct{}{}

// InjectPayloadIntoContext injeta o payload validado no contexto
func InjectPayloadIntoContext(ctx context.Context, payload *signetv1.SignetPayload) context.Context {
	return context.WithValue(ctx, contextKey, payload)
}

// PayloadFromContext extrai o payload do contexto, se presente
func PayloadFromContext(ctx context.Context) (*signetv1.SignetPayload, bool) {
	v := ctx.Value(contextKey)
	if v == nil {
		return nil, false
	}
	payload, ok := v.(*signetv1.SignetPayload)
	return payload, ok
}

// KeyResolverFunc define a assinatura para funções que resolvem uma chave pública
// com base no contexto e no kid do token.
type KeyResolverFunc func(ctx context.Context, kid string) (ed25519.PublicKey, error)
