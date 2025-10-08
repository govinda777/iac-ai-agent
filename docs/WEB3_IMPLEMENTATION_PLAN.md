# Plano de Implementação Web3 - IaC AI Agent

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Componentes a Implementar](#componentes-a-implementar)
3. [Caminho Crítico](#caminho-crítico)
4. [Implementação Frontend](#implementação-frontend)
5. [Implementação Backend](#implementação-backend)
6. [Contratos Smart](#contratos-smart)
7. [Integração Privy](#integração-privy)
8. [Validação e Testes](#validação-e-testes)
9. [Timeline](#timeline)

## 🎯 Visão Geral

Este documento detalha o plano para implementar a integração Web3 no IaC AI Agent, permitindo:
- Login com Privy
- Compra de NFT de acesso
- Compra de tokens
- Consulta ao serviço
- Recebimento de sugestões

A implementação segue o modelo "Mínimo Produto Viável" (MVP) com foco na funcionalidade completa usando mock APIs quando necessário.

## 🧩 Componentes a Implementar

### Frontend
1. **Login/Authentication**
   - Componente Privy Login
   - Integração com Wallet Connect
   - Dashboard de usuário autenticado
   
2. **NFT Purchase**
   - Página de tiers de acesso
   - Fluxo de compra com ETH direto
   - Fluxo de onramp via cartão/pix
   
3. **Token Purchase**
   - Loja de tokens
   - Fluxo de compra com ETH
   - Fluxo de onramp
   
4. **Consulta**
   - Interface para envio de código Terraform
   - Painel de controle de análise
   - Histórico de consultas
   
5. **Sugestões**
   - Visualizador de resultados
   - Exportação de relatórios
   - Detalhamento de sugestões

### Backend
1. **Privy API Integration**
   - Verificação de tokens
   - Embedded wallets
   - Onramp de fiat para crypto
   
2. **Base Network Integration**
   - NFT contract interaction
   - Token contract interaction
   - Transaction management
   
3. **API de Análise**
   - Validação de acesso
   - Gerenciamento de tokens
   - Rate limiting por tier

## 🚀 Caminho Crítico

O caminho crítico da aplicação segue o fluxo:

1. **Autenticação**
   - Login via Privy (wallet ou email)
   - Verificação de identidade
   
2. **Acesso**
   - Compra/verificação de NFT
   - Validação de permissões
   
3. **Tokens**
   - Compra de tokens IACAI
   - Verificação de saldo
   
4. **Consulta**
   - Envio de código Terraform
   - Processamento da análise
   - Débito de tokens
   
5. **Sugestões**
   - Exibição de resultados
   - Interação com sugestões

## 💻 Implementação Frontend

### Tecnologias
- React/Next.js
- Privy SDK (@privy-io/react-auth)
- Wagmi/Viem para interação com blockchain
- TailwindCSS para UI

### Componentes Principais
1. **AuthProvider**: Wrapper do Privy SDK
   ```jsx
   <PrivyProvider
     appId={process.env.PRIVY_APP_ID}
     config={{
       loginMethods: ['wallet', 'email'],
       appearance: {
         theme: 'light',
         accentColor: '#3182ce'
       },
       embeddedWallets: {
         createOnLogin: 'all-users',
         noPromptOnSignature: true
       }
     }}
   >
     <App />
   </PrivyProvider>
   ```

2. **NFT Purchase Flow**
   ```jsx
   // Componentes a implementar:
   - <NFTTierSelector />
   - <PurchaseWithETH />
   - <PurchaseWithFiat />
   - <OnrampModal />
   - <TransactionStatus />
   ```

3. **Token Store**
   ```jsx
   // Componentes a implementar:
   - <TokenPackageList />
   - <PurchaseTokens />
   - <TokenBalance />
   - <TransactionHistory />
   ```

4. **Analysis Interface**
   ```jsx
   // Componentes a implementar:
   - <CodeSubmission />
   - <AnalysisOptions />
   - <AnalysisStatus />
   - <ResultsViewer />
   ```

## 🔧 Implementação Backend

### Módulos a Completar
1. **Privy Client** (`internal/platform/web3/privy_client.go`)
   - Implementar login/verificação
   - Gerenciar usuários
   - Integrar embedded wallets

2. **NFT Access** (`internal/platform/web3/nft_access.go`)
   - Finalizar interações com contrato
   - Implementar verificação de acesso
   - Gerenciar tier permissions

3. **Bot Token** (`internal/platform/web3/bot_token.go`)
   - Implementar compra de tokens
   - Gerenciar saldo de tokens
   - Controlar débitos por operação

4. **Onramp** (`internal/platform/web3/privy_onramp.go`)
   - Integrar com provedores (Moonpay/Transak)
   - Processar webhook callbacks
   - Gerenciar transações fiat → crypto

5. **API REST** (`api/rest/handlers.go`)
   - Endpoints de autenticação
   - Endpoints de compra
   - Endpoints de análise

## 📝 Contratos Smart

Para o MVP, usaremos versões simplificadas dos contratos, já preparados para o caminho crítico:

1. **NFT Access Contract** (ERC-721)
   ```solidity
   // Simplificado:
   - mint(address to, uint8 tierId)
   - balanceOf(address owner)
   - ownerOf(uint256 tokenId)
   - getTier(uint256 tokenId)
   - upgradeTier(uint256 tokenId, uint8 newTierId)
   ```

2. **Bot Token Contract** (ERC-20)
   ```solidity
   // Simplificado:
   - buyTokens(uint8 packageId)
   - balanceOf(address owner)
   - transfer(address to, uint256 amount)
   - approve(address spender, uint256 amount)
   - deductForService(address from, uint256 amount)
   ```

## 🔒 Integração Privy

Etapas para integrar o Privy:

1. **Setup da aplicação**
   - Configurar APP_ID e APP_SECRET
   - Adicionar domínios permitidos

2. **Frontend SDK**
   - Instalar: `npm install @privy-io/react-auth`
   - Implementar login flow

3. **Backend API**
   - Implementar verificação de tokens
   - Gerenciar embedded wallets

4. **Onramp**
   - Configurar provedores
   - Implementar webhook handlers

## ✅ Validação e Testes

### Testes BDD (Behavior-Driven Development)
Complementar os testes existentes para cobrir o caminho crítico:

1. **Authentication Tests** (`test/bdd/features/user_onboarding.feature`)
   - Testar login com diversas wallets
   - Testar persistência de sessão
   - Testar embedded wallets

2. **NFT Purchase Tests** (`test/bdd/features/nft_purchase.feature`)
   - Testar compra com ETH
   - Testar onramp e minting
   - Testar verificação de acesso

3. **Token Purchase Tests** (`test/bdd/features/token_purchase.feature`)
   - Testar compra de pacotes
   - Testar saldo e débitos
   - Testar histórico de transações

4. **Analysis Flow Tests** (`test/bdd/features/bot_analysis.feature`)
   - Testar submissão de código
   - Testar processamento de análise
   - Testar recebimento de resultados

### Mocks para MVP
Para acelerar o desenvolvimento, implementar mocks para:

1. **Blockchain Interactions**
   - Simular transações
   - Simular confirmações
   - Retornar balances fixos

2. **Privy API Responses**
   - Simular autenticação
   - Simular criação de wallet
   - Simular transações de onramp

## ⏱️ Timeline

| Fase | Duração | Descrição |
|------|---------|-----------|
| **Fase 1** | 1 dia | Setup inicial e estrutura base |
| **Fase 2** | 1-2 dias | Implementar autenticação Privy |
| **Fase 3** | 2-3 dias | Implementar NFT e Token interfaces |
| **Fase 4** | 1-2 dias | Implementar consulta e análises |
| **Fase 5** | 1 dia | Testes e validação |
| **Total** | 6-9 dias | MVP completo |

## 📊 Métricas de Sucesso

Para o MVP, consideraremos bem-sucedido quando:

1. Usuários conseguirem completar o caminho crítico:
   - Login → Comprar NFT → Comprar Tokens → Fazer Consulta → Ver Sugestões

2. Todos os testes BDD passarem com sucesso

3. A jornada do usuário estiver fluida e sem bloqueios críticos
