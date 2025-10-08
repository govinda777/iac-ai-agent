# 🤖 IaC AI Agent - Agente Inteligente com Capacidades Modulares

## 🎯 Visão Geral

O **IaC AI Agent** é um agente inteligente com arquitetura modular que possui **habilidades** (capabilities) para se comunicar via diferentes canais e realizar análises de Infrastructure as Code. A arquitetura permite adicionar novas habilidades facilmente sem modificar o core do agente.

## 🏗️ Arquitetura Modular

### Componentes Principais

```mermaid
graph TB
    A[IaC AI Agent Core] --> B[WhatsApp Capability]
    A --> C[IaC Analysis Capability]
    A --> D[Telegram Capability]
    A --> E[Slack Capability]
    
    B --> F[WhatsApp Business API]
    C --> G[Terraform Analyzer]
    C --> H[Security Analyzer]
    C --> I[Cost Analyzer]
    
    A --> J[Message Router]
    A --> K[Capability Manager]
    A --> L[Logger]
    A --> M[Config Manager]
```

### Fluxo de Processamento

1. **Recepção**: Mensagem via webhook (WhatsApp, Telegram, etc.)
2. **Roteamento**: Agente identifica qual habilidade pode processar
3. **Processamento**: Habilidade específica processa a mensagem
4. **Análise IaC**: Se necessário, delega para habilidade de análise
5. **Resposta**: Retorna resposta formatada para o canal original

## 🚀 Habilidades Disponíveis

### ✅ WhatsApp Capability
- **Comunicação**: Integração com WhatsApp Business API
- **Comandos**: `/help`, `/status`, `/ping`, `/analyze`, `/security`, `/cost`
- **Análise IaC**: Delegação automática para análise de código Terraform
- **Rate Limiting**: Controle de taxa de requisições
- **Segurança**: Verificação de assinaturas e tokens

### ✅ IaC Analysis Capability
- **Análise Terraform**: Análise geral de código Terraform
- **Segurança**: Detecção de vulnerabilidades e problemas de segurança
- **Custos**: Otimização de custos e estimativas
- **Provedores**: Suporte para AWS, Azure, GCP, Kubernetes
- **Cache**: Cache de resultados para melhor performance

### 🔄 Habilidades Futuras
- **Telegram Capability**: Comunicação via Telegram
- **Slack Capability**: Integração com Slack
- **GitHub Capability**: Análise de repositórios GitHub
- **CI/CD Capability**: Integração com pipelines

## 📱 Exemplos de Uso

### Análise via WhatsApp

```
Usuário: /analyze
```hcl
resource "aws_instance" "web" {
  instance_type = "t3.micro"
  ami           = "ami-0c55b159cbfafe1d0"
}
```

Agente: ✅ Análise Terraform Concluída
*Código analisado:* 89 caracteres
*Problemas encontrados:* 2
*Sugestões:* 3

*Problemas:*
• Recurso sem tags adequadas
• Instance type pode ser otimizada

*Sugestões:*
• Adicionar tags para melhor organização
• Considerar usar instance types menores
• Implementar lifecycle rules
```

### Verificação de Segurança

```
Usuário: /security
```hcl
resource "aws_s3_bucket" "data" {
  bucket = "my-bucket"
}
```

Agente: 🔒 Análise de Segurança Concluída
*Vulnerabilidades encontradas:* 1
*Recomendações:* 2

*Vulnerabilidades:*
• Bucket S3 sem criptografia (Severidade: High)

*Recomendações de Segurança:*
• Habilitar criptografia no bucket S3
• Restringir regras do security group
```

### Análise de Custos

```
Usuário: /cost
```hcl
resource "aws_instance" "web" {
  instance_type = "t3.large"
  ami           = "ami-0c55b159cbfafe1d0"
}
```

Agente: 💰 Análise de Custos Concluída
*Custo estimado mensal:* $150.00
*Potencial de economia:* $30.00
*Otimizações sugeridas:* 2

*Otimizações:*
• Usar Spot Instances (Economia: $20.00/mês)
• Otimizar tamanho de storage (Economia: $10.00/mês)
```

## 🛠️ Instalação e Configuração

### 1. Pré-requisitos

```bash
# Go 1.21+
go version

# Docker e Docker Compose
docker --version
docker-compose --version

# Dependências opcionais
# PostgreSQL 15+
# Redis 7+
```

### 2. Configuração

```bash
# Copiar configuração
cp configs/agent-config.yaml.example configs/agent-config.yaml

# Editar configurações
vim configs/agent-config.yaml
```

### 3. Variáveis de Ambiente

```bash
# WhatsApp Configuration
WHATSAPP_WEBHOOK_URL=https://seu-dominio.com/webhook/whatsapp
WHATSAPP_VERIFY_TOKEN=seu_token_secreto_aqui
WHATSAPP_API_KEY=sua_chave_api_whatsapp

# Web3 Configuration
WALLET_ADDRESS=0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5
NFT_CONTRACT=nation.fun

# Database Configuration
DATABASE_URL=postgres://iac_agent:senha@localhost:5432/iac_ai_agent

# Redis Configuration
REDIS_URL=redis://localhost:6379/0

# Logging Configuration
LOG_LEVEL=info
LOG_FILE=/var/log/iac-ai-agent.log
```

### 4. Execução

```bash
# Execução local
go run cmd/agent/main.go

# Ou compilar e executar
go build -o iac-ai-agent cmd/agent/main.go
./iac-ai-agent

# Execução com Docker
docker-compose up -d
```

