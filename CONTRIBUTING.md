# 🤝 Guia de Contribuição — signet-go

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Obrigado por considerar contribuir com o signet-go! Este documento fornece diretrizes para contribuições.

## 🎯 Como Contribuir

### 📋 Tipos de Contribuição

- **🐛 Bug fixes**: Correções de bugs e problemas
- **✨ Novas funcionalidades**: Adições à API
- **📚 Documentação**: Melhorias na documentação
- **🧪 Testes**: Adição ou melhoria de testes
- **🔧 Exemplos**: Novos exemplos de uso
- **⚡ Performance**: Otimizações de performance
- **🔒 Segurança**: Melhorias de segurança

### 🚀 Processo de Contribuição

1. **Fork o repositório**
2. **Crie uma branch** para sua feature/fix
3. **Faça suas mudanças** seguindo as diretrizes
4. **Adicione/atualize testes** se necessário
5. **Atualize documentação** se necessário
6. **Commit suas mudanças** com mensagens claras
7. **Abra um Pull Request**

## 🛠️ Ambiente de Desenvolvimento

### Pré-requisitos

- **Go 1.21+**: Versão mínima suportada
- **Git**: Para controle de versão
- **Make** (opcional): Para comandos de build

### Setup Local

```bash
# Clone o repositório
git clone https://github.com/lucas-de-lima/signet-go.git
cd signet-go

# Instale dependências
go mod download

# Execute testes
go test ./...

# Execute linting
go vet ./...
```

## 📝 Diretrizes de Código

### 🐹 Padrões Go

- Siga as [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` para formatação
- Execute `go vet` antes de commitar
- Mantenha cobertura de testes alta (>90%)

### 🔒 Segurança

- **Nunca commite chaves privadas** ou secrets
- **Valide todas as entradas** de usuário
- **Use constantes de tempo** para comparações criptográficas
- **Documente vulnerabilidades** encontradas

### 📚 Documentação

- **Documente funções públicas** com GoDoc
- **Atualize READMEs** quando necessário
- **Adicione exemplos** para novas funcionalidades
- **Mantenha documentação sincronizada** com o código

### 🧪 Testes

```bash
# Execute todos os testes
go test ./...

# Execute testes com cobertura
go test -cover ./...

# Execute testes de benchmark
go test -bench=. ./...

# Execute testes de race condition
go test -race ./...
```

## 🏗️ Estrutura do Projeto

```
signet-go/
├── signet/              # 📦 Pacote principal
├── grpcinterceptor/     # 🔌 Interceptor gRPC
├── internal/           # 🔒 Código interno
├── proto/              # 📋 Definições protobuf
├── examples/           # 🔧 Exemplos de uso
└── docs/              # 📚 Documentação
```

## 📋 Checklist de Pull Request

### Antes de Abrir o PR

- [ ] **Código compila** sem erros
- [ ] **Testes passam** localmente
- [ ] **Linting passa** (`go vet`, `golint`)
- [ ] **Documentação atualizada** se necessário
- [ ] **Exemplos funcionam** se adicionados
- [ ] **Commits organizados** com mensagens claras

### Template de Pull Request

```markdown
## 📝 Descrição
Breve descrição das mudanças

## 🎯 Tipo de Mudança
- [ ] Bug fix
- [ ] Nova funcionalidade
- [ ] Breaking change
- [ ] Documentação

## 🧪 Testes
- [ ] Testes unitários adicionados/atualizados
- [ ] Testes de integração passam
- [ ] Exemplos testados

## 📚 Documentação
- [ ] GoDoc atualizado
- [ ] README atualizado se necessário
- [ ] Exemplos atualizados se necessário

## 🔒 Segurança
- [ ] Mudanças de segurança documentadas
- [ ] Vulnerabilidades conhecidas listadas
```

## 🚨 Reportando Bugs

### Template de Bug Report

```markdown
## 🐛 Descrição do Bug
Descrição clara e concisa do bug

## 🔄 Passos para Reproduzir
1. Vá para '...'
2. Clique em '...'
3. Veja o erro

## ✅ Comportamento Esperado
O que deveria acontecer

## 📱 Comportamento Atual
O que realmente acontece

## 🖥️ Ambiente
- OS: [ex: Windows 10]
- Go Version: [ex: 1.21.0]
- signet-go Version: [ex: v1.0.0]

## 📋 Informações Adicionais
Logs, screenshots, etc.
```

## 🏷️ Versionamento

Seguimos [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: Novas funcionalidades (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

## 📞 Comunicação

- **Issues**: Para bugs e feature requests
- **Discussions**: Para perguntas e discussões
- **Security**: Para vulnerabilidades (privado) - dev.lucasdelima@gmail.com
- **Email**: dev.lucasdelima@gmail.com
- **LinkedIn**: [dev-lucasdelima](https://www.linkedin.com/in/dev-lucasdelima/)

## 🎉 Reconhecimento

Contribuidores serão listados no README e releases. Obrigado por ajudar a tornar o signet-go melhor!

---

> **💡 Sua contribuição é valiosa para a comunidade! Juntos, podemos criar uma biblioteca de segurança ainda mais robusta e confiável.** 