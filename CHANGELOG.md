# Changelog

## [Unreleased]

### Added - Modo de Validação (2025-10-06)

#### Novos Métodos de Validação

**CheckovAnalyzer (`internal/agent/analyzer/checkov.go`)**
- `ValidateAndParseResult(jsonResult []byte)`: Valida e faz parse de resultado Checkov já executado
- `validateCheckovResult(result *models.CheckovResult)`: Valida estrutura e consistência dos dados Checkov
  - Verifica se números de checks não são negativos
  - Valida presença de CheckID em todos os checks
  - Alerta sobre inconsistências entre summary e lista de checks

**AnalysisService (`internal/services/analysis.go`)**
- `ValidatePreExistingResults(checkovJSON []byte, tfAnalysis *models.TerraformAnalysis)`: Valida resultados de análises já executadas externamente
  - Processa resultados Checkov sem executar a ferramenta
  - Valida análise Terraform fornecida
  - Executa análise IAM nos dados Terraform
  - Calcula scores e gera sugestões
  - Adiciona metadata `validation_mode: "pre_existing_results"`
- `validateTerraformAnalysis(analysis *models.TerraformAnalysis)`: Valida estrutura de análise Terraform
  - Verifica se totais não são negativos
  - Alerta sobre inconsistências entre contadores e listas

#### Benefícios
- ⚡ **Performance**: Não executa ferramentas externas pesadas
- 🔄 **Reutilização**: Aproveita análises já executadas em CI/CD
- 💰 **Economia**: Reduz tempo de pipeline
- 🔍 **Validação**: Garante consistência de dados externos
- 📊 **Score Unificado**: Calcula métricas independente da origem

#### Documentação
- `docs/VALIDATION_MODE.md`: Documentação completa do modo de validação
  - Exemplos de uso
  - Formatos de entrada
  - Validações realizadas
  - Casos de uso
  - Limitações e benefícios

#### Testes
- `test/unit/validation_test.go`: Suite completa de testes
  - ✅ Validação de resultado Checkov válido
  - ✅ Tratamento de JSON inválido
  - ✅ Tratamento de dados inválidos
  - ✅ Validação completa com ambos os resultados
  - ✅ Validação apenas com Terraform (sem Checkov)

### Fixed

**Correções de Compatibilidade**
- Corrigido erro de type assertion em `terraform.go` (linha 136)
- Corrigido parâmetros de `NewPRScorer()` para não receber argumentos
- Atualizado uso de `CalculateScore()` para receber `*models.AnalysisDetails`
- Corrigido métodos `IsApproved` → `ShouldApprove` e `GetRecommendation` removido
- Adicionado campo `minPassScore` ao `AnalysisService`
- Corrigido imports faltantes (`bytes` em `webhook/handlers.go`)
- Corrigido arquivos placeholder sem declaração de package

**API Handlers**
- `api/rest/handlers.go`: Atualizado construtor para nova API
- `internal/platform/webhook/handlers.go`: Atualizado construtor e imports
- `internal/services/analysis.go`: Adicionado método wrapper `Analyze()`

### Changed

**AnalysisService**
- Agora mantém `minPassScore` como campo interno
- `NewAnalysisService` aceita `(log *logger.Logger, minPassScore int)`
- Todos os métodos de análise agora retornam metadata adicional:
  - `score_level`: Nível do score (Excelente, Bom, Regular, etc.)
  - `score_summary`: Resumo textual formatado do score
  - `is_approved`: Boolean indicando se PR seria aprovado

**Metadata em Responses**
- Adicionado `validation_mode` quando usando `ValidatePreExistingResults`
- Score breakdown detalhado incluído em todas as respostas

## Uso

### Modo Normal (executa ferramentas)
```go
log := logger.New("info", "json")
analysisService := services.NewAnalysisService(log, 70)

// Analisa diretório (executa Checkov, Terraform, etc.)
response, err := analysisService.AnalyzeDirectory("./terraform")
```

### Modo Validação (apenas valida resultados)
```go
log := logger.New("info", "json")
analysisService := services.NewAnalysisService(log, 70)

// Lê resultados já executados
checkovJSON, _ := ioutil.ReadFile("checkov-results.json")
tfAnalysis := loadTerraformAnalysis() // de alguma fonte

// Valida sem re-executar
response, err := analysisService.ValidatePreExistingResults(checkovJSON, tfAnalysis)
```

## Compatibilidade

- ✅ Go 1.21+
- ✅ Backward compatible com análises existentes
- ✅ Novos métodos não quebram API existente
- ✅ Todos os testes passando (5/5 specs)

## Próximos Passos

- [ ] Endpoint REST dedicado para `/api/validate`
- [ ] Suporte a outros formatos (tfsec, terrascan)
- [ ] Cache de resultados validados
- [ ] Validação de análises CloudFormation
- [ ] Agregação de múltiplas análises

