# 🪝 Git Hooks

## O Que São Git Hooks?

Git hooks são scripts que são executados automaticamente em eventos específicos do Git (commit, push, merge, etc).

## Hooks Configurados

### Pre-Push Hook

**Arquivo**: `.githooks/pre-push`

**Quando executa**: Antes de cada `git push`

**O que faz**:
1. ✅ Executa linter (golangci-lint)
2. ✅ Executa testes unitários
3. ✅ Executa testes de integração
4. ✅ Verifica build

**Se qualquer etapa falhar**:
- ❌ Push é **bloqueado**
- 📋 Erros são exibidos
- 💡 Sugestões de correção são fornecidas

## Como Instalar

### Opção 1: Via Make (Recomendado)

```bash
make setup
# Ou apenas os hooks
make git-hooks
```

### Opção 2: Manual

```bash
# Dar permissão de execução
chmod +x .githooks/pre-push

# Configurar git para usar .githooks
git config core.hooksPath .githooks
```

## Como Funciona

### Fluxo Normal

```bash
$ git push origin main

🔍 ==================================
🔍 PRE-PUSH: Executando Testes
🔍 ==================================

📋 [1/4] Executando Linter...
✅ Linter passou!

🧪 [2/4] Executando Testes Unitários...
✅ Testes unitários passaram!

🔗 [3/4] Executando Testes de Integração...
✅ Testes de integração passaram!

🏗️  [4/4] Verificando Build...
✅ Build passou!

🎉 ==================================
🎉 TODOS OS TESTES PASSARAM!
🎉 ==================================

📦 Prosseguindo com push...

Enumerating objects: 5, done.
...
```

### Quando Há Erros

```bash
$ git push origin main

🔍 ==================================
🔍 PRE-PUSH: Executando Testes
🔍 ==================================

📋 [1/4] Executando Linter...
❌ Linter falhou! Corrija os erros antes do push.

💡 Dica: Execute 'make lint' para ver os erros

error: failed to push some refs to 'origin'
```

## Como Desabilitar (Temporariamente)

**⚠️ Use com MUITO cuidado!**

```bash
# Pular todos os hooks
git push --no-verify

# Ou
git push --no-verify origin main
```

**Quando usar `--no-verify`**:
- ✅ Hotfix urgente em produção
- ✅ WIP push para branch pessoal
- ❌ NUNCA em main/develop sem validação

## Como Desinstalar

```bash
# Voltar para hooks padrão do git
git config --unset core.hooksPath
```

## Problemas Comuns

### 1. "Permission denied"

```bash
# Solução
chmod +x .githooks/pre-push
```

### 2. "golangci-lint: command not found"

```bash
# Solução
make lint-install

# Ou manual
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
```

### 3. "Testes muito lentos"

```bash
# Para desenvolvimento rápido, use branches:
git checkout -b feature/my-feature
# Trabalhe normalmente
git push origin feature/my-feature

# Merge para main só após testes locais
make test
git checkout main
git merge feature/my-feature
git push origin main  # Aqui os hooks rodam
```

### 4. "Hook não está executando"

```bash
# Verificar configuração
git config core.hooksPath

# Deve retornar: .githooks

# Se não, reinstalar
make git-hooks
```

## Customização

### Adicionar Novas Verificações

Edite `.githooks/pre-push`:

```bash
# Exemplo: Adicionar check de TODO
echo "Verificando TODOs críticos..."
if grep -r "TODO: URGENT" .; then
    echo "❌ TODO urgente encontrado!"
    exit 1
fi
```

### Mudar Ordem de Execução

Reordene as seções no arquivo `pre-push`.

### Adicionar Outros Hooks

Crie novos arquivos em `.githooks/`:
- `pre-commit` - Antes de cada commit
- `post-commit` - Depois de cada commit
- `pre-merge-commit` - Antes de merge
- etc.

## Boas Práticas

### 1. Teste Localmente Antes de Push

```bash
make test-unit
make test-integration
make lint
```

### 2. Commits Pequenos e Frequentes

Evite commits grandes que demoram para testar.

### 3. Use Branches para Desenvolvimento

```bash
# Desenvolvimento
git checkout -b feature/new-feature
# Pushs rápidos sem hooks
git push origin feature/new-feature --no-verify

# Antes de merge para main
make test  # Valida localmente
git checkout main
git merge feature/new-feature
git push  # Hooks executam aqui
```

### 4. Mantenha Testes Rápidos

Testes que demoram muito desencorajam commits.

**Meta**: Pre-push < 2 minutos

### 5. CI/CD como Backup

Hooks locais podem ser pulados (`--no-verify`), então sempre tenha CI/CD configurado como backup.

## Próximos Hooks Planejados

### Pre-Commit (Futuro)

```bash
# Validações rápidas antes de commit
- Formatar código (gofmt)
- Verificar mensagens de commit
- Check de secrets (.env no commit)
```

### Commit-Msg (Futuro)

```bash
# Validar mensagem de commit
- Padrão: feat|fix|docs|chore|test|refactor
- Exemplo: "feat: adiciona login com Google"
```

---

**Última Atualização**: 2025-01-15  
**Versão**: 1.0.0  
**Status**: ✅ Ativo
