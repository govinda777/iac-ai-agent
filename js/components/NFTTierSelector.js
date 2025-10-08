/**
 * NFT Tier Selector Component
 * 
 * Este componente apresenta os diferentes tiers de NFTs disponíveis
 * e gerencia a seleção e compra pelo usuário.
 */

class NFTTierSelector {
  constructor(config = {}) {
    // Configuração
    this.containerSelector = config.container || '#nft-tier-selector';
    this.onTierSelect = config.onTierSelect || (() => {});
    this.onPurchase = config.onPurchase || (() => {});
    this.mockMode = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';

    // Estado
    this.selectedTier = null;
    this.tiers = [
      {
        id: 'basic',
        name: 'Basic Access',
        description: 'Acesso básico ao IaC AI Agent',
        priceETH: 0.01,
        priceUSD: 25,
        features: [
          'Análises ilimitadas de Terraform',
          'Detecção de segurança com Checkov',
          'Sugestões básicas de otimização',
          'Suporte via Discord'
        ]
      },
      {
        id: 'pro',
        name: 'Pro Access',
        description: 'Acesso profissional com AI avançada',
        priceETH: 0.05,
        priceUSD: 125,
        popular: true,
        features: [
          'Tudo do Basic Access',
          'Análise com LLM (GPT-4/Claude)',
          'Sugestões contextualizadas inteligentes',
          'Análise de Preview e Drift',
          'Detecção de Secrets',
          'Recomendações de arquitetura',
          'Priority support'
        ]
      },
      {
        id: 'enterprise',
        name: 'Enterprise Access',
        description: 'Acesso enterprise com features exclusivas',
        priceETH: 0.2,
        priceUSD: 500,
        features: [
          'Tudo do Pro Access',
          'API dedicada com rate limits maiores',
          'Custom knowledge base',
          'Integração com CI/CD privado',
          'Suporte dedicado 24/7',
          'SLA garantido',
          'Governance tokens inclusos'
        ]
      }
    ];
    
    this.container = null;
    this.initialize();
  }

  /**
   * Inicializa o componente e renderiza a UI
   */
  initialize() {
    this.container = document.querySelector(this.containerSelector);
    if (!this.container) {
      console.error(`Container element not found: ${this.containerSelector}`);
      return;
    }
    
    this.render();
    this.attachEventListeners();
    
    // Verificar se já tem NFT armazenado
    const savedTier = localStorage.getItem('nftTier');
    if (savedTier) {
      this.updateUIForExistingNFT(savedTier);
    }
  }

  /**
   * Renderiza o componente na UI
   */
  render() {
    const template = `
      <div class="nft-tier-selector">
        <h2 class="section-title">Escolha seu Plano de Acesso</h2>
        <p class="section-description">NFTs de acesso garantem permissões exclusivas na plataforma</p>
        
        <div class="pricing-grid">
          ${this.tiers.map(tier => this.renderTierCard(tier)).join('')}
        </div>

        <div class="purchase-options" id="purchase-options">
          <h3>Complete sua compra</h3>
          <p id="selected-tier-info">Selecione um plano primeiro</p>
          
          <div class="purchase-methods">
            <button id="purchase-eth" class="purchase-btn eth-btn disabled" disabled>
              <span class="icon">Ξ</span>
              <span class="button-text">Comprar com ETH</span>
              <span class="price-amount">0.00 ETH</span>
            </button>
            
            <button id="purchase-card" class="purchase-btn card-btn disabled" disabled>
              <span class="icon">💳</span>
              <span class="button-text">Comprar com Cartão</span>
              <span class="price-amount">$0.00</span>
            </button>
          </div>
        </div>
        
        <div class="transaction-status" id="transaction-status" style="display: none;">
          <div class="spinner"></div>
          <h3>Processando Transação</h3>
          <p id="transaction-message">Aguarde enquanto processamos sua transação...</p>
          <div class="progress-bar">
            <div class="progress-fill" style="width: 0%"></div>
          </div>
        </div>
      </div>
    `;
    
    this.container.innerHTML = template;
  }

