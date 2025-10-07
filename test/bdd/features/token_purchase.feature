# language: pt
Funcionalidade: Compra de Tokens do Bot (IACAI)
  Como um usuário com NFT de acesso
  Eu quero comprar tokens IACAI
  Para pagar por análises avançadas com LLM

  Contexto:
    Dado que estou autenticado com wallet "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
    E que possuo NFT "Pro Access"
    E que o contrato de token IACAI está ativo

  Cenário: Visualizar pacotes de tokens disponíveis
    Quando eu acesso a página "Buy Tokens"
    Então devo ver 4 pacotes disponíveis:
      | ID | Nome            | Tokens | Preço ETH | Preço USD | Desconto |
      | 1  | Starter Pack    | 100    | 0.005     | $10       | -        |
      | 2  | Power Pack      | 500    | 0.0225    | $45       | 10%      |
      | 3  | Pro Pack        | 1000   | 0.0425    | $85       | 15%      |
      | 4  | Enterprise Pack | 5000   | 0.1875    | $375      | 25%      |
    E devo ver meu saldo atual de tokens
    E devo ver a cotação atual: "1 IACAI = $0.10"

  Cenário: Comprar tokens com ETH na wallet
    Dado que minha wallet tem "0.1" ETH
    E meu saldo de tokens é 0 IACAI
    Quando eu seleciono o pacote "Power Pack" (500 tokens)
    E clico em "Buy with ETH"
    E confirmo a transação
    Então a transação deve ser processada
    E após confirmação meu saldo deve ser 500 IACAI
    E devo receber notificação "500 IACAI tokens adicionados!"

  Cenário: Comprar tokens com cartão de crédito via Privy Onramp
    Dado que minha wallet tem "0" ETH
    E eu seleciono o pacote "Pro Pack" (1000 tokens)
    Quando eu clico em "Buy with Card"
    Então devo ver a modal do Privy Onramp
    E devo ver:
      | Campo           | Valor                           |
      | Valor do pacote | 0.0425 ETH                      |
      | Preço estimado  | $85.00 USD                      |
      | Taxa de rede    | $2.50                           |
      | Taxa provedor   | $5.00                           |
      | Total a pagar   | $92.50                          |
    Quando eu completo o pagamento
    Então os tokens devem ser transferidos automaticamente
    E meu saldo final deve ser 1000 IACAI

  Cenário: Verificar saldo de tokens
    Dado que possuo 500 IACAI tokens
    Quando eu acesso "My Account"
    Então devo ver meu saldo: "500 IACAI"
    E devo ver o valor em USD: "~$50.00"
    E devo ver um botão "Buy More Tokens"

  Cenário: Gastar tokens em análise com LLM
    Dado que possuo 100 IACAI tokens
    E uma análise com LLM custa 5 tokens
    Quando eu faço uma análise de Terraform com LLM
    Então o sistema deve deduzir 5 tokens do meu saldo
    E meu saldo final deve ser 95 IACAI
    E devo ver no histórico: "Analysis with LLM: -5 IACAI"

  Cenário: Tentativa de análise sem tokens suficientes
    Dado que possuo 2 IACAI tokens
    E uma análise com LLM custa 5 tokens
    Quando eu tento fazer uma análise com LLM
    Então devo ver erro "Saldo insuficiente de tokens"
    E devo ver: "Você tem 2 IACAI, precisa de 5 IACAI"
    E devo ver botão "Buy Tokens"
    E a análise não deve ser executada

  Cenário: Histórico de transações de tokens
    Dado que fiz as seguintes operações:
      | Tipo              | Quantidade | Data       |
      | Compra (Pack 500) | +500       | 2024-01-10 |
      | Análise LLM       | -5         | 2024-01-11 |
      | Análise Full      | -15        | 2024-01-12 |
      | Compra (Pack 100) | +100       | 2024-01-15 |
    Quando eu acesso "Transaction History"
    Então devo ver todas as 4 transações
    E devo ver o saldo após cada operação
    E devo poder filtrar por tipo: "Compras" ou "Gastos"

  Cenário: Transferir tokens para outra wallet
    Dado que possuo 1000 IACAI tokens
    Quando eu acesso "Transfer Tokens"
    E insiro o endereço destino "0x8f3Cf7ad23Cd3CaDbD9735AFF958023239c6A063"
    E insiro a quantidade "100" tokens
    E confirmo a transferência
    Então os tokens devem ser transferidos
    E meu saldo deve ser 900 IACAI
    E a wallet destino deve receber 100 IACAI

  Cenário: Preços das operações em tokens
    Quando eu acesso "Pricing Info"
    Então devo ver a tabela de custos:
      | Operação            | Custo (IACAI) | Descrição                    |
      | Terraform Analysis  | 1             | Análise básica               |
      | Checkov Scan        | 2             | Scan de segurança            |
      | LLM Analysis        | 5             | Análise com IA               |
      | Preview Analysis    | 3             | Análise de terraform plan    |
      | Security Audit      | 10            | Auditoria completa           |
      | Cost Optimization   | 5             | Otimização de custos         |
      | Full Review         | 15            | Review completo com LLM      |

  Cenário: Desconto por volume na compra
    Dado que eu estou comprando o pacote "Enterprise Pack"
    Quando eu visualizo os detalhes do pacote
    Então devo ver:
      | Campo              | Valor                      |
      | Tokens             | 5000 IACAI                 |
      | Preço sem desconto | $500.00                    |
      | Desconto           | 25% (-$125.00)             |
      | Preço final        | $375.00                    |
      | Economia           | $125.00                    |
    E devo ver badge "BEST VALUE"
