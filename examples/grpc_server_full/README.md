# 🚀 Exemplo: Servidor gRPC Completo com Signet-go

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![gRPC](https://img.shields.io/badge/gRPC-1.0+-lightgrey.svg)](https://grpc.io)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Este exemplo demonstra como proteger um microsserviço gRPC de ponta a ponta usando signet-go, integrando KeyResolver com cache, métricas e validação avançada.

## 🏗️ Arquitetura

- **🖥️ Servidor gRPC**: inicializa com `GRPCAuthInterceptor`, `KeyResolver` com cache TTL, `MetricsRecorder` e validação de audiência
- **📱 Cliente gRPC**: gera token válido, injeta no metadata e faz chamada `GetProfile`
- **🧪 Simulação de falhas**: altere o `kid` ou a audiência no cliente para ver erros `PermissionDenied`/`Unauthenticated`

## 🚀 Como executar

### 1. Inicie o servidor

```bash
cd examples/grpc_server_full/server
go run .
```

O servidor ficará ouvindo em `localhost:50051`.

### 2. Em outro terminal, execute o cliente

```bash
cd examples/grpc_server_full/client
go run .
```

Você verá a resposta do endpoint protegido:

```
Perfil recebido: subject=user-xyz, audience=my-protected-service, roles=[admin], expira=2025-07-02 22:00:00 +0000 UTC
```

### 3. Simule falhas de autenticação

- **🔑 Token inválido**: altere o `kid` no cliente para um valor não registrado
- **🎯 Audiência errada**: altere o campo `WithAudience` no cliente
- O cliente exibirá o erro retornado pelo servidor (`PermissionDenied`, `Unauthenticated`)

### 4. 📊 Métricas

O servidor imprime métricas de validação no log (em produção, integre com Prometheus).

## 🔧 Componentes do Exemplo

### Servidor (`server/main.go`)
- **gRPC Server** com interceptor de autenticação
- **KeyResolver com cache TTL** para performance
- **MetricsRecorder** para observabilidade
- **Validação de audiência** obrigatória

### Cliente (`client/main.go`)
- **Geração de token** com claims específicos
- **Injeção no metadata** gRPC
- **Chamada protegida** ao endpoint
- **Tratamento de erros** de autenticação

### Proto (`proto/profile.proto`)
- **Definição do serviço** ProfileService
- **Endpoint protegido** GetProfile
- **Mensagens** Profile e GetProfileRequest

## 🛡️ Segurança Implementada

- **🔐 Autenticação obrigatória** em todos os endpoints
- **⏰ Validação temporal** (exp/iat)
- **🎯 Validação de audiência** específica
- **🔑 Rotação de chaves** via KeyResolver
- **📊 Observabilidade** completa

---

> **💡 Este exemplo mostra como aplicar as melhores práticas de segurança, performance e observabilidade do signet-go em um cenário real de microsserviços gRPC.** 