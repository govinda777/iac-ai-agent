# language: pt
Funcionalidade: Compra de NFT de Acesso
  Como um usuário autenticado
  Eu quero comprar um NFT de acesso
  Para poder usar o bot de análise de IaC

  Contexto:
    Dado que estou autenticado com wallet "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
    E que minha wallet está na Base Network
    E que os tiers de NFT estão disponíveis

  Cenário: Visualizar tiers de acesso disponíveis
    Quando eu acesso a página de "Pricing"
    Então devo ver 3 tiers disponíveis:
      | Tier ID | Nome             | Preço    | Benefícios                              |
      | 1       | Basic Access     | 0.01 ETH | Análises ilimitadas, Checkov, Suporte   |
      | 2       | Pro Access       | 0.05 ETH | + LLM, Preview, Drift, Priority Support |
      | 3       | Enterprise Access| 0.2 ETH  | + API dedicada, Custom KB, SLA 24/7     |
    E cada tier deve exibir o preço em ETH e USD

  Cenário: Comprar NFT Basic Access pagando com ETH na wallet
    Dado que minha wallet tem saldo de "0.1" ETH
    E eu estou visualizando o tier "Basic Access"
    Quando eu clico em "Buy with ETH"
    E confirmo a transação no MetaMask
    Então a transação deve ser enviada para a blockchain
    E devo ver o status "Transaction Pending"
    E após a confirmação devo receber o NFT
    E devo ver a mensagem "Acesso ativado! Você pode usar o bot agora"
    E meu saldo de NFT deve ser 1

  Cenário: Comprar NFT Pro Access usando Privy Onramp (sem ETH na wallet)
    Dado que minha wallet tem saldo de "0" ETH
    E eu estou visualizando o tier "Pro Access"
    Quando eu clico em "Buy with Card"
    Então devo ver a modal do Privy Onramp
    E devo ver o valor necessário: "0.05 ETH (~$125 USD)"
    Quando eu seleciono "Credit Card" como método de pagamento
    E insiro os dados do cartão:
      | Campo         | Valor                |
      | Número        | 4242 4242 4242 4242  |
      | Validade      | 12/25                |
      | CVV           | 123                  |
      | Nome          | John Doe             |
    E concluo o pagamento
    Então devo ver "Payment Processing"
    E em até 10 minutos o ETH deve chegar na wallet
    E o NFT deve ser automaticamente mintado
    E devo receber notificação "NFT Pro Access adquirido!"

  Cenário: Comprar NFT Pro Access usando PIX (Brasil)
    Dado que estou no Brasil
    E minha moeda preferida é "BRL"
    E eu estou visualizando o tier "Pro Access"
    Quando eu clico em "Buy with PIX"
    Então devo ver o valor em BRL: "R$ 625,00"
    E devo ver um QR Code do PIX
    Quando eu escanei o QR Code e pago
    Então devo ver "Aguardando confirmação do pagamento"
    E após confirmação do banco o ETH deve chegar na wallet
    E o NFT deve ser automaticamente mintado

  Cenário: Tentativa de compra com saldo insuficiente
    Dado que minha wallet tem saldo de "0.001" ETH
    E eu estou visualizando o tier "Pro Access" (0.05 ETH)
    Quando eu clico em "Buy with ETH"
    Então devo ver a mensagem de erro "Saldo insuficiente"
    E devo ver a opção "Buy with Card" sugerida
    E não deve ser possível prosseguir com ETH

  Cenário: Upgrade de tier Basic para Pro
    Dado que eu possuo NFT "Basic Access"
    Quando eu acesso "My Access"
    E clico em "Upgrade to Pro"
    Então devo ver o preço de upgrade: "0.04 ETH"
    E devo ver "(diferença entre tiers)"
    Quando eu confirmo o upgrade
    E pago a diferença
    Então meu NFT deve ser atualizado para "Pro Access"
    E devo ter acesso aos novos benefícios imediatamente

  Cenário: Transferir NFT de acesso para outra wallet
    Dado que eu possuo NFT "Pro Access" (Token ID: 123)
    Quando eu acesso "My Access"
    E clico em "Transfer NFT"
    E insiro o endereço destino "0x8f3Cf7ad23Cd3CaDbD9735AFF958023239c6A063"
    E confirmo a transferência
    Então o NFT deve ser transferido
    E eu não devo mais ter acesso ao bot
    E a wallet destino deve ganhar acesso

  Cenário: Verificar acesso antes de usar o bot
    Dado que possuo NFT "Pro Access"
    Quando eu tento fazer uma análise de Terraform
    Então o sistema deve verificar meu NFT
    E deve confirmar que tenho tier "Pro" ou superior
    E deve permitir a análise

  Cenário: Acesso negado sem NFT
    Dado que não possuo nenhum NFT de acesso
    Quando eu tento fazer uma análise de Terraform
    Então devo ver erro "Acesso negado"
    E devo ser redirecionado para página de pricing
    E devo ver a mensagem "Adquira um NFT de acesso para usar o bot"
