# 🤖 WhatsApp Agent - IaC AI Agent

Agente inteligente para análise de infraestrutura via WhatsApp, integrado com Web3 e sistema de tokens IACAI.

## 🚀 Características

- **Análise de Código Terraform**: Analisa código de infraestrutura e identifica problemas
- **Verificação de Segurança**: Detecta vulnerabilidades e problemas de segurança
- **Otimização de Custos**: Sugere melhorias para reduzir custos de infraestrutura
- **Autenticação Web3**: Integração com wallets Ethereum e NFTs Nation.fun
- **Sistema de Billing**: Cobrança automática de tokens IACAI
- **Logging Avançado**: Sistema completo de logs e monitoramento
- **Rate Limiting**: Controle de taxa de requisições
- **Docker Ready**: Containerização completa com Docker

## 📋 Pré-requisitos

- Go 1.21+
- Docker e Docker Compose
- PostgreSQL 15+
- Redis 7+
- Wallet Ethereum com NFT Nation.fun
- Chave API do WhatsApp Business

## 🛠️ Instalação

### 1. Clone o repositório

```bash
git clone https://github.com/iac-ai-agent/iac-ai-agent.git
cd iac-ai-agent
```

### 2. Configure as variáveis de ambiente

```bash
cp env.example .env
```

Edite o arquivo `.env` com suas configurações:

```env
# WhatsApp Configuration
WHATSAPP_WEBHOOK_URL=https://seu-dominio.com/webhook/whatsapp
WHATSAPP_VERIFY_TOKEN=seu_token_secreto_aqui
WHATSAPP_API_KEY=sua_chave_api_whatsapp

# Web3 Configuration
WALLET_ADDRESS=0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5
NFT_CONTRACT=nation.fun

# Database Configuration
DATABASE_URL=postgres://whatsapp_user:whatsapp_password@localhost:5432/whatsapp_agent

# Redis Configuration
REDIS_URL=redis://localhost:6379/0

# Logging Configuration
LOG_LEVEL=info
LOG_FILE=/var/log/whatsapp-agent.log
```

### 3. Execute com Docker Compose

```bash
docker-compose -f configs/docker-compose.whatsapp.yml up -d
```

### 4. Verifique se está funcionando

```bash
curl http://localhost:8080/webhook/whatsapp/health
```

## 📱 Comandos Disponíveis

### Comandos Gratuitos

- `/help` - Lista todos os comandos disponíveis
- `/status` - Status do agente e informações do sistema
- `/balance` - Verifica saldo de tokens IACAI
- `/usage` - Estatísticas de uso e histórico

### Comandos Pagos (1 token IACAI cada)

- `/analyze` - Analisa código Terraform
- `/security` - Verifica segurança do código
- `/cost` - Otimiza custos do código

### Exemplos de Uso

#### Análise de Código
```
/analyze
```hcl
resource "aws_instance" "web" {
  instance_type = "t3.micro"
  ami           = "ami-0c55b159cbfafe1d0"
}
```
```

#### Verificação de Segurança
```
/security
```hcl
resource "aws_s3_bucket" "data" {
  bucket = "my-bucket"
}
```
```

#### Análise de Custos
```
/cost
```hcl
resource "aws_instance" "web" {
  instance_type = "t3.large"
  ami           = "ami-0c55b159cbfafe1d0"
}
```
```

## 🔐 Autenticação Web3

O agente utiliza autenticação Web3 com as seguintes características:

- **Wallet Verification**: Verifica se a wallet possui NFT Nation.fun
- **Signature Validation**: Valida assinaturas digitais das mensagens
- **Token Balance**: Verifica saldo de tokens IACAI antes de processar comandos pagos
- **Lit Protocol**: Armazenamento seguro de chaves API

### Configuração da Wallet

1. Possua um NFT Nation.fun
2. Configure sua wallet no arquivo de configuração
3. O agente verificará automaticamente a propriedade do NFT

## 💰 Sistema de Billing

### Tokens IACAI

