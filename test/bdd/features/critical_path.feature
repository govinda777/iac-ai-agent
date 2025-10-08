# language: pt
Funcionalidade: Caminho Crítico do Produto
  Como um novo usuário
  Eu quero completar todo o fluxo principal do produto
  Para validar que o sistema funciona end-to-end

  Contexto:
    Dado que o serviço está disponível
    E que o serviço Privy está configurado
    E que os contratos na Base Network estão deployados

  @critical_path
  Cenário: Novo usuário completa o caminho crítico end-to-end
    # Etapa 1: Login com Privy
    Dado que sou um novo usuário sem conta
    Quando eu acesso a página inicial
    E clico em "Connect Wallet"
    E seleciono "MetaMask" como provedor
    E autorizo a conexão no MetaMask
    Então devo estar autenticado com sucesso
    E meu endereço de wallet deve estar visível no cabeçalho

    # Etapa 2: Compra de NFT de acesso
    Quando eu navego para a página "Pricing"
    E seleciono o tier "Pro Access"
    E clico em "Buy with ETH"
    E confirmo a transação na wallet
    Então a transação deve ser processada
    E após confirmação devo receber o NFT Pro Access
    E devo ver o status "Access Granted" no dashboard

    # Etapa 3: Compra de tokens IACAI
    Quando eu navego para "Buy Tokens"
    E seleciono o pacote "Power Pack" (500 tokens)
    E clico em "Buy with ETH" 
    E confirmo a transação na wallet
    Então os tokens devem ser adicionados à minha conta
    E meu saldo deve mostrar "500 IACAI"

    # Etapa 4: Consulta de análise de código
    Quando eu navego para "New Analysis"
    E submeto o seguinte código Terraform:
      """
      resource "aws_s3_bucket" "example" {
        bucket = "my-terraform-bucket"
        acl    = "public-read"  # Problema de segurança intencional
      }
      
      resource "aws_instance" "web" {
        ami           = "ami-0c55b159cbfafe1f0"
        instance_type = "t3.micro"
      }
      """
    E seleciono "Full Review" (15 tokens)
    E clico em "Analyze"
    Então o sistema deve processar minha análise
    E deve debitar 15 IACAI tokens da minha conta
    E meu saldo final deve ser 485 IACAI

    # Etapa 5: Recebimento de sugestões
    E devo receber um relatório de análise completo em até 30 segundos
    E o relatório deve conter as seguintes seções:
      | Seção                    | Status |
      | Executive Summary        | ✓      |
      | Security Issues          | ✓      |
      | Best Practices           | ✓      |
      | Cost Optimization        | ✓      |
      | Detailed Findings        | ✓      |
    E pelo menos uma sugestão de segurança sobre o "public-read" ACL
    E cada sugestão deve conter:
      - Descrição do problema
      - Impacto potencial
      - Código correto recomendado
      - Referências adicionais

  @critical_path @mock
  Cenário: Caminho crítico usando mocks e simulação
    # Etapa 1: Login com Privy (mock)
    Dado que o sistema está em modo de validação
    E a API do Privy está mockada
    Quando eu faço login com a wallet mockada "0xMockedWalletAddress"
    Então devo estar autenticado com sucesso
    E devo ver "Wallet conectada: 0xMocked...Address"

    # Etapa 2: Compra de NFT de acesso (mock)
    Quando eu seleciono o tier "Pro Access" para compra
    E confirmo a compra usando a função de simulação
    Então o sistema deve simular a transação com sucesso
    E devo receber um NFT Pro Access simulado
    E meu status de acesso deve ser "Pro Access"

    # Etapa 3: Compra de tokens IACAI (mock)
    Quando eu seleciono o pacote "Power Pack" (500 tokens)
    E confirmo a compra usando a função de simulação
    Então o sistema deve simular a transação com sucesso
    E meu saldo simulado deve ser "500 IACAI"

    # Etapa 4: Consulta de análise (real)
    Quando eu submeto código Terraform para análise:
      """
      resource "aws_s3_bucket" "example" {
        bucket = "my-terraform-bucket"
      }
      """
    E seleciono "Full Review"
    Então o sistema deve processar a análise real
    E debitar tokens do saldo simulado

    # Etapa 5: Recebimento de sugestões (real)
    E devo receber sugestões reais de melhoria
    E o relatório deve ser completo e detalhado
    E meu saldo final simulado deve refletir o custo da análise

  @critical_path @alternative
  Cenário: Caminho crítico com compra via cartão de crédito
    # Similar ao primeiro cenário, mas usando onramp de cartão de crédito
    # para comprar NFT e tokens em vez de ETH diretamente na wallet
    
    # Etapa 1: Login com Privy
    Dado que sou um novo usuário sem conta
    Quando eu me autentico com email e senha
    Então uma embedded wallet deve ser criada para mim
    
    # Etapa 2: Compra de NFT com cartão de crédito
    Quando eu navego para a página "Pricing"
    E seleciono o tier "Pro Access"
    E clico em "Buy with Card"
    E completo o processo de pagamento com cartão
    Então o NFT deve ser mintado para minha embedded wallet
    
    # Etapa 3-5: Continua como no primeiro cenário
    # ...

  @edge_cases
  Esquema do Cenário: Validação de acesso por tier nos diferentes serviços
    Dado que estou autenticado com wallet
    E possuo NFT "<Tier>" de acesso
    Quando tento acessar o serviço "<Serviço>"
    Então devo ver o resultado "<Resultado>"
    
    Exemplos:
      | Tier           | Serviço           | Resultado         |
      | Basic Access   | Análise Básica    | Sucesso           |
      | Basic Access   | Análise com LLM   | Acesso Negado     |
      | Basic Access   | Preview Analysis  | Acesso Negado     |
      | Pro Access     | Análise Básica    | Sucesso           |
      | Pro Access     | Análise com LLM   | Sucesso           |
      | Pro Access     | API Dedicada      | Acesso Negado     |
      | Enterprise     | Análise com LLM   | Sucesso           |
      | Enterprise     | API Dedicada      | Sucesso           |
      | Sem NFT        | Análise Básica    | Acesso Negado     |
