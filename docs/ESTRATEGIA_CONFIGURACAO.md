# 🔐 Estratégia de Configuração e Segurança - IaC AI Agent

Este documento explica a estratégia de configuração e segurança adotada no IaC AI Agent, que simplifica drasticamente o setup para usuários finais.

## 🎯 Filosofia da Configuração

### ✅ **Simplicidade para o Usuário**
- **Apenas 1 variável obrigatória**: `LLM_API_KEY`
- **Configurações hardcoded**: Privy App ID e Wallet Address
- **Secrets gerenciados automaticamente**: Git Secrets + Lit Protocol

### 🔒 **Segurança Máxima**
- **Secrets nunca em texto plano**: Todos os secrets são criptografados
- **Acesso controlado**: Apenas desenvolvedores autorizados têm acesso
- **Auditoria completa**: Rastreamento de quem tem acesso aos secrets

## 📋 Estratégia por Tipo de Configuração

### 🟢 **Hardcoded (Código)**

| Configuração | Valor | Localização | Motivo |
|-------------|-------|-------------|--------|
| **Privy App ID** | `cmgh6un8w007bl10ci0tgitwp` | `configs/app.yaml` | App público, não sensível |
| **Wallet Address** | `0x147e832418Cc06A501047019E956714271098b89` | `configs/app.yaml` | Endereço público, não sensível |
| **Base RPC URL** | `https://mainnet.base.org` | `configs/app.yaml` | URL pública |
| **Chain ID** | `8453` | `configs/app.yaml` | ID público |

### 🔐 **Git Secrets (Criptografado)**

| Secret | Uso | Como Acessar |
|--------|-----|--------------|
| **wallet_private_key** | Assinatura de transações | `config.GetWalletPrivateKey()` |
| **github_token** | Integração GitHub | `config.GetGitHubToken()` |
| **github_webhook_secret** | Validação webhooks | `config.GetGitHubWebhookSecret()` |

### 🌐 **Lit Protocol (Web3)**

| Secret | Uso | Como Acessar |
|--------|-----|--------------|
| **whatsapp_api_key** | Integração WhatsApp | `litClient.GetWhatsAppAPIKey()` |

### 🔑 **Variável de Ambiente (Usuário)**

| Variável | Uso | Obrigatória |
|----------|-----|-------------|
| **LLM_API_KEY** | OpenAI API | ✅ Sim |

## 🛠️ Implementação Técnica

### 1. Configuração Hardcoded

```yaml
# configs/app.yaml
web3:
  privy_app_id: "cmgh6un8w007bl10ci0tgitwp"  # Hardcoded
  wallet_address: "0x147e832418Cc06A501047019E956714271098b89"  # Hardcoded
```

### 2. Git Secrets

```go
// pkg/config/config.go
func (c *Config) GetWalletPrivateKey() (string, error) {
    return getGitSecret("wallet_private_key")
}

func getGitSecret(secretName string) (string, error) {
    cmd := exec.Command("git", "secret", "show", secretName)
    // Executa git secret show para descriptografar
}
```

### 3. Lit Protocol

```go
// internal/platform/web3/lit_protocol.go
func (lpc *LitProtocolClient) GetWhatsAppAPIKey() (string, error) {
    // Usa Lit Protocol para descriptografar secret
    // Condições de acesso: apenas owner da wallet
}
```

## 🔄 Fluxo de Configuração

### Para Desenvolvedores

1. **Setup inicial**:
   ```bash
   git secret init
   git secret tell developer@example.com
   echo "private_key_here" | git secret add wallet_private_key
   git secret hide -m
   ```

2. **Acesso aos secrets**:
   ```bash
   git secret reveal  # Descriptografa todos os secrets
   ```

### Para Usuários Finais

1. **Setup mínimo**:
   ```bash
   cp env.example .env
   echo "LLM_API_KEY=sk-proj-..." >> .env
   make run
   ```

2. **Configuração WhatsApp** (opcional):
   ```bash
   # Via API - armazena automaticamente no Lit Protocol
   curl -X POST /api/v1/whatsapp/configure \
     -d '{"api_key": "your_whatsapp_key"}'
   ```

## 🎯 Benefícios da Estratégia

### ✅ **Para Usuários**
- **Setup em 1 minuto**: Apenas API key do OpenAI
- **Zero configuração Web3**: Tudo já configurado
- **Secrets seguros**: Não precisam gerenciar chaves privadas

### ✅ **Para Desenvolvedores**
- **Secrets centralizados**: Todos em um lugar
- **Acesso controlado**: Apenas quem tem chave GPG
- **Auditoria completa**: Rastreamento de acesso

### ✅ **Para Segurança**
- **Nenhum secret em texto plano**: Tudo criptografado
- **Princípio do menor privilégio**: Acesso mínimo necessário
- **Rotação fácil**: Trocar secrets sem afetar usuários

## 🔧 Comandos Úteis

### Git Secrets

```bash
# Listar secrets
git secret list

# Adicionar novo secret
echo "new_secret" | git secret add secret_name

# Remover secret
git secret remove secret_name

# Re-criptografar após mudanças
git secret hide -m
```

### Lit Protocol

```bash
# Verificar se tem WhatsApp API key
curl /api/v1/whatsapp/status

# Configurar WhatsApp API key
curl -X POST /api/v1/whatsapp/configure \
  -d '{"api_key": "your_key"}'
```

## 🚨 Troubleshooting

### Erro: "git-secret não está disponível"

```bash
# Instalar git-secret
brew install git-secret  # macOS
sudo apt-get install git-secret  # Ubuntu
```

### Erro: "Secret não encontrado"

```bash
# Verificar secrets disponíveis
git secret list

# Descriptografar todos os secrets
git secret reveal
```

### Erro: "Lit Protocol error"

```bash
# Verificar conectividade com Base Network
curl https://mainnet.base.org

# Verificar se wallet está conectada
curl /api/v1/wallet/status
```

## 📚 Documentação Relacionada

- **[GIT_SECRETS_SETUP.md](GIT_SECRETS_SETUP.md)** - Configuração detalhada do Git Secrets
- **[WHATSAPP_INTEGRATION.md](WHATSAPP_INTEGRATION.md)** - Integração WhatsApp + Lit Protocol
- **[CONFIGURACAO_VARIAVEIS.md](CONFIGURACAO_VARIAVEIS.md)** - Detalhamento de variáveis
- **[SECURE_TOKEN_USAGE.md](SECURE_TOKEN_USAGE.md)** - Uso seguro de tokens

---

**Status**: ✅ Estratégia implementada e documentada  
**Versão**: 1.0.0  
**Última atualização**: 2025-01-15
