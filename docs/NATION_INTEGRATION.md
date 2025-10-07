# 🎨 Integração com Nation.fun

## 📋 Visão Geral

O IaC AI Agent utiliza exclusivamente o [Nation.fun](https://nation.fun) como provedor de LLM (Large Language Model). Esta integração permite:

1. **Acesso via NFT**: Validação de posse do NFT Nation.fun para acesso ao serviço
2. **LLM Especializado**: Modelo de linguagem especializado em análise de Infrastructure as Code
3. **Autenticação Segura**: Via wallet token e NFT contract address

## 🔧 Como Funciona

### Fluxo de Integração

```
┌───────────────┐      ┌───────────────┐      ┌───────────────┐
│   IaC Agent   │ ──→  │  Nation.fun   │ ──→  │  LLM Service  │
└───────────────┘      └───────────────┘      └───────────────┘
        │                     ↑                      │
        │                     │                      │
        └─────────────────────┼──────────────────────┘
                              │
                      ┌───────────────┐
                      │  NFT Access   │
                      │  Validation   │
                      └───────────────┘
```

1. O IaC AI Agent envia requisições para a API Nation.fun
2. A API valida a posse do NFT Nation.fun
3. Se válido, processa a requisição via LLM
4. Retorna a resposta estruturada

## 🔐 Requisitos

Para usar o Nation.fun como provedor LLM, você precisa:

1. **NFT Nation.fun**: Adquirir um NFT em [nation.fun](https://nation.fun)
2. **Wallet Token**: Token de autenticação da sua wallet
3. **API Key**: Chave de API do Nation.fun

## ⚙️ Configuração

### Variáveis de Ambiente Obrigatórias

```bash
# Nation.fun
LLM_PROVIDER=nation.fun
LLM_API_KEY=your_api_key
LLM_MODEL=nation-1

# NFT Access
NFT_CONTRACT_ADDRESS=0x...  # Endereço do contrato Nation.fun
WALLET_TOKEN=your_wallet_token
WALLET_ADDRESS=0x...  # Seu endereço de wallet
```

### Configuração YAML

```yaml
# LLM Configuration
llm:
  provider: "nation.fun"
  model: "nation-1"
  temperature: 0.2
  max_tokens: 4000
  # api_key: definido via env var

# Web3 Configuration
web3:
  nft_access_contract_address: "0x..."  # Endereço do contrato Nation.fun
  wallet_token: "your_wallet_token"
  wallet_address: "0x..."  # Seu endereço de wallet
```

## 📝 Uso da API

### Cliente Nation.fun

O cliente Nation.fun está implementado em:

```
internal/agent/llm/nation_client.go
```

Este cliente implementa a interface `LLMProvider` e é usado automaticamente quando o provedor é configurado como "nation.fun".

### Exemplo de Uso

```go
import (
    "github.com/gosouza/iac-ai-agent/internal/agent/llm"
    "github.com/gosouza/iac-ai-agent/internal/models"
)

// Criar cliente
client := llm.NewClient(cfg, log)

// Gerar resposta
req := &models.LLMRequest{
    Prompt:      "Analise este código Terraform: resource \"aws_s3_bucket\" \"example\" { ... }",
    Temperature: 0.2,
    MaxTokens:   2000,
}

resp, err := client.Generate(req)
if err != nil {
    log.Error("Erro ao gerar resposta", "error", err)
    return
}

fmt.Println(resp.Content)
```

### Respostas Estruturadas

Para obter respostas em formato estruturado:

```go
analysis := &models.LLMStructuredResponse{}
err := client.GenerateStructured(req, analysis)
if err != nil {
    log.Error("Erro ao gerar resposta estruturada", "error", err)
    return
}

fmt.Printf("Resumo Executivo: %s\n", analysis.ExecutiveSummary)
fmt.Printf("Problemas Críticos: %d\n", len(analysis.CriticalIssues))
```

## 🔍 Validação

Durante o startup, o sistema valida automaticamente:

1. Conexão com Nation.fun API
2. Validade do API Key
3. Posse do NFT Nation.fun
4. Validade do Wallet Token

Se qualquer validação falhar, a aplicação não inicia.

## 🚀 Próximos Passos

1. **Implementar Cache**: Armazenar respostas comuns para reduzir chamadas à API
2. **Streaming**: Implementar suporte a respostas em streaming
3. **Métricas**: Tracking detalhado de uso, tokens e custos

## ❓ FAQ

### Por que Nation.fun?

Nation.fun foi escolhido como provedor exclusivo por:
- Especialização em análise de código
- Modelo treinado em Infrastructure as Code
- Integração nativa com NFTs para acesso
- Comunidade ativa de desenvolvedores

### Como adquirir um NFT Nation.fun?

Visite [nation.fun](https://nation.fun) e conecte sua wallet para adquirir um NFT de acesso.

### Posso usar outros provedores?

Não. O sistema foi redesenhado para usar exclusivamente Nation.fun como provedor LLM.

### Quais modelos estão disponíveis?

- `nation-1`: Modelo padrão
- `nation-pro`: Modelo avançado (requer NFT Pro)
- `nation-enterprise`: Modelo enterprise (requer NFT Enterprise)

## 📚 Referências

- [Nation.fun Website](https://nation.fun)
- [Nation.fun API Docs](https://docs.nation.fun)
- [NFT Access Guide](https://docs.nation.fun/nft-access)

---

**Última Atualização**: 2025-01-15  
**Versão**: 1.0.0  
**Status**: ✅ Implementado
