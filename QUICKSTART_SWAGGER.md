# 🚀 Quickstart - Swagger UI

Guia rápido para começar a usar o IaC AI Agent com Swagger UI.

## ⚡ Início Rápido (5 minutos)

### 1. Gere a Documentação Swagger

```bash
make swagger
```

### 2. Execute a Aplicação

```bash
make run
```

Ou rode tudo de uma vez:

```bash
make run-swagger
```

### 3. Acesse o Swagger UI

Abra no navegador: **http://localhost:8080/swagger/**

---

## 🎯 Testando os Endpoints

### Health Check

1. No Swagger UI, clique em **GET /health**
2. Clique em **"Try it out"**
3. Clique em **"Execute"**
4. Veja a resposta: `{"status": "healthy", "service": "iac-ai-agent", "version": "1.0.0"}`

### Análise de Código Terraform

1. Clique em **POST /analyze**
2. Clique em **"Try it out"**
3. Cole este exemplo no Request body:

```json
{
  "repository": "my-org/infrastructure",
  "path": "terraform/",
  "content": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-test-bucket\"\n  acl    = \"public-read\"\n}\n\nresource \"aws_instance\" \"web\" {\n  ami           = \"ami-12345678\"\n  instance_type = \"m5.xlarge\"\n}"
}
```

4. Clique em **"Execute"**
5. Veja o resultado com:
   - Score geral
   - Análise de segurança (Checkov)
   - Análise de IAM
   - Sugestões de melhorias
   - Estimativas de custo

### Review de Pull Request

1. Clique em **POST /review**
2. Clique em **"Try it out"**
3. Cole este exemplo:

```json
{
  "repository": "my-repo",
  "pr_number": 42,
  "owner": "my-org"
}
```

4. Clique em **"Execute"**
5. Veja o resultado do review completo

---

## 📚 Usando via cURL

### Health Check

```bash
curl http://localhost:8080/health
```

### Analisar Código

```bash
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "repository": "my-org/infrastructure",
    "path": "terraform/",
    "content": "resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"my-test-bucket\"\n}"
  }'
```

### Review de PR

```bash
curl -X POST http://localhost:8080/review \
  -H "Content-Type: application/json" \
  -d '{
    "repository": "my-repo",
    "pr_number": 42,
    "owner": "my-org"
  }'
```

---

## 🎨 Features do Swagger UI

### 1. Documentação Interativa
- Veja todos os endpoints disponíveis
- Leia descrições detalhadas
- Visualize estruturas de request/response

### 2. Teste Direto no Navegador
- Clique em "Try it out"
- Preencha os campos
- Execute requisições reais
- Veja respostas formatadas

### 3. Validação Automática
- Schema validation built-in
- Exemplos pré-preenchidos
- Erros explicados

### 4. Exportação
- Baixe `swagger.json` ou `swagger.yaml`
- Importe no Postman
- Use em outras ferramentas

---

## 🔧 Comandos Úteis

```bash
# Gerar documentação Swagger
make swagger

# Rodar aplicação
make run

# Rodar com Swagger (gera docs + executa)
make run-swagger

# Instalar swag CLI manualmente
make swagger-install

# Limpar tudo (incluindo docs)
make clean

# Build
make build

# Testes
make test
```

---

## 📂 Arquivos Gerados

Após `make swagger`, você terá:

```
docs/
├── docs.go          # Código Go da documentação
├── swagger.json     # Spec OpenAPI em JSON
└── swagger.yaml     # Spec OpenAPI em YAML
```

---

## 🐛 Troubleshooting

### Swagger UI não carrega

```bash
# 1. Verifique se os docs foram gerados
ls docs/

# 2. Se não existirem, gere
make swagger

# 3. Rebuild
make build

# 4. Execute
make run
```

### Erro "docs package not found"

```bash
# Sempre gere a documentação antes de compilar
make swagger
make build
```

### Mudanças não aparecem

```bash
# Após qualquer mudança nos handlers, regenere
make swagger
make run
```

### Port 8080 já em uso

```bash
# Mude a porta no configs/app.yaml
# Ou use variável de ambiente
PORT=8081 make run
```

---

## 🎓 Próximos Passos

1. **Explore a UI**: Navegue pelos endpoints e modelos
2. **Teste com dados reais**: Use seus arquivos Terraform
3. **Integre com CI/CD**: Adicione análises automáticas nos PRs
4. **Configure webhooks**: Integre com GitHub para reviews automáticos

---

## 📖 Documentação Adicional

- [README Principal](docs/README.md)
- [Guia Completo do Swagger](docs/SWAGGER.md)
- [Arquitetura](docs/ARCHITECTURE.md)
- [Roadmap](docs/IMPLEMENTATION_ROADMAP.md)

---

## 💡 Dicas Pro

1. **Use Ctrl+F** no Swagger UI para buscar endpoints
2. **Salve exemplos** que funcionam para reutilizar
3. **Exporte swagger.json** para usar no Postman/Insomnia
4. **Marque a aba** do Swagger para acesso rápido
5. **Configure variáveis de ambiente** para diferentes ambientes

---

**Pronto! Você está pronto para usar o IaC AI Agent! 🎉**

Qualquer dúvida, consulte a [documentação completa](docs/README.md) ou abra uma issue.
