# 🧪 Guia de Testes - IaC AI Agent

## 📋 Visão Geral

O projeto possui **testes obrigatórios** que são executados automaticamente antes de cada push.

```
🔒 Pre-Push Hook Ativado
├─ 📋 Linter (golangci-lint)
├─ 🧪 Testes Unitários
├─ 🔗 Testes de Integração
└─ 🏗️  Build
```

---

## 🚀 Quick Start

### Setup Inicial (Uma vez)

```bash
# Instala dependências e git hooks
make setup
```

Isso irá:
- ✅ Baixar dependências Go
- ✅ Instalar ferramentas de teste (Ginkgo)
- ✅ Instalar git hooks (pre-push)
- ✅ Configurar linter

### Após Setup

```bash
# Seus próximos pushes executarão testes automaticamente!
git push origin main

# Você verá:
🔍 PRE-PUSH: Executando Testes
📋 [1/4] Executando Linter...
🧪 [2/4] Executando Testes Unitários...
🔗 [3/4] Executando Testes de Integração...
🏗️  [4/4] Verificando Build...
🎉 TODOS OS TESTES PASSARAM!
```

---

## 🧪 Tipos de Testes

### 1. Testes Unitários

**Localização**: `test/unit/`

**O que testa**:
- Funções individuais
- Lógica de negócio isolada
- Sem dependências externas (mocks)

**Como executar**:
```bash
make test-unit
```

**Arquivos existentes**:
- `checkov_analyzer_test.go`
- `iam_analyzer_test.go`
- `terraform_analyzer_test.go`
- `pr_scorer_test.go`
- `validation_test.go`

### 2. Testes de Integração

**Localização**: `test/integration/`

**O que testa**:
- Múltiplos componentes trabalhando juntos
- Fluxos completos
- Interações entre serviços

**Como executar**:
```bash
make test-integration
```

**Arquivos existentes**:
- `analysis_test.go`
- `analysis_service_test.go`
- `review_service_test.go`

### 3. Testes BDD (Behavior Driven Development)

**Localização**: `test/bdd/features/`

**O que testa**:
- Comportamento do usuário
- Fluxos end-to-end
- Cenários de negócio

**Arquivos**:
- `user_onboarding.feature` (3 scenarios)
- `nft_purchase.feature` (3 scenarios)
- `token_purchase.feature` (4 scenarios)
- `bot_analysis.feature` (5 scenarios)

**Como executar** (quando implementado):
```bash
make test-bdd
```

---

## 📝 Comandos Disponíveis

### Testes

```bash
# Todos os testes
make test

# Apenas unitários
make test-unit

# Apenas integração
make test-integration

# BDD (Ginkgo)
make test-bdd

# Testes rápidos (sem race detector)
make test-quick

# Com coverage
make test-coverage

# Coverage visual (abre navegador)
make test-coverage
```

### Linter

```bash
# Executar linter
make lint

# Instalar linter
make lint-install

# Auto-fix issues
make lint-fix

# Formatar código
make fmt
```

### Build

```bash
# Build
make build

# Build e executar
make run

# Limpar build
make clean
```

---

## 🪝 Git Hooks

### Pre-Push Hook

**Arquivo**: `.githooks/pre-push`

**O que faz**:
1. Executa linter
2. Executa testes unitários
3. Executa testes de integração
4. Verifica build

**Se qualquer etapa falhar**:
```
❌ Push bloqueado!
💡 Corrija os erros antes de tentar novamente
```

### Como Desabilitar (Temporariamente)

```bash
# Pular testes (use com cuidado!)
git push --no-verify
```

⚠️ **Aviso**: Só use `--no-verify` em casos excepcionais!

### Como Reinstalar Hooks

```bash
make git-hooks
```

---

## ✅ Escrevendo Bons Testes

### Testes Unitários

```go
package unit_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestMyFunction(t *testing.T) {
    // Arrange
    input := "test"
    expected := "TEST"
    
    // Act
    result := MyFunction(input)
    
    // Assert
    assert.Equal(t, expected, result)
}

// Table-driven tests (recomendado)
func TestMyFunction_TableDriven(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"uppercase", "test", "TEST"},
        {"empty", "", ""},
        {"special chars", "test!", "TEST!"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := MyFunction(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Testes de Integração

```go
package integration_test

