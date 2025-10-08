# language: pt
Funcionalidade: Testes de Integração Real vs Mock - Sistema Completo
  Como desenvolvedor
  Eu quero ter testes que validem tanto mocks quanto integração real
  Para garantir que o sistema funciona em todos os ambientes

  Contexto:
    Dado que o sistema suporta modo mock e modo real
    E que as configurações de ambiente estão definidas

  @integration @real @web3
  Cenário: Integração real com Privy.io - Validação de Configuração
    Dado que o sistema está em modo de integração real
    E que o app ID padrão "cmgh6un8w007bl10ci0tgitwp" está configurado
    E que a Base Network está conectada
    Quando eu verifico a configuração do sistema
    Então todas as configurações devem estar válidas:
      | Serviço        | Status | Detalhes                           |
      | Privy.io       | ✓      | App ID válido, conectado           |
      | Base Network   | ✓      | RPC funcionando, Chain ID 84531   |
      | Smart Contracts| ✓      | NFTs e Tokens deployados         |
      | Nation.fun     | ✓      | API key válida, conectado         |
    E devo ver "Sistema configurado e pronto para produção"

  @integration @real @authentication
  Cenário: Autenticação real com Privy.io
    Dado que o sistema está em modo de integração real
    E que tenho uma wallet real na Base Network
    Quando eu faço login com minha wallet real
    E o Privy.io valida minha autenticação
    Então devo estar autenticado com dados reais:
      | Campo           | Valor Real                    |
      | User ID         | privy_user_123456            |
      | Wallet Address  | 0x742d35Cc6634C0532925a3b8... |
      | Email           | user@example.com             |
      | Session Token   | eyJhbGciOiJIUzI1NiIsInR5cCI6... |
    E minha sessão deve ser válida por 24 horas

  @integration @real @blockchain
  Cenário: Operações reais na Base Network
    Dado que o sistema está em modo de integração real
    E que tenho ETH real na Base Network
    E que os contratos estão deployados
    Quando eu verifico meu saldo de NFT
    Então devo receber dados reais da blockchain:
      | Campo           | Valor Real                    |
      | NFT Balance     | 1                             |
      | Token ID        | 12345                         |
      | Tier            | Pro Access                    |
      | Contract Address| 0x147e832418Cc06A501047019E956714271098b89 |
    E o gas usado deve ser < 0.001 ETH

  @integration @real @nation_fun
  Cenário: Análise real com Nation.fun
    Dado que o sistema está em modo de integração real
    E que tenho API key válida do Nation.fun
    E que possuo tokens suficientes
    Quando eu submeto código Terraform real para análise
    Então o Nation.fun deve processar a análise
    E devo receber resultado real:
      | Campo           | Valor Real                    |
      | Analysis ID     | nation_analysis_789          |
      | Processing Time | 2.3 segundos                 |
      | Score           | 87/100                       |
      | Issues Found    | 3                            |
      | Suggestions     | 5                            |
    E os tokens devem ser debitados da minha conta real

  @mock @unit @fast
  Cenário: Testes rápidos com mocks - Fluxo completo
    Dado que o sistema está em modo mock
    E que o MockPrivyClient está configurado
    E que o MockNFTAccessManager está ativo
    E que o MockBotTokenManager está configurado
    Quando eu executo o fluxo completo:
      1. Login com wallet mockada
      2. Verificação de NFT mockado
      3. Compra de tokens mockada
      4. Análise com mocks
    Então todos os passos devem ser executados em < 1 segundo
    E devo receber dados mockados consistentes:
      | Serviço        | Dados Mockados                |
      | Privy          | mock_user_123                 |
      | NFT            | Pro Access (Tier 2)           |
      | Tokens         | 1000 IACAI                    |
      | Análise        | Score: 85, Issues: 2          |
    E nenhuma chamada externa deve ser feita

  @mock @error_simulation
  Cenário: Simulação de falhas com mocks
    Dado que o sistema está em modo mock
    E que as falhas simuladas estão habilitadas
    Quando eu tento fazer login
    E o MockPrivyClient simula falha de autenticação
    Então devo receber erro simulado:
      | Campo           | Valor                         |
      | Error Code     | PRIVY_AUTH_FAILED            |
      | Error Message   | "Token inválido ou expirado"  |
      | Retry After    | 30 segundos                   |
    E o sistema deve estar preparado para lidar com o erro

  @mock @performance @load
  Cenário: Teste de carga com mocks
    Dado que o sistema está em modo mock
    E que o MockAnalysisService está configurado
    Quando eu executo 100 análises simultâneas
    Então todas devem ser processadas em < 5 segundos
    E o sistema deve manter performance:
      | Métrica         | Valor Esperado |
      | Throughput      | > 20 req/s     |
      | Latência P95    | < 200ms       |
      | CPU Usage       | < 50%         |
      | Memory Usage    | < 100MB       |
    E nenhum erro deve ocorrer

  @integration @real @end_to_end
  Cenário: Fluxo completo end-to-end em ambiente real
    Dado que o sistema está em modo de integração real
    E que tenho wallet real com ETH suficiente
    E que todos os serviços externos estão funcionando
    Quando eu executo o fluxo completo:
      1. Login real com Privy.io
      2. Compra real de NFT na Base Network
      3. Compra real de tokens IACAI
      4. Análise real com Nation.fun
      5. Recebimento de resultado real
    Então todo o fluxo deve funcionar:
      | Etapa           | Status | Tempo | Custo Real |
      | Login           | ✓      | 2s    | $0.00      |
      | NFT Purchase    | ✓      | 30s   | $125.00    |
      | Token Purchase  | ✓      | 15s   | $45.00     |
      | Analysis        | ✓      | 5s    | $0.50      |
      | **Total**       | ✓      | 52s   | $170.50    |
    E todos os dados devem ser persistidos na blockchain

  @mock @integration @hybrid
  Cenário: Teste híbrido - Mock + Integração Real
    Dado que o sistema está em modo híbrido
    E que Privy.io está mockado
    E que Base Network está real
    E que Nation.fun está real
    Quando eu executo análise
    Então:
      - Autenticação deve usar mock (rápido)
      - Verificação de NFT deve usar blockchain real
      - Análise deve usar Nation.fun real
      - Resultado deve ser persistido na blockchain real
    E o tempo total deve ser < 10 segundos
    E apenas os serviços críticos devem ser testados

  @integration @real @monitoring
  Cenário: Monitoramento de serviços externos
    Dado que o sistema está em modo de integração real
    Quando eu verifico o status dos serviços externos
    Então devo receber métricas reais:
      | Serviço        | Status | Latência | Uptime | Rate Limit |
      | Privy.io       | ✓      | 150ms    | 99.9%  | 1000/min   |
      | Base Network   | ✓      | 200ms    | 99.8%  | 100/min    |
      | Nation.fun     | ✓      | 300ms    | 99.5%  | 50/min     |
      | OpenAI API     | ✓      | 500ms    | 99.9%  | 60/min     |
    E alertas devem ser gerados se algum serviço estiver degradado

  @mock @integration @ci_cd
  Cenário: Execução em pipeline CI/CD
    Dado que estou executando testes em pipeline CI/CD
    E que não tenho acesso a serviços externos
    Quando eu executo a suíte de testes
    Então:
      - Testes unitários devem usar mocks (100% cobertura)
      - Testes de integração devem ser pulados
      - Testes de smoke devem usar mocks
      - Deploy deve ser bloqueado se testes falharem
    E o pipeline deve completar em < 5 minutos
    E relatório de cobertura deve ser gerado

  @integration @real @production
  Cenário: Validação pré-produção
    Dado que estou preparando deploy para produção
    E que todos os serviços estão configurados
    Quando eu executo validação completa
    Então todos os serviços devem estar funcionando:
      | Validação           | Status | Detalhes                    |
      | Privy.io Config     | ✓      | App ID válido               |
      | Base Network        | ✓      | Contratos deployados        |
      | Nation.fun API      | ✓      | API key válida              |
      | OpenAI API          | ✓      | API key válida              |
      | Database            | ✓      | Conectado e migrado         |
      | Redis Cache         | ✓      | Conectado                   |
      | Monitoring          | ✓      | Alertas configurados        |
    E sistema deve estar pronto para produção
