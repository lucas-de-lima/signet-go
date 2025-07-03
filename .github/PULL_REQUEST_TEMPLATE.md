---
name: Pull Request Template
about: Template padrão para pull requests
title: '[TIPO] Descrição breve da mudança'
labels: ['triage']
assignees: ''
reviewers: ''
---

## 📝 Descrição
<!-- Descreva as mudanças feitas e o motivo -->

## 🎯 Tipo de Mudança
- [ ] 🐛 Bug fix (mudança que corrige um problema)
- [ ] ✨ Nova funcionalidade (mudança que adiciona funcionalidade)
- [ ] 💥 Breaking change (correção ou funcionalidade que faria com que a funcionalidade existente não funcionasse como esperado)
- [ ] 📚 Documentação (mudanças na documentação)
- [ ] 🔧 Refatoração (mudança que não corrige um bug nem adiciona uma funcionalidade)
- [ ] ⚡ Performance (mudança que melhora o desempenho)
- [ ] 🧪 Teste (adicionando testes ausentes ou corrigindo testes existentes)
- [ ] 🔒 Segurança (mudança relacionada à segurança)

## 🔗 Issues Relacionadas
<!-- Fecha issues automaticamente -->
Closes #(issue)

## 🧪 Testes
- [ ] ✅ Testes unitários adicionados/atualizados
- [ ] ✅ Testes de integração passam
- [ ] ✅ Testes de benchmark executados (se aplicável)
- [ ] ✅ Testes de race condition executados (`go test -race`)
- [ ] ✅ Exemplos testados (se aplicável)

## 📚 Documentação
- [ ] ✅ GoDoc atualizado para funções públicas
- [ ] ✅ README atualizado (se necessário)
- [ ] ✅ Exemplos atualizados (se necessário)
- [ ] ✅ CHANGELOG atualizado (se aplicável)

## 🔒 Segurança
- [ ] ✅ Mudanças de segurança documentadas
- [ ] ✅ Vulnerabilidades conhecidas listadas
- [ ] ✅ Chaves privadas/secrets não foram commitados
- [ ] ✅ Validação de entrada implementada (se aplicável)

## 🏗️ Mudanças Técnicas
<!-- Descreva mudanças arquiteturais ou técnicas importantes -->

## 📊 Impacto
<!-- Descreva o impacto das mudanças -->

### Compatibilidade
- [ ] ✅ Backward compatible
- [ ] ❌ Breaking changes (documentados)

### Performance
- [ ] ✅ Sem impacto na performance
- [ ] ⚡ Melhoria na performance
- [ ] 📉 Degradação na performance (justificada)

## 🔍 Checklist de Qualidade
- [ ] ✅ Código segue os padrões Go
- [ ] ✅ `go fmt` executado
- [ ] ✅ `go vet` executado
- [ ] ✅ `golint` executado (se disponível)
- [ ] ✅ Cobertura de testes mantida/improved
- [ ] ✅ Commits organizados com mensagens claras
- [ ] ✅ Branch atualizada com main/master

## 📋 Informações Adicionais
<!-- Qualquer informação adicional que possa ajudar os revisores -->

## 🎯 Próximos Passos
<!-- Se aplicável, descreva próximos passos ou dependências -->

---
*Por favor, preencha todas as seções relevantes. Pull requests sem informações suficientes podem ser rejeitados.* 