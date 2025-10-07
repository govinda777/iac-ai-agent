# 🚀 Exemplo de Startup - IaC AI Agent

## 📋 O Que Acontece ao Iniciar

Este documento mostra **exatamente** o que você verá ao executar `go run cmd/agent/main.go` pela primeira vez.

---

## ✅ Startup Bem-Sucedido (Primeira Execução)

```
	██╗ █████╗  ██████╗     █████╗ ██╗     █████╗  ██████╗ ███████╗███╗   ██╗████████╗
	██║██╔══██╗██╔════╝    ██╔══██╗██║    ██╔══██╗██╔════╝ ██╔════╝████╗  ██║╚══██╔══╝
	██║███████║██║         ███████║██║    ███████║██║  ███╗█████╗  ██╔██╗ ██║   ██║   
	██║██╔══██║██║         ██╔══██║██║    ██╔══██║██║   ██║██╔══╝  ██║╚██╗██║   ██║   
	██║██║  ██║╚██████╗    ██║  ██║██║    ██║  ██║╚██████╔╝███████╗██║ ╚████║   ██║   
	╚═╝╚═╝  ╚═╝ ╚═════╝    ╚═╝  ╚═╝╚═╝    ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝   
	
	Infrastructure as Code AI Agent v1.0.0
	Powered by LLM | Secured by Privy.io | Running on Base Network

2025-01-15T10:30:00Z INFO 🚀 Starting IaC AI Agent version=1.0.0

🔍 Executando validações de startup...

─────────────────────────────────────────────────
📋 Validando configuração básica...
✅ LLM_API_KEY encontrada
✅ PRIVY_APP_ID encontrada
✅ PRIVY_APP_SECRET encontrada
✅ WALLET_ADDRESS encontrada
─────────────────────────────────────────────────

─────────────────────────────────────────────────
🤖 Validando conexão com LLM...
INFO Testando conexão com LLM... provider=openai model=gpt-4

⏱️  Enviando mensagem de teste...
✅ LLM respondeu com sucesso
   Latency: 2.3s
   Tokens Used: 8
   Model: gpt-4-0125-preview
   Response: "OK"

✅ LLM validado com sucesso
─────────────────────────────────────────────────

─────────────────────────────────────────────────
🔐 Validando credenciais Privy.io...
INFO Privy credentials configuradas app_id=app_ckse...

✅ Privy.io validado com sucesso
─────────────────────────────────────────────────

─────────────────────────────────────────────────
🌐 Validando conexão com Base Network...
INFO Conectando ao RPC... rpc_url=https://goerli.base.org

✅ Conectado ao RPC
   Chain ID: 84531 (Base Goerli)
   Latest Block: 15432876
   Sync Status: ✅ Synced

✅ Base Network validado com sucesso
─────────────────────────────────────────────────

─────────────────────────────────────────────────
🎨 Validando posse do NFT Nation.fun...
INFO Verificando NFT... wallet=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb

✅ NFT detectado!
   Contract: 0x...
   Balance: 1 NFT
   Token ID: #1337

✅ NFT Nation.fun validado com sucesso
─────────────────────────────────────────────────

─────────────────────────────────────────────────
📦 Inicializando serviços...

✅ Agent Service inicializado
   Templates carregados: 4
   - general-purpose
   - security-specialist
   - cost-optimizer
   - architecture-advisor
─────────────────────────────────────────────────

─────────────────────────────────────────────────
🤖 Verificando agente padrão...
INFO Buscando agente existente... owner=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb

ℹ️  Nenhum agente encontrado para este wallet
✨ Criando novo agente automaticamente...

🎨 Usando template: general-purpose
   ✓ Configurando LLM (GPT-4, temp=0.2)
   ✓ Habilitando análises: Checkov, IAM, Cost, Drift, Preview, Secrets
   ✓ Configurando personalidade: Professional, Encouraging, Balanced
   ✓ Carregando knowledge base: Terraform Expert, AWS Expert
   ✓ Definindo limites: 100 req/h, $10/dia
   ✓ Inicializando métricas

✅ Novo agente criado automaticamente!
   ID: agent-f7b3c9e1-2a4d-4b5e-9c8f-1d6e3a7b2c5d
   Name: IaC Agent - 0x742d35
   Template: General Purpose
   Status: Active

INFO Agente configurado
   id=agent-f7b3c9e1-2a4d-4b5e-9c8f-1d6e3a7b2c5d
   name="IaC Agent - 0x742d35"
   owner=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
   status=active

✅ Agente pronto!
─────────────────────────────────────────────────

🎉 Todas as validações passaram! Aplicação pronta para iniciar.

============================================================
📊 RELATÓRIO DE VALIDAÇÃO DE STARTUP
============================================================
✅ Status: PASSOU

📋 Checklist de Validações:
  ✅ LLM Connection
  ✅ Privy.io Credentials
  ✅ Base Network
  ✅ Nation.fun NFT
  ✅ Default Agent

🤖 Agent Details:
  ID: agent-f7b3c9e1-2a4d-4b5e-9c8f-1d6e3a7b2c5d
  Name: IaC Agent - 0x742d35

============================================================

✅ Validação completa - Aplicação iniciando...

─────────────────────────────────────────────────
🌐 Configurando servidor HTTP...

✅ Rotas configuradas
✅ Middleware instalado
✅ Swagger UI pronto

🚀 Servidor HTTP iniciado
   Address: localhost:8080
   Environment: development

📚 Swagger UI: http://localhost:8080/swagger/
❤️  Health Check: http://localhost:8080/health

✨ Aplicação pronta para receber requisições!
Press Ctrl+C to shutdown gracefully

```

