# 🔧 Exemplos de Produção — signet-go

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Este diretório contém exemplos práticos e completos de como usar o signet-go em cenários reais de produção.

## 📚 Índice dos Exemplos

### 🚀 [Servidor gRPC Completo](grpc_server_full/)
**Proteção completa de microsserviços gRPC**
- Interceptor de autenticação
- KeyResolver com cache TTL
- Métricas e observabilidade
- Cliente e servidor funcionais

### 🔄 [KeyResolver com Cache](keyresolver_cache/)
**Performance e resiliência para produção**
- Cache em memória com TTL
- Mitigação de ataques DoS
- Comparação de latência cache hit/miss
- Padrão recomendado para alta escala

### 🔑 [Demonstração KeyResolver](keyresolver_demo/)
**Fundamentos do KeyResolver**
- Rotação de chaves básica
- Suporte a múltiplas chaves
- Fallback para tokens antigos
- Conceitos fundamentais

### 📊 [Métricas Prometheus](metrics_prometheus/)
**Observabilidade completa**
- Integração com Prometheus
- Métricas de validação
- Endpoint `/metrics`
- Queries e dashboards

## 🎯 Como usar os exemplos

### 1. Escolha o exemplo
```bash
cd examples/[nome-do-exemplo]
```

### 2. Execute o exemplo
```bash
go run .
```

### 3. Siga as instruções
Cada exemplo tem seu próprio README com instruções detalhadas.

## 🏗️ Arquitetura dos Exemplos

```
examples/
├── grpc_server_full/     # 🚀 Exemplo completo gRPC
│   ├── server/          # Servidor protegido
│   ├── client/          # Cliente funcional
│   └── proto/           # Definições protobuf
├── keyresolver_cache/    # 🔄 Cache TTL
├── keyresolver_demo/     # 🔑 Fundamentos
└── metrics_prometheus/   # 📊 Observabilidade
```

## 💡 Ordem de Aprendizado Recomendada

1. **🔑 [keyresolver_demo](keyresolver_demo/)**: Entenda os fundamentos
2. **🔄 [keyresolver_cache](keyresolver_cache/)**: Otimize para produção
3. **📊 [metrics_prometheus](metrics_prometheus/)**: Adicione observabilidade
4. **🚀 [grpc_server_full](grpc_server_full/)**: Aplique tudo junto

## 🛠️ Pré-requisitos

- **Go 1.21+**: Para todos os exemplos
- **gRPC**: Para o exemplo completo
- **Prometheus**: Para métricas (opcional)

## 🔧 Personalização

Todos os exemplos são **auto-contidos** e podem ser facilmente adaptados:

- **Altere as chaves**: Use suas próprias chaves Ed25519
- **Modifique claims**: Personalize subject, audience, roles
- **Integre com sua infraestrutura**: KMS, Redis, banco de dados
- **Adicione suas métricas**: Prometheus, OpenTelemetry, etc.

## 📖 Documentação Relacionada

- **[README Principal](../README.md)**: Visão geral do projeto
- **[GoDoc Reference](../GODOC-REFERENCE.md)**: Documentação da API
- **[Especificação Signet](../SPECIFICATION-v1.0.md)**: Padrão técnico

---

> **💡 Estes exemplos demonstram as melhores práticas para usar signet-go em ambientes de produção seguros e escaláveis.**

---

## 👨‍💻 Autor

**Lucas de Lima**
- 📧 Email: dev.lucasdelima@gmail.com
- 💼 LinkedIn: [dev-lucasdelima](https://www.linkedin.com/in/dev-lucasdelima/)
- 🚀 Software Engineer | Backend, Full Stack and Mobile Development 