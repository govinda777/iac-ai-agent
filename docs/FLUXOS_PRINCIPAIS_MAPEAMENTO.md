# 📊 Mapeamento dos Principais Fluxos - IaC AI Agent

## 🎯 Visão Geral

Este documento mapeia os principais fluxos de negócio do sistema IaC AI Agent, identificados através da análise do código e documentação existente.

## 🔄 Fluxos Principais Identificados

### 1. **Fluxo de Autenticação Web3**
**Objetivo**: Autenticar usuários via Web3 usando Privy.io

**Etapas**:
1. Usuário conecta wallet (MetaMask, Coinbase Wallet, etc.)
2. Privy.io valida a conexão
3. Sistema verifica token de autenticação
4. Usuário é autenticado e recebe sessão

**APIs Envolvidas**:
- `POST /api/auth/web3/verify` - Verificar token Privy
- `POST /api/auth/web3/check-access` - Verificar acesso a operações

**Componentes**:
- `Web3AuthService`
- `PrivyClient`
- `AuthenticatedUser`

---

### 2. **Fluxo de Compra de NFT de Acesso**
**Objetivo**: Permitir que usuários comprem NFTs para acessar diferentes tiers do sistema

**Etapas**:
1. Usuário navega para página de pricing
2. Seleciona tier desejado (Basic, Pro, Enterprise)
3. Escolhe método de pagamento (ETH direto ou cartão via onramp)
4. Confirma transação na wallet
5. NFT é mintado automaticamente
6. Acesso é concedido baseado no tier

**Tiers Disponíveis**:
- **Basic Access** (Tier 1): $25 - Análise básica
- **Pro Access** (Tier 2): $125 - Análise com LLM + Preview
- **Enterprise Access** (Tier 3): $500 - API dedicada + Suporte 24/7

**Componentes**:
- `NFTAccessManager`
- `PrivyOnrampManager`
- `BaseClient`

---

### 3. **Fluxo de Compra de Tokens IACAI**
**Objetivo**: Permitir que usuários comprem tokens para pagar por análises

**Etapas**:
1. Usuário navega para "Buy Tokens"
2. Seleciona pacote de tokens (Starter, Power, Enterprise)
3. Confirma compra com ETH ou cartão
4. Tokens são transferidos para wallet do usuário
5. Saldo é atualizado

**Pacotes Disponíveis**:
- **Starter Pack**: 100 tokens - $25
- **Power Pack**: 500 tokens - $100
- **Enterprise Pack**: 1000 tokens - $150

**Componentes**:
- `BotTokenManager`
- `PrivyOnrampManager`

---

### 4. **Fluxo de Análise de Código**
**Objetivo**: Analisar código Terraform e fornecer sugestões de melhoria

**Etapas**:
1. Usuário submete código Terraform
2. Sistema verifica acesso (NFT tier)
3. Sistema verifica saldo de tokens
4. Executa análises baseadas no tier:
   - Análise básica (Checkov)
   - Análise com LLM (GPT-4)
   - Preview Analysis
   - Security Audit
5. Debita tokens do saldo
6. Retorna relatório estruturado

**Tipos de Análise**:
- **Basic Analysis**: Checkov scan + Terraform validation
- **LLM Analysis**: Análise com IA + sugestões
- **Full Review**: Análise completa + otimizações
- **Preview Analysis**: Análise de terraform plan

**Componentes**:
- `AnalysisService`
- `TerraformAnalyzer`
- `CheckovAnalyzer`
- `IAMAnalyzer`
- `LLMService`

---

### 5. **Fluxo de Criação Automática de Agentes**
**Objetivo**: Criar automaticamente um agente personalizado para cada usuário

**Etapas**:
1. Sistema inicia e verifica se usuário tem agente
2. Se não tem agente, cria automaticamente usando template "General Purpose"
3. Configura agente com:
   - LLM (GPT-4)
   - Todas as análises habilitadas
   - Personalidade profissional
   - Limites padrão
4. Agente fica disponível para uso

**Templates Disponíveis**:
- **General Purpose**: Análise completa (padrão)
- **Security Specialist**: Foco em segurança
- **Cost Optimizer**: Otimização de custos
- **Architecture Advisor**: Análise arquitetural