- **Comandos Gratuitos**: `/help`, `/status`, `/balance`, `/usage`
- **Comandos Pagos**: `/analyze`, `/security`, `/cost` (1 token cada)
- **Cobrança Automática**: Tokens são debitados automaticamente após processamento
- **Histórico**: Todas as transações são registradas

### Verificação de Saldo

```bash
curl -X GET http://localhost:8080/api/billing/balance/0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5
```

## 📊 Monitoramento

### Métricas Disponíveis

- **Requisições por minuto/hora**
- **Taxa de sucesso/erro**
- **Tempo de resposta médio**
- **Uso de tokens por usuário**
- **Comandos mais utilizados**

### Dashboards

- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Health Check**: http://localhost:8080/webhook/whatsapp/health

### Logs

```bash
# Ver logs em tempo real
docker-compose -f configs/docker-compose.whatsapp.yml logs -f whatsapp-agent

# Ver logs específicos
docker-compose -f configs/docker-compose.whatsapp.yml exec whatsapp-agent tail -f /var/log/whatsapp-agent.log
```

## 🧪 Testes

### Testes Unitários

```bash
go test ./internal/agent/whatsapp/...
```

### Testes de Integração

```bash
go test ./test/integration/...
```

### Testes com Docker

```bash
docker-compose -f configs/docker-compose.whatsapp.yml run --rm whatsapp-agent go test ./...
```

## 🚀 Deploy em Produção

### 1. Configurar SSL

```bash
# Gerar certificados SSL
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
        proxy_pass http://whatsapp-agent:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 3. Configurar WhatsApp Business API

1. Acesse o Facebook Developer Console
2. Configure o webhook: `https://seu-dominio.com/webhook/whatsapp`
3. Configure o token de verificação
4. Teste a integração

### 4. Deploy com Kubernetes

```bash
kubectl apply -f k8s/whatsapp-agent-deployment.yaml
kubectl apply -f k8s/whatsapp-agent-service.yaml
kubectl apply -f k8s/whatsapp-agent-ingress.yaml
```

## 🔧 Desenvolvimento

### Estrutura do Projeto

```
internal/agent/whatsapp/
├── agent.go          # Lógica principal do agente
├── commands.go       # Implementação dos comandos
├── config.go         # Configuração do agente
├── logging.go        # Sistema de logging
├── templates.go      # Templates de resposta
└── types.go          # Tipos e estruturas

api/rest/
└── whatsapp_handlers.go  # Handlers REST para webhook

cmd/whatsapp-agent/
└── main.go               # Ponto de entrada da aplicação
```

### Adicionando Novos Comandos

1. Adicione o comando em `AvailableCommands()` em `commands.go`
2. Implemente o handler do comando
3. Adicione testes unitários
4. Atualize a documentação

### Exemplo de Novo Comando

```go
"deploy": {
    Name:            "deploy",
    Description:    "Simula deploy do código",
    Pattern:        `^/deploy\s*(.*)`,
    Handler:        handleDeployCommand,
    RequiresPayment: true,
    TokenCost:      2,
},
```

## 📚 Documentação Adicional

- [Guia de Desenvolvimento](docs/WHATSAPP_AGENT_DEVELOPMENT_GUIDE.md)
- [Arquitetura do Sistema](docs/ARCHITECTURE.md)
- [Configuração de Ambiente](docs/CONFIGURACAO_AMBIENTE.md)
- [Integração Web3](docs/WEB3_INTEGRATION_GUIDE.md)

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🆘 Suporte

- **Issues**: [GitHub Issues](https://github.com/iac-ai-agent/iac-ai-agent/issues)
- **Discord**: [Discord Server](https://discord.gg/iac-ai-agent)
- **Email**: support@iac-ai-agent.com

## 🏆 Roadmap

- [ ] Suporte a múltiplos idiomas
- [ ] Integração com GitHub/GitLab
- [ ] Análise de código Python/Node.js
- [ ] Dashboard web para gerenciamento
- [ ] API REST completa
- [ ] Suporte a webhooks personalizados
- [ ] Integração com CI/CD
- [ ] Análise de custos em tempo real

---

**Desenvolvido com ❤️ pela equipe IaC AI Agent**
