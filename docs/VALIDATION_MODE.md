# Modo de Validação - Análise de Resultados Pré-existentes

## Visão Geral

O IaC AI Agent agora suporta um **modo de validação** onde você pode fornecer resultados de análises já executadas externamente, ao invés de executar as ferramentas (Checkov, Terraform, etc.) novamente.

Este modo é útil quando:
- Você já executou o Checkov em seu pipeline CI/CD e quer apenas validar e processar os resultados
- Quer economizar tempo evitando re-execução de ferramentas pesadas
- Precisa validar consistência de resultados de análises externas
- Quer apenas calcular scores e gerar sugestões baseadas em dados existentes

## Como Funciona

### 1. Validação de Resultado Checkov

O `CheckovAnalyzer` agora possui um método `ValidateAndParseResult` que:
- Recebe um JSON de resultado Checkov já executado
- Valida a estrutura e consistência dos dados
- Converte para o modelo interno `SecurityAnalysis`
- Retorna erros se os dados estiverem inválidos

#### Exemplo de Uso

```go
import (
    "encoding/json"
    "io/ioutil"
    
    "github.com/gosouza/iac-ai-agent/internal/agent/analyzer"
    "github.com/gosouza/iac-ai-agent/pkg/logger"
)

// Cria o analyzer
log := logger.NewLogger("info")
checkovAnalyzer := analyzer.NewCheckovAnalyzer(log)

// Lê resultado Checkov de um arquivo
checkovJSON, err := ioutil.ReadFile("checkov-results.json")
if err != nil {
    log.Fatal("Erro ao ler arquivo", "error", err)
}

// Valida e converte o resultado
securityAnalysis, err := checkovAnalyzer.ValidateAndParseResult(checkovJSON)
if err != nil {
    log.Error("Resultado Checkov inválido", "error", err)
    return
}

log.Info("Resultado validado com sucesso",
    "passed", securityAnalysis.ChecksPassed,
    "failed", securityAnalysis.ChecksFailed,
    "total_issues", securityAnalysis.TotalIssues)
```

### 2. Validação Completa de Análise

O `AnalysisService` possui o método `ValidatePreExistingResults` que:
- Aceita resultados Checkov (JSON) e análise Terraform
- Valida ambos os resultados
- Executa análise IAM nos dados Terraform fornecidos
- Calcula score de qualidade
- Gera sugestões de melhorias
- **NÃO executa nenhuma ferramenta externa**

#### Exemplo de Uso

```go
import (
    "encoding/json"
    "io/ioutil"
    
    "github.com/gosouza/iac-ai-agent/internal/models"
    "github.com/gosouza/iac-ai-agent/internal/services"
    "github.com/gosouza/iac-ai-agent/pkg/logger"
)

// Cria o serviço
log := logger.NewLogger("info")
analysisService := services.NewAnalysisService(log, 70) // minPassScore = 70

// Lê resultado Checkov
checkovJSON, _ := ioutil.ReadFile("checkov-results.json")

// Cria análise Terraform (ou lê de outro lugar)
tfAnalysis := &models.TerraformAnalysis{
    TotalResources: 10,
    TotalModules:   2,
    TotalVariables: 5,
    TotalOutputs:   3,
    Resources: []models.TerraformResource{
        {
            Type: "aws_s3_bucket",
            Name: "my-bucket",
            File: "main.tf",
            Line: 10,
        },
        // ... outros recursos
    },
    Variables: []models.TerraformVariable{
        {
            Name:        "bucket_name",
            Description: "Name of the S3 bucket",
            Type:        "string",
        },
        // ... outras variáveis
    },
    // ... outros dados
}

// Valida resultados pré-existentes
response, err := analysisService.ValidatePreExistingResults(checkovJSON, tfAnalysis)
if err != nil {
    log.Error("Erro ao validar resultados", "error", err)
    return
}

// Processa resposta
log.Info("Validação concluída",
    "score", response.Score,
    "is_approved", response.Metadata["is_approved"],
    "score_level", response.Metadata["score_level"],
    "suggestions", len(response.Suggestions))

// Exibe score
fmt.Printf("Score Total: %d/100\n", response.Score)
fmt.Printf("Nível: %s\n", response.Metadata["score_level"])
fmt.Printf("Aprovado: %v\n", response.Metadata["is_approved"])
fmt.Printf("\n%s\n", response.Metadata["score_summary"])
```

## Formato de Entrada

### Resultado Checkov (JSON)

O JSON do Checkov deve seguir o formato padrão:

