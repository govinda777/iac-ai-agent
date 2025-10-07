# 🤖 Configuração do Agente Nation.fun

## 📋 Visão Geral

O IaC AI Agent utiliza o [Nation.fun](https://nation.fun) como provedor LLM exclusivo, com um agente específico configurado como padrão. Este documento explica como configurar e usar o agente Nation.fun.

## 🔗 Agente Padrão

Por padrão, o sistema está configurado para usar o seguinte agente Nation.fun:

```
Agent Address: 0x147e832418Cc06A501047019E956714271098b89
URL: https://nation.fun/agent/0x147e832418Cc06A501047019E956714271098b89?chat=true
```

Este agente é especializado em análise de Infrastructure as Code (IaC) e será usado automaticamente caso não seja possível criar um agente personalizado.

## ⚙️ Configuração

### Variáveis de Ambiente

```bash
# Configuração do agente Nation.fun
DEFAULT_AGENT_ADDRESS=0x147e832418Cc06A501047019E956714271098b89  # Opcional, usa o padrão se não especificado
```

### Configuração YAML

```yaml
web3:
  # Nation.fun Configuration
  default_agent_address: "0x147e832418Cc06A501047019E956714271098b89"  # Agente padrão Nation.fun
```

## 🔧 Como Funciona

1. O sistema tenta usar o agente especificado em `DEFAULT_AGENT_ADDRESS`
2. Se não for especificado, usa o agente padrão `0x147e832418Cc06A501047019E956714271098b89`
3. Todas as requisições de LLM são enviadas para este agente via API Nation.fun
4. O agente processa a análise de código Terraform e retorna resultados estruturados

## 🎯 Características do Agente

O agente Nation.fun padrão é especializado em:

- Análise de código Terraform
- Detecção de problemas de segurança
- Sugestões de otimização de custos
- Recomendações de boas práticas
- Insights arquiteturais

## 🔄 Alterando o Agente Padrão

Para usar um agente diferente:

1. **Via variável de ambiente**:
   ```bash
   export DEFAULT_AGENT_ADDRESS=0xSeuEnderecoDeAgente
   ```

2. **Via arquivo de configuração**:
   ```yaml
   web3:
     default_agent_address: "0xSeuEnderecoDeAgente"
   ```

## 🧪 Testando o Agente

Você pode testar o agente diretamente acessando:
[https://nation.fun/agent/0x147e832418Cc06A501047019E956714271098b89?chat=true](https://nation.fun/agent/0x147e832418Cc06A501047019E956714271098b89?chat=true)

## 📚 Exemplos de Uso

### Análise Básica

```go
// O agente padrão será usado automaticamente
req := &models.LLMRequest{
    Prompt: "Analise este código Terraform: resource \"aws_s3_bucket\" \"example\" { ... }",
}

resp, err := llmClient.Generate(req)
```

### Análise Estruturada

```go
// O agente padrão será usado para gerar resposta estruturada
analysis := &models.LLMStructuredResponse{}
req := &models.LLMRequest{
    Prompt: "Analise este código Terraform: resource \"aws_s3_bucket\" \"example\" { ... }",
    ResponseFormat: "json",
}

err := llmClient.GenerateStructured(req, analysis)
```

## 🔍 Validação

Durante o startup, o sistema valida:

1. Conexão com Nation.fun API
2. Validade do API Key
3. Posse do NFT Nation.fun
4. Validade do Wallet Token
5. Existência do agente especificado

Se qualquer validação falhar, a aplicação não inicia.

## ❓ FAQ

### Por que usar um agente específico?

O agente `0x147e832418Cc06A501047019E956714271098b89` foi especialmente treinado para análise de Infrastructure as Code, com conhecimento específico sobre Terraform, AWS, Azure, GCP e boas práticas de segurança.

### Posso criar meu próprio agente?

Sim, você pode criar seu próprio agente no Nation.fun e especificá-lo via `DEFAULT_AGENT_ADDRESS`. No entanto, recomendamos usar o agente padrão que já está otimizado para análise de IaC.

### Como o agente é diferente de um LLM genérico?

O agente Nation.fun foi treinado especificamente em:
- Código Terraform e HCL
- Padrões de segurança em cloud
- Otimização de custos em AWS, Azure e GCP
- Arquitetura de infraestrutura como código

---

**Última Atualização**: 2025-01-15  
**Versão**: 1.0.0  
**Status**: ✅ Implementado
