# 🤖 IaC AI Agent

> Agente de IA para análise, revisão e otimização de código Infrastructure as Code (Terraform) com autenticação Web3 e pagamentos on-chain.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Privy](https://img.shields.io/badge/Auth-Privy.io-6366F1)](https://privy.io)
[![Base Network](https://img.shields.io/badge/L2-Base-0052FF)](https://base.org)
[![Nation.fun](https://img.shields.io/badge/Community-Nation.fun-FF6B6B)](https://nation.fun)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 🚨 ANTES DE COMEÇAR - LEIA ISTO!

> **A aplicação NÃO VAI INICIAR sem estas 3 coisas configuradas:**

| Requisito | O que é | Onde obter |
|-----------|---------|------------|
| 🎨 **Nation.fun NFT** | NFT de membership da Nation.fun | https://nation.fun/ |
| 🔐 **Privy.io Account** | Credenciais de autenticação Web3 | https://privy.io |
| 🤖 **OpenAI API Key** | Chave de API do LLM | https://platform.openai.com/api-keys |

📖 **Setup Completo**: Leia [`SETUP.md`](SETUP.md) ou [`docs/ENVIRONMENT_VARIABLES.md`](docs/ENVIRONMENT_VARIABLES.md)

---

## 🎯 O Que é?

O **IaC AI Agent** é um bot inteligente que analisa código Terraform e fornece:

- ✅ **Análise de Segurança** (Checkov)
- ✅ **Análise com LLM** (GPT-4/Claude) - Sugestões contextualizadas
- ✅ **Detecção de Drift** e Preview Analysis
- ✅ **Otimização de Custos** com estimativas
- ✅ **Recomendações de Arquitetura**
- ✅ **Best Practices** e IAM Analysis
- ✅ **Scoring Automático** de qualidade

### 🤖 Sistema de Agentes (NOVO!)

O IaC AI Agent possui um **sistema de agentes inteligentes** que:

- ✨ **Cria automaticamente um agente** quando você inicia pela primeira vez
- 🎨 **4 templates pré-definidos**: General Purpose, Security, Cost, Architecture
- 🧠 **Personalidade customizável**: Ajuste tom, verbosidade, estilo
- 📊 **Conhecimento especializado**: Expertise em AWS, Azure, GCP, Terraform
- 🔧 **Limites configuráveis**: Rate limits, custos, timeouts
- 📈 **Métricas de uso**: Performance, custos, qualidade

```
🤖 Verificando agente padrão...
ℹ️  Nenhum agente encontrado
✨ Criando novo agente automaticamente...
✅ Novo agente criado: IaC Agent - 0x742d35
```

📖 **Documentação completa**: [`docs/AGENT_SYSTEM.md`](docs/AGENT_SYSTEM.md)

### 🔐 Web3 Native

- **Autenticação via Privy.io**: Login com wallet (MetaMask, Coinbase) ou email
- **NFTs de Acesso** (Base Network): 3 tiers de acesso permanente
- **Token IACAI** (ERC-20): Pague por análises com tokens on-chain
- **Privy Onramp**: Compre crypto com cartão/PIX sem ter wallet

---

## 🏗️ Arquitetura

```
┌─────────────────────────────────────────────────────┐
│                    Frontend                          │
│  Privy SDK + Wagmi + Next.js                        │
└─────────────────────────────────────────────────────┘
                       ↓
┌─────────────────────────────────────────────────────┐
│              Backend (Go)                            │
│                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────┐│
│  │ API REST     │  │ Web3 Platform│  │ LLM       ││
│  │ • Handlers   │  │ • Privy      │  │ • OpenAI  ││
│  │ • Auth       │  │ • Base       │  │ • Claude  ││
│  └──────────────┘  │ • NFT/Token  │  └───────────┘│
│                     │ • Onramp     │                │
│  ┌──────────────┐  └──────────────┘  ┌───────────┐│
│  │ Analyzers    │                     │ Knowledge ││
│  │ • Terraform  │                     │ Base      ││
│  │ • Checkov    │                     │ • Rules   ││
│  │ • IAM        │                     │ • Modules ││
│  │ • Preview    │                     │ • Patterns││
│  └──────────────┘                     └───────────┘│
└─────────────────────────────────────────────────────┘
                       ↓
┌─────────────────────────────────────────────────────┐
│           Base Network (L2 Ethereum)                 │
│                                                      │
│  ┌────────────────────┐    ┌────────────────────┐  │
│  │ NFT Access (ERC-721)│    │ IACAI Token (ERC-20)│  │
│  │ • Basic: 0.01 ETH  │    │ • Packages          │  │
│  │ • Pro: 0.05 ETH    │    │ • Payments          │  │
│  │ • Enterprise: 0.2  │    │ • Transfers         │  │
│  └────────────────────┘    └────────────────────┘  │
└─────────────────────────────────────────────────────┘
```

---

## ⚡ Quick Start

### 🔴 REQUISITOS OBRIGATÓRIOS

Antes de iniciar, você PRECISA ter:

1. **Nation.fun NFT** - Compre em https://nation.fun/
2. **Privy.io Account** - Crie em https://privy.io
3. **OpenAI API Key** - Obtenha em https://platform.openai.com/api-keys

```bash
# 1. Clone
git clone https://github.com/gosouza/iac-ai-agent
cd iac-ai-agent

# 2. Configure variáveis OBRIGATÓRIAS
cp .env.example .env

# Edite .env e adicione:
# - PRIVY_APP_ID=app_xxx
# - PRIVY_APP_SECRET=xxx
# - WALLET_ADDRESS=0x... (com Nation.fun NFT)
# - WALLET_PRIVATE_KEY=0x...
# - NATION_NFT_CONTRACT=0x...
# - LLM_API_KEY=sk-...

# 3. Execute
go run cmd/agent/main.go

# A aplicação vai validar TUDO antes de iniciar!
# ✅ LLM Connection
# ✅ Privy.io Credentials
# ✅ Base Network
# ✅ Nation.fun NFT Ownership

# 4. Teste
curl http://localhost:8080/health
```

📖 **Documentação de Variáveis**: [ENVIRONMENT_VARIABLES.md](docs/ENVIRONMENT_VARIABLES.md)  
📖 **Guia completo**: [QUICKSTART.md](docs/QUICKSTART.md)  
📖 **Integração Nation.fun**: [NATION_FUN_INTEGRATION.md](docs/NATION_FUN_INTEGRATION.md)

---

## 🎫 Sistema de Acesso (NFTs)

### Tiers Disponíveis

| Tier | Preço | Benefícios |
|------|-------|------------|
| **Basic** | 0.01 ETH (~$25) | Análises ilimitadas, Checkov, Suporte Discord |
| **Pro** | 0.05 ETH (~$125) | + LLM, Preview, Drift, Priority Support |
| **Enterprise** | 0.2 ETH (~$500) | + API dedicada, Custom KB, SLA 24/7 |

### Como Funciona?

1. **Compra NFT** → Acesso permanente
2. **Compra Tokens (IACAI)** → Pague por análises
3. **Use o Bot** → Tokens são debitados automaticamente

### Métodos de Pagamento

- ✅ **Cartão de Crédito/Débito** (via Privy Onramp)
- ✅ **PIX** (Brasil)
- ✅ **ETH na wallet**
- ✅ **Apple Pay / Google Pay**

---

## 💎 Tokens IACAI

### Pacotes Disponíveis

| Pacote | Tokens | Preço | Desconto |
|--------|--------|-------|----------|
| Starter | 100 | 0.005 ETH ($10) | - |
| Power | 500 | 0.0225 ETH ($45) | 10% |
| Pro | 1000 | 0.0425 ETH ($85) | 15% |
| Enterprise | 5000 | 0.1875 ETH ($375) | 25% |

### Tabela de Custos

| Operação | Custo (IACAI) |
|----------|---------------|
| Terraform Analysis | 1 |
| Checkov Scan | 2 |
| LLM Analysis | 5 |
| Preview Analysis | 3 |
| Security Audit | 10 |
| Cost Optimization | 5 |
| Full Review | 15 |

---

## 📚 Documentação

### Para Começar

- 📖 [Quick Start](docs/QUICKSTART.md) - Setup em 5 minutos
- 🎯 [Objetivo do Projeto](docs/OBJECTIVE.md) - Visão completa
- 🏗️ [Arquitetura](docs/ARCHITECTURE.md) - Design técnico

### Para Desenvolvedores

- 🔌 [Guia de Integração Web3](docs/WEB3_INTEGRATION_GUIDE.md) - Privy + Base
- 📝 [Resumo de Implementação](docs/IMPLEMENTATION_SUMMARY.md) - O que foi feito
- 🗺️ [Roadmap](docs/IMPLEMENTATION_ROADMAP.md) - Próximos passos

### Para Executivos

- 📊 [Executive Summary](docs/EXECUTIVE_SUMMARY.md) - Visão executiva
- 📈 [Análise do Projeto](docs/PROJECT_ANALYSIS.md) - Status atual

### Índice Completo

- 📚 [INDEX.md](docs/INDEX.md) - Todos os documentos

---

## 🧪 Testes BDD

Testes completos em Gherkin (português) cobrindo todos os fluxos:

```bash
# Instalar Godog
go install github.com/cucumber/godog/cmd/godog@latest

# Executar todos os testes
godog test/bdd/features/

# Testes disponíveis:
# ✓ user_onboarding.feature     - Autenticação Privy
# ✓ nft_purchase.feature         - Compra de NFT
# ✓ token_purchase.feature       - Compra de tokens
# ✓ bot_analysis.feature         - Uso do bot
```

### Exemplo de Cenário BDD

```gherkin
Cenário: Comprar NFT Pro Access usando Privy Onramp
  Dado que minha wallet tem saldo de "0" ETH
  E eu estou visualizando o tier "Pro Access"
  Quando eu clico em "Buy with Card"
  Então devo ver a modal do Privy Onramp
  E devo ver o valor necessário: "0.05 ETH (~$125 USD)"
  Quando eu seleciono "Credit Card" como método de pagamento
  E insiro os dados do cartão
  E concluo o pagamento
  Então devo ver "Payment Processing"
  E em até 10 minutos o ETH deve chegar na wallet
  E o NFT deve ser automaticamente mintado
  E devo receber notificação "NFT Pro Access adquirido!"
```

---

## 🚀 Deployment

### Backend (Docker)

```bash
# Build
docker build -t iacai-agent .

# Run
docker run -p 8080:8080 \
  -e PRIVY_APP_ID=xxx \
  -e PRIVY_APP_SECRET=xxx \
  -e LLM_API_KEY=xxx \
  -e BASE_RPC_URL=https://mainnet.base.org \
  iacai-agent
```

### Smart Contracts

```bash
cd contracts
npm install
npx hardhat run scripts/deploy.ts --network base
```

Contratos deployados na **Base Mainnet** (Chain ID 8453):
- NFT Access: `0x...` (a ser deployado)
- IACAI Token: `0x...` (a ser deployado)

---

## 🛠️ Stack Tecnológica

### Backend
- **Linguagem**: Go 1.21+
- **Frameworks**: Standard library, Gorilla Mux
- **LLM**: OpenAI GPT-4, Anthropic Claude
- **Security**: Checkov integration

### Web3
- **Auth**: Privy.io SDK
- **Blockchain**: Base Network (L2 Ethereum)
- **Wallets**: MetaMask, Coinbase Wallet, Embedded Wallets
- **Onramp**: MoonPay, Transak (via Privy)
- **Contracts**: Solidity 0.8.20, OpenZeppelin

### Frontend (Sugerido)
- **Framework**: Next.js 14
- **Auth**: `@privy-io/react-auth`
- **Web3**: Wagmi, Viem
- **UI**: Tailwind CSS, shadcn/ui

---

## 📊 Métricas e Monitoramento

```yaml
Business:
  - NFT mints por dia/tier
  - Token purchases
  - Revenue (ETH/USD)
  - Active users por tier
  - Análises executadas

Technical:
  - API response time
  - LLM latency
  - Blockchain tx success rate
  - Onramp conversion rate
  - Error rate
```

---

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor:

1. Fork o repositório
2. Crie uma branch (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Add nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

### Guidelines

- Escreva testes BDD para novas features
- Mantenha cobertura de testes > 80%
- Siga Go best practices
- Documente APIs públicas
- Atualize README se necessário

---

## 📝 License

MIT License - veja [LICENSE](LICENSE) para detalhes.

---

## 🔗 Links Úteis

- **Privy Docs**: https://docs.privy.io
- **Base Network**: https://docs.base.org
- **OpenAI API**: https://platform.openai.com/docs
- **Checkov**: https://www.checkov.io
- **Terraform**: https://www.terraform.io

---

## 📞 Suporte

- **Issues**: [GitHub Issues](https://github.com/gosouza/iac-ai-agent/issues)
- **Email**: support@iacai.com
- **Discord**: (em breve)
- **Twitter**: [@iacaiagent](https://twitter.com/iacaiagent)

---

## 🎯 Roadmap

- [x] Análise básica de Terraform
- [x] Integração Checkov
- [x] LLM Analysis (GPT-4/Claude)
- [x] Autenticação Web3 (Privy)
- [x] NFTs de acesso (Base Network)
- [x] Token IACAI (ERC-20)
- [x] Privy Onramp
- [x] Testes BDD completos
- [ ] Preview Analysis
- [ ] Drift Detection
- [ ] Dashboard Web
- [ ] Integração CI/CD
- [ ] Mobile App
- [ ] Governance DAO

---

## 🌟 Star History

Se você gostou do projeto, dê uma ⭐ no GitHub!

---

**Status**: 🚀 Pronto para produção  
**Versão**: 1.0.0  
**Última Atualização**: 2025-01-15

---

Made with ❤️ by the IaC AI Agent Team
