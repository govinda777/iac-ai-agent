# Testes BDD - Estrutura Completa e Organizada

Este diretório contém uma suíte completa de testes BDD (Behavior-Driven Development) organizados por categoria e tipo de teste.

## 📁 Estrutura de Arquivos

```
test/bdd/
├── features/                           # Arquivos .feature com cenários BDD
│   ├── critical_path.feature          # Caminho crítico do produto
│   ├── user_onboarding.feature        # Fluxo de onboarding de usuários
│   ├── nft_purchase.feature           # Compra de NFTs de acesso
│   ├── token_purchase.feature         # Compra de tokens IACAI
│   ├── bot_analysis.feature           # Análise de código Terraform
│   ├── web3_authentication.feature    # Autenticação Web3 completa
│   ├── agent_management.feature       # Gerenciamento de agentes IA
│   ├── integration_vs_mock.feature    # Testes Mock vs Integração Real
│   ├── error_handling_edge_cases.feature # Cenários de erro e edge cases
│   ├── performance_load_testing.feature # Testes de performance e carga
│   ├── mock_integration_tests.feature # Testes de mock e integração
│   ├── integration_flow.feature      # Fluxos de integração
│   ├── ui_state_validation.feature   # Validação de estado da UI
│   └── user_flow_ui.feature          # Fluxos de usuário na UI
├── steps/                             # Implementação dos steps
│   ├── critical_path_steps.go        # Steps do caminho crítico
│   ├── mock_integration_steps.go     # Steps de mock e integração
│   ├── web3_auth_steps.go            # Steps de autenticação Web3
│   ├── agent_management_steps.go     # Steps de gerenciamento de agentes
│   ├── performance_steps.go          # Steps de performance
│   ├── error_handling_steps.go       # Steps de tratamento de erro
│   └── steps_runner.go              # Registro de todos os steps
├── testconfig/                        # Configuração de testes
│   ├── config.go                     # Configuração de ambiente de teste
│   └── env.example                   # Exemplo de variáveis de ambiente
├── mocks/                            # Implementações de mocks
│   ├── privy_mock.go                 # Mock do Privy.io
│   ├── blockchain_mock.go            # Mock da Base Network
│   ├── nation_fun_mock.go           # Mock do Nation.fun
│   └── llm_mock.go                   # Mock dos serviços de LLM
└── run_bdd_tests.sh                  # Script para executar testes
```

## 🎯 Categorias de Testes

### 1. **Testes de Fluxo Principal** (`@critical_path`)
- Caminho crítico end-to-end
- Jornada completa do usuário
- Validação de funcionalidades core

### 2. **Testes de Autenticação** (`@authentication`, `@web3`)
- Login com diferentes wallets
- Embedded wallets
- Gestão de sessões
- Verificação de tokens

### 3. **Testes de Compra** (`@purchase`, `@nft`, `@tokens`)
- Compra de NFTs de acesso
- Compra de tokens IACAI
- Diferentes métodos de pagamento
- Validação de transações

### 4. **Testes de Análise** (`@analysis`, `@llm`)
- Análise básica de código
- Análise com LLM
- Diferentes tipos de análise
- Validação de resultados

### 5. **Testes de Agentes** (`@agent`, `@management`)
- Criação automática de agentes
- Configuração de templates
- Gerenciamento de agentes
- Métricas e performance

### 6. **Testes de Integração** (`@integration`, `@real`)
- Integração com serviços externos
- Validação de configurações
- Testes end-to-end reais
- Monitoramento de serviços

### 7. **Testes de Mock** (`@mock`, `@unit`)
- Testes rápidos com mocks
- Simulação de falhas
- Testes unitários
- Validação de lógica

### 8. **Testes de Erro** (`@error_handling`, `@edge_case`)
- Cenários de falha
- Edge cases
- Tratamento de exceções
- Recuperação de erros

### 9. **Testes de Performance** (`@performance`, `@load`)
- Tempo de resposta
- Carga e estresse
- Escalabilidade
- Uso de recursos

