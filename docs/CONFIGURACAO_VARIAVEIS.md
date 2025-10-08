# Configuração de Variáveis de Ambiente - IaC AI Agent

Este documento explica como configurar as variáveis de ambiente necessárias para o projeto IaC AI Agent.

## Passo 1: Criar o arquivo .env

Crie um arquivo chamado `.env` na raiz do projeto. Você pode usar o arquivo `.env-example` como base:

```bash
cp .env-example .env
```

## Passo 2: Configurar variáveis obrigatórias

Edite o arquivo `.env` e preencha as seguintes variáveis obrigatórias:

### Privy.io (Autenticação Web3)

```bash
PRIVY_APP_ID=cmgh6un8w007bl10ci0tgitwp  # Já configurado
PRIVY_APP_SECRET=                        # Obter no dashboard do Privy.io
```

### Wallet e NFT da Nation.fun

```bash
WALLET_ADDRESS=0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5  # Já configurado
# WALLET_PRIVATE_KEY=                    # Opcional: usada apenas para geração automática do token
# WALLET_TOKEN=                          # Opcional: configuração manual ou gerado automaticamente
NATION_NFT_CONTRACT=                     # Endereço do contrato NFT da Nation.fun (padrão usado se não definido)
```

### LLM API (Nation.fun ou OpenAI)

```bash
LLM_PROVIDER=nation.fun                  # Usar nation.fun, openai ou anthropic
LLM_API_KEY=                             # Chave de API do provedor LLM
```

## Passo 3: Configurar app.yaml (opcional)

O arquivo `configs/app.yaml` já foi atualizado com as seguintes configurações:

```yaml
web3:
  privy_app_id: "cmgh6un8w007bl10ci0tgitwp"
  wallet_address: "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5"
  # outras configurações...
```

## Passo 4: Verificação de segurança

Certifique-se de que o arquivo `.env` está no `.gitignore` para evitar que informações sensíveis sejam expostas.

## Obtendo as variáveis necessárias

### 1. Privy App Secret

1. Acesse https://dashboard.privy.io/
2. Faça login na sua conta
3. Vá em Settings → API Keys
4. Copie o App Secret

### 2. Wallet Private Key

**ATENÇÃO: Nunca compartilhe sua chave privada!**

- **OPCIONAL**: Usada apenas para geração automática do token
- Se preferir, use a abordagem manual e fornecer o WALLET_TOKEN diretamente
- Exporte da sua carteira MetaMask/Coinbase Wallet apenas se necessário
- Certifique-se de que este arquivo está no .gitignore

### 3. Wallet Token

- Token de autenticação para a API Nation.fun
- Você tem duas opções:
  1. **Geração automática**: Fornecendo WALLET_PRIVATE_KEY
  2. **Configuração manual**: Obtendo o token diretamente da plataforma
- O sistema tentará fazer o melhor mesmo se não for possível gerar o token

### 3. Nation NFT Contract

1. Acesse https://nation.fun/
2. Navegue até sua Nation
3. Clique para ver o contrato no Basescan
4. Copie o endereço do contrato

### 4. LLM API Key

#### Nation.fun
1. Acesse https://nation.fun/api-access
2. Gere ou copie sua API key

#### OpenAI (alternativa)
1. Acesse https://platform.openai.com/api-keys
2. Crie uma nova chave API

#### Anthropic (alternativa)
1. Acesse https://console.anthropic.com/
2. Obtenha sua API key

## Validação

Após configurar todas as variáveis, execute a aplicação para verificar se está tudo correto:

```bash
go run cmd/agent/main.go
```

O sistema fará a validação automática das variáveis durante o startup.
