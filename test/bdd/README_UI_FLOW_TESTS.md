# Testes BDD - Fluxo de UI

Este diretório contém testes BDD (Behavior Driven Development) específicos para validar o fluxo completo de UI do usuário no IaC AI Agent.

## Estrutura dos Testes

### Features (Cenários)
- **`user_flow_ui.feature`** - Fluxo completo de UI do usuário
- **`ui_state_validation.feature`** - Validação de estado da interface
- **`integration_flow.feature`** - Integração completa do fluxo

### Steps (Implementações)
- **`ui_flow_steps.go`** - Implementação dos steps para testes de fluxo de UI
- **`ui_test_config.go`** - Configuração específica para testes de UI

### Scripts
- **`run_ui_flow_tests.sh`** - Script principal para executar os testes
- **`ui_tests.env`** - Configurações de ambiente para os testes

## Fluxo Testado

Os testes cobrem o fluxo completo de 5 etapas:

1. **Conectar Wallet** - Autenticação via Privy.io
2. **Adquirir NFT** - Compra de NFT de acesso (Basic/Pro/Enterprise)
3. **Comprar Tokens** - Aquisição de tokens IACAI para pagamento
4. **Submeter Código** - Envio de código Terraform para análise
5. **Receber Sugestões** - Obtenção de análise e recomendações

## Tipos de Teste

### @ui_flow
Testes do fluxo principal de UI, incluindo:
- Navegação entre etapas
- Validação de estados intermediários
- Feedback visual para o usuário
- Tratamento de erros

### @ui_state
Testes de validação de estado da interface:
- Estados de autenticação
- Status de NFT e tokens
- Habilitação progressiva de funcionalidades
- Consistência de dados

### @integration
Testes de integração completa:
- Integração entre todos os sistemas
- Performance do fluxo integrado
- Segurança durante o fluxo
- Escalabilidade

### @mock
Testes usando simulações:
- Autenticação mockada
- Transações simuladas
- APIs mockadas
- Dados de teste

## Executando os Testes

### Pré-requisitos
```bash
# Instalar godog se não estiver instalado
go install github.com/cucumber/godog/cmd/godog@latest

# Verificar se o sistema está rodando
curl http://localhost:8080/health
```

### Executar Todos os Testes
```bash
./run_ui_flow_tests.sh
```

### Executar Testes Específicos
```bash
# Apenas testes mock
./run_ui_flow_tests.sh --mode mock

# Apenas testes de integração
./run_ui_flow_tests.sh --mode integration

# Apenas testes de fluxo de UI
./run_ui_flow_tests.sh --mode ui-flow

# Apenas testes de estado de UI
./run_ui_flow_tests.sh --mode ui-state
```

### Executar com Relatório
```bash
./run_ui_flow_tests.sh --report
```

## Configuração

### Variáveis de Ambiente
As configurações são definidas no arquivo `ui_tests.env`:

```bash
# Carregar configurações
source ui_tests.env

# Ou definir manualmente
export UI_TEST_MOCK_MODE="true"
export UI_TEST_BASE_URL="http://localhost:8080"
export UI_TEST_HEADLESS="true"
```

### Modos de Teste

#### Modo Mock (Padrão)
- Usa simulações para Web3 e APIs
- Mais rápido e confiável
- Ideal para desenvolvimento

#### Modo Integração
- Usa serviços reais
- Mais lento mas mais realista
- Ideal para validação final

## Cenários de Teste

### Fluxo Completo
```gherkin
Cenário: Usuário completa o fluxo completo através da UI
  Dado que acesso a página inicial do IaC AI Agent
  Quando eu clico no botão "Conectar Wallet"
  E seleciono "MetaMask" como provedor
  E autorizo a conexão no MetaMask
  Então devo estar autenticado com sucesso
  # ... continua com todas as etapas
```

### Validação de Estado
```gherkin
Cenário: Estados de autenticação refletidos na UI
  Dado que não estou autenticado
  Quando acesso qualquer página do sistema
  Então devo ver no cabeçalho:
    - Botão "Conectar Wallet" visível
    - Seção de perfil do usuário oculta
```

### Tratamento de Erros
```gherkin
Cenário: Validação de saldo insuficiente de tokens
  Dado que meu saldo de tokens é "2 IACAI"
  Quando eu seleciono "Full Review" (15 tokens)
  E clico em "Analisar Código"
  Então devo ver a mensagem "Saldo insuficiente"
```

## Relatórios

### Relatório JUnit
Os testes geram relatórios JUnit em formato XML:
```
test/bdd/reports/junit_user_flow_ui.xml
test/bdd/reports/junit_ui_state_validation.xml
test/bdd/reports/junit_integration_flow.xml
```

### Relatório HTML
Quando executado com `--report`, gera relatório HTML:
```
test/bdd/reports/ui_flow_test_report.html
```

### Screenshots
Screenshots são capturados durante os testes:
```
test/bdd/screenshots/
├── login_success.png
├── nft_purchase.png
├── token_purchase.png
├── analysis_submission.png
└── results_display.png
```

## Debugging

### Logs Detalhados
```bash
export LOG_LEVEL="debug"
./run_ui_flow_tests.sh
```

### Modo Não-Headless
```bash
export UI_TEST_HEADLESS="false"
./run_ui_flow_tests.sh
```

### Screenshots em Caso de Falha
Os screenshots são automaticamente capturados quando um teste falha.

## Integração com CI/CD

### GitHub Actions
```yaml
- name: Run UI Flow Tests
  run: |
    cd test/bdd
    ./run_ui_flow_tests.sh --mode mock --report
```

### Jenkins
```groovy
stage('UI Flow Tests') {
    steps {
        sh 'cd test/bdd && ./run_ui_flow_tests.sh --mode mock'
    }
    post {
        always {
            publishHTML([
                allowMissing: false,
                alwaysLinkToLastBuild: true,
                keepAll: true,
                reportDir: 'test/bdd/reports',
                reportFiles: 'ui_flow_test_report.html',
                reportName: 'UI Flow Test Report'
            ])
        }
    }
}
```

## Troubleshooting

### Problemas Comuns

#### Teste Falha por Timeout
```bash
# Aumentar timeout
export UI_TEST_TIMEOUT="60"
```

#### Erro de Conexão
```bash
# Verificar se o sistema está rodando
curl http://localhost:8080/health

# Usar modo mock
export UI_TEST_MOCK_MODE="true"
```

#### Screenshots Não Capturados
```bash
# Verificar permissões
chmod 755 test/bdd/screenshots/

# Verificar se o diretório existe
mkdir -p test/bdd/screenshots/
```

## Contribuindo

### Adicionando Novos Cenários
1. Crie um novo arquivo `.feature` ou adicione ao existente
2. Implemente os steps necessários em `ui_flow_steps.go`
3. Execute os testes para validar
4. Atualize a documentação

### Padrões de Nomenclatura
- Features: `descriptive_name.feature`
- Steps: `action_expected_result`
- Tags: `@category @subcategory`

### Boas Práticas
- Use dados de teste consistentes
- Mantenha cenários independentes
- Documente cenários complexos
- Valide estados intermediários
