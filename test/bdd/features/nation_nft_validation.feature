# language: pt
Funcionalidade: Validação de NFT Pass do Nation em Tempo de Execução
  Como um sistema IaC AI Agent
  Eu quero validar se a carteira padrão possui NFT Pass do Nation
  Para garantir acesso autorizado e funcionalidade completa

  Contexto:
    Dado que o sistema está configurado com a carteira padrão "0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5"
    E que NATION_NFT_REQUIRED está definido como true
    E que a API do Nation.fun está disponível
    E que o contrato NationPassNFT está deployado

  @nation_nft @startup_validation @critical
  Cenário: Validação de NFT Pass na inicialização da aplicação
    Dado que a aplicação está sendo inicializada
    E que a carteira padrão está configurada
    Quando o sistema executa a validação de NFT Pass do Nation
    Então a carteira deve ser verificada contra a carteira padrão autorizada
    E a API do Nation.fun deve ser consultada para verificar NFT
    E o sistema deve confirmar que a carteira possui NFT Pass válido
    E um teste de conectividade deve ser enviado para o Nation.fun
    E a resposta do teste deve ser coletada com sucesso
    E a aplicação deve inicializar normalmente
    E deve ser logado "NFT Pass do Nation validado com sucesso"

  @nation_nft @wallet_validation @security
  Cenário: Validação de carteira não autorizada
    Dado que uma carteira diferente da padrão está sendo verificada
    E que o endereço é "0x1234567890123456789012345678901234567890"
    Quando o sistema tenta validar NFT Pass do Nation
    Então deve ser retornado erro "wallet não autorizada para acesso"
    E nenhuma consulta à API do Nation.fun deve ser feita
    E a validação deve falhar imediatamente
    E deve ser logado "wallet não autorizada"

  @nation_nft @api_integration @success
  Cenário: Validação bem-sucedida com NFT Pass válido
    Dado que a carteira padrão está sendo verificada
    E que a API do Nation.fun retorna NFT válido
    E que o NFT possui tier "PRO"
    E que o NFT está ativo e não expirado
    Quando o sistema valida NFT Pass do Nation
    Então deve ser retornado sucesso
    E o token ID deve ser capturado
    E o tier deve ser identificado como "PRO"
    E o status deve ser "ativo"
    E deve ser logado "NFT Pass do Nation validado com sucesso"

  @nation_nft @api_integration @failure
  Cenário: Validação falha por NFT inexistente
    Dado que a carteira padrão está sendo verificada
    E que a API do Nation.fun retorna que não possui NFT
    Quando o sistema valida NFT Pass do Nation
    Então deve ser retornado erro "carteira não possui NFT Pass do Nation válido"
    E a aplicação não deve inicializar
    E deve ser logado "Falha na validação de NFT na inicialização"

  @nation_nft @test_connectivity @integration
  Cenário: Envio de teste para Nation.fun com resposta bem-sucedida
    Dado que a validação de NFT foi bem-sucedida
    E que a conectividade com Nation.fun está funcionando
    Quando o sistema envia teste "Teste de conectividade na inicialização"
    Então a API do Nation.fun deve receber o teste
    E deve retornar test_id válido
    E deve retornar status "success"
    E deve retornar timestamp atual
    E deve ser logado "Teste de conectividade com Nation.fun bem-sucedido"

  @nation_nft @test_connectivity @failure
  Cenário: Falha no teste de conectividade com Nation.fun
    Dado que a validação de NFT foi bem-sucedida
    E que a API do Nation.fun está indisponível
    Quando o sistema tenta enviar teste de conectividade
    Então deve ser retornado erro de conectividade
    E deve ser logado "Falha no teste de conectividade"
    E a aplicação deve continuar inicializando normalmente
    E não deve falhar por causa do teste

  @nation_nft @runtime_validation @security
  Cenário: Validação em tempo de execução durante operação
    Dado que a aplicação está rodando normalmente
    E que um usuário tenta executar operação protegida
    Quando o sistema verifica permissões em tempo de execução
    Então a carteira padrão deve ser revalidada
    E o NFT Pass deve ser verificado novamente
    E se válido, a operação deve prosseguir
    E se inválido, a operação deve ser negada
    E deve ser logado o resultado da validação

  @nation_nft @configuration @environment
  Cenário: Configuração incorreta de variáveis de ambiente
    Dado que WALLET_ADDRESS não está configurado
    E que NATION_NFT_REQUIRED está definido como true
    Quando o sistema tenta inicializar
    Então deve ser retornado erro "WALLET_ADDRESS não configurado"
    E a aplicação não deve inicializar
    E deve ser logado erro de configuração

  @nation_nft @configuration @optional
  Cenário: Validação opcional quando NATION_NFT_REQUIRED é false
    Dado que NATION_NFT_REQUIRED está definido como false
    E que WALLET_ADDRESS está configurado
    Quando o sistema inicializa
    Então a validação de NFT deve ser pulada
    E a aplicação deve inicializar normalmente
    E deve ser logado "Validação de NFT opcional - pulada"

  @nation_nft @api_error_handling @resilience
  Cenário: Tratamento de erro da API do Nation.fun
    Dado que a carteira padrão está sendo verificada
    E que a API do Nation.fun retorna erro HTTP 500
    Quando o sistema tenta validar NFT Pass
    Então deve ser retornado erro "API retornou status 500"
    E a aplicação não deve inicializar
    E deve ser logado erro detalhado da API

  @nation_nft @timeout @resilience
  Cenário: Timeout na comunicação com API do Nation.fun
    Dado que a carteira padrão está sendo verificada
    E que a API do Nation.fun demora mais de 30 segundos para responder
    Quando o sistema tenta validar NFT Pass
    Então deve ser retornado erro de timeout
    E a aplicação não deve inicializar
    E deve ser logado "timeout na comunicação com Nation.fun"

  @nation_nft @json_parsing @error_handling
  Cenário: Resposta JSON inválida da API do Nation.fun
    Dado que a carteira padrão está sendo verificada
    E que a API do Nation.fun retorna JSON malformado
    Quando o sistema tenta validar NFT Pass
    Então deve ser retornado erro "erro ao parsear resposta JSON"
    E a aplicação não deve inicializar
    E deve ser logado erro de parsing

  @nation_nft @nft_expired @business_logic
  Cenário: NFT Pass expirado
    Dado que a carteira padrão está sendo verificada
    E que a API do Nation.fun retorna NFT expirado
    E que expires_at está no passado
    Quando o sistema valida NFT Pass do Nation
    Então deve ser retornado erro "NFT Pass expirado"
    E a aplicação não deve inicializar
    E deve ser logado "NFT Pass do Nation expirado"

  @nation_nft @nft_inactive @business_logic
  Cenário: NFT Pass inativo
    Dado que a carteira padrão está sendo verificada
    E que a API do Nation.fun retorna NFT inativo
    E que is_active é false
    Quando o sistema valida NFT Pass do Nation
    Então deve ser retornado erro "NFT Pass inativo"
    E a aplicação não deve inicializar
    E deve ser logado "NFT Pass do Nation inativo"
