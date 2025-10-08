# language: pt
Funcionalidade: Onboarding de Usuário via Privy
  Como um novo usuário
  Eu quero me autenticar usando minha wallet
  Para acessar o IaC AI Agent Bot

  Contexto:
    Dado que o serviço Privy está disponível
    E que a Base Network está acessível
    E que os contratos de NFT e Token estão deployados
    E que o frontend está corretamente configurado com Privy SDK

  Cenário: Novo usuário faz login com Metamask
    Dado que sou um novo usuário sem conta
    Quando eu clico em "Connect Wallet"
    E seleciono "MetaMask" como provider
    E aprovo a conexão no MetaMask
    Então devo estar autenticado
    E meu endereço da wallet deve estar visível
    E devo ver a mensagem "Bem-vindo ao IaC AI Agent"
    E devo ver as opções de compra de acesso

  Cenário: Novo usuário faz login com Coinbase Wallet
    Dado que sou um novo usuário sem conta
    Quando eu clico em "Connect Wallet"
    E seleciono "Coinbase Wallet" como provider
    E aprovo a conexão no Coinbase Wallet
    Então devo estar autenticado
    E meu endereço da wallet deve estar visível

  Cenário: Usuário cria embedded wallet via Privy
    Dado que sou um novo usuário sem wallet
    Quando eu clico em "Create Wallet"
    E concluo o processo de autenticação com email
    Então uma embedded wallet deve ser criada automaticamente
    E devo estar autenticado
    E devo ver meu endereço de wallet

  Cenário: Usuário vincula email à conta Privy
    Dado que estou autenticado com wallet
    Quando eu clico em "Link Email"
    E insiro meu email "user@example.com"
    E confirmo o código de verificação recebido
    Então meu email deve estar vinculado à conta
    E devo poder fazer login com email ou wallet

  Cenário: Tentativa de acesso sem autenticação
    Dado que não estou autenticado
    Quando eu tento acessar "/dashboard"
    Então devo ser redirecionado para "/login"
    E devo ver a mensagem "Por favor, conecte sua wallet"

  Cenário: Sessão expirada
    Dado que estou autenticado há mais de 24 horas
    Quando eu tento fazer uma análise
    Então devo receber erro "Token expirado"
    E devo ser redirecionado para login
    E após re-autenticar devo voltar para a página original
