package signet

import (
	"crypto/ed25519"
	"testing"
	"time"
)

// Testa se o NewPayload define iat e exp corretamente (segurança por padrão)
func TestNewPayloadDefaults(t *testing.T) {
	builder := NewPayload()
	payload, err := builder.Build()
	if err != nil {
		t.Fatalf("esperava payload válido, erro: %v", err)
	}
	if payload.Iat <= 0 || payload.Exp <= 0 {
		t.Error("iat e exp devem ser positivos")
	}
	if payload.Exp <= payload.Iat {
		t.Error("exp deve ser maior que iat")
	}
	if payload.Exp-payload.Iat != 15*60 {
		t.Errorf("exp padrão deve ser 15 minutos após iat, obtido: %d", payload.Exp-payload.Iat)
	}
}

// Testa encadeamento dos métodos do builder
func TestPayloadBuilderChaining(t *testing.T) {
	builder := NewPayload().
		WithSubject("user-123").
		WithAudience("service-x").
		WithRole("admin").
		WithCustomClaim("foo", "bar")
	payload, err := builder.Build()
	if err != nil {
		t.Fatalf("erro ao construir payload: %v", err)
	}
	if payload.Sub != "user-123" {
		t.Error("subject não definido corretamente")
	}
	if payload.Aud != "service-x" {
		t.Error("audience não definido corretamente")
	}
	if len(payload.Roles) != 1 || payload.Roles[0] != "admin" {
		t.Error("role não definido corretamente")
	}
	if v, ok := payload.CustomClaims["foo"]; !ok || v != "bar" {
		t.Error("custom claim não definido corretamente")
	}
}

// Testa validação de expiração e iat inválidos
func TestPayloadBuilderInvalidTimestamps(t *testing.T) {
	builder := NewPayload().WithIssuedAt(100).WithExpiration(50)
	_, err := builder.Build()
	if err != ErrInvalidExpIat {
		t.Errorf("esperava ErrInvalidExpIat, obteve: %v", err)
	}
	builder = NewPayload().WithIssuedAt(0).WithExpiration(0)
	_, err = builder.Build()
	if err != ErrInvalidPayload {
		t.Errorf("esperava ErrInvalidPayload, obteve: %v", err)
	}
}

// Testa assinatura e validação round-trip (fluxo completo)
func TestSignAndParse_RoundTrip(t *testing.T) {
	// Gera par de chaves Ed25519
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("erro ao gerar chave: %v", err)
	}
	// Cria payload completo
	builder := NewPayload().
		WithSubject("user-abc").
		WithAudience("api-x").
		WithRole("dev").
		WithCustomClaim("k1", "v1").
		WithSessionID([]byte("sessao-xyz")).
		WithExpiration(time.Now().Unix() + 60).
		WithIssuedAt(time.Now().Unix())
	pl, err := builder.Build()
	if err != nil {
		t.Fatalf("erro ao construir payload: %v", err)
	}
	// Assina o payload
	tokenBytes, err := builder.Sign(priv)
	if err != nil {
		t.Fatalf("erro ao assinar: %v", err)
	}
	// Valida o token
	parsed, err := Parse(tokenBytes, pub)
	if err != nil {
		t.Fatalf("erro ao validar token: %v", err)
	}
	// Verifica equivalência dos campos principais
	if parsed.Sub != pl.Sub || parsed.Aud != pl.Aud || parsed.Iat != pl.Iat || parsed.Exp != pl.Exp {
		t.Error("payload retornado não bate com o original")
	}
	if string(parsed.Sid) != string(pl.Sid) {
		t.Error("sid não bate")
	}
	if len(parsed.Roles) != len(pl.Roles) || parsed.Roles[0] != pl.Roles[0] {
		t.Error("roles não batem")
	}
	if v, ok := parsed.CustomClaims["k1"]; !ok || v != "v1" {
		t.Error("custom claim não bate")
	}
}

// Testa rejeição de token expirado
func TestParse_ExpiredToken(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	agora := time.Now().Unix()
	// exp < agora, mas exp > iat
	payload := NewPayload().WithIssuedAt(agora - 20).WithExpiration(agora - 10)
	tokenBytes, err := payload.Sign(priv)
	if err != nil {
		t.Fatalf("erro ao assinar: %v", err)
	}
	_, err = Parse(tokenBytes, pub)
	if err != ErrTokenExpired {
		t.Errorf("esperava ErrTokenExpired, obteve: %v", err)
	}
}

// Testa rejeição de token com iat no futuro
func TestParse_IatNoFuturo(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	agora := time.Now().Unix()
	// iat > agora, exp > iat
	payload := NewPayload().WithIssuedAt(agora + 3600).WithExpiration(agora + 3700)
	tokenBytes, err := payload.Sign(priv)
	if err != nil {
		t.Fatalf("erro ao assinar: %v", err)
	}
	_, err = Parse(tokenBytes, pub)
	if err != ErrTokenNotYetValid {
		t.Errorf("esperava ErrTokenNotYetValid, obteve: %v", err)
	}
}

// Testa rejeição por assinatura inválida
func TestParse_AssinaturaInvalida(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	payload := NewPayload()
	tokenBytes, err := payload.Sign(priv)
	if err != nil {
		t.Fatalf("erro ao assinar: %v", err)
	}
	// Corrompe a assinatura
	tokenBytes[len(tokenBytes)-1] ^= 0xFF
	_, err = Parse(tokenBytes, pub)
	if err != ErrInvalidSignature {
		t.Errorf("esperava ErrInvalidSignature, obteve: %v", err)
	}
}

// Testa validação de audiência
func TestParse_AudienceMismatch(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	payload := NewPayload().WithAudience("servico-a")
	tokenBytes, err := payload.Sign(priv)
	if err != nil {
		t.Fatalf("erro ao assinar: %v", err)
	}
	_, err = Parse(tokenBytes, pub, WithExpectedAudience("servico-b"))
	if err != ErrAudienceMismatch {
		t.Errorf("esperava ErrAudienceMismatch, obteve: %v", err)
	}
}
