package core

import (
	"crypto/ed25519"
	"errors"
	"testing"
)

// TestSignAndVerifyHappyPath testa o fluxo "feliz": assina e verifica com sucesso.
// Garante que a assinatura gerada é válida para os dados e chave correspondente.
func TestSignAndVerifyHappyPath(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("erro ao gerar chave: %v", err)
	}
	data := []byte("dados de teste signet-go")
	sig, err := sign(priv, data)
	if err != nil {
		t.Fatalf("erro ao assinar: %v", err)
	}
	if err := verify(pub, data, sig); err != nil {
		t.Errorf("verificação falhou no happy path: %v", err)
	}
}

// TestVerifyWrongKey garante que a verificação falha ao usar uma chave pública incorreta.
func TestVerifyWrongKey(t *testing.T) {
	_, priv, _ := ed25519.GenerateKey(nil)
	pub2, _, _ := ed25519.GenerateKey(nil)
	data := []byte("dados de teste signet-go")
	sig, _ := sign(priv, data)
	if err := verify(pub2, data, sig); err == nil {
		t.Error("verificação deveria falhar com chave pública incorreta")
	}
}

// TestVerifyCorruptedData garante que a verificação falha se os dados forem alterados após a assinatura.
func TestVerifyCorruptedData(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	data := []byte("dados de teste signet-go")
	sig, _ := sign(priv, data)
	corrupted := make([]byte, len(data))
	copy(corrupted, data)
	corrupted[0] ^= 0xFF // altera um byte
	if err := verify(pub, corrupted, sig); err == nil {
		t.Error("verificação deveria falhar com dados corrompidos")
	}
}

// TestVerifyCorruptedSignature garante que a verificação falha se a assinatura for alterada.
func TestVerifyCorruptedSignature(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	data := []byte("dados de teste signet-go")
	sig, _ := sign(priv, data)
	corrupted := make([]byte, len(sig))
	copy(corrupted, sig)
	corrupted[0] ^= 0xFF // altera um byte
	if err := verify(pub, data, corrupted); err == nil {
		t.Error("verificação deveria falhar com assinatura corrompida")
	}
}

// TestSignInputErrors cobre os casos de erro de entrada para a função sign.
func TestSignInputErrors(t *testing.T) {
	_, priv, _ := ed25519.GenerateKey(nil)
	data := []byte("dados")
	// Chave privada nula
	_, err := sign(nil, data)
	if !errors.Is(err, ErrInvalidPrivateKey) {
		t.Errorf("esperado erro ErrInvalidPrivateKey, mas obteve: %v", err)
	}
	// Chave privada tamanho errado
	_, err = sign(priv[:10], data)
	if !errors.Is(err, ErrInvalidPrivateKey) {
		t.Errorf("esperado erro ErrInvalidPrivateKey, mas obteve: %v", err)
	}
	// Dados nulos
	_, err = sign(priv, nil)
	if !errors.Is(err, ErrNilData) {
		t.Errorf("esperado erro ErrNilData, mas obteve: %v", err)
	}
}

// TestVerifyInputErrors cobre os casos de erro de entrada para a função verify.
func TestVerifyInputErrors(t *testing.T) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	data := []byte("dados")
	sig, _ := sign(priv, data)
	// Chave pública nula
	err := verify(nil, data, sig)
	if !errors.Is(err, ErrInvalidPublicKey) {
		t.Errorf("esperado erro ErrInvalidPublicKey, mas obteve: %v", err)
	}
	// Chave pública tamanho errado
	err = verify(pub[:10], data, sig)
	if !errors.Is(err, ErrInvalidPublicKey) {
		t.Errorf("esperado erro ErrInvalidPublicKey, mas obteve: %v", err)
	}
	// Dados nulos
	err = verify(pub, nil, sig)
	if !errors.Is(err, ErrNilData) {
		t.Errorf("esperado erro ErrNilData, mas obteve: %v", err)
	}
	// Assinatura nula
	err = verify(pub, data, nil)
	if !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("esperado erro ErrInvalidSignature, mas obteve: %v", err)
	}
	// Assinatura tamanho errado
	err = verify(pub, data, sig[:10])
	if !errors.Is(err, ErrInvalidSignature) {
		t.Errorf("esperado erro ErrInvalidSignature, mas obteve: %v", err)
	}
}

// BenchmarkSign mede o tempo e as alocações da operação de assinatura para um payload típico.
func BenchmarkSign(b *testing.B) {
	_, priv, _ := ed25519.GenerateKey(nil)
	payload := make([]byte, 256)
	for i := range payload {
		payload[i] = byte(i % 256)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sign(priv, payload)
		if err != nil {
			b.Fatalf("erro ao assinar: %v", err)
		}
	}
}

// BenchmarkVerify mede o tempo e as alocações da operação de verificação para um payload típico.
func BenchmarkVerify(b *testing.B) {
	pub, priv, _ := ed25519.GenerateKey(nil)
	payload := make([]byte, 256)
	for i := range payload {
		payload[i] = byte(i % 256)
	}
	sig, err := sign(priv, payload)
	if err != nil {
		b.Fatalf("erro ao assinar: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := verify(pub, payload, sig); err != nil {
			b.Fatalf("verificação falhou: %v", err)
		}
	}
}
