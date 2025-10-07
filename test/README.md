# Testes BDD - IaC AI Agent

Este diretório contém os testes unitários e de integração em BDD (Behavior-Driven Development) para o IaC AI Agent.

## Estrutura de Testes

```
test/
├── unit/                    # Testes unitários BDD
│   ├── suite_test.go       # Suite de testes unitários
│   ├── terraform_analyzer_test.go
│   ├── pr_scorer_test.go
│   ├── checkov_analyzer_test.go
│   └── iam_analyzer_test.go
├── integration/             # Testes de integração BDD
│   ├── suite_test.go       # Suite de testes de integração
│   ├── analysis_service_test.go
│   └── review_service_test.go
└── mocks/                   # Mocks para testes
    └── mocks.go
```

## Frameworks Utilizados

- **Ginkgo v2**: Framework BDD para Go
- **Gomega**: Biblioteca de matchers para assertions

## Executando os Testes

### Todos os Testes

```bash
# Executar todos os testes
make test

# Ou usando go test diretamente
go test ./test/...
```

### Testes Unitários

```bash
# Executar apenas testes unitários
make test-unit

# Ou usando ginkgo
ginkgo ./test/unit/

# Com verbose
ginkgo -v ./test/unit/
```

### Testes de Integração

```bash
# Executar apenas testes de integração
make test-integration

# Ou usando ginkgo
ginkgo ./test/integration/

# Com verbose
ginkgo -v ./test/integration/
```

### Executando Testes Específicos

```bash
# Executar um arquivo de teste específico
ginkgo ./test/unit/terraform_analyzer_test.go

# Executar testes que correspondem a um padrão
ginkgo --focus="TerraformAnalyzer" ./test/unit/

# Pular testes específicos
ginkgo --skip="slow" ./test/...
```

### Com Coverage

```bash
# Executar testes com coverage
make test-coverage

# Ou usando ginkgo
ginkgo -cover -coverprofile=coverage.out ./test/...

# Ver relatório de coverage
go tool cover -html=coverage.out
```

## Instalando Ginkgo CLI

Para executar os testes usando o CLI do Ginkgo:

```bash
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

## Estrutura dos Testes BDD

### Organização Ginkgo

```go
var _ = Describe("ComponentName", func() {
    var (
        // Variáveis de teste
    )

    BeforeEach(func() {
        // Setup antes de cada teste
    })

    AfterEach(func() {
        // Cleanup depois de cada teste
    })

    Describe("Feature Description", func() {
        Context("when specific condition", func() {
            It("should behave in expected way", func() {
                // Test implementation
                Expect(result).To(Equal(expected))
            })
        })
    })
})
```

### Matchers Gomega Comuns

```go
// Igualdade
Expect(actual).To(Equal(expected))
Expect(actual).NotTo(Equal(unexpected))

// Valores booleanos
Expect(actual).To(BeTrue())
Expect(actual).To(BeFalse())

// Nil/Not Nil
Expect(actual).To(BeNil())
Expect(actual).NotTo(BeNil())

// Comparações numéricas
Expect(actual).To(BeNumerically(">", 5))
Expect(actual).To(BeNumerically(">=", 10))
Expect(actual).To(BeNumerically("~", 5.0, 0.1)) // Aproximadamente

// Strings
Expect(actual).To(ContainSubstring("expected"))
Expect(actual).To(HavePrefix("prefix"))
Expect(actual).To(MatchRegexp("pattern"))

// Collections
Expect(slice).To(HaveLen(5))
Expect(slice).To(ContainElement(item))
Expect(slice).To(ContainElements(item1, item2))
Expect(slice).To(BeEmpty())

// Errors
Expect(err).To(HaveOccurred())
Expect(err).NotTo(HaveOccurred())
```

## Testes Unitários

### TerraformAnalyzer

Testa análise de código Terraform:
- Parsing de recursos, módulos, variáveis, outputs
- Detecção de erros de sintaxe
- Verificação de best practices
- Análise de diretórios

### PRScorer

Testa cálculo de scores de qualidade:
- Score de segurança
- Score de best practices
- Score de documentação
- Score de maintainability
- Score de performance
- Média ponderada

### CheckovAnalyzer

Testa integração com Checkov:
- Conversão de resultados
- Determinação de severidade
- Geração de recomendações
- Contabilização de issues

### IAMAnalyzer

Testa análise de IAM:
- Detecção de políticas permissivas
- Identificação de acesso público
- Análise de roles
- Geração de recomendações

## Testes de Integração

### AnalysisService

Testa orquestração completa de análise:
- Análise de conteúdo Terraform
- Análise de diretórios
- Integração entre analyzers
- Geração de sugestões
- Cálculo de scores
- Validação de análises

### ReviewService

Testa processo de review de PRs:
- Review de PRs
- Review de arquivos
- Determinação de status
- Geração de comentários
- Aprovação automática
- Integração com AnalysisService

## Ambiente de Testes

### Variáveis de Ambiente

```bash
# Log level para testes
export LOG_LEVEL=info

# Desabilitar Checkov em testes (se não instalado)
export CHECKOV_ENABLED=false
```

### Pré-requisitos Opcionais

- **Checkov**: Para testes de integração com análise de segurança
  ```bash
  pip install checkov
  ```

## Continuous Integration

Os testes podem ser executados em pipelines CI/CD:

```yaml
# Exemplo GitHub Actions
- name: Run Tests
  run: |
    go test -v ./test/...
    
- name: Run Tests with Coverage
  run: |
    go test -v -coverprofile=coverage.out ./test/...
    go tool cover -func=coverage.out
```

## Debugging Testes

### Modo Verbose

```bash
ginkgo -v ./test/unit/
```

### Foco em Teste Específico

```bash
# Adicione FIt, FDescribe, FContext no código
FIt("should test specific behavior", func() {
    // test
})
```

### Pular Testes Temporariamente

```bash
# Adicione XIt, XDescribe, XContext no código
XIt("should test something (pending)", func() {
    // test
})
```

### Debug com Delve

```bash
dlv test ./test/unit/
```

## Boas Práticas

1. **Nomes Descritivos**: Use nomes que descrevem o comportamento esperado
2. **Um Assertion por It**: Quando possível, teste um comportamento por vez
3. **Setup/Cleanup**: Use BeforeEach/AfterEach para preparar e limpar
4. **Independência**: Testes não devem depender uns dos outros
5. **Determinístico**: Testes devem produzir mesmo resultado sempre
6. **Rápidos**: Testes unitários devem ser muito rápidos
7. **Isolados**: Use mocks para dependências externas

## Cobertura de Testes

Meta de cobertura: **80%+**

Para visualizar cobertura:

```bash
make test-coverage
go tool cover -html=coverage.out -o coverage.html
open coverage.html  # macOS
xdg-open coverage.html  # Linux
```

## Troubleshooting

### Testes Falhando

1. Verifique as dependências: `go mod tidy`
2. Limpe o cache: `go clean -testcache`
3. Execute com verbose: `ginkgo -v`
4. Verifique logs: `LOG_LEVEL=debug ginkgo -v`

### Ginkgo não Encontrado

```bash
go install github.com/onsi/ginkgo/v2/ginkgo@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

## Recursos

- [Ginkgo Documentation](https://onsi.github.io/ginkgo/)
- [Gomega Documentation](https://onsi.github.io/gomega/)
- [Go Testing Package](https://golang.org/pkg/testing/)

