# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/),
e este projeto adere ao [Semantic Versioning](https://semver.org/lang/pt-BR/).

## [Unreleased]

### Adicionado
- Templates de issues e pull requests
- Código de conduta para contribuidores
- Badges de status no README
- Workflow de CI melhorado com cobertura de testes
- Documentação completa de exemplos

### Alterado
- Melhorada formatação de todos os READMEs
- Atualizada documentação GoDoc

## [1.0.0] - 2024-01-XX

### Adicionado
- **API Principal**: Implementação completa da especificação Signet v1.0
  - `signet.NewPayload()`: Criação fluente de payloads
  - `signet.Parse()`: Validação rigorosa de tokens
  - `signet.WithAudience()`: Validação de audiência
  - `signet.RequireRole()`: Validação de papéis
  - `signet.WithRevocationCheck()`: Suporte a revogação (perfil STATEFUL)
  - `signet.WithMetricsRecorder()`: Observabilidade integrada

- **KeyResolver**: Sistema de rotação de chaves
  - `KeyResolverFunc`: Interface para resolução dinâmica de chaves
  - Suporte a múltiplas chaves simultâneas
  - Fallback para tokens antigos (sem `kid`)

- **Interceptor gRPC**: Proteção automática de endpoints
  - `grpcinterceptor.GRPCAuthInterceptor()`: Interceptor de autenticação
  - Injeção automática de payload no contexto
  - Mapeamento de erros para status gRPC

- **Observabilidade**: Métricas e monitoramento
  - Interface `MetricsRecorder` para integração com Prometheus/OpenTelemetry
  - Razões padronizadas para falhas de validação
  - Propagação de contexto para tracing distribuído

- **Exemplos de Produção**:
  - Servidor gRPC completo com autenticação
  - KeyResolver com cache TTL para performance
  - Demonstração de rotação de chaves
  - Integração com métricas Prometheus

### Segurança
- Criptografia Ed25519 por padrão
- Validação temporal obrigatória (exp/iat)
- Verificação de integridade obrigatória
- Tratamento seguro de erros com wrapping
- Validação rigorosa de claims

### Performance
- Serialização binária com Protocol Buffers
- API fluente para construção de payloads
- Cache TTL para KeyResolver
- Otimizações de memória e CPU

### Documentação
- README completo com exemplos práticos
- Documentação GoDoc detalhada
- Exemplos de produção funcionais
- Guia de contribuição
- Código de conduta

---

## Convenções de Versionamento

- **MAJOR**: Breaking changes na API
- **MINOR**: Novas funcionalidades (backward compatible)
- **PATCH**: Bug fixes e melhorias (backward compatible)

## Links

- [Especificação Signet](SPECIFICATION-v1.0.md)
- [Documentação GoDoc](GODOC-REFERENCE.md)
- [Exemplos de Produção](examples/) 