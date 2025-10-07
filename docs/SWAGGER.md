# Swagger UI - IaC AI Agent

Este projeto inclui uma interface web Swagger UI para facilitar a interação com a API.

## 🚀 Como Usar

### Gerar Documentação Swagger

```bash
# Opção 1: Usando Makefile
make swagger

# Opção 2: Usando script
./scripts/generate-swagger.sh

# Opção 3: Comando direto
swag init -g cmd/agent/main.go -o docs --parseDependency --parseInternal
```

### Rodar a Aplicação com Swagger

```bash
# Gera docs e inicia a aplicação
make run-swagger

# Ou gere manualmente e depois rode
make swagger
make run
```

### Acessar a Interface

Após iniciar a aplicação, acesse:

**Swagger UI:** http://localhost:8080/swagger/

## 📚 Endpoints Disponíveis

### Health Check
- **GET** `/health` - Verifica status do serviço

### Análise
- **POST** `/analyze` - Analisa código Terraform

### Review
- **POST** `/review` - Executa review de Pull Request

### Informações
- **GET** `/` - Informações da API

## 🎨 Features do Swagger UI

- ✅ **Documentação Interativa** - Teste endpoints diretamente no navegador
- ✅ **Exemplos de Request/Response** - Veja estruturas esperadas
- ✅ **Validação de Schema** - Validação automática de payloads
- ✅ **Try it Out** - Execute requests reais contra a API
- ✅ **Download Specs** - Exporte documentação em JSON/YAML

## 📝 Exemplo de Uso

### 1. Análise de Código Terraform

Acesse: http://localhost:8080/swagger/

1. Clique em **POST /analyze**
2. Clique em **Try it out**
3. Cole o exemplo de request:

```json
{
  "repository": "my-org/my-repo",
  "path": "infrastructure/",
  "content": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-bucket\"\n}"
}
```

4. Clique em **Execute**
5. Veja o resultado com score, análise de segurança e sugestões

### 2. Review de Pull Request

1. Clique em **POST /review**
2. Clique em **Try it out**
3. Cole o exemplo:

```json
{
  "repository": "my-repo",
  "pr_number": 123,
  "owner": "my-org"
}
```

4. Execute e veja o resultado completo do review

## 🔧 Desenvolvimento

### Adicionar Novos Endpoints

1. **Adicione anotações Swagger no handler:**

```go
// HandleNewEndpoint processa nova requisição
// @Summary Descrição curta
// @Description Descrição detalhada do que o endpoint faz
// @Tags nome-da-tag
// @Accept json
// @Produce json
// @Param request body models.MyRequest true "Descrição do request"
// @Success 200 {object} models.MyResponse "Descrição do sucesso"
// @Failure 400 {object} models.ErrorResponse "Erro de validação"
// @Router /novo-endpoint [post]
func (h *Handler) HandleNewEndpoint(w http.ResponseWriter, r *http.Request) {
    // implementação
}
```

2. **Regenere a documentação:**

```bash
make swagger
```

3. **Reinicie a aplicação:**

```bash
make run
```

### Estrutura das Anotações

- `@title` - Título da API
- `@version` - Versão
- `@description` - Descrição geral
- `@host` - Host da API
- `@BasePath` - Caminho base
- `@Summary` - Resumo do endpoint
- `@Description` - Descrição detalhada
- `@Tags` - Agrupa endpoints relacionados
- `@Accept` - Tipo de content aceito
- `@Produce` - Tipo de content retornado
- `@Param` - Parâmetros do endpoint
- `@Success` - Resposta de sucesso
- `@Failure` - Respostas de erro
- `@Router` - Rota e método HTTP

## 🎯 Melhores Práticas

1. **Sempre documente novos endpoints** - Mantenha a documentação atualizada
2. **Use exemplos claros** - Facilita o entendimento
3. **Documente todos os erros possíveis** - Ajuda na depuração
4. **Agrupe endpoints com Tags** - Organiza a UI
5. **Regenere após mudanças** - `make swagger` sempre que alterar handlers

## 🐛 Troubleshooting

### Swagger UI não carrega

```bash
# Verifique se a documentação foi gerada
ls docs/

# Deve ter: docs.go, swagger.json, swagger.yaml
```

### Erro "docs package not found"

```bash
# Regenere a documentação
make swagger

# Reconstrua a aplicação
make build
```

### Mudanças não aparecem

```bash
# Sempre regenere após alterar handlers
make swagger
make run
```

## 📦 Arquivos Gerados

Após rodar `make swagger`, os seguintes arquivos são criados:

```
docs/
├── docs.go          # Documentação em Go (importado pelo main.go)
├── swagger.json     # Spec OpenAPI em JSON
└── swagger.yaml     # Spec OpenAPI em YAML
```

## 🔗 Links Úteis

- [Swagger Documentation](https://swagger.io/docs/)
- [Swaggo GitHub](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)

## 💡 Dicas

- Use **Ctrl+F** no Swagger UI para buscar endpoints rapidamente
- Clique em **"Models"** no final da página para ver todas as estruturas de dados
- Use **"Authorize"** se adicionar autenticação no futuro
- Exporte a spec (swagger.json) para usar em ferramentas como Postman

---

**Precisa de ajuda?** Consulte a [documentação principal](README.md)
