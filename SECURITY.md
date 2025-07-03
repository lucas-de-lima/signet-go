# Política de Segurança

## 🛡️ Reportando uma Vulnerabilidade

Agradecemos que você queira reportar uma vulnerabilidade de segurança. Sua contribuição é importante para manter o signet-go seguro para todos.

### 📧 Como Reportar

**NÃO** abra uma issue pública para vulnerabilidades de segurança. Em vez disso:

1. **Envie um email** para dev.lucasdelima@gmail.com
2. **Use o assunto**: `[SECURITY] signet-go - [Descrição breve]`
3. **Inclua detalhes** sobre a vulnerabilidade
4. **Aguarde resposta** da equipe de segurança

### 📋 Informações Necessárias

Para ajudar na investigação, inclua:

- **Descrição detalhada** da vulnerabilidade
- **Passos para reproduzir** o problema
- **Impacto potencial** da vulnerabilidade
- **Sugestões de correção** (se aplicável)
- **Informações do ambiente** (versão Go, OS, etc.)

### ⏱️ Processo de Resposta

1. **Confirmação**: Você receberá confirmação em 48 horas
2. **Investigação**: A equipe investigará a vulnerabilidade
3. **Correção**: Desenvolveremos e testaremos uma correção
4. **Disclosure**: Coordenaremos a divulgação pública
5. **Release**: Publicaremos uma versão corrigida

## 🔒 Práticas de Segurança

### Criptografia
- **Ed25519**: Algoritmo padrão para assinaturas
- **Validação temporal**: Verificação obrigatória de exp/iat
- **Integridade**: Verificação obrigatória de assinatura
- **Constantes de tempo**: Comparações seguras contra timing attacks

### Validação de Entrada
- **Claims obrigatórios**: Validação rigorosa de sub, aud, exp, iat
- **Claims opcionais**: Validação quando presentes
- **Sanitização**: Limpeza de dados de entrada
- **Limites**: Proteção contra payloads maliciosos

### Gerenciamento de Chaves
- **Rotação**: Suporte a múltiplas chaves via KeyResolver
- **Validação**: Verificação de formato e tamanho de chaves
- **Cache seguro**: TTL para mitigar ataques DoS
- **Fallback**: Suporte a tokens antigos sem `kid`

### Observabilidade
- **Métricas**: Monitoramento de falhas de validação
- **Logs seguros**: Sem exposição de dados sensíveis
- **Tracing**: Propagação de contexto para debugging
- **Alertas**: Detecção de anomalias

## 🚨 Vulnerabilidades Conhecidas

### Nenhuma vulnerabilidade conhecida atualmente

Se você descobrir uma vulnerabilidade, siga o processo de reporte acima.

## 📚 Recursos de Segurança

### Documentação
- [Especificação Signet](SPECIFICATION-v1.0.md): Padrão de segurança
- [GoDoc Reference](GODOC-REFERENCE.md): Documentação da API
- [Exemplos de Produção](examples/): Implementações seguras

### Ferramentas
- **go vet**: Análise estática de código
- **gosec**: Análise de segurança específica para Go
- **CodeQL**: Análise de segurança avançada (GitHub)

### Boas Práticas
- **Princípio do menor privilégio**: Tokens com claims mínimos
- **Validação rigorosa**: Sempre validar claims obrigatórios
- **Rotação de chaves**: Implementar KeyResolver com cache
- **Monitoramento**: Usar MetricsRecorder para observabilidade
- **Atualizações**: Manter dependências atualizadas

## 🤝 Coordenação de Vulnerabilidades

### Equipe de Segurança
- **Mantenedores**: Revisão e correção de vulnerabilidades
- **Contribuidores**: Reporte e teste de correções
- **Comunidade**: Feedback e validação

### Processo de Disclosure
1. **Responsible disclosure**: Coordenação com reportador
2. **CVE assignment**: Solicitação de CVE quando apropriado
3. **Public disclosure**: Divulgação coordenada
4. **Patch release**: Lançamento de versão corrigida

---

> **💡 A segurança é uma responsabilidade compartilhada. Obrigado por ajudar a manter o signet-go seguro!** 