# Fluxos Operacionais e de Diagnóstico

Este documento detalha os fluxos internos para inicialização, verificação de saúde (health check) e diagnóstico de agentes.

---

## 🚀 Start App

Fluxo executado na inicialização da aplicação para garantir que o ambiente está pronto.

1.  **Subir a API**: Iniciar o servidor web e expor os endpoints.
2.  **Procurar Agentes na Nation**: Verificar se a wallet padrão possui agentes registrados no sistema "Nation".
3.  **Criar Novo Agente**: Se nenhum agente for encontrado, criar um novo com o template padrão.
4.  **Teste de LLM**: Realizar uma chamada de teste para o Large Language Model (LLM) para validar a conectividade e a API Key.

---

## ❤️ Health Check

Verificação de saúde geral da aplicação, focada na integração com a wallet padrão.

1.  **Dados da Wallet Padrão**:
    -   **Lista de Transações**: Obter o histórico de transações recentes.
    -   **Lista de NFTs**: Listar todos os NFTs pertencentes à wallet.
    -   **Saldo**: Verificar o saldo de criptomoedas (e.g., ETH, MATIC).
    -   **Assinatura**: Validar se a aplicação consegue assinar uma transação ou mensagem em nome da wallet.

---

## 🩺 Health Agent

Verificação de saúde específica para a funcionalidade dos Agentes.

1.  **Verificar NFT da Nation**: A wallet padrão possui o NFT específico que concede acesso ao "Nation"?
2.  **Verificar Agentes na Base**: A wallet possui Agentes registrados no nosso sistema (Base)?
3.  **Teste do Primeiro Agente**: Conseguimos executar uma operação de teste com o primeiro agente encontrado?
4.  **Teste de Criação de Agente**: Conseguimos criar um novo agente usando o template pré-definido?