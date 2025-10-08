# Configuração Git Secrets

## Por que usar Git Secrets?

**Segurança**: Secrets nunca devem ser armazenados em variáveis de ambiente ou arquivos de configuração em texto plano. Git Secrets permite criptografar secrets diretamente no repositório Git.

## Instalação

### macOS
```bash
brew install git-secret
```

### Ubuntu/Debian
```bash
sudo apt-get install git-secret
```

### Outros sistemas
Consulte: https://git-secret.io/installation

## Configuração Inicial

### 1. Inicializar Git Secrets no repositório
```bash
git secret init
```

### 2. Adicionar sua chave GPG
```bash
# Gerar chave GPG se não tiver
gpg --gen-key

# Adicionar sua chave ao Git Secrets
git secret tell your-email@example.com
```

### 3. Adicionar secrets
```bash
# GitHub Token
echo "ghp_your_github_token_here" | git secret add github_token

# GitHub Webhook Secret
echo "your_webhook_secret_here" | git secret add github_webhook_secret

# Chave privada da carteira
echo "0x1234567890abcdef..." | git secret add wallet_private_key

# WhatsApp API Key
echo "your_whatsapp_api_key_here" | git secret add whatsapp_api_key
```

### 4. Commit dos arquivos criptografados
```bash
git add .secret/*
git commit -m "Add encrypted secrets"
```

## Uso na Aplicação

A aplicação agora acessa secrets através de métodos específicos:

```go
config := &Config{}

// Obter GitHub Token
token, err := config.GetGitHubToken()
if err != nil {
    log.Fatal("Erro ao obter GitHub token:", err)
}

// Obter chave privada da carteira
privateKey, err := config.GetWalletPrivateKey()
if err != nil {
    log.Fatal("Erro ao obter chave privada:", err)
}

// Obter WhatsApp API Key
whatsappKey, err := config.GetWhatsAppAPIKey()
if err != nil {
    log.Fatal("Erro ao obter WhatsApp API key:", err)
}
```

## Comandos Úteis

### Listar secrets disponíveis
```bash
git secret list
```

### Decriptar todos os secrets (para desenvolvimento)
```bash
git secret reveal
```

### Re-criptografar após modificações
```bash
git secret hide -m
```

### Remover um secret
```bash
git secret remove secret_name
```

## Segurança

- ✅ **Secrets criptografados** no repositório
- ✅ **Não há variáveis de ambiente** com secrets
- ✅ **Acesso controlado** via chaves GPG
- ✅ **Auditoria completa** de quem tem acesso

## Troubleshooting

### Erro: "git-secret não está disponível"
```bash
# Instalar git-secret
brew install git-secret  # macOS
sudo apt-get install git-secret  # Ubuntu
```

### Erro: "Secret não encontrado"
```bash
# Verificar se o secret existe
git secret list

# Adicionar o secret se necessário
echo "your_secret_value" | git secret add secret_name
```

### Erro de permissão GPG
```bash
# Verificar chaves GPG
gpg --list-keys

# Adicionar chave ao Git Secrets
git secret tell your-email@example.com
```
