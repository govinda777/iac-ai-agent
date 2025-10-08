# Testes BDD - Mock e Integração Real

Este diretório contém testes BDD (Behavior-Driven Development) que validam tanto o comportamento mockado quanto a integração real com serviços externos.

## Estrutura

```
test/bdd/
├── features/                    # Arquivos .feature com cenários BDD
│   ├── critical_path.feature   # Caminho crítico do produto
│   └── mock_integration_tests.feature  # Testes de mock e integração
├── steps/                       # Implementação dos steps
│   ├── critical_path_steps.go   # Steps do caminho crítico
│   ├── mock_integration_steps.go # Steps de mock e integração
│   └── steps_runner.go          # Registro de todos os steps
├── testconfig/                  # Configuração de testes
│   ├── config.go               # Configuração de ambiente de teste
│   └── env.example             # Exemplo de variáveis de ambiente
└── run_bdd_tests.sh            # Script para executar testes
```

## Modos de Teste

### 1. Modo Mock (`mock`)
- Usa apenas mocks e simulações
- Não faz chamadas para serviços externos
- Ideal para testes unitários e desenvolvimento local
- Execução rápida e determinística

### 2. Modo Integração (`integration`)
- Usa serviços reais para validação
- Testa integração com Privy, Base Network, Nation.fun
- Requer configuração de API keys reais
- Valida configurações e conectividade

### 3. Modo Real (`real`)
- Testes end-to-end completos
- Usa todos os serviços em produção
- Valida fluxo completo do usuário
- Requer ambiente de teste configurado

## Configuração

### Variáveis de Ambiente

Copie o arquivo de exemplo e configure as variáveis:

```bash
cp test/bdd/testconfig/env.example .env.test
```

Principais variáveis:

```bash
# Modo de teste
export TEST_MODE=mock  # mock, integration, real

# Configurações para integração real
export TEST_REAL_PRIVY_APP_ID=cmgh6un8w007bl10ci0tgitwp
export TEST_REAL_BASE_RPC_URL=https://goerli.base.org
export TEST_REAL_CONTRACT_ADDR=0x147e832418Cc06A501047019E956714271098b89
export TEST_REAL_NATION_API_KEY=your_api_key_here

# Configurações para mocks
export TEST_MOCK_USER_ID=mock_user_123
export TEST_MOCK_WALLET_ADDRESS=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
```

### App ID Padrão

O sistema usa o App ID padrão do Privy para testes de integração:
- **App ID**: `cmgh6un8w007bl10ci0tgitwp`
- **Network**: Base Goerli Testnet (Chain ID: 84531)
- **Contract**: `0x147e832418Cc06A501047019E956714271098b89`

## Executando Testes

### Usando o Script

```bash
# Testes com mocks
./test/bdd/run_bdd_tests.sh mock

# Testes de integração
./test/bdd/run_bdd_tests.sh integration

# Testes reais
./test/bdd/run_bdd_tests.sh real

# Todos os modos
./test/bdd/run_bdd_tests.sh all

# Com opções específicas
./test/bdd/run_bdd_tests.sh mock --verbose --coverage
```

### Usando Godog Diretamente

```bash
# Instalar godog
go install github.com/cucumber/godog/cmd/godog@latest

# Executar com tags específicas
godog --tags="@mock" test/bdd/features/

# Executar com configuração específica
TEST_MODE=integration godog test/bdd/features/
```

### Tags Disponíveis

- `@mock` - Testes usando mocks
- `@integration` - Testes de integração
- `@real` - Testes com serviços reais
- `@unit` - Testes unitários
- `@performance` - Testes de performance
- `@error_handling` - Testes de tratamento de erro
- `@data_validation` - Testes de validação de dados

## Cenários de Teste

### Testes com Mocks

```gherkin
@mock @unit
Cenário: Fluxo completo usando mocks - Login e Autenticação
  Dado que o sistema está em modo mock
  E que o MockPrivyClient está configurado
  Quando eu faço login com a wallet "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
  Então devo estar autenticado com sucesso
  E meu user ID deve ser "mock_user_123"
```

### Testes de Integração

