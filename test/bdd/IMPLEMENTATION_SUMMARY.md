# Configuração de Testes BDD - Mock e Integração Real

## Resumo

Este projeto agora possui uma estrutura completa de testes BDD que permite validar tanto o comportamento mockado quanto a integração real com serviços externos.

## Arquivos Criados

### 1. Mocks (`test/mocks/`)
- **`web3_mocks.go`** - Mocks para todos os serviços Web3 (Privy, NFT, Tokens, Onramp)
- **`test_environment.go`** - Ambiente de teste configurável com mocks e serviços reais

### 2. Testes BDD (`test/bdd/`)
- **`features/mock_integration_tests.feature`** - Cenários de teste para mocks e integração
- **`steps/mock_integration_steps.go`** - Implementação dos steps de teste
- **`testconfig/config.go`** - Configuração de ambiente de teste
- **`testconfig/env.example`** - Exemplo de variáveis de ambiente
- **`run_bdd_tests.sh`** - Script para executar testes com diferentes configurações
- **`README.md`** - Documentação completa dos testes

## Funcionalidades Implementadas

### ✅ Testes com Mocks
- Login e autenticação mockados
- Compra de NFT simulada
- Compra de tokens simulada
- Análise de código mockada
- Onramp com cartão simulado
- Tratamento de erros mockados

### ✅ Testes de Integração Real
- Validação de configuração do sistema
- Integração com Privy usando App ID padrão
- Validação de contratos na Base Network
- Integração com Nation.fun para LLM
- Testes end-to-end completos

### ✅ Configuração Flexível
- Alternância entre modo mock e integração
- Configuração via variáveis de ambiente
- App ID padrão configurado
- Suporte a diferentes cenários de teste

## Como Usar

### 1. Executar Testes com Mocks
```bash
# Configurar ambiente
export TEST_MODE=mock

# Executar testes
./test/bdd/run_bdd_tests.sh mock
```

### 2. Executar Testes de Integração
```bash
# Configurar ambiente
export TEST_MODE=integration
export TEST_REAL_NATION_API_KEY=your_api_key

# Executar testes
./test/bdd/run_bdd_tests.sh integration
```

### 3. Executar Todos os Testes
```bash
./test/bdd/run_bdd_tests.sh all
```

## App ID Padrão Configurado

O sistema está configurado para usar o App ID padrão do Privy:
- **App ID**: `cmgh6un8w007bl10ci0tgitwp`
- **Network**: Base Goerli Testnet (Chain ID: 84531)
- **Contract**: `0x147e832418Cc06A501047019E956714271098b89`

## Tags de Teste Disponíveis

- `@mock` - Testes usando mocks
- `@integration` - Testes de integração
- `@real` - Testes com serviços reais
- `@unit` - Testes unitários
- `@performance` - Testes de performance
- `@error_handling` - Testes de tratamento de erro
- `@data_validation` - Testes de validação de dados

## Exemplos de Cenários

### Cenário Mock
```gherkin
@mock @unit
Cenário: Fluxo completo usando mocks - Login e Autenticação
  Dado que o sistema está em modo mock
  E que o MockPrivyClient está configurado
  Quando eu faço login com a wallet "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
  Então devo estar autenticado com sucesso
  E meu user ID deve ser "mock_user_123"
```

### Cenário de Integração
```gherkin
@integration @real
Cenário: Integração real com Privy - Validação de Configuração
  Dado que o sistema está em modo de integração real
  E que o app ID padrão "cmgh6un8w007bl10ci0tgitwp" está configurado
  Quando eu verifico a configuração do sistema
  Então todas as configurações devem estar válidas
```

## Benefícios

1. **Desenvolvimento Rápido**: Testes com mocks executam rapidamente
2. **Validação Real**: Testes de integração validam serviços externos
3. **Flexibilidade**: Alternância fácil entre modos
4. **Cobertura Completa**: Testa tanto comportamento quanto integração
5. **CI/CD Ready**: Configuração para pipelines automatizados

## Próximos Passos

1. **Configurar API Keys**: Adicionar chaves reais para testes de integração
2. **Executar Testes**: Rodar os testes para validar funcionamento
3. **Integrar CI/CD**: Adicionar aos pipelines de integração contínua
4. **Expandir Cenários**: Adicionar mais cenários conforme necessário

## Troubleshooting

### Problemas Comuns
1. **Testes de integração falham**: Verificar API keys e conectividade
2. **Mocks não funcionam**: Verificar configuração de TEST_MODE
3. **Timeout**: Aumentar TEST_TIMEOUT_SECONDS

### Logs Úteis
```bash
# Logs detalhados
TEST_DEBUG=true ./test/bdd/run_bdd_tests.sh mock --verbose

# Logs de rede
TEST_LOG_LEVEL=debug ./test/bdd/run_bdd_tests.sh integration
```

## Conclusão

A estrutura de testes BDD está completa e pronta para uso. Ela permite validar tanto o comportamento mockado quanto a integração real com o App ID padrão do Privy, fornecendo uma base sólida para desenvolvimento e validação do sistema.
