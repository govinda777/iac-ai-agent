# language: pt
Funcionalidade: Validação de Estado da Interface do Usuário
  Como um desenvolvedor de testes
  Eu quero validar que a UI reflete corretamente o estado do usuário
  Para garantir uma experiência consistente e confiável

  Contexto:
    Dado que o sistema está funcionando
    E que a interface web está carregada
    E que os serviços de autenticação estão disponíveis

  @ui_state @authentication
  Cenário: Estados de autenticação refletidos na UI
    Dado que não estou autenticado
    Quando acesso qualquer página do sistema
    Então devo ver no cabeçalho:
      - Botão "Conectar Wallet" visível
      - Seção de perfil do usuário oculta
      - Endereço de wallet não visível
    
    Quando faço login com sucesso
    Então devo ver no cabeçalho:
      - Botão "Conectar Wallet" oculto
      - Seção de perfil do usuário visível
      - Endereço de wallet truncado visível
      - Botão de logout visível
    
    Quando faço logout
    Então devo voltar ao estado inicial
    E todos os dados sensíveis devem ser limpos

  @ui_state @nft_status
  Cenário: Estados de NFT refletidos na UI
    Dado que estou autenticado mas não possuo NFT
    Quando acesso a página inicial
    Então devo ver:
      - Status "NFT necessário" no passo 2
      - Botão "Ver Planos" habilitado
      - Seção de análise desabilitada
    
    Quando adquiro NFT "Basic Access"
    Então devo ver:
      - Status "NFT Basic adquirido" no passo 2
      - Funcionalidades básicas habilitadas
      - Funcionalidades avançadas desabilitadas
    
    Quando adquiro NFT "Pro Access"
    Então devo ver:
      - Status "NFT Pro adquirido" no passo 2
      - Todas as funcionalidades habilitadas
      - Acesso a análises com LLM

  @ui_state @token_balance
  Cenário: Saldo de tokens refletido na UI
    Dado que estou autenticado com NFT
    E meu saldo de tokens é "0 IACAI"
    
    Quando acesso qualquer página
    Então devo ver:
      - "0 IACAI" no cabeçalho
      - "0 IACAI" na seção de status
      - Botões de análise desabilitados
    
    Quando adquiro "100 tokens"
    Então devo ver:
      - "100 IACAI" no cabeçalho
      - "100 IACAI" na seção de status
      - Botões de análise básica habilitados
    
    Quando uso "5 tokens" em uma análise
    Então devo ver:
      - "95 IACAI" no cabeçalho
      - "95 IACAI" na seção de status
      - Saldo atualizado imediatamente

  @ui_state @progressive_enablement
  Cenário: Habilitação progressiva de funcionalidades
    Dado que não estou autenticado
    Quando acesso a página inicial
    Então devo ver:
      - Passo 1: habilitado (login)
      - Passos 2-5: visualmente desabilitados
      - Botões de compra desabilitados
      - Formulário de análise desabilitado
    
    Quando faço login
    Então devo ver:
      - Passo 1: concluído
      - Passo 2: habilitado (NFT)
      - Passos 3-5: visualmente desabilitados
    
    Quando adquiro NFT
    Então devo ver:
      - Passo 2: concluído
      - Passo 3: habilitado (tokens)
      - Passos 4-5: visualmente desabilitados
    
    Quando adquiro tokens
    Então devo ver:
      - Passo 3: concluído
      - Passos 4-5: habilitados
      - Formulário de análise habilitado

  @ui_state @error_states
  Cenário: Estados de erro refletidos na UI
    Dado que estou autenticado
    E possuo NFT "Basic Access"
    E meu saldo é "0 tokens"
    
    Quando tento fazer análise "Full Review" (15 tokens)
    Então devo ver:
      - Mensagem de erro "Saldo insuficiente"
      - Botão de análise desabilitado
      - Link para comprar tokens destacado
    
    Quando tento fazer análise "LLM" (requer Pro Access)
    Então devo ver:
      - Mensagem de erro "Tier insuficiente"
      - Opção desabilitada no dropdown
      - Texto explicativo sobre requisitos

  @ui_state @loading_states
  Cenário: Estados de carregamento na UI
    Dado que estou autenticado com NFT e tokens
    Quando submeto código para análise
    Então devo ver:
      - Botão "Analisar Código" com spinner
      - Texto "Processando análise..."
      - Formulário desabilitado
      - Indicador de progresso
    
    Quando a análise é concluída
    Então devo ver:
      - Botão volta ao estado normal
      - Resultados exibidos
      - Formulário habilitado novamente
      - Saldo de tokens atualizado

  @ui_state @responsive_states
  Cenário: Estados responsivos da UI
    Dado que acesso o site em desktop (1200px+)
    Então devo ver:
      - Fluxo de 5 passos em layout horizontal
      - Cards de status lado a lado
      - Navegação completa visível
    
    Quando redimensiono para tablet (768px-1199px)
    Então devo ver:
      - Fluxo adaptado para 2 colunas
      - Cards empilhados verticalmente
      - Navegação condensada
    
    Quando redimensiono para mobile (< 768px)
    Então devo ver:
      - Fluxo em coluna única
      - Cards em largura total
      - Menu hambúrguer ativo

  @ui_state @persistence
  Cenário: Persistência de estado na UI
    Dado que estou autenticado com NFT e tokens
    E tenho "100 tokens" de saldo
    Quando recarregar a página
    Então devo ver:
      - Estado de autenticação mantido
      - NFT status mantido
      - Saldo de tokens correto
      - Progresso do fluxo mantido
    
    Quando fecho e reabro o navegador
    Então devo ver:
      - Sessão mantida (se configurado)
      - Estado atualizado corretamente
      - Dados sincronizados com blockchain

  @ui_state @real_time_updates
  Cenário: Atualizações em tempo real da UI
    Dado que estou autenticado
    E tenho a página aberta
    Quando recebo tokens de outra fonte
    Então o saldo deve ser atualizado automaticamente
    E devo ver notificação "Novos tokens recebidos"
    
    Quando minha transação de NFT é confirmada
    Então o status deve mudar automaticamente
    E devo ver notificação "NFT adquirido com sucesso"
    
    Quando minha análise é concluída
    Então os resultados devem aparecer automaticamente
    E devo ver notificação "Análise concluída"

  @ui_state @validation_feedback
  Cenário: Feedback de validação na UI
    Dado que estou na seção de análise
    Quando deixo o campo de código vazio
    E clico em "Analisar Código"
    Então devo ver:
      - Campo destacado em vermelho
      - Mensagem "Por favor, insira código Terraform"
      - Botão permanece desabilitado
    
    Quando insiro código inválido
    E clico em "Analisar Código"
    Então devo ver:
      - Mensagem "Código Terraform inválido"
      - Sugestões de correção
      - Botão habilitado para tentar novamente
    
    Quando insiro código válido
    Então devo ver:
      - Campo destacado em verde
      - Contador de linhas atualizado
      - Botão habilitado
