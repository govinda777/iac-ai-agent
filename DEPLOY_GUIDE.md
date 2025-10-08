# 🚀 Guia de Deploy Rápido - IaC AI Agent

## Deploy em Produção (5 minutos)

### 1. Pré-requisitos
```bash
# Verificar se Docker está instalado
docker --version
docker-compose --version

# Verificar se está no diretório do projeto
pwd
# Deve mostrar: /path/to/iac-ai-agent
```

### 2. Configuração
```bash
# Copiar arquivo de configuração
cp configs/production.env.example .env.prod

# Editar configurações obrigatórias
nano .env.prod
```

**Variáveis obrigatórias:**
- `LLM_API_KEY`: Sua chave da OpenAI
- `LLM_PROVIDER`: openai
- `LLM_MODEL`: gpt-4

### 3. Deploy Automático
```bash
# Executar script de deploy
./scripts/deploy-production.sh production
```

### 4. Verificação
```bash
# Verificar status
docker-compose -f configs/docker-compose.prod.yml ps

# Verificar logs
docker-compose -f configs/docker-compose.prod.yml logs iac-ai-agent

# Testar health check
curl http://localhost:8080/health
```

### 5. URLs Importantes
- **Aplicação**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **API Docs**: http://localhost:8080/swagger

---

## Deploy Manual (Passo a Passo)

### 1. Construir Imagem
```bash
docker build -t iacai-agent:prod -f deployments/Dockerfile.prod .
```

### 2. Iniciar Serviços
```bash
docker-compose -f configs/docker-compose.prod.yml up -d
```

### 3. Verificar Status
```bash
docker-compose -f configs/docker-compose.prod.yml ps
```

---

## Troubleshooting

### Problema: "LLM_API_KEY não configurada"
**Solução:**
```bash
# Editar arquivo de configuração
nano .env.prod

# Adicionar:
LLM_API_KEY=sk-proj-sua-chave-aqui
```

### Problema: "Serviços não ficam saudáveis"
**Solução:**
```bash
# Verificar logs
docker-compose -f configs/docker-compose.prod.yml logs

# Reiniciar serviços
docker-compose -f configs/docker-compose.prod.yml restart
```

### Problema: "Porta 8080 já em uso"
**Solução:**
```bash
# Verificar o que está usando a porta
lsof -i :8080

# Parar serviço conflitante ou mudar porta no .env.prod
PORT=8081
```

---

## Comandos Úteis

### Gerenciamento de Serviços
```bash
# Parar todos os serviços
docker-compose -f configs/docker-compose.prod.yml down

# Reiniciar aplicação
docker-compose -f configs/docker-compose.prod.yml restart iac-ai-agent

# Ver logs em tempo real
docker-compose -f configs/docker-compose.prod.yml logs -f iac-ai-agent
```

### Limpeza
```bash
# Remover containers parados
docker-compose -f configs/docker-compose.prod.yml down --remove-orphans

# Limpar volumes não utilizados
docker volume prune

# Limpar imagens não utilizadas
docker image prune
```

### Monitoramento
```bash
# Verificar uso de recursos
docker stats

# Verificar saúde dos containers
docker-compose -f configs/docker-compose.prod.yml ps
```

---

## Próximos Passos

1. **Configurar SSL/TLS** para HTTPS
2. **Configurar domínio** personalizado
3. **Configurar monitoramento** (Prometheus/Grafana)
4. **Configurar backup** automático
5. **Configurar CI/CD** pipeline

---

## Suporte

- **Documentação**: `/docs/`
- **Issues**: GitHub Issues
- **Logs**: `docker-compose logs iac-ai-agent`
- **Health Check**: `curl http://localhost:8080/health`
