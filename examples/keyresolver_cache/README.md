# Exemplo: KeyResolver Resiliente com Cache TTL

Este exemplo demonstra como implementar um KeyResolverFunc performático e resiliente, usando cache em memória com TTL para mitigar problemas de performance e ataques de negação de serviço (DoS) ao buscar chaves públicas.

## O Problema

Em produção, a fonte de chaves públicas pode ser lenta (ex: chamada de rede, banco, JWKS). Buscar a chave a cada validação de token pode degradar o sistema e abrir brecha para DoS.

## A Solução

- **Cache em memória seguro para concorrência (sync.Map)**
- **TTL configurável**: cada chave é armazenada por um tempo limitado
- **Cache Hit**: resposta instantânea
- **Cache Miss**: busca na fonte lenta, armazena no cache

## Como funciona o exemplo

- Um provider simula uma fonte de chaves lenta (100ms de latência artificial)
- O CachingKeyResolver implementa o padrão de cache TTL
- O main.go gera um token, valida duas vezes:
  - Primeira: cache miss (lento)
  - Segunda: cache hit (rápido)
- O tempo de cada validação é impresso para comparação

## Como executar

```sh
cd examples/keyresolver_cache
go run .
```

## Saída esperada

```
Validação (cache miss): 100.123456ms
Validação (cache hit):  123.456µs
Demonstração concluída. Veja o README.md para detalhes.
```

## Padrão recomendado para produção

- Use cache com TTL (ex: 5 minutos) para cada chave pública
- Sempre busque de fonte confiável e valide o formato da chave
- Monitore métricas de cache hit/miss para ajustar o TTL

---

Este padrão é fundamental para ambientes de alta escala e segurança. 