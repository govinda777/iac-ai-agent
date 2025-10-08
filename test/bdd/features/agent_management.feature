# language: pt
Funcionalidade: Fluxo de Criação Automática de Agentes
  Como um usuário autenticado
  Eu quero que um agente seja criado automaticamente para mim
  Para poder usar o sistema de análise inteligente

  Contexto:
    Dado que estou autenticado com wallet "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
    E que o sistema está inicializado
    E que os templates de agente estão disponíveis

  @agent_creation @automatic
  Cenário: Criação automática de agente na primeira visita
    Dado que sou um novo usuário sem agente
    Quando eu acesso o sistema pela primeira vez
    Então o sistema deve verificar se tenho um agente
    E deve detectar que não tenho agente
    E deve criar automaticamente um agente usando template "General Purpose"
    E o agente deve ser configurado com:
      | Configuração           | Valor                    |
      | Template               | General Purpose          |
      | LLM Provider           | OpenAI                   |
      | LLM Model              | GPT-4                    |
      | Temperature            | 0.7                      |
      | Max Tokens             | 4000                     |
      | Enable Checkov          | true                     |
      | Enable IAM Analysis     | true                     |
      | Enable Cost Analysis    | true                     |
      | Enable Drift Detection  | true                     |
      | Enable Preview Analysis | true                     |
      | Enable Secrets Scanning | true                     |
    E devo receber notificação "Agente criado automaticamente!"
    E devo poder usar o agente imediatamente

  @agent_creation @templates
  Cenário: Visualizar templates de agente disponíveis
    Quando eu acesso "Create Agent"
    Então devo ver 4 templates disponíveis:
      | ID   | Nome                | Descrição                           | Recomendado |
      | 1    | General Purpose     | Análise completa e versátil         | ✓           |
      | 2    | Security Specialist | Foco em segurança e compliance      | -           |
      | 3    | Cost Optimizer      | Otimização de custos e recursos     | -           |
      | 4    | Architecture Advisor| Análise arquitetural e design      | -           |
    E cada template deve mostrar:
      - Descrição detalhada
      - Casos de uso
      - Configurações padrão
      - Tags relevantes

  @agent_creation @custom
  Cenário: Criar agente customizado
    Dado que estou na página "Create Agent"
    Quando eu seleciono "Custom Agent"
    E configuro:
      | Campo                | Valor                    |
      | Nome                 | Meu Agente DevOps        |
      | Descrição            | Agente especializado     |
      | LLM Provider         | Anthropic                |
      | LLM Model            | Claude-3-Opus            |
      | Temperature          | 0.5                      |
      | Enable Checkov       | true                     |
      | Enable LLM Analysis  | true                     |
      | Enable Cost Analysis | false                    |
    E clico em "Create Agent"
    Então o agente deve ser criado com minhas configurações
    E devo ver "Agente 'Meu Agente DevOps' criado com sucesso!"
    E o agente deve estar ativo e pronto para uso

  @agent_creation @personality
  Cenário: Configurar personalidade do agente
    Dado que estou criando um agente customizado
    Quando eu configuro a personalidade:
      | Campo              | Valor        |
      | Estilo             | Professional |
      | Tom                | Encorajador  |
      | Verbosidade        | Balanceada   |
      | Usa Emojis         | true         |
      | Explica Raciocínio | true         |
      | Dá Exemplos        | true         |
      | Idioma             | pt-BR        |
    E salvo o agente
    Então o agente deve ter a personalidade configurada
    E as respostas devem seguir o estilo definido

  @agent_creation @limits
  Cenário: Configurar limites do agente
    Dado que estou criando um agente
    Quando eu configuro os limites:
      | Limite                    | Valor |
      | Requests por hora         | 50    |
      | Requests por dia          | 500   |
      | Tokens por request        | 2000  |
      | Custo máximo por request  | $0.25 |
      | Custo máximo por dia      | $5.00 |
    E salvo o agente
    Então o agente deve respeitar os limites configurados
    E devo receber alertas quando os limites forem atingidos

  @agent_creation @knowledge
  Cenário: Adicionar conhecimento especializado ao agente
    Dado que estou criando um agente
    Quando eu adiciono conhecimento:
      | Tipo                | Conteúdo                                    |
      | Regras de Negócio    | Sempre sugerir uso de módulos Terraform    |
      | Padrões de Segurança | Nunca usar ACLs públicas em S3            |
      | Otimizações          | Preferir instâncias spot quando possível   |
      | Documentação         | AWS Well-Architected Framework             |
    E salvo o agente
    Então o agente deve incorporar o conhecimento
    E as análises devem seguir as regras definidas

  @agent_creation @error_handling
  Cenário: Falha na criação de agente por configuração inválida
    Dado que estou criando um agente
    Quando eu configuro:
      | Campo        | Valor Inválido |
      | Nome         | ""             |
      | Max Tokens   | -100           |
      | Temperature  | 2.5            |
    E tento salvar
    Então devo ver erros de validação:
      - "Nome é obrigatório"
      - "Max Tokens deve ser positivo"
      - "Temperature deve estar entre 0 e 1"
    E o agente não deve ser criado
    E devo poder corrigir os erros

  @agent_management @list
  Cenário: Listar agentes do usuário
    Dado que tenho 3 agentes criados
    Quando eu acesso "My Agents"
    Então devo ver lista com:
      | Nome                | Status | Template         | Última Atividade |
      | IaC Agent - 0x742d  | Active | General Purpose | Hoje 14:30       |
      | Security Bot        | Active | Security Spec    | Ontem 16:45      |
      | Cost Optimizer      | Inactive| Cost Optimizer  | Semana passada   |
    E devo poder:
      - Ver detalhes de cada agente
      - Editar configurações
      - Ativar/desativar
      - Deletar agente

  @agent_management @update
  Cenário: Atualizar configurações do agente
    Dado que tenho um agente "Meu Agente DevOps"
    Quando eu acesso "Edit Agent"
    E modifico:
      | Campo        | Valor Anterior | Valor Novo |
      | Temperature  | 0.7            | 0.5        |
      | Max Tokens   | 4000           | 2000       |
      | Enable Cost  | true           | false      |
    E salvo as alterações
    Então o agente deve ser atualizado
    E devo ver "Agente atualizado com sucesso!"
    E as novas configurações devem estar ativas

  @agent_management @delete
  Cenário: Deletar agente
    Dado que tenho um agente "Test Agent"
    E que não é meu único agente
    Quando eu acesso "Delete Agent"
    E confirmo a exclusão
    Então o agente deve ser removido
    E devo ver "Agente deletado com sucesso"
    E o agente não deve mais aparecer na lista

  @agent_management @metrics
  Cenário: Visualizar métricas do agente
    Dado que tenho um agente ativo há 30 dias
    Quando eu acesso "Agent Metrics"
    Então devo ver estatísticas:
      | Métrica                    | Valor |
      | Total de Análises          | 150   |
      | Análises Hoje             | 5     |
      | Análises Esta Semana      | 25    |
      | Análises Este Mês         | 80    |
      | Score Médio               | 82.5  |
      | Tempo Médio de Resposta   | 3.2s  |
      | Tokens Gastos Total       | 750   |
      | Custo Total               | $15.50|
    E devo ver gráficos de uso ao longo do tempo
