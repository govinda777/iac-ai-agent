# language: pt
Funcionalidade: Análise de Infraestrutura como Código
  Como um usuário com acesso ao bot
  Eu quero analisar meu código Terraform
  Para receber sugestões de melhorias

  Contexto:
    Dado que estou autenticado com a wallet "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
    E possuo NFT "Pro Access" (Token ID: 12345)
    E possuo 100 IACAI tokens
    E minha sessão está ativa

  Cenário: Análise básica de código Terraform
    Dado que tenho o seguinte código Terraform:
      """
      resource "aws_s3_bucket" "example" {
        bucket = "my-bucket"
      }
      """
    Quando eu submeto o código para análise básica
    Então o sistema deve:
      | Ação                          | Status |
      | Verificar acesso NFT          | ✓      |
      | Parsear código Terraform      | ✓      |
      | Executar análise estática     | ✓      |
      | Verificar saldo de tokens     | ✓      |
      | Deduzir 1 token               | ✓      |
      | Retornar resultado            | ✓      |
    E devo receber um relatório com:
      | Seção                | Conteúdo                         |
      | Score                | 75/100                           |
      | Recursos encontrados | 1 (aws_s3_bucket)                |
      | Sugestões            | 3 recomendações                  |
      | Tempo de análise     | < 5 segundos                     |
    E meu saldo final deve ser 99 IACAI

  Cenário: Análise com LLM (análise inteligente)
    Dado que tenho código Terraform válido
    Quando eu submeto para análise com LLM
    Então o sistema deve:
      | Ação                            | Status |
      | Verificar tier Pro ou superior  | ✓      |
      | Verificar saldo (5 tokens)      | ✓      |
      | Executar análise básica         | ✓      |
      | Consultar Knowledge Base        | ✓      |
      | Enviar para LLM (GPT-4/Claude)  | ✓      |
      | Deduzir 5 tokens                | ✓      |
      | Retornar análise enriquecida    | ✓      |
    E devo receber resposta estruturada com:
      | Seção                    | Presente |
      | Executive Summary        | ✓        |
      | Critical Issues          | ✓        |
      | Improvements             | ✓        |
      | Best Practices           | ✓        |
      | Architectural Insights   | ✓        |
      | Priority Actions         | ✓        |
      | Quick Wins               | ✓        |
    E cada sugestão deve ter código de exemplo
    E meu saldo final deve ser 95 IACAI

  Cenário: Análise de Checkov (segurança)
    Dado que tenho código Terraform com problemas de segurança
    Quando eu submeto para análise de segurança
    Então o sistema deve:
      | Ação                        | Status |
      | Executar Checkov            | ✓      |
      | Classificar por severidade  | ✓      |
      | Deduzir 2 tokens            | ✓      |
    E devo receber relatório de segurança com:
      | Severidade | Quantidade |
      | Critical   | 2          |
      | High       | 5          |
      | Medium     | 3          |
      | Low        | 1          |
    E cada issue deve ter:
      - Check ID
      - Descrição do problema
      - Impacto de negócio
      - Como corrigir (código)
      - Referências externas

  Cenário: Análise completa (Full Review)
    Dado que tenho um projeto Terraform completo
    Quando eu submeto para Full Review
    Então o sistema deve executar:
      | Análise               | Custo (tokens) | Status |
      | Terraform parsing     | incluído       | ✓      |
      | Checkov security      | incluído       | ✓      |
      | IAM analysis          | incluído       | ✓      |
      | LLM enrichment        | incluído       | ✓      |
      | Cost optimization     | incluído       | ✓      |
      | Architecture review   | incluído       | ✓      |
      | **Total**             | **15 tokens**  | ✓      |
    E devo receber relatório completo
    E meu saldo final deve ser 85 IACAI
    E o tempo total deve ser < 30 segundos

  Cenário: Análise bloqueada por falta de tokens
    Dado que possuo apenas 2 IACAI tokens
    Quando eu tento fazer análise com LLM (custo: 5 tokens)
    Então devo ver mensagem:
      """
      ⚠️ Saldo insuficiente de tokens
      
      Você tem: 2 IACAI
      Necessário: 5 IACAI
      Faltam: 3 IACAI
      
      [Buy Tokens]
      """
    E a análise não deve ser executada
    E meu saldo deve permanecer 2 IACAI

  Cenário: Análise bloqueada por tier insuficiente
    Dado que possuo NFT "Basic Access"
    E análise com LLM requer tier "Pro" ou superior
    Quando eu tento fazer análise com LLM
    Então devo ver mensagem:
      """
      ⚠️ Tier insuficiente
      
      Seu tier: Basic Access
      Requerido: Pro Access ou superior
      
      [Upgrade NFT]
      """
    E devo ver botão para fazer upgrade
    E a análise não deve ser executada

  Cenário: Rate limiting por tier
    Dado que possuo NFT "Basic Access"
    E o limite é 10 análises por hora
    Quando eu faço 10 análises consecutivas
    Então todas devem ser processadas
    Mas quando eu tento fazer a 11ª análise
    Então devo ver mensagem:
      """
      ⚠️ Limite de rate atingido
      
      Tier: Basic Access
      Limite: 10 análises/hora
      Próxima análise disponível em: 45 minutos
      
      Upgrade para Pro Access para rate limit maior
      """

  Cenário: Histórico de análises
    Dado que fiz 5 análises hoje
    Quando eu acesso "Analysis History"
    Então devo ver lista com:
      | Data/Hora        | Tipo         | Score | Custo | Status    |
      | 2024-01-15 10:30 | Full Review  | 85    | 15    | Completed |
      | 2024-01-15 11:45 | LLM Analysis | 78    | 5     | Completed |
      | 2024-01-15 12:00 | Basic        | 90    | 1     | Completed |
      | 2024-01-15 14:20 | Security     | 65    | 2     | Completed |
      | 2024-01-15 15:10 | LLM Analysis | 82    | 5     | Completed |
    E devo poder clicar em cada análise para ver detalhes
    E devo poder baixar o relatório em PDF

  Cenário: Análise via API
    Dado que possuo API key válida
    E meu tier é "Enterprise Access"
    Quando eu faço POST para "/api/v1/analyze" com:
      ```json
      {
        "code": "resource \"aws_s3_bucket\" \"example\" {...}",
        "type": "llm_analysis"
      }
      ```
    Então devo receber resposta 200 OK
    E o corpo deve conter análise completa em JSON
    E os tokens devem ser debitados automaticamente
