# language: pt
Funcionalidade: Testes BDD com Mocks e Integração Real
  Como desenvolvedor
  Eu quero executar testes BDD que validem tanto o comportamento mockado quanto a integração real
  Para garantir que o sistema funciona corretamente em ambos os cenários

  Contexto:
    Dado que o ambiente de teste está configurado
    E que os mocks estão disponíveis
    E que o app ID padrão está configurado

  @mock @unit
  Cenário: Fluxo completo usando mocks - Login e Autenticação
    Dado que o sistema está em modo mock
    E que o MockPrivyClient está configurado
    Quando eu faço login com a wallet "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
    Então devo estar autenticado com sucesso
    E meu user ID deve ser "mock_user_123"
    E meu email deve ser "test@example.com"
    E minha wallet deve estar validada

  @mock @unit
  Cenário: Fluxo completo usando mocks - Compra de NFT
    Dado que estou autenticado em modo mock
    E que o MockNFTAccessManager está configurado
    Quando eu solicito mint de NFT tier "pro"
    Para a wallet "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
    Então o NFT deve ser mintado com sucesso
    E o token ID deve ser "mock_token_123"
    E o transaction hash deve ser "0xmock_tx_hash_123"
    E o status deve ser "success"

  @mock @unit
  Cenário: Fluxo completo usando mocks - Compra de Tokens
    Dado que estou autenticado em modo mock
    E que o MockBotTokenManager está configurado
    E meu saldo inicial é 1000 tokens
    Quando eu solicito mint de 500 tokens
    Para a wallet "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
    Então os tokens devem ser mintados com sucesso
    E meu saldo final deve ser 1500 tokens
    E o transaction hash deve ser "0xmock_token_tx_123"

  @mock @unit
  Cenário: Fluxo completo usando mocks - Análise de Código
    Dado que estou autenticado em modo mock
    E que o MockAnalysisService está configurado
    E meu saldo de tokens é 1000
    Quando eu submeto código Terraform para análise:
      """
      resource "aws_s3_bucket" "example" {
        bucket = "my-terraform-bucket"
        acl    = "public-read"
      }
      """
    E seleciono análise tipo "Full Review"
    Então a análise deve ser processada com sucesso
    E o score deve ser 85
    E deve haver 2 issues encontrados
    E deve haver pelo menos 1 issue de severidade "HIGH"
    E deve haver pelo menos 1 issue de severidade "MEDIUM"

  @mock @unit
  Cenário: Fluxo completo usando mocks - Onramp com Cartão
    Dado que estou autenticado em modo mock
    E que o MockPrivyOnrampManager está configurado
    Quando eu crio uma sessão de onramp:
      | Campo | Valor |
      | user_id | mock_user_123 |
      | wallet_address | 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb |
      | purpose | nft_access |
      | target_item_id | pro_tier |
      | source_currency | USD |
    Então a sessão deve ser criada com sucesso
    E o session ID deve ser "mock_session_123"
    E o quote deve ter source_amount "125.00"
    E o quote deve ter target_amount "0.05"
    E o provider deve ser "moonpay"

  @integration @real
  Cenário: Integração real com Privy - Validação de Configuração
    Dado que o sistema está em modo de integração real
    E que o app ID padrão "cmgh6un8w007bl10ci0tgitwp" está configurado
    E que a Base Network está configurada (Chain ID: 84531)
    E que o contrato NFT está deployado em "0x147e832418Cc06A501047019E956714271098b89"
    Quando eu verifico a configuração do sistema
    Então todas as configurações devem estar válidas
    E o Privy App ID deve estar correto
    E a RPC URL da Base deve estar acessível
    E o contrato NFT deve estar deployado

  @integration @real
  Cenário: Integração real com Privy - Validação de Usuário
    Dado que o sistema está em modo de integração real
    E que o PrivyClient está configurado com app ID real
    Quando eu busco informações do usuário "test_user_real"
    Então a resposta deve ser válida
    E deve conter informações do usuário
    E deve conter wallets associadas
    E deve conter dados de autenticação

  @integration @real
  Cenário: Integração real com Base Network - Verificação de Contrato
    Dado que o sistema está em modo de integração real
    E que o contrato NFT está deployado na Base Network
    Quando eu verifico o status do contrato
    Então o contrato deve estar ativo
    E deve responder a chamadas de leitura
    E deve ter as funções necessárias implementadas
    E deve estar no endereço correto

  @integration @real
  Cenário: Integração real com Nation.fun - Validação de LLM
    Dado que o sistema está em modo de integração real
    E que o Nation.fun está configurado
    E que a API key está válida
    Quando eu envio uma requisição de teste para o LLM
    Então a resposta deve ser válida
    E deve conter análise estruturada
    E deve seguir o formato esperado
    E deve processar código Terraform

  @integration @real
  Cenário: Integração real end-to-end - Fluxo Completo Simulado
    Dado que o sistema está em modo de integração real
    E que todos os serviços estão configurados
    E que uso dados de teste válidos
    Quando eu executo o fluxo completo:
      | Etapa | Ação | Parâmetros |
      | 1 | Login | wallet: 0xTestWallet |
      | 2 | Verificar NFT | tier: pro |
      | 3 | Verificar Tokens | amount: 1000 |
      | 4 | Análise | code: terraform_basic |
    Então cada etapa deve ser executada com sucesso
    E os resultados devem ser consistentes
    E não deve haver erros de integração
    E os logs devem mostrar comunicação real com serviços

  @mock @error_handling
  Cenário: Tratamento de erros em modo mock - Falha de Autenticação
    Dado que o sistema está em modo mock
    E que o MockPrivyClient está configurado para falhar
    Quando eu tento fazer login com wallet inválida
    Então deve retornar erro de autenticação
    E o erro deve ter código "AUTHENTICATION_FAILED"
    E a mensagem deve conter "Mock authentication failed"

  @mock @error_handling
  Cenário: Tratamento de erros em modo mock - Falha de Mint NFT
    Dado que estou autenticado em modo mock
    E que o MockNFTAccessManager está configurado para falhar
    Quando eu tento mintar NFT
    Então deve retornar erro de mint
    E o erro deve ter código "MINT_FAILED"
    E a mensagem deve conter "Mock mint failed"

  @mock @error_handling
  Cenário: Tratamento de erros em modo mock - Falha de Análise
    Dado que estou autenticado em modo mock
    E que o MockAnalysisService está configurado para falhar
    Quando eu submeto código para análise
    Então deve retornar erro de análise
    E o erro deve ter código "ANALYSIS_FAILED"
    E a mensagem deve conter "Mock analysis failed"

  @integration @performance
  Cenário: Teste de performance - Múltiplas requisições simultâneas
    Dado que o sistema está em modo de integração real
    E que o ambiente está configurado para testes de carga
    Quando eu executo 10 requisições simultâneas de análise
    Então todas as requisições devem ser processadas
    E o tempo médio de resposta deve ser menor que 30 segundos
    E não deve haver timeouts
    E os recursos devem ser utilizados eficientemente

  @mock @data_validation
  Cenário: Validação de dados mockados - Estrutura de Resposta
    Dado que o sistema está em modo mock
    E que todos os mocks estão configurados
    Quando eu executo operações que retornam dados
    Então os dados devem ter estrutura válida
    E devem conter todos os campos obrigatórios
    E os tipos de dados devem estar corretos
    E os valores devem estar dentro dos ranges esperados

  @integration @data_validation
  Cenário: Validação de dados reais - Estrutura de Resposta
    Dado que o sistema está em modo de integração real
    E que os serviços externos estão respondendo
    Quando eu executo operações que retornam dados reais
    Então os dados devem ter estrutura válida
    E devem conter todos os campos obrigatórios
    E os tipos de dados devem estar corretos
    E os valores devem estar dentro dos ranges esperados
    E devem ser consistentes com a documentação da API