---

## 🔄 Startup Subsequente (Agente Já Existe)

Na próxima vez que você executar, o agente já existirá:

```
	Infrastructure as Code AI Agent v1.0.0
	Powered by LLM | Secured by Privy.io | Running on Base Network

2025-01-15T11:00:00Z INFO 🚀 Starting IaC AI Agent version=1.0.0

🔍 Executando validações de startup...

[... validações anteriores ...]

─────────────────────────────────────────────────
🤖 Verificando agente padrão...
INFO Buscando agente existente... owner=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb

✅ Agente existente encontrado!
   ID: agent-f7b3c9e1-2a4d-4b5e-9c8f-1d6e3a7b2c5d
   Name: IaC Agent - 0x742d35
   Status: Active
   Requests Today: 45
   Cost Today: $2.35

✅ Agente pronto!
─────────────────────────────────────────────────

[... resto do startup ...]
```

---

## ❌ Startup com Falha (LLM Inválido)

Se a API key do LLM estiver errada:

```
	Infrastructure as Code AI Agent v1.0.0
	Powered by LLM | Secured by Privy.io | Running on Base Network

2025-01-15T12:00:00Z INFO 🚀 Starting IaC AI Agent version=1.0.0

🔍 Executando validações de startup...

─────────────────────────────────────────────────
🤖 Validando conexão com LLM...
INFO Testando conexão com LLM... provider=openai model=gpt-4

⏱️  Enviando mensagem de teste...

❌ ERRO: Falha ao comunicar com LLM
   Error: invalid API key
   Hint: Verifique sua LLM_API_KEY em .env

─────────────────────────────────────────────────

============================================================
📊 RELATÓRIO DE VALIDAÇÃO DE STARTUP
============================================================
❌ Status: FALHOU

📋 Checklist de Validações:
  ❌ LLM Connection
  ⏭️  Privy.io Credentials (não executado)
  ⏭️  Base Network (não executado)
  ⏭️  Nation.fun NFT (não executado)
  ⏭️  Default Agent (não executado)

❌ Erros Encontrados:
  ❌ LLM validation failed: invalid API key

============================================================

💥 APLICAÇÃO NÃO PODE INICIAR - Validação falhou
Por favor, corrija os erros acima e tente novamente.

Erros críticos:
  - LLM validation failed: invalid API key

panic: Startup validation failed
```

---

## ❌ Startup com Falha (NFT Faltando)

Se a wallet não tiver o NFT Nation.fun:

```
	Infrastructure as Code AI Agent v1.0.0
	Powered by LLM | Secured by Privy.io | Running on Base Network

2025-01-15T13:00:00Z INFO 🚀 Starting IaC AI Agent version=1.0.0

🔍 Executando validações de startup...

[... LLM, Privy, Base Network passam ...]

─────────────────────────────────────────────────
🎨 Validando posse do NFT Nation.fun...
INFO Verificando NFT... wallet=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb

❌ NFT NÃO ENCONTRADO!
   Wallet: 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
   Balance: 0 NFTs
   
🛒 AÇÃO NECESSÁRIA:
   Compre um NFT Nation.fun em: https://nation.fun/
   
   Para usar a aplicação, você precisa possuir um NFT Nation.fun
   na wallet configurada no WALLET_ADDRESS.

─────────────────────────────────────────────────

============================================================
📊 RELATÓRIO DE VALIDAÇÃO DE STARTUP
============================================================
❌ Status: FALHOU

📋 Checklist de Validações:
  ✅ LLM Connection
  ✅ Privy.io Credentials
  ✅ Base Network
  ❌ Nation.fun NFT
  ⏭️  Default Agent (não executado)

❌ Erros Encontrados:
  ❌ Nation.fun NFT validation failed: wallet 0x742d35... não possui NFT

⚠️ Avisos:
  ⚠️ Verificação real de NFT Nation.fun implementada

============================================================

💥 APLICAÇÃO NÃO PODE INICIAR - Validação falhou
Por favor, corrija os erros acima e tente novamente.

Erros críticos:
  - Nation.fun NFT validation failed: wallet 0x742d35... não possui NFT

🛒 Compre Nation.fun NFT: https://nation.fun/

panic: Startup validation failed
```