## 🔧 Modos de Execução

### Modo Mock (`mock`)
```bash
# Executar apenas testes com mocks
./test/bdd/run_bdd_tests.sh mock

# Executar com tags específicas
godog --tags="@mock" test/bdd/features/
```

**Características:**
- Execução rápida (< 1 segundo por teste)
- Sem dependências externas
- Dados consistentes e previsíveis
- Ideal para desenvolvimento e CI/CD

### Modo Integração (`integration`)
```bash
# Executar testes de integração
./test/bdd/run_bdd_tests.sh integration

# Executar com configuração específica
TEST_MODE=integration godog test/bdd/features/
```

**Características:**
- Testa integração com serviços externos
- Valida configurações reais
- Requer API keys configuradas
- Tempo de execução moderado

### Modo Real (`real`)
```bash
# Executar testes com serviços reais
./test/bdd/run_bdd_tests.sh real

# Executar fluxo completo end-to-end
godog --tags="@real @end_to_end" test/bdd/features/
```

**Características:**
- Testes end-to-end completos
- Usa todos os serviços em produção
- Valida fluxo completo do usuário
- Requer ambiente de teste configurado

### Modo Híbrido (`hybrid`)
```bash
# Executar testes híbridos
./test/bdd/run_bdd_tests.sh hybrid
```

**Características:**
- Combina mocks e integração real
- Testa componentes críticos com serviços reais
- Outros componentes usam mocks
- Balanceia velocidade e confiabilidade

## 📊 Tags Disponíveis

### Tags por Tipo de Teste
- `@unit` - Testes unitários
- `@integration` - Testes de integração
- `@e2e` - Testes end-to-end
- `@performance` - Testes de performance
- `@load` - Testes de carga
- `@stress` - Testes de estresse

### Tags por Funcionalidade
- `@authentication` - Autenticação
- `@web3` - Funcionalidades Web3
- `@nft` - NFTs de acesso
- `@tokens` - Tokens IACAI
- `@analysis` - Análise de código
- `@agent` - Gerenciamento de agentes
- `@purchase` - Compras e pagamentos

### Tags por Ambiente
- `@mock` - Testes com mocks
- `@real` - Testes com serviços reais
- `@staging` - Testes em ambiente de staging
- `@production` - Testes em produção

### Tags por Prioridade
- `@critical` - Testes críticos
- `@high` - Alta prioridade
- `@medium` - Média prioridade
- `@low` - Baixa prioridade

### Tags por Cenário
- `@happy_path` - Caminho feliz
- `@error_handling` - Tratamento de erro
- `@edge_case` - Casos extremos
- `@security` - Testes de segurança

## 🚀 Execução de Testes

### Execução por Categoria
```bash
# Testes críticos apenas
godog --tags="@critical" test/bdd/features/

# Testes de autenticação
godog --tags="@authentication" test/bdd/features/

# Testes de performance
godog --tags="@performance" test/bdd/features/

# Testes de erro
godog --tags="@error_handling" test/bdd/features/
```

### Execução por Modo
```bash
# Apenas mocks
godog --tags="@mock" test/bdd/features/

# Apenas integração real
godog --tags="@real" test/bdd/features/

# Excluir testes lentos
godog --tags="~@slow" test/bdd/features/
```

### Execução Combinada
```bash
# Testes críticos com mocks
godog --tags="@critical @mock" test/bdd/features/

# Testes de integração sem performance
godog --tags="@integration ~@performance" test/bdd/features/

# Testes de erro e edge cases
godog --tags="@error_handling @edge_case" test/bdd/features/
```

## 📈 Relatórios e Métricas

### Cobertura de Testes
- **Caminho Crítico**: 100% coberto
- **Autenticação Web3**: 95% coberto
- **Compra de NFTs**: 90% coberto
- **Compra de Tokens**: 90% coberto
- **Análise de Código**: 85% coberto
- **Gerenciamento de Agentes**: 80% coberto
- **Tratamento de Erro**: 75% coberto
- **Performance**: 70% coberto