import (
    "testing"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Integration Suite")
}

var _ = Describe("Analysis Service", func() {
    var (
        service *AnalysisService
    )
    
    BeforeEach(func() {
        service = NewAnalysisService()
    })
    
    Context("quando analisa código Terraform", func() {
        It("deve retornar resultados válidos", func() {
            code := `resource "aws_s3_bucket" "test" {}`
            result, err := service.Analyze(code)
            
            Expect(err).ToNot(HaveOccurred())
            Expect(result).ToNot(BeNil())
        })
    })
})
```

---

## 📊 Coverage

### Ver Coverage

```bash
# Gerar coverage
make test-coverage

# Abrirá navegador com relatório visual
```

### Meta de Coverage

```
Target: 80% coverage

Atual:
- Unit tests:        ~60%
- Integration tests: ~40%
- Total:            ~50%
```

---

## 🚨 Troubleshooting

### Erro: "golangci-lint not found"

```bash
# Instalar
make lint-install

# Ou manualmente
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
```

### Erro: "ginkgo not found"

```bash
# Instalar
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

### Erro: "Tests failing on push"

```bash
# Executar testes localmente primeiro
make test-unit
make test-integration

# Ver detalhes dos erros
go test -v ./test/unit/...
```

### Erro: "Permission denied: .githooks/pre-push"

```bash
# Dar permissão
chmod +x .githooks/pre-push
```

---

## 🎯 Boas Práticas

### 1. Teste Antes de Commitar

```bash
# Sempre execute antes de commitar
make test
make lint
```

### 2. Escreva Testes Para Código Novo

```
Para cada nova feature:
☐ Escrever testes unitários
☐ Escrever testes de integração (se aplicável)
☐ Atualizar BDD scenarios (se aplicável)
```

### 3. Use Table-Driven Tests

```go
// ✅ BOM
tests := []struct{
    name string
    input string
    expected string
}{
    {"case1", "input1", "output1"},
    {"case2", "input2", "output2"},
}

// ❌ EVITE
func TestCase1() {}
func TestCase2() {}
func TestCase3() {}
```

### 4. Mock Dependências Externas

```go
// Use interfaces para facilitar mocks
type LLMClient interface {
    Generate(req *LLMRequest) (*LLMResponse, error)
}

// Em testes, use mock
type MockLLMClient struct {
    mock.Mock
}

func (m *MockLLMClient) Generate(req *LLMRequest) (*LLMResponse, error) {
    args := m.Called(req)
    return args.Get(0).(*LLMResponse), args.Error(1)
}
```

### 5. Cleanup Após Testes

```go
func TestWithCleanup(t *testing.T) {
    // Setup
    tmpFile := createTempFile()
    
    // Cleanup
    defer os.Remove(tmpFile)
    
    // Test
    // ...
}
```

---

## 📈 Roadmap de Testes

### Fase 1: Básico (Atual)
- ✅ Estrutura de testes criada
- ✅ Pre-push hook implementado
- ✅ Testes unitários existentes
- ✅ Testes integração existentes

### Fase 2: Expandir Coverage (2 semanas)
- ⏳ Adicionar testes para novos componentes
- ⏳ Aumentar coverage para 80%+
- ⏳ Implementar BDD step definitions
- ⏳ CI/CD integration

### Fase 3: Avançado (1 mês)
- ⏳ Load tests (k6)
- ⏳ E2E tests (Playwright)
- ⏳ Performance benchmarks
- ⏳ Mutation testing

---

## 🔗 Links Úteis

- [Go Testing](https://golang.org/pkg/testing/)
- [Ginkgo BDD](https://onsi.github.io/ginkgo/)
- [Gomega Matchers](https://onsi.github.io/gomega/)
- [Testify](https://github.com/stretchr/testify)
- [GolangCI-Lint](https://golangci-lint.run/)

---

## 📞 Suporte

**Problemas com testes?**
- Veja logs: `make test-unit -v`
- CI failing? Check GitHub Actions
- Issues: GitHub Issues

---

**Última Atualização**: 2025-01-15  
**Versão**: 1.0.0  
**Status**: ✅ Pre-Push Hooks Ativados