## 🔧 Desenvolvimento

### Estrutura do Projeto

```
internal/agent/
├── core/
│   ├── agent.go          # Core do agente
│   └── types.go          # Tipos e interfaces
├── capabilities/
│   ├── whatsapp_capability.go    # Habilidade WhatsApp
│   ├── iac_analysis_capability.go # Habilidade IaC
│   └── telegram_capability.go     # Habilidade Telegram (futuro)
└── ...

api/rest/
└── agent_handler.go      # Handler REST

cmd/agent/
└── main.go               # Ponto de entrada

configs/
└── agent-config.yaml     # Configuração
```

### Adicionando Nova Habilidade

1. **Criar estrutura da habilidade**:

```go
type MyCapability struct {
    id          string
    name        string
    description string
    version     string
    status      string
    // ... outros campos
}

func (m *MyCapability) GetID() string { return m.id }
func (m *MyCapability) GetName() string { return m.name }
// ... implementar interface Capability
```

2. **Registrar no agente**:

```go
myCapability := capabilities.NewMyCapability()
agent.RegisterCapability(myCapability)
```

3. **Configurar no YAML**:

```yaml
capabilities:
  my-capability:
    enabled: true
    # ... configurações específicas
```

### Exemplo de Nova Habilidade

```go
// Telegram Capability
type TelegramCapability struct {
    id          string
    name        string
    description string
    version     string
    status      string
    
    botToken    string
    webhookURL  string
    logger      *core.Logger
}

func (t *TelegramCapability) CanHandle(message *core.Message) bool {
    return message.Source == "telegram"
}

func (t *TelegramCapability) ProcessMessage(ctx context.Context, message *core.Message) (*core.Response, error) {
    // Implementar processamento específico do Telegram
    // Pode delegar para outras habilidades se necessário
    return &core.Response{
        To:   message.From,
        Text: "Resposta do Telegram",
        Type: "text",
    }, nil
}
```

## 🧪 Testes

### Testes Unitários

```bash
# Testar core do agente
go test ./internal/agent/core/...

# Testar habilidades
go test ./internal/agent/capabilities/...

# Testar handlers
go test ./api/rest/...
```

### Testes de Integração

```bash
# Testar integração completa
go test ./test/integration/...
```

### Testes com Docker

```bash
# Executar testes em container
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## 📊 Monitoramento

### Health Check

```bash
# Status do agente
curl http://localhost:8080/agent/status

# Saúde das habilidades
curl http://localhost:8080/agent/health

# Lista de habilidades
curl http://localhost:8080/agent/capabilities
```

### Métricas Prometheus

```bash
# Métricas do agente
curl http://localhost:8080/metrics
```

### Logs

```bash
# Logs em tempo real
tail -f /var/log/iac-ai-agent.log

# Logs com Docker
docker-compose logs -f agent
```

## 🚀 Deploy em Produção

### 1. Configurar SSL

```bash
# Gerar certificados
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes
```

### 2. Configurar Nginx

```nginx
server {
    listen 443 ssl;
    server_name seu-dominio.com;
    
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location /webhook/whatsapp {
        proxy_pass http://agent:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    location /agent/ {
        proxy_pass http://agent:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 3. Deploy com Kubernetes

```bash
# Aplicar configurações
kubectl apply -f k8s/agent-deployment.yaml
kubectl apply -f k8s/agent-service.yaml
kubectl apply -f k8s/agent-ingress.yaml
```

## 🔒 Segurança

### Autenticação Web3
- Verificação de wallet Ethereum
- Validação de NFT Nation.fun
- Assinatura digital de mensagens

### Rate Limiting
- Limite por usuário
- Limite por habilidade
- Controle de burst

### Armazenamento Seguro
- Criptografia AES-256
- Integração com Lit Protocol
- Armazenamento distribuído

## 📈 Roadmap

### Próximas Habilidades
- [ ] **Telegram Capability** - Comunicação via Telegram
- [ ] **Slack Capability** - Integração com Slack
- [ ] **GitHub Capability** - Análise de repositórios
- [ ] **CI/CD Capability** - Integração com pipelines
- [ ] **Email Capability** - Comunicação via email

### Melhorias de Análise IaC
- [ ] **Multi-cloud Analysis** - Análise cross-cloud
- [ ] **Compliance Checking** - Verificação de conformidade
- [ ] **Cost Optimization** - Otimização avançada de custos
- [ ] **Security Scanning** - Escaneamento de segurança
- [ ] **Performance Analysis** - Análise de performance

### Melhorias de Performance
- [ ] **Async Processing** - Processamento assíncrono
- [ ] **Caching Layer** - Camada de cache avançada
- [ ] **Load Balancing** - Balanceamento de carga
- [ ] **Auto-scaling** - Escalabilidade automática

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-habilidade`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova habilidade'`)
4. Push para a branch (`git push origin feature/nova-habilidade`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🆘 Suporte

- **Issues**: [GitHub Issues](https://github.com/iac-ai-agent/iac-ai-agent/issues)
- **Discord**: [Discord Server](https://discord.gg/iac-ai-agent)
- **Email**: support@iac-ai-agent.com

---

**🎉 Agente Inteligente com Capacidades Modulares - Pronto para Uso!**

O IaC AI Agent agora possui uma arquitetura modular que permite fácil extensão e manutenção, com habilidades específicas para diferentes canais de comunicação e análise de Infrastructure as Code.
