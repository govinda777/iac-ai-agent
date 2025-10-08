# language: pt
Funcionalidade: Fluxo de Autenticação Web3 Completo
  Como um usuário
  Eu quero me autenticar via Web3
  Para acessar o sistema IaC AI Agent

  Contexto:
    Dado que o sistema está configurado com Privy.io
    E que a Base Network está conectada
    E que os contratos estão deployados

  @authentication @web3
  Cenário: Login com MetaMask pela primeira vez
    Dado que sou um novo usuário
    E que tenho MetaMask instalado
    Quando eu acesso a página inicial
    E clico em "Connect Wallet"
    E seleciono "MetaMask"
    E autorizo a conexão no MetaMask
    E confirmo a assinatura da mensagem
    Então devo estar autenticado com sucesso
    E meu endereço de wallet deve estar visível
    E devo ver "Welcome! Connect your wallet to get started"
    E uma sessão deve ser criada automaticamente

  @authentication @web3 @embedded_wallet
  Cenário: Login com email criando embedded wallet
    Dado que sou um novo usuário sem wallet
    Quando eu acesso a página inicial
    E clico em "Get Started"
    E seleciono "Continue with Email"
    E insiro meu email "user@example.com"
    E confirmo o código de verificação
    Então uma embedded wallet deve ser criada
    E devo estar autenticado
    E devo ver "Embedded wallet created successfully"
    E devo poder fazer login futuramente com email

  @authentication @web3 @coinbase
  Cenário: Login com Coinbase Wallet
    Dado que tenho Coinbase Wallet instalado
    Quando eu acesso a página inicial
    E clico em "Connect Wallet"
    E seleciono "Coinbase Wallet"
    E autorizo a conexão
    Então devo estar autenticado
    E meu endereço Coinbase deve estar visível
    E devo ter acesso às funcionalidades básicas

  @authentication @web3 @error_handling
  Cenário: Falha na autenticação por wallet não conectada
    Dado que não tenho nenhuma wallet conectada
    Quando eu tento acessar uma página protegida
    Então devo ser redirecionado para login
    E devo ver mensagem "Please connect your wallet to continue"
    E devo ver opções de conexão disponíveis

  @authentication @web3 @session_management
  Cenário: Renovação automática de sessão
    Dado que estou autenticado há 23 horas
    E minha sessão expira em 1 hora
    Quando eu faço uma nova ação no sistema
    Então minha sessão deve ser renovada automaticamente
    E devo continuar autenticado
    E não devo ver mensagens de expiração

  @authentication @web3 @logout
  Cenário: Logout completo do sistema
    Dado que estou autenticado
    Quando eu clico em "Logout"
    E confirmo o logout
    Então minha sessão deve ser invalidada
    E devo ser redirecionado para página inicial
    E não devo ter acesso a páginas protegidas
    E minha wallet deve ser desconectada

  @authentication @web3 @multiple_wallets
  Cenário: Trocar entre múltiplas wallets
    Dado que tenho 2 wallets conectadas
    E estou usando a wallet "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
    Quando eu clico em "Switch Wallet"
    E seleciono "0x8f3Cf7ad23Cd3CaDbD9735AFF958023239c6A063"
    Então devo estar autenticado com a nova wallet
    E meu NFT e tokens devem ser verificados para a nova wallet
    E o contexto deve ser atualizado corretamente
