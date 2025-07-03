# 🔑 Exemplo: Demonstração de KeyResolver

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Este exemplo demonstra o uso básico do `KeyResolverFunc` para rotação de chaves públicas em tokens Signet.

## 🎯 Objetivo

Mostrar como implementar um `KeyResolverFunc` simples que:
- Resolve chaves públicas baseado no `kid` (Key ID) do token
- Suporta múltiplas chaves para rotação
- Fornece fallback para tokens antigos (sem `kid`)

## 🔧 Como funciona

1. **Gera múltiplas chaves** com IDs diferentes (`v1`, `v2`)
2. **Cria KeyResolver** que mapeia `kid` → chave pública
3. **Gera tokens** com diferentes `kid`s
4. **Valida tokens** usando o KeyResolver
5. **Demonstra fallback** para tokens sem `kid`

## 🚀 Como executar

```bash
cd examples/keyresolver_demo
go run .
```

## 📊 Saída esperada

```
Token v1 validado com sucesso
Token v2 validado com sucesso
Token sem kid validado com sucesso (fallback)
Token com kid desconhecido rejeitado
```

## 💡 Conceitos Demonstrados

- **🔑 Rotação de chaves**: múltiplas chaves ativas simultaneamente
- **🆔 Key ID (kid)**: identificador para seleção de chave
- **🔄 Fallback**: suporte para tokens antigos sem `kid`
- **❌ Rejeição**: tokens com `kid` desconhecido são rejeitados

## 🏭 Padrão para Produção

```go
keyResolver := func(ctx context.Context, kid string) (ed25519.PublicKey, error) {
    if kid == "" {
        return defaultKey, nil // fallback para tokens antigos
    }
    
    key, exists := keyStore.Get(kid)
    if !exists {
        return nil, signet.ErrUnknownKeyID
    }
    
    return key, nil
}
```

---

> **💡 Este exemplo é ideal para entender os fundamentos do KeyResolver antes de implementar cache e outras otimizações.** 