"# signet-go" 

## Observabilidade e Métricas

A biblioteca signet-go permite instrumentação de métricas de validação de tokens via a interface `MetricsRecorder`, sem acoplamento a sistemas específicos (Prometheus, OpenTelemetry, etc).

### Como usar

Implemente a interface:

```go
import (
    "context"
    "github.com/lucas-de-lima/signet-go/signet"
)

type MeuRecorder struct{}
func (r *MeuRecorder) IncrementTokenValidation(ctx context.Context, success bool, failureReason string) {
    // Exemplo: enviar para Prometheus, logar, etc.
    fmt.Printf("Token validado: sucesso=%v, motivo=%s\n", success, failureReason)
}
```

E injete no Parse:

```go
recorder := &MeuRecorder{}
payload, err := signet.Parse(tokenBytes, publicKey, signet.WithMetricsRecorder(recorder))
```

- `success` indica se a validação foi bem-sucedida.
- `failureReason` é uma string padronizada (ex: `signet.ReasonTokenExpired`, `signet.ReasonInvalidSignature`).

Assim, você pode criar dashboards, alertas e monitorar a saúde do seu sistema de autenticação em tempo real. 
