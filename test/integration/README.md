# Testes de Integração NFT Access - Nation Pass

## 🎯 Objetivo

Estes testes vão além dos mocks comuns e validam efetivamente:
- ✅ Acesso real a NFTs Nation Pass
- ✅ Interação com contratos reais na Base Network
- ✅ Validação de wallets configuradas
- ✅ Fluxo completo de integração

## 🚀 Como Executar

### 1. Configuração Básica

```bash
# Copie o arquivo de exemplo
cp env.example .env

# Configure suas variáveis
nano .env
```

### 2. Variáveis Obrigatórias

```bash
# ============================================
# 🔴 OBRIGATÓRIAS PARA TESTES DE INTEGRAÇÃO
# ============================================

# Wallet com Nation Pass NFT
WALLET_ADDRESS=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb

# Contrato Nation Pass na Base Network
NATION_NFT_CONTRACT=0x147e832418Cc06A501047019E956714271098b89

# Base Network RPC
BASE_RPC_URL=https://mainnet.base.org

# Token de autenticação (opcional)
WALLET_TOKEN=your_wallet_token_here
```

### 3. Executar Testes

```bash
# Testes básicos de integração
INTEGRATION_TESTS=true go test ./test/integration/ -v

# Testes com contratos reais
REAL_CONTRACT_TESTS=true go test ./test/integration/ -v

# Testes específicos Nation Pass
NATION_PASS_TESTS=true go test ./test/integration/ -v

# Todos os testes de integração
INTEGRATION_TESTS=true REAL_CONTRACT_TESTS=true NATION_PASS_TESTS=true go test ./test/integration/ -v
```

## 📋 Tipos de Testes

### 1. **Testes Básicos de Integração** (`INTEGRATION_TESTS=true`)
- ✅ Validação de wallet real configurada
- ✅ Listagem de NFTs da wallet
- ✅ Validação de acesso para diferentes tiers
- ✅ Conexão com Base Network
- ✅ Obtenção de tiers disponíveis

### 2. **Testes com Contratos Reais** (`REAL_CONTRACT_TESTS=true`)
- ✅ Validação do contrato Nation Pass
- ✅ Verificação de código do contrato
- ✅ Verificação de saldo ETH da wallet
- ✅ Validação de endereços

### 3. **Testes Específicos Nation Pass** (`NATION_PASS_TESTS=true`)
- ✅ Verificação de posse do NFT Nation Pass
- ✅ Validação de acesso por tier
- ✅ Teste de integração completa
- ✅ Listagem de funcionalidades disponíveis

## 🔍 O que os Testes Validam

### Validação de Wallet Real
```go
// Verifica se a wallet configurada tem acesso
nft, err := nftManager.CheckAccess(ctx, walletAddress)
if err != nil {
    // Wallet não tem NFT - esperado se não comprou
} else {
    // Wallet tem NFT - valida tier e status
}
```

### Validação de Contrato Real
```go
// Conecta com Base Network e valida contrato
client, err := ethclient.Dial(cfg.Web3.BaseRPCURL)
code, err := client.CodeAt(ctx, contractAddr, nil)
assert.NotEmpty(t, code, "Contract should have code")
```

### Validação de Acesso por Tier
```go
// Testa acesso para diferentes tiers
for tier := uint8(1); tier <= 3; tier++ {
    hasAccess, err := nftManager.ValidateAccess(ctx, walletAddress, tier)
    // Valida se tem acesso suficiente para o tier
}
```

## 🎨 Cenários de Teste

### Cenário 1: Wallet com Nation Pass
```
✅ Wallet configurada com NFT Nation Pass
✅ Validação de acesso bem-sucedida
✅ Listagem de NFTs funcionando
✅ Acesso validado para tiers apropriados
```

### Cenário 2: Wallet sem Nation Pass
```
❌ Wallet sem NFT Nation Pass
❌ Validação de acesso falha (esperado)
✅ Instruções para obter NFT são mostradas
✅ Testes continuam normalmente
```

### Cenário 3: Contrato Inválido
```
❌ Endereço de contrato inválido
❌ Falha na validação do contrato
✅ Testes são pulados com mensagem clara
```

## 🔧 Configuração Avançada

### Testes com Múltiplas Wallets
```bash
# Teste com wallet específica
WALLET_ADDRESS=0x123... go test ./test/integration/ -v

# Teste com contrato específico
NATION_NFT_CONTRACT=0x456... go test ./test/integration/ -v
```

### Testes com RPC Personalizado
```bash
# Usar RPC personalizado
BASE_RPC_URL=https://your-base-rpc.com go test ./test/integration/ -v
```

## 📊 Relatórios de Teste

### Saída Esperada - Wallet com NFT
```
=== RUN   TestNFTAccessIntegration/RealWalletNationPassValidation/CheckWalletAccess
    nft_access_integration_test.go:45: Testing with real wallet: 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
    nft_access_integration_test.go:58: Wallet has NFT access - Tier: Pro Access, Token ID: 123
--- PASS: TestNFTAccessIntegration/RealWalletNationPassValidation/CheckWalletAccess (0.15s)
```

### Saída Esperada - Wallet sem NFT
```
=== RUN   TestNFTAccessIntegration/RealWalletNationPassValidation/CheckWalletAccess
    nft_access_integration_test.go:45: Testing with real wallet: 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
    nft_access_integration_test.go:47: Wallet does not have NFT access: wallet não possui NFT de acesso
--- PASS: TestNFTAccessIntegration/RealWalletNationPassValidation/CheckWalletAccess (0.12s)
```

## 🚨 Troubleshooting

### Erro: "WALLET_ADDRESS not configured"
```bash
# Configure a variável de ambiente
export WALLET_ADDRESS=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
```

### Erro: "Failed to connect to Base Network"
```bash
# Verifique se o RPC está funcionando
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  https://mainnet.base.org
```

### Erro: "Contract code should not be empty"
```bash
# Verifique se o endereço do contrato está correto
# O contrato Nation Pass deve estar na Base Network
```

## 🎯 Próximos Passos

1. **Implementar ABI do Contrato**: Para chamadas específicas do contrato
2. **Adicionar Testes de Mint**: Para testar criação de novos NFTs
3. **Implementar Testes de Transfer**: Para testar transferência de NFTs
4. **Adicionar Testes de Upgrade**: Para testar upgrade de tiers
5. **Implementar Rate Limiting**: Para testar limites por tier

## 📚 Documentação Relacionada

- [Nation.fun Integration Guide](../docs/NATION_INTEGRATION.md)
- [Web3 Implementation Plan](../docs/WEB3_IMPLEMENTATION_PLAN.md)
- [Environment Variables](../docs/ENVIRONMENT_VARIABLES.md)
- [Setup Guide](../SETUP.md)
