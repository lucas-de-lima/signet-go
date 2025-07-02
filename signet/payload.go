package signet

import (
	"crypto/ed25519"
	"errors"
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
		return nil, err
	}
	// 2. Serializar o payload para bytes canônicos (protobuf)
	payloadBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, ErrInvalidPayload
	}
	// 3. Assinar os bytes do payload usando Ed25519
	signature, err := core.Sign(privateKey, payloadBytes)
	if err != nil {
		return nil, err
	}
	// 4. Construir o SignetToken final
	token := &signetv1.SignetToken{
		Payload:   payloadBytes,
		Signature: signature,
	}
	// 5. Serializar o token final
	tokenBytes, err := proto.Marshal(token)
	if err != nil {
		return nil, ErrInvalidPayload
	}
	return tokenBytes, nil
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

// Parse deserializa e valida rigorosamente um SignetToken.
// Fluxo seguro e obrigatório:
// 1. Deserializa o token (SignetToken)
// 2. Verifica a assinatura ANTES de processar o payload
// 3. Deserializa o payload
// 4. Valida expiração, iat e audiência (por padrão)
// 5. Retorna o payload válido ou erro sentinela específico
func Parse(tokenBytes []byte, publicKey ed25519.PublicKey, options ...ValidationOption) (*signetv1.SignetPayload, error) {
	// Configurar opções de validação
	config := &validationConfig{}
	for _, option := range options {
		option(config)
	}
	// 1. Deserializar o token
	var token signetv1.SignetToken
	if err := proto.Unmarshal(tokenBytes, &token); err != nil {
		return nil, ErrInvalidPayload
	}
	// 2. Verificar se os campos obrigatórios estão presentes
	if token.Payload == nil || token.Signature == nil {
		return nil, ErrInvalidPayload
	}
	// 3. Verificar a assinatura ANTES de deserializar o payload
	if err := core.Verify(publicKey, token.Payload, token.Signature); err != nil {
		return nil, ErrInvalidSignature
	}
	// 4. Deserializar o payload
	var payload signetv1.SignetPayload
	if err := proto.Unmarshal(token.Payload, &payload); err != nil {
		return nil, ErrInvalidPayload
	}
	// 5. Validações temporais (a menos que explicitamente puladas)
	now := time.Now().Unix()
	if !config.skipExpirationCheck {
		if payload.Exp <= now {
			return nil, ErrTokenExpired
		}
	}
	if !config.skipIssuedAtCheck {
		if payload.Iat > now {
			return nil, ErrTokenNotYetValid
		}
	}
	// 6. Validação de audiência declarativa (WithAudience)
	if config.expectedAudience != "" && payload.Aud != config.expectedAudience {
		return nil, ErrAudienceMismatch
	}
	// 7. Validação de papéis declarativa (RequireRole)
	if len(config.requiredRoles) > 0 {
		rolesMap := make(map[string]struct{}, len(payload.Roles))
		for _, r := range payload.Roles {
			rolesMap[r] = struct{}{}
		}
		for _, required := range config.requiredRoles {
			if _, ok := rolesMap[required]; !ok {
				return nil, ErrMissingRequiredRole
			}
		}
	}
	// 8. Validação de revogação declarativa (WithRevocationCheck)
	if config.revocationChecker != nil && len(payload.Sid) > 0 {
		if config.revocationChecker(payload.Sid) {
			return nil, ErrTokenRevoked
		}
	}
	return &payload, nil
}