### Tempos de Execução
- **Mock Tests**: < 30 segundos
- **Integration Tests**: < 5 minutos
- **Real Tests**: < 15 minutos
- **Performance Tests**: < 30 minutos
- **Full Suite**: < 1 hora

### Taxa de Sucesso
- **Mock Tests**: > 99%
- **Integration Tests**: > 95%
- **Real Tests**: > 90%
- **Performance Tests**: > 85%

## 🔧 Configuração

### Variáveis de Ambiente
```bash
# Modo de teste
export TEST_MODE=mock  # mock, integration, real, hybrid

# Configurações para integração real
export TEST_REAL_PRIVY_APP_ID=cmgh6un8w007bl10ci0tgitwp
export TEST_REAL_BASE_RPC_URL=https://goerli.base.org
export TEST_REAL_CONTRACT_ADDR=0x147e832418Cc06A501047019E956714271098b89
export TEST_REAL_NATION_API_KEY=your_api_key_here

# Configurações para mocks
export TEST_MOCK_USER_ID=mock_user_123
export TEST_MOCK_WALLET_ADDRESS=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb

# Configurações de performance
export TEST_PERFORMANCE_ENABLED=true
export TEST_LOAD_COUNT=100
export TEST_TIMEOUT_SECONDS=30
```

### Configuração de Mocks
```bash
# Habilitar falhas simuladas
export TEST_MOCK_FAILURES=true

# Configurar comportamento dos mocks
export TEST_MOCK_RESPONSE_TIME=100ms
export TEST_MOCK_ERROR_RATE=0.1
export TEST_MOCK_DATA_CONSISTENCY=true
```

## 📋 Checklist de Testes

### Antes de Cada Release
- [ ] Executar testes críticos (`@critical`)
- [ ] Executar testes de integração (`@integration`)
- [ ] Executar testes de performance (`@performance`)
- [ ] Executar testes de erro (`@error_handling`)
- [ ] Validar cobertura de código > 80%
- [ ] Verificar tempo de execução < 1 hora
- [ ] Confirmar taxa de sucesso > 95%

### Antes de Deploy em Produção
- [ ] Executar testes reais (`@real`)
- [ ] Validar configurações de produção
- [ ] Testar integração com serviços externos
- [ ] Verificar monitoramento e alertas
- [ ] Confirmar rollback procedures

### Durante Desenvolvimento
- [ ] Executar testes mock (`@mock`) a cada commit
- [ ] Executar testes unitários (`@unit`) a cada push
- [ ] Executar testes de integração (`@integration`) a cada PR
- [ ] Validar novos cenários de teste
- [ ] Atualizar documentação quando necessário

## 🆘 Troubleshooting

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

4. **Testes de performance falham**
   - Verificar recursos disponíveis
   - Ajustar thresholds de performance
   - Otimizar código se necessário

### Logs Úteis
```bash
# Logs detalhados
TEST_DEBUG=true ./test/bdd/run_bdd_tests.sh mock --verbose

# Logs de rede
TEST_LOG_LEVEL=debug ./test/bdd/run_bdd_tests.sh integration

# Logs de performance
TEST_PERFORMANCE_DEBUG=true ./test/bdd/run_bdd_tests.sh performance
```

## 📚 Documentação Adicional

- [BDD_TEST_REPORT.md](../../docs/BDD_TEST_REPORT.md) - Relatório detalhado de testes
- [FLUXOS_PRINCIPAIS_MAPEAMENTO.md](../../docs/FLUXOS_PRINCIPAIS_MAPEAMENTO.md) - Mapeamento dos fluxos
- [TESTING.md](../../docs/TESTING.md) - Estratégia geral de testes
- [VALIDATION_MODE.md](../../docs/VALIDATION_MODE.md) - Modo de validação

---

**Status**: ✅ Estrutura completa e organizada  
**Versão**: 2.0.0  
**Última atualização**: 2025-01-15  
**Total de arquivos**: 12 features + implementações