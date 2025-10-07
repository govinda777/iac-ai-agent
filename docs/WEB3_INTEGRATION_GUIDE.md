# Guia de Integração Web3 - IaC AI Agent

## 📋 Índice

1. [Visão Geral](#visão-geral)
2. [Arquitetura](#arquitetura)
3. [Privy.io Setup](#privyio-setup)
4. [Base Network Setup](#base-network-setup)
5. [Smart Contracts](#smart-contracts)
6. [Fluxos de Usuário](#fluxos-de-usuário)
7. [Testes BDD](#testes-bdd)
8. [Deploy](#deploy)

---

## 🎯 Visão Geral

O IaC AI Agent integra com **Privy.io** para autenticação Web3 e **Base Network** (L2 Ethereum) para pagamentos e acesso via NFT/Tokens.

### Componentes Principais

1. **Privy.io**: Autenticação Web3 simplificada
   - Login com wallet (MetaMask, Coinbase Wallet, etc)
   - Embedded wallets
   - Onramp para compra de crypto com fiat

2. **Base Network**: Blockchain L2 de baixo custo
   - NFTs de acesso (3 tiers)
   - Token ERC-20 (IACAI) para pagamentos
   - Transações baratas (~$0.01)

3. **Smart Contracts**:
   - NFT de Acesso (ERC-721)
   - Token IACAI (ERC-20)

---

## 🏗️ Arquitetura

```
┌─────────────────────────────────────────────────────────────┐
│                         Frontend                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Privy SDK    │  │ Wagmi/Viem   │  │ Base Wagmi   │      │
│  │ (Auth)       │  │ (Blockchain) │  │ (Base Chain) │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────┐
│                      Backend (Go)                            │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │ internal/platform/web3/                              │  │
│  │                                                       │  │
│  │  • privy_client.go      - Privy API integration     │  │
│  │  • base_client.go       - Base Network client       │  │
│  │  • nft_access.go        - NFT management            │  │
│  │  • bot_token.go         - Token management          │  │
│  │  • privy_onramp.go      - Fiat to crypto            │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                             ↓
┌─────────────────────────────────────────────────────────────┐
│                    Base Network (L2)                         │
│                                                              │
│  ┌──────────────────┐         ┌──────────────────┐         │
│  │ NFT Access       │         │ IACAI Token      │         │
│  │ (ERC-721)        │         │ (ERC-20)         │         │
│  │                  │         │                  │         │
│  │ - Basic Tier     │         │ - Buy packages   │         │
│  │ - Pro Tier       │         │ - Pay analyses   │         │
│  │ - Enterprise     │         │ - Transfer       │         │
│  └──────────────────┘         └──────────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔐 Privy.io Setup

### 1. Criar Conta no Privy

1. Acesse https://privy.io
2. Crie uma conta/app
3. Anote `APP_ID` e `APP_SECRET`

### 2. Configurar App

No dashboard do Privy:

```yaml
Settings:
  App Name: IaC AI Agent
  App Domain: https://iacai.yourdomain.com
  
Allowed Origins:
  - http://localhost:3000
  - https://iacai.yourdomain.com
  
Login Methods:
  ✓ Wallet (MetaMask, Coinbase Wallet, WalletConnect)
  ✓ Email
  ✓ Social (optional)
  
Embedded Wallets:
  ✓ Enable embedded wallets
  Default Chain: Base (8453)
  
Onramp:
  ✓ Enable fiat onramp
  Providers: MoonPay, Transak
  Supported Currencies: USD, EUR, BRL, GBP
  Supported Crypto: ETH, USDC
```

### 3. Variáveis de Ambiente

```bash
# .env
PRIVY_APP_ID=your-app-id
PRIVY_APP_SECRET=your-app-secret
```

### 4. Frontend Integration (React)

```tsx
// app/providers.tsx
import { PrivyProvider } from '@privy-io/react-auth';
import { base } from 'wagmi/chains';

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <PrivyProvider
      appId={process.env.NEXT_PUBLIC_PRIVY_APP_ID!}
      config={{
        loginMethods: ['wallet', 'email'],
        appearance: {
          theme: 'dark',
          accentColor: '#6366F1',
        },
        embeddedWallets: {
          createOnLogin: 'users-without-wallets',
        },
        defaultChain: base,
        supportedChains: [base],
      }}
    >
      {children}
    </PrivyProvider>
  );
}
```

---

## 🌐 Base Network Setup

### 1. Network Info

```yaml
Base Mainnet:
  Chain ID: 8453
  RPC URL: https://mainnet.base.org
  Currency: ETH
  Block Explorer: https://basescan.org
  
Base Goerli Testnet:
  Chain ID: 84531
  RPC URL: https://goerli.base.org
  Currency: ETH
  Block Explorer: https://goerli.basescan.org
```

### 2. Adicionar ao MetaMask

```javascript
// Add Base Network
await window.ethereum.request({
  method: 'wallet_addEthereumChain',
  params: [{
    chainId: '0x2105', // 8453 in hex
    chainName: 'Base',
    nativeCurrency: {
      name: 'Ether',
      symbol: 'ETH',
      decimals: 18
    },
    rpcUrls: ['https://mainnet.base.org'],
    blockExplorerUrls: ['https://basescan.org']
  }]
});
```

### 3. Obter ETH para Gas

**Opção 1: Bridge de Ethereum L1**
- https://bridge.base.org
- Bridge ETH de Ethereum Mainnet para Base

**Opção 2: Comprar via Onramp**
- Privy Onramp (MoonPay/Transak)
- Compre ETH direto na Base Network

---

## 📜 Smart Contracts

### NFT de Acesso (ERC-721)

```solidity
// contracts/IACaiAccessNFT.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract IACaiAccessNFT is ERC721, Ownable {
    uint256 private _nextTokenId;
    
    struct Tier {
        uint8 id;
        string name;
        uint256 price;
        uint256 maxSupply;
        uint256 currentSupply;
    }
    
    mapping(uint8 => Tier) public tiers;
    mapping(uint256 => uint8) public tokenTiers;
    
    event NFTMinted(address indexed to, uint256 tokenId, uint8 tier);
    event TierUpgraded(uint256 tokenId, uint8 fromTier, uint8 toTier);
    
    constructor() ERC721("IaC AI Access", "IACAI-ACCESS") Ownable(msg.sender) {
        // Define tiers
        tiers[1] = Tier(1, "Basic", 0.01 ether, 10000, 0);
        tiers[2] = Tier(2, "Pro", 0.05 ether, 5000, 0);
        tiers[3] = Tier(3, "Enterprise", 0.2 ether, 1000, 0);
    }
    
    function mint(uint8 tierId) external payable {
        Tier storage tier = tiers[tierId];
        require(tier.id > 0, "Invalid tier");
        require(msg.value >= tier.price, "Insufficient payment");
        require(tier.currentSupply < tier.maxSupply, "Tier sold out");
        
        uint256 tokenId = _nextTokenId++;
        tokenTiers[tokenId] = tierId;
        tier.currentSupply++;
        
        _safeMint(msg.sender, tokenId);
        emit NFTMinted(msg.sender, tokenId, tierId);
        
        // Refund excess
        if (msg.value > tier.price) {
            payable(msg.sender).transfer(msg.value - tier.price);
        }
    }
    
    function upgrade(uint256 tokenId, uint8 newTierId) external payable {
        require(ownerOf(tokenId) == msg.sender, "Not token owner");
        
        uint8 currentTierId = tokenTiers[tokenId];
        require(newTierId > currentTierId, "Can only upgrade");
        
        Tier storage newTier = tiers[newTierId];
        Tier storage currentTier = tiers[currentTierId];
        
        uint256 priceDiff = newTier.price - currentTier.price;
        require(msg.value >= priceDiff, "Insufficient payment");
        
        tokenTiers[tokenId] = newTierId;
        emit TierUpgraded(tokenId, currentTierId, newTierId);
        
        if (msg.value > priceDiff) {
            payable(msg.sender).transfer(msg.value - priceDiff);
        }
    }
    
    function getTier(uint256 tokenId) external view returns (uint8) {
        return tokenTiers[tokenId];
    }
    
    function withdraw() external onlyOwner {
        payable(owner()).transfer(address(this).balance);
    }
}
```

### Token IACAI (ERC-20)

```solidity
// contracts/IACaiToken.sol
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract IACaiToken is ERC20, Ownable {
    struct TokenPackage {
        uint8 id;
        uint256 tokenAmount;
        uint256 price;
        uint8 discountPercent;
    }
    
    mapping(uint8 => TokenPackage) public packages;
    
    event TokensPurchased(address indexed buyer, uint256 amount, uint256 price);
    
    constructor() ERC20("IaC AI Token", "IACAI") Ownable(msg.sender) {
        // Mint initial supply
        _mint(address(this), 1_000_000 * 10**decimals());
        
        // Define packages
        packages[1] = TokenPackage(1, 100 * 10**decimals(), 0.005 ether, 0);
        packages[2] = TokenPackage(2, 500 * 10**decimals(), 0.0225 ether, 10);
        packages[3] = TokenPackage(3, 1000 * 10**decimals(), 0.0425 ether, 15);
        packages[4] = TokenPackage(4, 5000 * 10**decimals(), 0.1875 ether, 25);
    }
    
    function buyTokens(uint8 packageId) external payable {
        TokenPackage storage package = packages[packageId];
        require(package.id > 0, "Invalid package");
        require(msg.value >= package.price, "Insufficient payment");
        
        _transfer(address(this), msg.sender, package.tokenAmount);
        emit TokensPurchased(msg.sender, package.tokenAmount, package.price);
        
        if (msg.value > package.price) {
            payable(msg.sender).transfer(msg.value - package.price);
        }
    }
    
    function withdraw() external onlyOwner {
        payable(owner()).transfer(address(this).balance);
    }
}
```

### Deploy Contracts

```bash
# Install dependencies
npm install --save-dev hardhat @openzeppelin/contracts

# Create hardhat config
npx hardhat init

# Configure for Base
# hardhat.config.ts
import { HardhatUserConfig } from "hardhat/config";

const config: HardhatUserConfig = {
  solidity: "0.8.20",
  networks: {
    base: {
      url: "https://mainnet.base.org",
      accounts: [process.env.DEPLOYER_PRIVATE_KEY!],
      chainId: 8453,
    },
    baseGoerli: {
      url: "https://goerli.base.org",
      accounts: [process.env.DEPLOYER_PRIVATE_KEY!],
      chainId: 84531,
    },
  },
};

export default config;

# Deploy
npx hardhat run scripts/deploy.ts --network base
```

---

## 👤 Fluxos de Usuário

### 1. Onboarding + Compra de NFT via Onramp

```
1. Usuário clica "Get Started"
2. Privy modal abre
3. Usuário escolhe "Continue with Email"
4. Embedded wallet é criada automaticamente
5. Usuário vê tiers de NFT disponíveis
6. Seleciona "Pro Access" ($125)
7. Clica "Buy with Card"
8. Privy Onramp abre (MoonPay)
9. Insere dados do cartão
10. Pagamento aprovado
11. ETH chega na wallet (~5 min)
12. Backend detecta ETH na wallet
13. Backend minta NFT automaticamente
14. Usuário recebe notificação
15. Acesso liberado!
```

### 2. Compra de Tokens

```
1. Usuário autenticado acessa "Buy Tokens"
2. Seleciona pacote (ex: Pro Pack - 1000 tokens)
3. Escolhe método: "ETH" ou "Card"
4. Se Card: Privy Onramp → ETH na wallet → Auto-compra
5. Se ETH: Aprova transação → Tokens transferidos
6. Saldo atualizado
```

### 3. Fazer Análise

```
1. Usuário cola código Terraform
2. Seleciona tipo de análise
3. Sistema verifica:
   - Tem NFT válido? ✓
   - Tier suficiente? ✓
   - Tokens suficientes? ✓
   - Rate limit OK? ✓
4. Análise executada
5. Tokens debitados
6. Resultado exibido
```

---

## 🧪 Testes BDD

Os testes BDD estão em `test/bdd/features/`:

1. **user_onboarding.feature**: Autenticação e onboarding
2. **nft_purchase.feature**: Compra de NFT (com ETH e com Onramp)
3. **token_purchase.feature**: Compra de tokens
4. **bot_analysis.feature**: Uso do bot com tokens

### Executar Testes

```bash
# Instalar Godog (BDD framework para Go)
go get github.com/cucumber/godog/cmd/godog

# Executar todos os testes
godog test/bdd/features/

# Executar teste específico
godog test/bdd/features/nft_purchase.feature

# Com tags
godog --tags=@smoke test/bdd/features/
```

---

## 🚀 Deploy

### 1. Backend

```bash
# Build
docker build -t iacai-agent .

# Run with env vars
docker run -p 8080:8080 \
  -e PRIVY_APP_ID=xxx \
  -e PRIVY_APP_SECRET=xxx \
  -e BASE_RPC_URL=https://mainnet.base.org \
  -e NFT_CONTRACT_ADDRESS=0x... \
  -e TOKEN_CONTRACT_ADDRESS=0x... \
  iacai-agent
```

### 2. Frontend

```bash
# Deploy to Vercel
vercel deploy --prod

# Environment variables
NEXT_PUBLIC_PRIVY_APP_ID=xxx
NEXT_PUBLIC_API_URL=https://api.iacai.com
```

### 3. Monitoring

```yaml
Health Checks:
  - /health (API health)
  - /health/privy (Privy connectivity)
  - /health/base (Base Network connectivity)
  
Metrics:
  - NFT mints per day
  - Token purchases per day
  - Active users by tier
  - Analysis requests per tier
  - Revenue (ETH/USD)
```

---

## 📞 Suporte

- **Documentação Privy**: https://docs.privy.io
- **Documentação Base**: https://docs.base.org
- **Issues**: https://github.com/your-org/iac-ai-agent/issues

---

**Status**: ✅ Pronto para implementação  
**Última atualização**: 2025-01-15