**Componentes**:
- `AgentService`
- `Agent` model
- `AgentTemplate`

---

### 6. **Fluxo de Review de Pull Request**
**Objetivo**: Analisar Pull Requests do GitHub automaticamente

**Etapas**:
1. Webhook recebe notificação de PR
2. Sistema baixa código do PR
3. Executa análise completa
4. Gera comentário com sugestões
5. Posta comentário no PR

**Componentes**:
- `ReviewService`
- `AnalysisService`
- GitHub Webhook integration

---

### 7. **Fluxo de Controle de Acesso por Tier**
**Objetivo**: Controlar quais operações cada usuário pode realizar

**Lógica de Acesso**:
- **Sem NFT**: Apenas view_docs, basic_analysis
- **Basic (Tier 1)**: + terraform_analysis, checkov_scan
- **Pro (Tier 2)**: + llm_analysis, preview_analysis, security_audit
- **Enterprise (Tier 3)**: + cost_optimization, priority_support, full_review

**Componentes**:
- `Web3AuthService.determineAllowedOperations()`
- `NFTAccessManager.CheckAccess()`

---

### 8. **Fluxo de Rate Limiting**
**Objetivo**: Controlar uso do sistema baseado no tier do usuário

**Limites por Tier**:
- **Público**: 5 requests/hora
- **Basic**: Configurável (padrão: 20/hora)
- **Pro**: Configurável (padrão: 100/hora)
- **Enterprise**: Configurável (padrão: 1000/hora)

**Componentes**:
- `Web3AuthService.CheckRateLimit()`

---

## 🔗 Integrações Externas

### Privy.io
- **Autenticação Web3**: Login com wallets
- **Embedded Wallets**: Criação automática de wallets
- **Onramp**: Compra de crypto com cartão

### Base Network
- **Smart Contracts**: NFTs e tokens ERC-20/ERC-721
- **Transações**: Mint, transfer, balance check
- **RPC**: Conectividade com blockchain

### Nation.fun
- **Análise de Código**: Serviço de análise externa
- **API Integration**: Chamadas para análise

### OpenAI/Anthropic
- **LLM Services**: GPT-4, Claude para análise
- **Análise Inteligente**: Sugestões baseadas em IA

---

## 📊 Métricas de Negócio

### Métricas de Usuário
- NFT mints por dia (por tier)
- Token purchases por dia
- Active users por tier
- Retention rate por tier

### Métricas de Uso
- Análises executadas por tipo
- Tokens gastos por operação
- Rate limit hits por tier
- Upgrade rate (Basic → Pro → Enterprise)

### Métricas Financeiras
- Revenue em ETH e USD
- Average revenue per user (ARPU)
- Customer lifetime value (CLV)
- Conversion rate (visitante → comprador)

---

## 🧪 Cenários de Teste Identificados

### Cenários Críticos
1. **Novo usuário completa jornada completa**
2. **Compra NFT via cartão de crédito**
3. **Análise de código com diferentes tiers**
4. **Rate limiting funciona corretamente**
5. **Criação automática de agente**

### Cenários de Erro
1. **Saldo insuficiente de tokens**
2. **Tier incorreto para operação**
3. **Rate limit excedido**
4. **Falha na autenticação Web3**
5. **Timeout em análise**

### Cenários de Edge Case
1. **Usuário com múltiplos NFTs**
2. **Transação blockchain falha**
3. **Serviço externo indisponível**
4. **Código Terraform inválido**
5. **Webhook GitHub malformado**

---

## 🎯 Próximos Passos

1. **Criar testes BDD** para cada fluxo mapeado
2. **Implementar cenários de erro** nos testes
3. **Adicionar testes de performance** para rate limiting
4. **Criar testes de integração** com serviços externos
5. **Implementar testes de regressão** para mudanças de contrato

---

**Status**: ✅ Mapeamento completo dos fluxos principais  
**Versão**: 1.0.0  
**Última atualização**: 2025-01-15  
**Total de fluxos mapeados**: 8 principais + integrações