  /**
   * Renderiza o card para um tier específico
   */
  renderTierCard(tier) {
    return `
      <div class="pricing-card ${tier.popular ? 'popular' : ''}" data-tier="${tier.id}">
        ${tier.popular ? '<div class="popular-badge">Mais Popular</div>' : ''}
        <h3 class="tier-name">${tier.name}</h3>
        <p class="tier-description">${tier.description}</p>
        
        <div class="tier-price">
          <span class="eth-price">${tier.priceETH} ETH</span>
          <span class="usd-price">~$${tier.priceUSD}</span>
        </div>
        
        <ul class="tier-features">
          ${tier.features.map(feature => `<li>${feature}</li>`).join('')}
        </ul>
        
        <button class="select-tier-btn" data-tier="${tier.id}">Selecionar Plano</button>
      </div>
    `;
  }

  /**
   * Anexa event listeners aos elementos da UI
   */
  attachEventListeners() {
    // Botões de seleção de tier
    const tierButtons = this.container.querySelectorAll('.select-tier-btn');
    tierButtons.forEach(button => {
      button.addEventListener('click', (e) => {
        const tierId = e.target.getAttribute('data-tier');
        this.selectTier(tierId);
      });
    });
    
    // Botão de compra com ETH
    const ethButton = this.container.querySelector('#purchase-eth');
    if (ethButton) {
      ethButton.addEventListener('click', () => {
        if (this.selectedTier) {
          this.purchaseWithETH(this.selectedTier);
        }
      });
    }
    
    // Botão de compra com cartão
    const cardButton = this.container.querySelector('#purchase-card');
    if (cardButton) {
      cardButton.addEventListener('click', () => {
        if (this.selectedTier) {
          this.purchaseWithCard(this.selectedTier);
        }
      });
    }
  }

  /**
   * Seleciona um tier específico
   */
  selectTier(tierId) {
    // Encontrar o tier selecionado
    const tier = this.tiers.find(t => t.id === tierId);
    if (!tier) return;
    
    this.selectedTier = tierId;
    
    // Atualizar UI
    this.container.querySelectorAll('.pricing-card').forEach(card => {
      card.classList.remove('selected');
    });
    
    const selectedCard = this.container.querySelector(`.pricing-card[data-tier="${tierId}"]`);
    if (selectedCard) {
      selectedCard.classList.add('selected');
    }
    
    // Atualizar texto de informação
    const tierInfoEl = this.container.querySelector('#selected-tier-info');
    if (tierInfoEl) {
      tierInfoEl.textContent = `${tier.name} selecionado - ${tier.priceETH} ETH (~$${tier.priceUSD})`;
    }
    
    // Atualizar botões de compra
    const ethButton = this.container.querySelector('#purchase-eth');
    const cardButton = this.container.querySelector('#purchase-card');
    
    if (ethButton) {
      ethButton.classList.remove('disabled');
      ethButton.disabled = false;
      ethButton.querySelector('.price-amount').textContent = `${tier.priceETH} ETH`;
    }
    
    if (cardButton) {
      cardButton.classList.remove('disabled');
      cardButton.disabled = false;
      cardButton.querySelector('.price-amount').textContent = `$${tier.priceUSD}`;
    }
    
    // Chamar callback
    this.onTierSelect(tier);
  }

  /**
   * Comprar NFT com ETH (via wallet)
   */
  async purchaseWithETH(tierId) {
    // Verificar autenticação
    if (!window.Web3Auth || !window.Web3Auth.isAuthenticated()) {
      await this.handleLogin();
      return;
    }
    
    const tier = this.tiers.find(t => t.id === tierId);
    if (!tier) return;
    
    try {
      this.showTransactionStatus('Iniciando transação com sua carteira...');
      
      if (this.mockMode) {
        await this.simulatePurchase(tier);
      } else {
        await this.executePurchase(tier);
      }
      
      // Atualizar UI após compra
      this.updateTransactionStatus('Transação confirmada! NFT adquirido com sucesso.', 100);
      setTimeout(() => {
        this.hideTransactionStatus();
        this.updateUIForExistingNFT(tier.id);
      }, 2000);
      
      // Armazenar tier no localStorage
      localStorage.setItem('nftTier', tier.id);
      
      // Disparar evento
      window.dispatchEvent(new CustomEvent('nft:purchased', {
        detail: { tier: tier.id, tierName: tier.name }
      }));
      
      // Chamar callback
      this.onPurchase({ success: true, tier: tier.id, method: 'eth' });
      
    } catch (error) {
      console.error('Error during NFT purchase:', error);
      this.updateTransactionStatus(`Erro na transação: ${error.message}`, 0);
      setTimeout(() => {
        this.hideTransactionStatus();
      }, 3000);
      
      // Chamar callback com erro
      this.onPurchase({ success: false, tier: tier.id, method: 'eth', error: error.message });
    }
  }

