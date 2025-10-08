# language: pt
Funcionalidade: Cenários de Erro e Edge Cases
  Como desenvolvedor
  Eu quero testar cenários de erro e casos extremos
  Para garantir robustez do sistema

  Contexto:
    Dado que o sistema está configurado
    E que os serviços externos podem falhar

  @error_handling @blockchain
  Cenário: Falha na transação blockchain
    Dado que estou tentando comprar NFT
    E que minha wallet tem ETH suficiente
    Quando eu confirmo a transação
    E a transação falha na blockchain (gas limit exceeded)
    Então devo receber erro específico:
      | Campo           | Valor                         |
      | Error Code      | BLOCKCHAIN_TX_FAILED         |
      | Error Message   | "Transação falhou: Gas limit exceeded" |
      | Transaction Hash| 0x1234567890abcdef...        |
      | Gas Used        | 21000                         |
      | Gas Limit       | 21000                         |
    E meu ETH não deve ser debitado
    E devo poder tentar novamente com gas limit maior

  @error_handling @rate_limit
  Cenário: Rate limit excedido em serviço externo
    Dado que estou fazendo muitas requisições
    E que atingi o rate limit do Nation.fun
    Quando eu tento fazer nova análise
    Então devo receber erro de rate limit:
      | Campo           | Valor                         |
      | Error Code      | RATE_LIMIT_EXCEEDED          |
      | Error Message   | "Rate limit excedido. Tente novamente em 60 segundos" |
      | Retry After     | 60 segundos                   |
      | Current Usage   | 50/50 requests per minute     |
    E a análise deve ser agendada para retry automático
    E devo receber notificação quando estiver disponível

  @error_handling @network
  Cenário: Falha de conectividade de rede
    Dado que estou tentando fazer análise
    E que há problemas de conectividade
    Quando eu submeto código para análise
    E o timeout é atingido (30 segundos)
    Então devo receber erro de timeout:
      | Campo           | Valor                         |
      | Error Code      | NETWORK_TIMEOUT              |
      | Error Message   | "Timeout na conexão com serviço externo" |
      | Timeout Duration| 30 segundos                   |
      | Service         | Nation.fun API               |
    E a análise deve ser marcada como falhada
    E devo poder tentar novamente

  @error_handling @authentication
  Cenário: Token de autenticação expirado
    Dado que estou autenticado há 25 horas
    E que minha sessão expira em 24 horas
    Quando eu tento fazer uma ação protegida
    Então devo receber erro de autenticação:
      | Campo           | Valor                         |
      | Error Code      | AUTH_TOKEN_EXPIRED           |
      | Error Message   | "Sessão expirada. Faça login novamente" |
      | Expired At      | 2024-01-15T10:30:00Z          |
      | Current Time    | 2024-01-16T11:30:00Z          |
    E devo ser redirecionado para login
    E minha sessão deve ser invalidada

  @error_handling @insufficient_funds
  Cenário: Saldo insuficiente de tokens
    Dado que possuo apenas 2 IACAI tokens
    E que uma análise custa 5 tokens
    Quando eu tento fazer análise com LLM
    Então devo receber erro de saldo insuficiente:
      | Campo           | Valor                         |
      | Error Code      | INSUFFICIENT_TOKENS           |
      | Error Message   | "Saldo insuficiente de tokens" |
      | Current Balance | 2 IACAI                       |
      | Required        | 5 IACAI                       |
      | Shortage        | 3 IACAI                       |
    E devo ver botão "Buy Tokens" destacado
    E a análise não deve ser executada

  @error_handling @invalid_input
  Cenário: Código Terraform inválido
    Dado que tenho tokens suficientes
    E que possuo NFT de acesso
    Quando eu submeto código Terraform inválido:
      """
      resource "aws_s3_bucket" "example" {
        bucket = "my-bucket"
        # Sintaxe inválida - falta fechamento
      }
      """
    Então devo receber erro de validação:
      | Campo           | Valor                         |
      | Error Code      | INVALID_TERRAFORM             |
      | Error Message   | "Código Terraform inválido"   |
      | Line Number     | 3                             |
      | Error Details   | "Expected closing brace"      |
      | Suggestions     | ["Verificar sintaxe", "Usar terraform fmt"] |
    E meus tokens não devem ser debitados
    E devo poder corrigir e tentar novamente

  @edge_case @multiple_nfts
  Cenário: Usuário com múltiplos NFTs
    Dado que possuo 2 NFTs diferentes:
      | Token ID | Tier           | Status |
      | 12345    | Basic Access   | Active |
      | 67890    | Pro Access     | Active |
    Quando eu verifico meu acesso
    Então o sistema deve usar o tier mais alto (Pro Access)
    E devo ter acesso a todas as funcionalidades do Pro
    E devo ver notificação "Usando tier mais alto: Pro Access"

  @edge_case @concurrent_analysis
  Cenário: Múltiplas análises simultâneas
    Dado que possuo tokens suficientes
    E que tenho rate limit adequado
    Quando eu submeto 3 análises simultaneamente
    Então todas devem ser processadas:
      | Análise | Status | Tempo | Tokens Gastos |
      | 1       | ✓      | 3.2s  | 5             |
      | 2       | ✓      | 3.5s  | 5             |
      | 3       | ✓      | 3.1s  | 5             |
    E meu saldo final deve refletir todas as deduções
    E todas devem aparecer no histórico

  @edge_case @large_codebase
  Cenário: Análise de código muito grande
    Dado que tenho código Terraform com 10.000 linhas
    E que possuo tokens suficientes
    Quando eu submeto para análise
    Então o sistema deve:
      - Dividir em chunks menores
      - Processar cada chunk
      - Consolidar resultados
      - Retornar análise completa
    E o tempo total deve ser < 60 segundos
    E o custo deve ser calculado proporcionalmente

  @edge_case @special_characters
  Cenário: Código com caracteres especiais
    Dado que tenho código com caracteres especiais:
      """
      resource "aws_s3_bucket" "exemplo-ção" {
        bucket = "meu-bucket-único"
        tags = {
          "Nome" = "Recurso com Acentos"
          "Descrição" = "Teste de caracteres especiais"
        }
      }
      """
    Quando eu submeto para análise
    Então o sistema deve processar corretamente
    E não deve haver problemas de encoding
    E a análise deve ser retornada normalmente

  @edge_case @empty_code
  Cenário: Submissão de código vazio
    Dado que tenho tokens suficientes
    Quando eu submeto código vazio ""
    Então devo receber erro de validação:
      | Campo           | Valor                         |
      | Error Code      | EMPTY_CODE                    |
      | Error Message   | "Código não pode estar vazio"   |
    E meus tokens não devem ser debitados
    E devo ver sugestão "Cole seu código Terraform aqui"

  @edge_case @malformed_json
  Cenário: Requisição com JSON malformado
    Dado que estou fazendo requisição para API
    Quando eu envio JSON malformado:
      ```json
      {
        "code": "resource \"aws_s3_bucket\" \"example\" {
        "type": "llm_analysis"
      ```
    Então devo receber erro de parsing:
      | Campo           | Valor                         |
      | Error Code      | INVALID_JSON                  |
      | Error Message   | "JSON malformado"             |
      | HTTP Status     | 400 Bad Request               |
    E a requisição deve ser rejeitada
    E devo receber exemplo de JSON válido

  @edge_case @service_degradation
  Cenário: Serviço externo com performance degradada
    Dado que o Nation.fun está com latência alta (10s)
    E que tenho timeout configurado para 30s
    Quando eu submeto análise
    Então a análise deve ser processada
    E devo receber warning sobre latência:
      | Campo           | Valor                         |
      | Warning         | "Serviço com latência alta"   |
      | Processing Time | 10.5 segundos                 |
      | Normal Time     | 2-3 segundos                 |
    E o resultado deve ser válido
    E devo ser notificado sobre a degradação

  @edge_case @partial_failure
  Cenário: Falha parcial em análise
    Dado que estou fazendo Full Review
    E que Checkov funciona mas LLM falha
    Quando eu submeto código
    Então devo receber resultado parcial:
      | Componente      | Status | Resultado            |
      | Terraform Parse | ✓      | Válido               |
      | Checkov Scan    | ✓      | 3 issues encontrados |
      | LLM Analysis    | ✗      | Serviço indisponível |
      | Cost Analysis   | ✓      | 2 otimizações        |
    E devo ser notificado sobre componentes falhados
    E o custo deve ser ajustado proporcionalmente
