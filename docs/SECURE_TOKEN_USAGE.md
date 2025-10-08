# 🔐 Uso Seguro de Tokens - IaC AI Agent

## ⚠️ ALERTA DE SEGURANÇA IMPORTANTE

Por razões de segurança, **NÃO UTILIZAMOS MAIS CHAVES PRIVADAS DIRETAMENTE** no código ou variáveis de ambiente. Este documento explica as alternativas seguras para autenticação.

## 🛡️ Práticas Recomendadas

### ❌ O que NÃO fazer:

- **NUNCA armazene chaves privadas** em código-fonte, arquivos de configuração ou variáveis de ambiente
- **NUNCA inclua chaves privadas** em repositórios Git, mesmo em arquivos `.env`
- **NUNCA compartilhe chaves privadas** por e-mail, chat ou qualquer meio não seguro
- **NUNCA use a mesma chave privada** para desenvolvimento e produção

### ✅ O que fazer:

- **USE tokens pré-gerados** para autenticação com Nation.fun
- **USE serviços de assinatura externos** como AWS KMS, HashiCorp Vault ou similares
- **USE diferentes tokens** para ambientes de desenvolvimento e produção
- **USE rotação periódica** de tokens de acesso

## 🔄 Alternativas Seguras para Autenticação

### 1. Token Pré-gerado (Recomendado para Produção)

```bash
# No arquivo .env
WALLET_ADDRESS=0x17eDfB8a794ec4f13190401EF7aF1c17f3cc90c5
WALLET_TOKEN=nft_v1_abc123...  # Token gerado externamente
```

### 2. Serviço de Assinatura Externo

Utilize um serviço como AWS KMS ou HashiCorp Vault para gerar assinaturas sem expor a chave privada:

```go
// Exemplo com AWS KMS (pseudocódigo)
signature, err := awsKMS.Sign(
    keyID: "alias/my-eth-key",
    message: messageHash,
    signingAlgorithm: "ECDSA_SHA_256"
)
```

### 3. HSM (Hardware Security Module)

Para segurança máxima, utilize um HSM para armazenar chaves e realizar operações criptográficas:

- Ledger/Trezor para desenvolvimento
- AWS CloudHSM ou similares para produção

## 🌐 Tokens para Nation.fun

O IaC AI Agent requer autenticação com Nation.fun, que pode ser feita de duas formas:

1. **Token Pré-gerado**: Obtenha um token válido da Nation.fun e configure-o diretamente
2. **Modo Desenvolvimento**: O sistema gera um token temporário para desenvolvimento local

### Obtenção de Token Válido

Para obter um token válido para produção:

1. Acesse o dashboard da Nation.fun
2. Vá para "API Access" → "Generate Token"
3. Selecione a wallet que possui o NFT
4. Copie o token gerado e configure-o em `WALLET_TOKEN`

## 🔍 Verificação de Segurança

Para verificar se sua configuração está segura:

```bash
# Execute o verificador de segurança
go run cmd/security/check.go

# Saída esperada
✅ Não foram encontradas chaves privadas expostas
✅ Token de autenticação configurado corretamente
✅ Permissões de arquivo .env corretas
```

## 🆘 Resposta a Incidentes

Se você acidentalmente expôs uma chave privada:

1. **Revogue imediatamente** qualquer token associado
2. **Transfira** todos os ativos para uma nova wallet
3. **Altere todas as senhas** e tokens de acesso
4. **Documente o incidente** e medidas tomadas

---

**Lembre-se**: A segurança de suas chaves privadas é sua responsabilidade. Uma chave privada exposta pode resultar na perda total de ativos digitais.
