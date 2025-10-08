# 📝 Integração Notion - Guia Completo

## 📋 Visão Geral

A integração com Notion permite que o IaC AI Agent gerencie agentes inteligentes diretamente na plataforma Notion, oferecendo:

- **Criação automática de agentes** na inicialização da aplicação
- **Gerenciamento centralizado** de agentes via API Notion
- **Sincronização** entre o sistema local e Notion
- **Validação automática** de configuração na startup

## 🚀 Funcionalidades Implementadas

### 1. **Cliente Notion** (`internal/platform/notion/client.go`)
- Comunicação completa com API Notion
- Operações CRUD para agentes
- Tratamento de erros e rate limiting
- Verificação de disponibilidade do serviço

### 2. **Serviço de Agentes** (`internal/services/notion_agent_service.go`)
- Gerenciamento de agentes Notion
- Criação automática de agente padrão
- Validação de configuração
- Interface unificada para operações

### 3. **Validação de Startup** (`internal/startup/validator.go`)
- Verificação automática na inicialização
- Criação de agente se não existir
- Validação de configuração Notion
- Relatório detalhado de status

### 4. **API REST** (`api/rest/notion_handlers.go`)
- Endpoints para gerenciamento de agentes
- Operações CRUD via HTTP
- Tratamento de erros padronizado

## ⚙️ Configuração

### 1. **Variáveis de Ambiente**

```bash
# Obrigatório para funcionalidade completa
NOTION_API_KEY=secret_xxxxxxxxxxxxxxxxxxxxxx

# Opcionais (com valores padrão)
NOTION_BASE_URL=https://api.notion.com/v1
NOTION_AGENT_NAME="IaC AI Agent"
NOTION_AGENT_DESCRIPTION="Intelligent Infrastructure as Code Analysis Agent"
NOTION_ENABLE_AGENT_CREATION=true
NOTION_AUTO_CREATE_ON_STARTUP=true
NOTION_MAX_REQUESTS_PER_MINUTE=60
```

### 2. **Arquivo de Configuração** (`configs/app.yaml`)

```yaml
notion:
  # api_key: ""             # Definir via env var NOTION_API_KEY
  base_url: "https://api.notion.com/v1"
  
  # Agent Configuration
  agent_name: "IaC AI Agent"
  agent_description: "Intelligent Infrastructure as Code Analysis Agent"
  
  # Features
  enable_agent_creation: true
  auto_create_on_startup: true
  
  # Rate Limiting
  max_requests_per_minute: 60
```

### 3. **Git Secrets** (Recomendado)

```bash
# Adicionar API key como secret
git secret add notion_api_key

# Editar o secret
git secret edit notion_api_key

# Verificar se foi adicionado
git secret list
```

## 🔧 Como Obter API Key do Notion

1. **Acesse** [Notion Integrations](https://www.notion.so/my-integrations)
2. **Clique** em "New integration"
3. **Preencha** os dados:
   - Name: `IaC AI Agent`
   - Associated workspace: Selecione seu workspace
4. **Copie** a "Internal Integration Token" (começa com `secret_`)
5. **Configure** a variável de ambiente ou Git Secret

## 🚀 Uso

### 1. **Inicialização Automática**

A aplicação automaticamente:
- Verifica se o serviço Notion está disponível
- Cria um agente padrão se não existir
- Valida a configuração
- Reporta o status na inicialização

### 2. **API Endpoints**

```bash
# Obter/criar agente padrão
GET /api/notion/agent/default

# Listar todos os agentes
GET /api/notion/agents

# Obter agente específico
GET /api/notion/agent?id=AGENT_ID

# Criar novo agente
POST /api/notion/agent
{
  "name": "Meu Agente",
  "description": "Descrição do agente",
  "capabilities": ["terraform_analysis", "security_review"]
}

# Remover agente
DELETE /api/notion/agent?id=AGENT_ID

# Status do serviço
GET /api/notion/status

# Capacidades padrão
GET /api/notion/capabilities
```

### 3. **Validação Manual**

```bash
# Executar validação de startup
go run ./cmd/agent --validate-config

# Testar integração
./scripts/test-notion-integration.sh
```

## 🧪 Testes

### 1. **Testes Unitários**

```bash
go test ./internal/services/... -v
```

### 2. **Testes de Integração**

```bash
go test ./test/integration/notion_integration_test.go -v
```

### 3. **Script de Teste Completo**

```bash
./scripts/test-notion-integration.sh
```

## 📊 Monitoramento

### 1. **Logs de Startup**

A aplicação reporta o status da integração Notion:

```
📝 Validando integração com Notion...
✅ Notion validado com sucesso
📝 Notion Agent Details:
  ID: agent_1234567890
  Name: IaC AI Agent
```

### 2. **Verificação de Status**

```bash
curl http://localhost:8080/api/notion/status
```

Resposta:
```json
{
  "available": true,
  "service": "notion"
}
```

## 🔒 Segurança

### 1. **API Key**
- Nunca commite a API key no código
- Use Git Secrets ou variáveis de ambiente
- Rotacione a key periodicamente

### 2. **Rate Limiting**
- Configurado para 60 requests/minuto por padrão
- Ajustável via configuração
- Tratamento automático de limites

### 3. **Validação**
- Validação de entrada em todos os endpoints
- Sanitização de dados
- Tratamento seguro de erros

## 🐛 Troubleshooting

### 1. **Erro: "NOTION_API_KEY não configurado"**
```bash
# Verificar se a variável está definida
echo $NOTION_API_KEY

# Ou verificar Git Secrets
git secret show notion_api_key
```

### 2. **Erro: "Serviço Notion não está disponível"**
- Verificar conectividade com internet
- Verificar se a API key é válida
- Verificar se o workspace Notion está ativo

### 3. **Erro: "Agente não encontrado"**
- Verificar se `NOTION_AGENT_NAME` está correto
- Verificar se `NOTION_ENABLE_AGENT_CREATION=true`
- Verificar permissões da API key

### 4. **Rate Limiting**
- Reduzir `NOTION_MAX_REQUESTS_PER_MINUTE`
- Implementar retry com backoff
- Verificar logs para padrões de uso

## 📚 Recursos Adicionais

- [Notion API Documentation](https://developers.notion.com/)
- [Notion Integration Guide](https://developers.notion.com/docs/getting-started)
- [Rate Limits](https://developers.notion.com/reference/request-limits)

## 🔄 Próximos Passos

1. **Implementar sincronização bidirecional**
2. **Adicionar webhooks Notion**
3. **Implementar cache de agentes**
4. **Adicionar métricas de uso**
5. **Implementar backup/restore de agentes**