```json
{
  "summary": {
    "passed": 15,
    "failed": 3,
    "skipped": 0,
    "parsing_errors": 0,
    "resource_count": 18
  },
  "results": {
    "passed_checks": [],
    "failed_checks": [
      {
        "check_id": "CKV_AWS_1",
        "check_name": "Ensure S3 bucket has encryption enabled",
        "description": "S3 bucket should have encryption enabled",
        "severity": "HIGH",
        "resource": "aws_s3_bucket.example",
        "file_path": "main.tf",
        "file_line_range": [10, 15],
        "guideline": "https://docs.bridgecrew.io/docs/s3_14"
      }
    ],
    "skipped_checks": []
  }
}
```

### Análise Terraform

A estrutura `TerraformAnalysis` deve conter:

```go
&models.TerraformAnalysis{
    TotalResources: 10,
    TotalModules:   2,
    TotalVariables: 5,
    TotalOutputs:   3,
    Resources: []models.TerraformResource{...},
    Variables: []models.TerraformVariable{...},
    Outputs: []models.TerraformOutput{...},
    Modules: []models.TerraformModule{...},
    BestPracticeWarnings: []string{...},
    SyntaxErrors: []string{...},
}
```

## Validações Realizadas

### Validação Checkov

- ✅ JSON bem formatado
- ✅ Estrutura `summary` presente
- ✅ Números de checks não negativos
- ✅ Consistência entre `summary.failed` e número de `failed_checks`
- ✅ Cada check possui `check_id` não vazio
- ⚠️ Warning se check não possui nome

### Validação Terraform

- ✅ Análise não nula
- ✅ `TotalResources` não negativo
- ✅ `TotalModules` não negativo
- ⚠️ Warning se houver inconsistência entre `TotalResources` e tamanho de `Resources`

## Diferenças do Modo Normal

| Aspecto | Modo Normal | Modo Validação |
|---------|-------------|----------------|
| Executa Checkov | ✅ Sim | ❌ Não |
| Executa Terraform Parser | ✅ Sim | ❌ Não |
| Valida JSON/Estruturas | ⚠️ Parcial | ✅ Completo |
| Análise IAM | ✅ Sim | ✅ Sim |
| Calcula Score | ✅ Sim | ✅ Sim |
| Gera Sugestões | ✅ Sim | ✅ Sim |
| Velocidade | 🐢 Lento | ⚡ Rápido |
| Metadata | - | `validation_mode: "pre_existing_results"` |

## Casos de Uso

### 1. Pipeline CI/CD

```bash
# 1. Execute Checkov no pipeline
checkov -d . -o json > checkov-results.json

# 2. Execute análise Terraform
terraform show -json > tfstate.json

# 3. Valide resultados no IaC AI Agent
curl -X POST http://localhost:8080/api/validate \
  -H "Content-Type: application/json" \
  -d @validation-request.json
```

### 2. Análise Offline

Quando você tem resultados salvos e quer reprocessá-los:

```go
// Lê resultados salvos
checkovData := loadCheckovResults("./saved-results/checkov.json")
tfData := loadTerraformAnalysis("./saved-results/terraform.json")

// Reprocessa sem re-executar ferramentas
response, _ := analysisService.ValidatePreExistingResults(checkovData, tfData)
```

### 3. Agregação de Múltiplas Análises

```go
var allSuggestions []models.Suggestion
var totalScore int

for _, dir := range directories {
    checkovData := loadCheckovForDirectory(dir)
    tfData := loadTerraformForDirectory(dir)
    
    response, _ := analysisService.ValidatePreExistingResults(checkovData, tfData)
    
    allSuggestions = append(allSuggestions, response.Suggestions...)
    totalScore += response.Score
}

avgScore := totalScore / len(directories)
```

## Limitações

1. **Não detecta novos problemas**: Apenas valida dados já coletados
2. **Dependente de formato**: Requer formato exato do Checkov
3. **Sem contexto de código**: Não tem acesso ao código fonte para análises mais profundas
4. **IAM limitado**: Análise IAM depende dos dados Terraform fornecidos

## Benefícios

- ⚡ **Performance**: Não executa ferramentas pesadas
- 🔄 **Reutilização**: Aproveita análises já existentes
- 💰 **Economia**: Reduz tempo de CI/CD
- 🔍 **Consistência**: Valida e normaliza resultados de diferentes fontes
- 📊 **Score unificado**: Calcula métricas consistentes independente da origem dos dados

## Próximos Passos

- Adicionar suporte a outros formatos (tfsec, terrascan, etc.)
- Validação de análises CloudFormation
- API REST endpoint dedicado para modo de validação
- Cache inteligente de resultados validados