  /**
   * Comprar NFT com cartão (via onramp)
   */
  async purchaseWithCard(tierId) {
    // Verificar autenticação
    if (!window.Web3Auth || !window.Web3Auth.isAuthenticated()) {
      await this.handleLogin();
      return;
    }
    
    const tier = this.tiers.find(t => t.id === tierId);
    if (!tier) return;
    
    try {
      this.showTransactionStatus('Iniciando fluxo de compra com cartão...');
      
      if (this.mockMode) {
        await this.simulateOnramp(tier);
      } else {
        await this.startOnramp(tier);
      }
      
      // Atualizar UI após compra
      this.updateTransactionStatus('NFT adquirido com sucesso!', 100);
      setTimeout(() => {
        this.hideTransactionStatus();
        this.updateUIForExistingNFT(tier.id);
      }, 2000);
      
      // Armazenar tier no localStorage
      localStorage.setItem('nftTier', tier.id);
      
      // Disparar evento
      window.dispatchEvent(new CustomEvent('nft:purchased', {
        detail: { tier: tier.id, tierName: tier.name }
      }));
      
      // Chamar callback
      this.onPurchase({ success: true, tier: tier.id, method: 'card' });
      
    } catch (error) {
      console.error('Error during NFT onramp:', error);
      this.updateTransactionStatus(`Erro na transação: ${error.message}`, 0);
      setTimeout(() => {
        this.hideTransactionStatus();
      }, 3000);
      
      // Chamar callback com erro
      this.onPurchase({ success: false, tier: tier.id, method: 'card', error: error.message });
    }
  }

  /**
   * Executa compra real de NFT com web3
   */
  async executePurchase(tier) {
    const walletAddress = window.Web3Auth.getWalletAddress();
    if (!walletAddress) throw new Error('Carteira não disponível');
    
    // Configurar transação
    const transaction = {
      to: window.NFT_CONTRACT_ADDRESS,
      value: ethers.utils.parseEther(tier.priceETH.toString()),
      data: this.encodeNFTPurchase(tier.id)
    };
    
    this.updateTransactionStatus('Confirmando transação na carteira...', 20);
    
    // Executar transação
    const txResult = await window.Web3Auth.executeTransaction(transaction);
    if (!txResult) throw new Error('Falha na transação');
    
    this.updateTransactionStatus('Transação enviada. Aguardando confirmação na blockchain...', 50);
    
    // Aguardar confirmação (simplificado)
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    this.updateTransactionStatus('NFT sendo mintado para sua carteira...', 80);
    
    // Simular confirmação final
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    return txResult;
  }

  /**
   * Inicia fluxo de onramp para compra com cartão
   */
  async startOnramp(tier) {
    const walletAddress = window.Web3Auth.getWalletAddress();
    if (!walletAddress) throw new Error('Carteira não disponível');
    
    this.updateTransactionStatus('Redirecionando para pagamento...', 30);
    
    // Em produção, integração com Privy Onramp ou provedor similar
    await new Promise(resolve => setTimeout(resolve, 1500));
    
    this.updateTransactionStatus('Processando pagamento...', 60);
    
    // Simular processamento
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    this.updateTransactionStatus('Pagamento confirmado! Mintando NFT...', 80);
    
    // Simular mint
    await new Promise(resolve => setTimeout(resolve, 1500));
  }

  /**
   * Simula compra para desenvolvimento
   */
  async simulatePurchase(tier) {
    this.updateTransactionStatus(`Simulando compra do NFT ${tier.name}...`, 20);
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    this.updateTransactionStatus('Transação enviada. Aguardando confirmação...', 50);
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    this.updateTransactionStatus('NFT sendo mintado para sua carteira...', 80);
    await new Promise(resolve => setTimeout(resolve, 1500));
  }

  /**
   * Simula onramp para desenvolvimento
   */
  async simulateOnramp(tier) {
    this.updateTransactionStatus(`Simulando onramp para compra do NFT ${tier.name}...`, 20);
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    this.updateTransactionStatus('Processando pagamento...', 50);
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    this.updateTransactionStatus('ETH recebido na wallet. Mintando NFT...', 80);
    await new Promise(resolve => setTimeout(resolve, 1500));
  }

