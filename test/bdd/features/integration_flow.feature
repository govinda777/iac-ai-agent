# language: pt
Funcionalidade: Integração Completa do Fluxo de Usuário
  Como um usuário do sistema
  Eu quero que todas as partes do fluxo funcionem integradas
  Para ter uma experiência fluida e sem interrupções

  Contexto:
    Dado que o sistema está funcionando completamente
    E que todos os serviços estão integrados
    E que a interface web está operacional

  @integration @end_to_end
  Cenário: Fluxo completo com validação de integração
    # Setup inicial
    Dado que o sistema está em modo de produção
    E que os contratos estão deployados na Base Network
    E que o Privy está configurado corretamente
    E que a API backend está funcionando
    
    # Fluxo completo do usuário
    Quando um novo usuário acessa o sistema
    E completa todo o fluxo de onboarding:
      - Login com MetaMask
      - Compra de NFT Pro Access
      - Compra de 500 tokens IACAI
      - Submissão de código Terraform
      - Recebimento de análise completa
    
    Então todas as integrações devem funcionar:
      | Integração              | Status Esperado |
      | Privy Authentication    | ✓ Sucesso       |
      | Base Network Contract   | ✓ Transação confirmada |
      | NFT Minting            | ✓ NFT recebido   |
      | Token Purchase         | ✓ Tokens creditados |
      | Backend API            | ✓ Análise processada |
      | Frontend Updates       | ✓ UI atualizada  |
    
    E o usuário deve ter:
      - Wallet conectada e verificada
      - NFT Pro Access em sua carteira
      - Saldo de tokens atualizado
      - Relatório de análise completo
      - Histórico de transações visível

  @integration @error_recovery
  Cenário: Recuperação de erros durante o fluxo
    Dado que estou no meio do fluxo de compra de NFT
    E a transação falha por rede instável
    
    Quando o sistema detecta a falha
    Então devo ver:
      - Mensagem de erro clara
      - Opção "Tentar novamente"
      - Opção "Comprar com cartão"
      - Estado anterior preservado
    
    Quando escolho "Tentar novamente"
    E a transação é bem-sucedida
    Então o fluxo deve continuar normalmente
    E todos os estados devem ser atualizados
    
    Quando escolho "Comprar com cartão"
    Então devo ser redirecionado para Privy Onramp
    E após pagamento o NFT deve ser mintado
    E o fluxo deve continuar

  @integration @data_consistency
  Cenário: Consistência de dados entre sistemas
    Dado que estou autenticado
    E possuo NFT "Pro Access"
    E tenho "100 tokens" de saldo
    
    Quando faço uma análise que custa "5 tokens"
    Então os dados devem ser consistentes em:
      | Sistema                | Valor Esperado |
      | Frontend UI           | 95 tokens      |
      | Backend Database      | 95 tokens      |
      | Blockchain Contract   | 95 tokens      |
      | User Session          | 95 tokens      |
    
    Quando recarregar a página
    Então todos os valores devem permanecer consistentes
    E não deve haver discrepâncias

  @integration @performance
  Cenário: Performance do fluxo integrado
    Dado que o sistema está sob carga normal
    
    Quando executo o fluxo completo
    Então os tempos de resposta devem ser:
      | Operação               | Tempo Máximo |
      | Login com Privy       | 3 segundos   |
      | Compra de NFT         | 30 segundos  |
      | Compra de tokens      | 30 segundos  |
      | Submissão de análise  | 2 segundos   |
      | Processamento análise | 30 segundos  |
      | Atualização UI        | 1 segundo    |
    
    E não deve haver:
      - Timeouts de transação
      - Perda de dados
      - Estados inconsistentes
      - Erros de sincronização

  @integration @security
  Cenário: Segurança durante o fluxo integrado
    Dado que estou executando o fluxo completo
    
    Então todas as comunicações devem ser:
      - Criptografadas (HTTPS/TLS)
      - Autenticadas com tokens válidos
      - Validadas no backend
      - Logadas para auditoria
    
    E os dados sensíveis devem ser:
      - Protegidos em trânsito
      - Não expostos no frontend
      - Criptografados no banco
      - Limpos após logout
    
    Quando há tentativa de acesso não autorizado
    Então o sistema deve:
      - Bloquear a requisição
      - Registrar o evento
      - Notificar o usuário
      - Manter a segurança

  @integration @scalability
  Cenário: Escalabilidade do fluxo integrado
    Dado que múltiplos usuários executam o fluxo simultaneamente
    
    Quando há 10 usuários simultâneos
    Então o sistema deve manter:
      - Tempos de resposta aceitáveis
      - Consistência de dados
      - Disponibilidade de serviços
      - Qualidade da experiência
    
    Quando há 100 usuários simultâneos
    Então o sistema deve:
      - Escalar automaticamente
      - Distribuir carga adequadamente
      - Manter performance
      - Não degradar funcionalidades

  @integration @monitoring
  Cenário: Monitoramento do fluxo integrado
    Dado que o sistema está sendo monitorado
    
    Quando executo o fluxo completo
    Então devem ser coletadas métricas de:
      - Tempo de cada etapa
      - Taxa de sucesso/erro
      - Uso de recursos
      - Satisfação do usuário
    
    E devem ser gerados alertas para:
      - Falhas de transação
      - Tempos de resposta altos
      - Erros de integração
      - Problemas de segurança
    
    Quando há problemas
    Então o sistema deve:
      - Detectar automaticamente
      - Notificar administradores
      - Tentar recuperação
      - Preservar dados do usuário

  @integration @backup_recovery
  Cenário: Backup e recuperação durante o fluxo
    Dado que estou no meio de uma transação importante
    
    Quando há falha de sistema
    Então os dados devem ser:
      - Salvos em backup automático
      - Recuperáveis após falha
      - Consistentes após restauração
      - Não perdidos pelo usuário
    
    Quando o sistema é restaurado
    Então o usuário deve poder:
      - Continuar de onde parou
      - Ver estado correto
      - Completar transações pendentes
      - Não perder progresso

  @integration @cross_browser
  Cenário: Compatibilidade entre navegadores
    Dado que o fluxo deve funcionar em diferentes navegadores
    
    Quando executo o fluxo em Chrome
    Então todas as funcionalidades devem funcionar
    
    Quando executo o fluxo em Firefox
    Então todas as funcionalidades devem funcionar
    
    Quando executo o fluxo em Safari
    Então todas as funcionalidades devem funcionar
    
    Quando executo o fluxo em Edge
    Então todas as funcionalidades devem funcionar
    
    E a experiência deve ser consistente entre todos

  @integration @mobile_integration
  Cenário: Integração em dispositivos móveis
    Dado que acesso o sistema em dispositivo móvel
    
    Quando executo o fluxo completo
    Então todas as integrações devem funcionar:
      - Login com wallet móvel
      - Transações via app móvel
      - Interface responsiva
      - Performance adequada
    
    E a experiência deve ser:
      - Intuitiva em touch
      - Rápida e fluida
      - Compatível com apps de wallet
      - Otimizada para mobile
