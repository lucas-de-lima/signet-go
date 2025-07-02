package signet

import (
	"context"
	"crypto/ed25519"
	"errors"
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
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("erro ao gerar chave: %v", err)
	}
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
	tokenBytes, err := builder.Sign(priv)
	if err != nil {
		t.Fatalf("erro ao assinar: %v", err)
	}
	keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
		return pub, nil
	}
	parsed, err := Parse(context.Background(), tokenBytes, keyResolver)
	if err != nil {
		t.Fatalf("erro ao validar token: %v", err)
	}
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
	payload := NewPayload().WithIssuedAt(agora - 20).WithExpiration(agora - 10)
	tokenBytes, err := payload.Sign(priv)
	if err != nil {
		t.Fatalf("erro ao assinar: %v", err)
	}
	keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
		return pub, nil
	}
	_, err = Parse(context.Background(), tokenBytes, keyResolver)
	if !errors.Is(err, ErrTokenExpired) {
		t.Errorf("esperava ErrTokenExpired, obteve: %v", err)
	}
}

// Testa rejeição de token com iat no futuro
func TestParse_IatNoFuturo(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	agora := time.Now().Unix()
	payload := NewPayload().WithIssuedAt(agora + 3600).WithExpiration(agora + 3700)
	tokenBytes, err := payload.Sign(priv)
	if err != nil {
		t.Fatalf("erro ao assinar: %v", err)
	}
	keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
		return pub, nil
	}
	_, err = Parse(context.Background(), tokenBytes, keyResolver)
	if !errors.Is(err, ErrTokenNotYetValid) {
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
	tokenBytes[len(tokenBytes)-1] ^= 0xFF
	keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
		return pub, nil
	}
	_, err = Parse(context.Background(), tokenBytes, keyResolver)
	if !errors.Is(err, ErrInvalidSignature) {
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
	keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
		return pub, nil
	}
	_, err = Parse(context.Background(), tokenBytes, keyResolver, WithExpectedAudience("servico-b"))
	if !errors.Is(err, ErrAudienceMismatch) {
		t.Errorf("esperava ErrAudienceMismatch, obteve: %v", err)
	}
}

// Testa validação declarativa de audiência (WithAudience)
func TestParse_WithAudience(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	payload := NewPayload().WithAudience("servico-a")
	tokenBytes, err := payload.Sign(priv)
	if err != nil {
		t.Fatalf("erro ao assinar: %v", err)
	}
	keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
		return pub, nil
	}
	// Sucesso: audiência correta
	_, err = Parse(context.Background(), tokenBytes, keyResolver, WithAudience("servico-a"))
	if err != nil {
		t.Errorf("não deveria falhar para audiência correta: %v", err)
	}
	// Falha: audiência incorreta
	_, err = Parse(context.Background(), tokenBytes, keyResolver, WithAudience("servico-b"))
	if !errors.Is(err, ErrAudienceMismatch) {
		t.Errorf("esperava ErrAudienceMismatch, obteve: %v", err)
	}
	// Falha: claim aud ausente
	payloadSemAud := NewPayload()
	tokenBytes, err = payloadSemAud.Sign(priv)
	if err != nil {
		t.Fatalf("erro ao assinar: %v", err)
	}
	_, err = Parse(context.Background(), tokenBytes, keyResolver, WithAudience("servico-x"))
	if !errors.Is(err, ErrAudienceMismatch) {
		t.Errorf("esperava ErrAudienceMismatch para aud ausente, obteve: %v", err)
	}
}

// Testa validação declarativa de papéis (RequireRole/RequireRoles) usando table-driven
func TestParse_RequireRole(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	payload := NewPayload().WithRole("admin").WithRole("auditor")
	tokenBytes, _ := payload.Sign(priv)
	keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
		return pub, nil
	}

	testCases := []struct {
		name          string
		options       []ValidationOption
		expectedError error
	}{
		{"Sucesso: papel único presente", []ValidationOption{RequireRole("admin")}, nil},
		{"Sucesso: múltiplos papéis presentes (RequireRole)", []ValidationOption{RequireRole("admin"), RequireRole("auditor")}, nil},
		{"Sucesso: múltiplos papéis presentes (RequireRoles)", []ValidationOption{RequireRoles("admin", "auditor")}, nil},
		{"Falha: papel ausente", []ValidationOption{RequireRole("user")}, ErrMissingRequiredRole},
		{"Falha: um papel presente, um ausente", []ValidationOption{RequireRole("admin"), RequireRole("user")}, ErrMissingRequiredRole},
		{"Falha: múltiplos papéis ausentes (RequireRoles)", []ValidationOption{RequireRoles("user", "auditor", "admin", "root")}, ErrMissingRequiredRole},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(context.Background(), tokenBytes, keyResolver, tc.options...)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("esperava erro '%v', mas obteve '%v'", tc.expectedError, err)
			}
		})
	}
}

// Testa validação declarativa de revogação (WithRevocationCheck) usando table-driven
func TestParse_WithRevocationCheck(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	sid := []byte("sid-123")
	payload := NewPayload().WithSessionID(sid)
	tokenBytes, _ := payload.Sign(priv)
	payloadStateless := NewPayload()
	tokenBytesStateless, _ := payloadStateless.Sign(priv)
	keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
		return pub, nil
	}

	testCases := []struct {
		name          string
		token         []byte
		checker       func([]byte) bool
		expectedError error
	}{
		{"Sucesso: sid não revogado", tokenBytes, func(_ []byte) bool { return false }, nil},
		{"Falha: sid revogado", tokenBytes, func(s []byte) bool { return string(s) == "sid-123" }, ErrTokenRevoked},
		{"Sucesso: token STATELESS ignora checker", tokenBytesStateless, func(_ []byte) bool { return true }, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse(context.Background(), tc.token, keyResolver, WithRevocationCheck(tc.checker))
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("esperava erro '%v', mas obteve '%v'", tc.expectedError, err)
			}
		})
	}
}
