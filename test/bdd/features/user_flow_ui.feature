# language: pt
Funcionalidade: Fluxo Completo de UI do Usuário
  Como um usuário do IaC AI Agent
  Eu quero seguir um fluxo claro e intuitivo na interface
  Para completar todas as etapas necessárias para analisar meu código Terraform

  Contexto:
    Dado que o sistema está funcionando corretamente
    E que a interface web está carregada
    E que os serviços Web3 estão disponíveis

  @ui_flow @critical_path
  Cenário: Usuário completa o fluxo completo através da UI
    # Etapa 1: Verificação do estado inicial
    Dado que acesso a página inicial do IaC AI Agent
    Então devo ver a seção "Seu Caminho para Análises Inteligentes"
    E devo ver 5 passos numerados:
      | Passo | Título                    | Status Esperado |
      | 1     | Conecte sua Wallet       | Não conectado   |
      | 2     | Adquira um NFT de Acesso | NFT necessário  |
      | 3     | Compre Tokens IACAI      | Tokens necessários |
      | 4     | Submeta seu Código       | Pronto para análise |
      | 5     | Receba Sugestões         | Aguardando análise |
    
    # Verificação do status atual
    E no card "Status Atual" devo ver:
      | Item   | Status Esperado |
      | Wallet | Não conectada  |
      | NFT    | Não adquirido  |
      | Tokens | 0 IACAI        |
    
    E o "Próximo Passo" deve ser "Conecte sua wallet para começar"
    E o botão de ação deve mostrar "Conectar Wallet"

    # Etapa 2: Login com wallet
    Quando eu clico no botão "Conectar Wallet" no passo 1
    E seleciono "MetaMask" como provedor
    E autorizo a conexão no MetaMask
    Então devo estar autenticado com sucesso
    E o status do passo 1 deve mudar para "Conectado"
    E meu endereço de wallet deve aparecer no cabeçalho
    E o "Próximo Passo" deve ser "Adquira um NFT de acesso"
    E o botão de ação deve mostrar "Ver Planos"

    # Etapa 3: Compra de NFT
    Quando eu clico no botão "Ver Planos" no passo 2
    E sou redirecionado para a seção de preços
    E seleciono o tier "Pro Access" (0.05 ETH)
    E clico em "Comprar NFT"
    E confirmo a transação na wallet
    Então a transação deve ser processada
    E após confirmação devo receber o NFT Pro Access
    E o status do passo 2 deve mudar para "NFT Pro adquirido"
    E o "Próximo Passo" deve ser "Compre tokens IACAI"
    E o botão de ação deve mostrar "Comprar Tokens"

    # Etapa 4: Compra de tokens
    Quando eu clico no botão "Comprar Tokens" no passo 3
    E seleciono o pacote "Power Pack" (500 tokens)
    E clico em "Comprar com ETH"
    E confirmo a transação na wallet
    Então os tokens devem ser adicionados à minha conta
    E meu saldo deve mostrar "500 IACAI"
    E o status do passo 3 deve mudar para "500 tokens disponíveis"
    E o "Próximo Passo" deve ser "Submeta seu código"
    E o botão de ação deve mostrar "Analisar Código"

    # Etapa 5: Submissão de código
    Quando eu clico no botão "Analisar Código" no passo 4
    E sou redirecionado para a seção de análise
    E cole o seguinte código Terraform:
      """
      resource "aws_s3_bucket" "example" {
        bucket = "my-terraform-bucket"
        acl    = "public-read"
      }
      
      resource "aws_instance" "web" {
        ami           = "ami-0c55b159cbfafe1f0"
        instance_type = "t3.micro"
      }
      """
    E seleciono "Análise com LLM" (5 tokens)
    E vejo que o custo é "5 tokens" e meu saldo é "500 IACAI"
    E clico em "Analisar Código"
    Então o sistema deve processar minha análise
    E deve debitar 5 IACAI tokens da minha conta
    E meu saldo deve ser atualizado para "495 IACAI"
    E o status do passo 5 deve mudar para "Análise concluída"

    # Etapa 6: Recebimento de resultados
    E em até 30 segundos devo receber um relatório de análise
    E o relatório deve conter as seguintes seções:
      | Seção                    | Status |
      | Executive Summary        | ✓      |
      | Security Issues          | ✓      |
      | Best Practices           | ✓      |
      | LLM Analysis             | ✓      |
      | Detailed Findings        | ✓      |
    E pelo menos uma sugestão sobre o "public-read" ACL
    E cada sugestão deve conter código correto recomendado
    E o "Próximo Passo" deve ser "Análise concluída com sucesso"

  @ui_flow @validation
  Cenário: Validação de estados intermediários da UI
    Dado que estou na página inicial
    E não estou autenticado
    
    # Teste de validação de acesso sem login
    Quando eu tento clicar em "Comprar NFT" sem estar logado
    Então devo ver a mensagem "Conecte sua wallet primeiro"
    E devo ser redirecionado para o login
    
    Quando eu tento clicar em "Comprar Tokens" sem estar logado
    Então devo ver a mensagem "Conecte sua wallet primeiro"
    
    Quando eu tento submeter código sem estar logado
    Então o formulário deve estar desabilitado
    E devo ver a mensagem "Conecte sua wallet para analisar código"

  @ui_flow @validation
  Cenário: Validação de acesso por tier de NFT
    Dado que estou autenticado com wallet "0x1234...5678"
    E possuo NFT "Basic Access"
    
    Quando eu acesso a seção de análise
    E seleciono "Análise com LLM" (requer Pro Access)
    Então devo ver a mensagem "Seu tier (basic) não permite este tipo de análise"
    E a opção deve estar desabilitada
    E devo ver "Requer pro ou superior"
    
    Quando eu seleciono "Análise Básica" (permitida para Basic)
    Então a opção deve estar habilitada
    E o custo deve mostrar "1 token"

  @ui_flow @validation
  Cenário: Validação de saldo insuficiente de tokens
    Dado que estou autenticado com wallet "0x1234...5678"
    E possuo NFT "Pro Access"
    E meu saldo de tokens é "2 IACAI"
    
    Quando eu acesso a seção de análise
    E seleciono "Full Review" (15 tokens)
    E clico em "Analisar Código"
    Então devo ver a mensagem "Saldo insuficiente de tokens. Necessário: 15, Disponível: 2"
    E o botão deve estar desabilitado
    E devo ver um link para "Comprar mais tokens"

  @ui_flow @error_handling
  Cenário: Tratamento de erros durante transações
    Dado que estou autenticado com wallet "0x1234...5678"
    E não possuo NFT de acesso
    
    Quando eu tento comprar um NFT
    E a transação falha por saldo insuficiente
    Então devo ver a mensagem "Transação falhou: Saldo insuficiente"
    E devo ver opções para "Tentar novamente" ou "Comprar com cartão"
    
    Quando eu tento comprar tokens sem NFT
    Então devo ver a mensagem "É necessário ter um NFT de acesso para comprar tokens"
    E devo ser redirecionado para a seção de NFTs

  @ui_flow @progressive_disclosure
  Cenário: Revelação progressiva de funcionalidades
    Dado que acesso a página inicial
    
    # Estado inicial - apenas login visível
    Então devo ver apenas o passo 1 habilitado
    E os passos 2-5 devem estar visualmente desabilitados
    E devo ver tooltips explicativos em cada passo
    
    # Após login - NFT e tokens habilitados
    Quando eu faço login com sucesso
    Então os passos 2 e 3 devem ficar habilitados
    E devo ver indicadores visuais de progresso
    
    # Após NFT - análise habilitada
    Quando eu adquiro um NFT
    Então o passo 4 deve ficar habilitado
    E devo ver o botão "Analisar Código" ativo
    
    # Após tokens - análise completa habilitada
    Quando eu adquiro tokens suficientes
    Então o passo 5 deve ficar habilitado
    E todos os tipos de análise devem estar disponíveis

  @ui_flow @responsive
  Cenário: Fluxo em dispositivos móveis
    Dado que acesso o site em um dispositivo móvel
    E a tela tem largura menor que 768px
    
    Então o fluxo de 5 passos deve ser exibido verticalmente
    E cada passo deve ocupar a largura total da tela
    E os botões devem ter tamanho adequado para toque
    E o texto deve ser legível sem zoom
    
    Quando eu rolo a página
    Então os passos devem ter scroll suave
    E o indicador de progresso deve ser visível

  @ui_flow @accessibility
  Cenário: Acessibilidade do fluxo
    Dado que acesso a página inicial
    
    Então todos os elementos devem ter:
      - Contraste adequado (WCAG AA)
      - Texto alternativo em imagens
      - Labels associados a campos
      - Navegação por teclado funcional
    
    Quando eu navego usando apenas o teclado
    Então devo conseguir:
      - Tab entre todos os elementos interativos
      - Ativar botões com Enter ou Space
      - Navegar pelos passos com setas
      - Acessar todas as funcionalidades

  @ui_flow @performance
  Cenário: Performance do fluxo
    Dado que acesso a página inicial
    
    Então a página deve carregar em menos de 3 segundos
    E as transições entre passos devem ser suaves (< 300ms)
    E as validações devem ser instantâneas
    
    Quando eu completo cada etapa
    Então as atualizações de UI devem ser imediatas
    E não deve haver travamentos ou delays visíveis