---

## ⚠️ Startup com Avisos (Warnings)

Se houver avisos não-críticos:

```
	Infrastructure as Code AI Agent v1.0.0
	Powered by LLM | Secured by Privy.io | Running on Base Network

2025-01-15T14:00:00Z INFO 🚀 Starting IaC AI Agent version=1.0.0

🔍 Executando validações de startup...

[... todas validações passam ...]

─────────────────────────────────────────────────
🤖 Validando conexão com LLM...
INFO Testando conexão com LLM... provider=openai model=gpt-4

⏱️  Enviando mensagem de teste...
✅ LLM respondeu com sucesso
   Latency: 12.5s ⚠️ (esperado < 10s)
   
⚠️ AVISO: Latência LLM alta
   A resposta demorou 12.5 segundos
   Esperado: < 10 segundos
   Possível causa: Região distante ou throttling
   Recomendação: Considere usar servidor mais próximo

✅ LLM validado com sucesso
─────────────────────────────────────────────────

[... resto das validações ...]

============================================================
📊 RELATÓRIO DE VALIDAÇÃO DE STARTUP
============================================================
✅ Status: PASSOU

📋 Checklist de Validações:
  ✅ LLM Connection
  ✅ Privy.io Credentials
  ✅ Base Network
  ✅ Nation.fun NFT
  ✅ Default Agent

⚠️ Avisos:
  ⚠️ LLM latency alta: 12.5s (esperado < 10s)

============================================================

✅ Validação completa - Aplicação iniciando...
```

---

## 📊 Logs Detalhados (Debug Mode)

Com `LOGGING_LEVEL=debug`:

```
2025-01-15T15:00:00Z DEBUG Loading configuration... path=configs/app.yaml
2025-01-15T15:00:00Z DEBUG Environment variables loaded count=23
2025-01-15T15:00:00Z DEBUG Config merged from env and yaml
2025-01-15T15:00:00Z INFO 🚀 Starting IaC AI Agent version=1.0.0

2025-01-15T15:00:00Z DEBUG Creating LLM client...
2025-01-15T15:00:00Z DEBUG LLM client created provider=openai model=gpt-4
2025-01-15T15:00:00Z DEBUG Building test prompt...
2025-01-15T15:00:00Z DEBUG Sending request to LLM... prompt_length=56
2025-01-15T15:00:00Z DEBUG Request sent, waiting response...
2025-01-15T15:00:00Z DEBUG Response received tokens=8 latency=2.3s
2025-01-15T15:00:00Z DEBUG Parsing response...
2025-01-15T15:00:00Z DEBUG Response parsed successfully content="OK"

2025-01-15T15:00:00Z DEBUG Connecting to Base Network... rpc=https://goerli.base.org
2025-01-15T15:00:00Z DEBUG ethclient.Dial() called
2025-01-15T15:00:00Z DEBUG Connection established
2025-01-15T15:00:00Z DEBUG Getting chain ID...
2025-01-15T15:00:00Z DEBUG Chain ID received: 84531
2025-01-15T15:00:00Z DEBUG Validating chain ID... expected=84531 got=84531
2025-01-15T15:00:00Z DEBUG Getting latest block...
2025-01-15T15T:00:00Z DEBUG Block number: 15432876

[... modo verboso completo ...]
```

---

## 🎯 Resumo

### ✅ Sucesso
- Banner exibido
- 5 validações executadas
- Agente criado automaticamente (primeira vez)
- Servidor HTTP iniciado
- Pronto para requisições

### ❌ Falha
- Valida até encontrar erro crítico
- Exibe mensagem clara do problema
- Sugere ação para correção
- Aplicação não inicia (`panic`)

### ⚠️ Avisos
- Aplicação inicia normalmente
- Avisos exibidos no relatório
- Funcionamento não comprometido
- Recomendações fornecidas

---

## 📝 Comandos Úteis

### Ver Logs em Tempo Real

```bash
# Executar com logs detalhados
LOGGING_LEVEL=debug go run cmd/agent/main.go

# Salvar logs em arquivo
go run cmd/agent/main.go 2>&1 | tee startup.log

# Apenas erros
LOGGING_LEVEL=error go run cmd/agent/main.go
```

### Forçar Recriação de Agente

```bash
# Limpar storage (in-memory não persiste)
# Basta reiniciar a aplicação

# Em produção com DB, deletar agente:
curl -X DELETE http://localhost:8080/api/v1/agents/{id}
```

### Health Check

```bash
# Verificar se aplicação está rodando
curl http://localhost:8080/health

# Response:
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": "2h30m15s",
  "validations": {
    "llm": true,
    "privy": true,
    "base_network": true,
    "nation_nft": true,
    "default_agent": true
  }
}
```

---

**Data**: 2025-01-15  
**Versão**: 1.0.0  
**Status**: ✅ Produção
