package core

import (
	"crypto/ed25519"
	"errors"
)

var (
	ErrInvalidPrivateKey  = errors.New("chave privada inválida: nula ou com tamanho incorreto")
	ErrInvalidPublicKey   = errors.New("chave pública inválida: nula ou com tamanho incorreto")
	ErrInvalidSignature   = errors.New("assinatura inválida: nula ou com tamanho incorreto")
	ErrNilData            = errors.New("dados para operação criptográfica não podem ser nulos")
	ErrVerificationFailed = errors.New("verificação da assinatura falhou")
)

// Sign gera uma assinatura Ed25519 para os dados fornecidos usando a chave privada.
// Retorna a assinatura ou um erro descritivo em caso de falha.
func Sign(privateKey ed25519.PrivateKey, data []byte) ([]byte, error) {
	if privateKey == nil || len(privateKey) != ed25519.PrivateKeySize {
		return nil, ErrInvalidPrivateKey
	}
	if data == nil {
		return nil, ErrNilData
	}
	sig := ed25519.Sign(privateKey, data)
	return sig, nil
}

// Verify verifica se a assinatura é válida para os dados e chave pública fornecidos.
// Retorna nil se válido, ou um erro descritivo caso contrário.
func Verify(publicKey ed25519.PublicKey, data, signature []byte) error {
	if publicKey == nil || len(publicKey) != ed25519.PublicKeySize {
		return ErrInvalidPublicKey
	}
	if data == nil {
		return ErrNilData
	}
	if signature == nil || len(signature) != ed25519.SignatureSize {
		return ErrInvalidSignature
	}
	if !ed25519.Verify(publicKey, data, signature) {
		return ErrVerificationFailed
	}
	return nil
}
