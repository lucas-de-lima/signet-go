# Exemplo: Servidor gRPC Completo com Signet-go

Este exemplo demonstra como proteger um microsserviço gRPC de ponta a ponta usando signet-go, integrando KeyResolver com cache, métricas e validação avançada.

## Arquitetura

- **Servidor gRPC**: inicializa com GRPCAuthInterceptor, KeyResolver com cache TTL, MetricsRecorder e validação de audiência.
- **Cliente gRPC**: gera token válido, injeta no metadata e faz chamada GetProfile.
- **Simulação de falhas**: altere o kid ou a audiência no cliente para ver erros PermissionDenied/Unauthenticated.

## Como executar

### 1. Inicie o servidor

```sh
cd examples/grpc_server_full/server
go run .
```

O servidor ficará ouvindo em `localhost:50051`.

### 2. Em outro terminal, execute o cliente

```sh
cd examples/grpc_server_full/client
go run .
```

Você verá a resposta do endpoint protegido:

```
Perfil recebido: subject=user-xyz, audience=my-protected-service, roles=[admin], expira=2025-07-02 22:00:00 +0000 UTC
```

### 3. Simule falhas de autenticação

- **Token inválido**: altere o `kid` no cliente para um valor não registrado.
- **Audiência errada**: altere o campo `WithAudience` no cliente.
- O cliente exibirá o erro retornado pelo servidor (PermissionDenied, Unauthenticated).

### 4. Métricas

O servidor imprime métricas de validação no log (em produção, integre com Prometheus).

---

Este exemplo mostra como aplicar as melhores práticas de segurança, performance e observabilidade do signet-go em um cenário real de microsserviços gRPC. 