```gherkin
@integration @real
Cenário: Integração real com Privy - Validação de Configuração
  Dado que o sistema está em modo de integração real
  E que o app ID padrão "cmgh6un8w007bl10ci0tgitwp" está configurado
  Quando eu verifico a configuração do sistema
  Então todas as configurações devem estar válidas
```

## Mocks Disponíveis

### MockPrivyClient
- Simula autenticação e validação de usuários
- Configurável para falhas simuladas
- Retorna dados de teste consistentes

### MockNFTAccessManager
- Simula mint e verificação de NFTs
- Configurável para diferentes tiers
- Simula transações blockchain

### MockBotTokenManager
- Simula mint e transferência de tokens
- Gerencia saldos mockados
- Simula operações de token

### MockAnalysisService
- Simula análise de código Terraform
- Retorna resultados estruturados
- Configurável para diferentes cenários

### MockPrivyOnrampManager
- Simula compra com cartão de crédito
- Gera cotações mockadas
- Simula fluxo de pagamento

## Validação de Integração

### Configuração do Sistema
- Verifica se todas as configurações estão válidas
- Valida conectividade com serviços externos
- Confirma deploy de contratos

### Privy Integration
- Valida App ID e configurações
- Testa autenticação de usuários
- Verifica validação de wallets

### Base Network Integration
- Confirma conectividade com RPC
- Valida deploy de contratos
- Testa operações blockchain

### Nation.fun Integration
- Valida API key e conectividade
- Testa análise de código
- Verifica formato de resposta

## Tratamento de Erros

### Falhas Simuladas
```bash
# Habilitar falhas simuladas
export TEST_MOCK_FAILURES=true
```

### Cenários de Erro
- Falha de autenticação
- Falha de mint NFT
- Falha de análise
- Timeout de serviços

## Performance e Carga

### Testes de Performance
```bash
# Executar testes de performance
export TEST_PERFORMANCE=true
export TEST_LOAD_COUNT=50
./test/bdd/run_bdd_tests.sh mock --performance
```

### Métricas Monitoradas
- Tempo de resposta
- Throughput
- Uso de recursos
- Taxa de erro

## Debugging

### Logs Detalhados
```bash
export TEST_DEBUG=true
export TEST_LOG_LEVEL=debug
```

### Configuração de Debug
- Logs de requisições HTTP
- Traces de operações blockchain
- Detalhes de mocks
- Métricas de performance

## CI/CD Integration

### GitHub Actions
```yaml
- name: Run Mock Tests
  run: ./test/bdd/run_bdd_tests.sh mock --coverage

- name: Run Integration Tests
  run: ./test/bdd/run_bdd_tests.sh integration
  env:
    TEST_REAL_NATION_API_KEY: ${{ secrets.NATION_API_KEY }}
```

### Relatórios
- Cobertura de código
- Relatórios de performance
- Logs de execução
- Métricas de qualidade

## Troubleshooting

### Problemas Comuns

1. **Testes de integração falham**
   - Verificar API keys configuradas
   - Confirmar conectividade de rede
   - Validar configurações de ambiente

2. **Mocks não funcionam**
   - Verificar configuração de TEST_MODE
   - Confirmar inicialização dos mocks
   - Validar dados de teste

3. **Timeout em testes**
   - Aumentar TEST_TIMEOUT_SECONDS
   - Verificar conectividade de rede
   - Otimizar operações lentas

### Logs Úteis
```bash
# Logs detalhados
TEST_DEBUG=true ./test/bdd/run_bdd_tests.sh mock --verbose

# Logs de rede
TEST_LOG_LEVEL=debug ./test/bdd/run_bdd_tests.sh integration
```

## Contribuição

### Adicionando Novos Testes

1. Criar cenário no arquivo `.feature`
2. Implementar steps correspondentes
3. Adicionar mocks se necessário
4. Configurar tags apropriadas
5. Documentar no README

### Padrões de Código

- Usar português para descrições de cenários
- Implementar steps de forma clara e reutilizável
- Documentar configurações de mock
- Seguir convenções de nomenclatura

### Testes de Regressão

- Executar todos os modos antes de merge
- Validar tanto mocks quanto integração
- Verificar performance não degradou
- Confirmar cobertura de código mantida