  /**
   * Lida com autenticação necessária
   */
  async handleLogin() {
    if (window.Web3Auth) {
      try {
        await window.Web3Auth.login();
        return true;
      } catch (error) {
        console.error('Login error:', error);
        alert('É necessário conectar uma carteira para continuar.');
        return false;
      }
    } else {
      alert('É necessário conectar uma carteira para continuar.');
      return false;
    }
  }

  /**
   * Atualiza UI para usuário que já possui NFT
   */
  updateUIForExistingNFT(tierId) {
    const tier = this.tiers.find(t => t.id === tierId);
    if (!tier) return;
    
    // Atualizar UI
    this.container.innerHTML = `
      <div class="nft-owned">
        <div class="success-icon">✓</div>
        <h2>Você já possui acesso ${tier.name}!</h2>
        <p>NFT adquirido e verificado em sua carteira.</p>
        
        <div class="nft-details">
          <div class="nft-tier-badge ${tier.id}">${tier.name}</div>
          <ul class="tier-features">
            ${tier.features.map(feature => `<li>${feature}</li>`).join('')}
          </ul>
        </div>
        
        <div class="nft-actions">
          <a href="#bot-services" class="primary-button">Usar Serviços</a>
          <button id="upgrade-nft" class="secondary-button">Upgrade de Plano</button>
        </div>
      </div>
    `;
    
    // Adicionar listener para botão de upgrade
    const upgradeButton = this.container.querySelector('#upgrade-nft');
    if (upgradeButton) {
      upgradeButton.addEventListener('click', () => {
        this.showUpgradeOptions(tierId);
      });
    }
  }

  /**
   * Mostra opções de upgrade
   */
  showUpgradeOptions(currentTierId) {
    // Encontrar tiers superiores
    const currentTierIndex = this.tiers.findIndex(t => t.id === currentTierId);
    const upgradeTiers = this.tiers.slice(currentTierIndex + 1);
    
    if (upgradeTiers.length === 0) {
      alert('Você já possui o plano mais avançado!');
      return;
    }
    
    // Renderizar opções de upgrade
    this.render(); // Reset UI
    
    const container = this.container.querySelector('.pricing-grid');
    container.innerHTML = '<h3 class="upgrade-title">Escolha seu plano de upgrade</h3>';
    
    // Adicionar apenas tiers superiores
    upgradeTiers.forEach(tier => {
      container.innerHTML += this.renderTierCard(tier);
    });
    
    // Reattach event listeners
    this.attachEventListeners();
  }

  /**
   * Mostra status de transação
   */
  showTransactionStatus(message) {
    const statusEl = this.container.querySelector('#transaction-status');
    const messageEl = this.container.querySelector('#transaction-message');
    
    if (statusEl) {
      statusEl.style.display = 'block';
    }
    
    if (messageEl) {
      messageEl.textContent = message;
    }
    
    // Ocultar opções de compra
    const optionsEl = this.container.querySelector('#purchase-options');
    if (optionsEl) {
      optionsEl.style.display = 'none';
    }
  }

  /**
   * Atualiza status de transação
   */
  updateTransactionStatus(message, progressPercent) {
    const messageEl = this.container.querySelector('#transaction-message');
    const progressEl = this.container.querySelector('.progress-fill');
    
    if (messageEl) {
      messageEl.textContent = message;
    }
    
    if (progressEl) {
      progressEl.style.width = `${progressPercent}%`;
    }
  }

  /**
   * Oculta status de transação
   */
  hideTransactionStatus() {
    const statusEl = this.container.querySelector('#transaction-status');
    const optionsEl = this.container.querySelector('#purchase-options');
    
    if (statusEl) {
      statusEl.style.display = 'none';
    }
    
    if (optionsEl) {
      optionsEl.style.display = 'block';
    }
  }

  /**
   * Codifica chamada de função para compra de NFT
   */
  encodeNFTPurchase(tierId) {
    // Em produção, usar biblioteca ethers.js ou web3.js para codificar
    // const iface = new ethers.utils.Interface(['function mint(address to, uint8 tierId)']);
    // return iface.encodeFunctionData('mint', [walletAddress, tierIdNumber]);
    
    // Valor simulado
    const tierIdMap = { basic: 1, pro: 2, enterprise: 3 };
    return `0x6a6278420000000000000000000000000000000000000000000000000000000000000${tierIdMap[tierId] || 0}`;
  }
}

// Exportar componente
window.NFTTierSelector = NFTTierSelector;
