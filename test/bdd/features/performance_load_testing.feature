# language: pt
Funcionalidade: Testes de Performance e Carga
  Como desenvolvedor
  Eu quero testar performance e capacidade do sistema
  Para garantir escalabilidade e responsividade

  Contexto:
    Dado que o sistema está configurado
    E que os recursos estão disponíveis

  @performance @response_time
  Cenário: Tempo de resposta de endpoints críticos
    Dado que o sistema está rodando
    Quando eu faço requisições para endpoints críticos
    Então os tempos de resposta devem estar dentro dos limites:
      | Endpoint                    | Método | Limite | Status Esperado |
      | /health                     | GET    | < 100ms| 200 OK          |
      | /api/auth/web3/verify       | POST   | < 500ms| 200 OK          |
      | /api/auth/web3/check-access | POST   | < 300ms| 200 OK          |
      | /analyze                    | POST   | < 5s   | 200 OK          |
      | /review                     | POST   | < 10s  | 200 OK          |
    E todos os endpoints devem responder consistentemente

  @performance @concurrent_users
  Cenário: Sistema com múltiplos usuários simultâneos
    Dado que o sistema está configurado para produção
    Quando 50 usuários fazem login simultaneamente
    E cada usuário executa 5 análises
    Então o sistema deve processar todas as requisições:
      | Métrica           | Valor Esperado | Valor Real |
      | Usuários Ativos   | 50             | 50         |
      | Análises Totais   | 250            | 250        |
      | Taxa de Sucesso   | > 99%          | 99.6%      |
      | Tempo Médio       | < 3s           | 2.8s       |
      | Tempo P95         | < 5s           | 4.2s       |
    E nenhum usuário deve ser bloqueado
    E todos devem receber resultados válidos

  @performance @memory_usage
  Cenário: Uso de memória durante operações intensivas
    Dado que o sistema está monitorado
    Quando eu executo 100 análises consecutivas
    Então o uso de memória deve permanecer estável:
      | Métrica           | Valor Inicial | Valor Final | Limite |
      | Memory Usage      | 50MB          | 65MB        | < 200MB|
      | Memory Growth     | -             | +15MB       | < 50MB |
      | GC Frequency      | Normal        | Normal      | -      |
      | Memory Leaks      | 0             | 0           | 0      |
    E não deve haver vazamentos de memória
    E o garbage collector deve funcionar normalmente

  @performance @database_performance
  Cenário: Performance do banco de dados
    Dado que o banco está configurado
    Quando eu executo operações de banco intensivas:
      - 1000 inserções de análise
      - 500 consultas de histórico
      - 100 atualizações de usuário
    Então as operações devem ser eficientes:
      | Operação          | Tempo Médio | Limite | Throughput |
      | Insert Analysis   | 5ms         | < 10ms | 200 ops/s  |
      | Query History     | 2ms         | < 5ms  | 500 ops/s  |
      | Update User       | 3ms         | < 8ms  | 300 ops/s  |
    E não deve haver locks ou deadlocks
    E as conexões devem ser gerenciadas eficientemente

  @load @stress_test
  Cenário: Teste de estresse - carga máxima
    Dado que o sistema está preparado para teste de estresse
    Quando eu aplico carga máxima sustentada:
      - 200 usuários simultâneos
      - 1000 análises por minuto
      - Operações contínuas por 30 minutos
    Então o sistema deve manter estabilidade:
      | Métrica           | Valor Esperado | Valor Real |
      | CPU Usage         | < 80%          | 75%        |
      | Memory Usage      | < 500MB        | 420MB      |
      | Response Time     | < 10s          | 8.5s       |
      | Error Rate        | < 1%           | 0.3%       |
      | Throughput        | > 1000/min     | 1200/min   |
    E o sistema deve se recuperar automaticamente
    E não deve haver falhas catastróficas

  @load @spike_test
  Cenário: Teste de pico - carga súbita
    Dado que o sistema está em operação normal
    Quando eu aplico carga súbita:
      - De 10 usuários para 200 usuários em 1 minuto
      - Pico de 500 análises simultâneas
      - Duração de 5 minutos
    Então o sistema deve lidar com o pico:
      | Métrica           | Durante Pico | Após Pico |
      | Response Time     | < 15s       | < 3s      |
      | Error Rate        | < 5%        | < 1%      |
      | Queue Length      | < 100       | < 10      |
      | Recovery Time     | -           | < 2min    |
    E deve voltar ao normal após o pico
    E nenhum dado deve ser perdido

  @performance @cache_performance
  Cenário: Performance do sistema de cache
    Dado que o Redis está configurado
    Quando eu executo operações com cache:
      - 1000 consultas de configuração
      - 500 consultas de histórico
      - 200 consultas de métricas
    Então o cache deve melhorar performance:
      | Operação          | Sem Cache | Com Cache | Melhoria |
      | Config Query      | 50ms      | 2ms       | 96%      |
      | History Query     | 100ms     | 5ms       | 95%      |
      | Metrics Query      | 200ms     | 10ms      | 95%      |
    E a taxa de hit do cache deve ser > 90%
    E o cache deve ser invalidado corretamente

  @performance @external_services
  Cenário: Performance de serviços externos
    Dado que todos os serviços externos estão monitorados
    Quando eu executo análises que dependem de serviços externos
    Então os serviços devem responder dentro dos limites:
      | Serviço        | Latência | Limite | Disponibilidade |
      | Privy.io       | 150ms    | < 500ms| 99.9%          |
      | Base Network   | 200ms    | < 1s   | 99.8%          |
      | Nation.fun     | 300ms    | < 2s   | 99.5%          |
      | OpenAI API     | 500ms    | < 3s   | 99.9%          |
    E timeouts devem ser configurados adequadamente
    E retry logic deve funcionar corretamente

  @load @endurance_test
  Cenário: Teste de resistência - operação contínua
    Dado que o sistema está configurado
    Quando eu executo operações contínuas por 24 horas:
      - Login/logout cíclico
      - Análises regulares
      - Operações de banco
      - Limpeza de cache
    Então o sistema deve manter estabilidade:
      | Métrica           | Início | 12h | 24h | Limite |
      | Memory Usage      | 50MB   | 55MB| 60MB| < 100MB|
      | CPU Usage         | 20%    | 25% | 30% | < 50%  |
      | Response Time     | 2s     | 2.1s| 2.2s| < 5s   |
      | Error Rate        | 0%     | 0.1%| 0.2%| < 1%   |
    E não deve haver degradação significativa
    E logs devem ser rotacionados corretamente

  @performance @optimization
  Cenário: Otimização de queries e operações
    Dado que o sistema está em produção
    Quando eu identifico operações lentas
    E aplico otimizações:
      - Índices de banco otimizados
      - Queries reescritas
      - Cache implementado
      - Pool de conexões ajustado
    Então as melhorias devem ser mensuráveis:
      | Operação          | Antes | Depois | Melhoria |
      | User Lookup       | 100ms | 10ms   | 90%      |
      | Analysis History  | 500ms | 50ms   | 90%      |
      | Config Load       | 200ms | 5ms    | 97.5%    |
      | Metrics Query     | 300ms | 20ms   | 93.3%    |
    E a experiência do usuário deve melhorar
    E os recursos devem ser utilizados mais eficientemente

  @load @scalability
  Cenário: Teste de escalabilidade horizontal
    Dado que o sistema está configurado para escalar
    Quando eu aumento a carga gradualmente:
      - 10 usuários → 50 usuários → 100 usuários → 200 usuários
      - Cada nível por 10 minutos
    Então o sistema deve escalar adequadamente:
      | Usuários | Response Time | CPU Usage | Memory | Status |
      | 10       | 1s           | 20%       | 50MB   | ✓      |
      | 50       | 2s           | 40%       | 80MB   | ✓      |
      | 100      | 3s           | 60%       | 120MB  | ✓      |
      | 200      | 5s           | 80%       | 180MB  | ✓      |
    E novos recursos devem ser provisionados automaticamente
    E o load balancer deve distribuir carga corretamente

  @performance @monitoring
  Cenário: Monitoramento de performance em tempo real
    Dado que o sistema está monitorado
    Quando eu executo operações normais
    Então as métricas devem ser coletadas:
      | Métrica           | Valor | Threshold | Status |
      | Response Time     | 2.5s  | < 5s      | ✓      |
      | Throughput        | 100/min| > 50/min  | ✓      |
      | Error Rate        | 0.5%  | < 2%      | ✓      |
      | CPU Usage         | 45%   | < 80%     | ✓      |
      | Memory Usage      | 120MB | < 500MB   | ✓      |
    E alertas devem ser gerados quando thresholds são excedidos
    E dashboards devem mostrar dados em tempo real
