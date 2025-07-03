# 🔄 Exemplo: KeyResolver Resiliente com Cache TTL

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Este exemplo demonstra como implementar um `KeyResolverFunc` performático e resiliente, usando cache em memória com TTL para mitigar problemas de performance e ataques de negação de serviço (DoS) ao buscar chaves públicas.

## 🎯 O Problema

Em produção, a fonte de chaves públicas pode ser lenta (ex: chamada de rede, banco, JWKS). Buscar a chave a cada validação de token pode:

- ⚡ **Degradar o sistema**: latência acumulativa
- 🚨 **Abrir brecha para DoS**: ataques de negação de serviço
- 💰 **Aumentar custos**: chamadas desnecessárias a serviços externos

## ✅ A Solução

- **🔄 Cache em memória seguro para concorrência** (`sync.Map`)
- **⏰ TTL configurável**: cada chave é armazenada por um tempo limitado
- **⚡ Cache Hit**: resposta instantânea
- **🐌 Cache Miss**: busca na fonte lenta, armazena no cache

## 🔧 Como funciona o exemplo

1. **Provider simula fonte lenta** (100ms de latência artificial)
2. **CachingKeyResolver implementa** o padrão de cache TTL
3. **main.go gera um token** e valida duas vezes:
   - **Primeira**: cache miss (lento)
   - **Segunda**: cache hit (rápido)
4. **Tempo de cada validação** é impresso para comparação

## 🚀 Como executar

```bash
cd examples/keyresolver_cache
go run .
```

## 📊 Saída esperada

```
Validação (cache miss): 100.123456ms
Validação (cache hit):  123.456µs
Demonstração concluída. Veja o README.md para detalhes.
```

## 🏭 Padrão recomendado para produção

- **⏰ Use cache com TTL** (ex: 5 minutos) para cada chave pública
- **🔒 Sempre busque de fonte confiável** e valide o formato da chave
- **📈 Monitore métricas** de cache hit/miss para ajustar o TTL
- **🔄 Implemente fallback** para tokens antigos (sem `kid`)

---

> **💡 Este padrão é fundamental para ambientes de alta escala e segurança.** 