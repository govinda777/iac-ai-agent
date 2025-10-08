/**
 * Token Package Selector Component
 * 
 * Este componente apresenta os diferentes pacotes de tokens disponíveis
 * e gerencia a seleção e compra pelo usuário.
 */

class TokenPackageSelector {
  constructor(config = {}) {
    // Configuração
    this.containerSelector = config.container || '#token-package-selector';
    this.onPackageSelect = config.onPackageSelect || (() => {});
    this.onPurchase = config.onPurchase || (() => {});
    this.mockMode = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';

    // Estado
    this.selectedPackage = null;
    this.packages = [
      {
        id: 'starter',
        name: 'Starter Pack',
        description: 'Pacote inicial para começar',
        tokens: 100,
        priceETH: 0.005,
        priceUSD: 10,
        discount: 0
      },
      {
        id: 'power',
        name: 'Power Pack',
        description: 'Pacote popular com 10% de desconto',
        tokens: 500,
        priceETH: 0.0225,
        priceUSD: 45,
        discount: 10,
        popular: true
      },
      {
        id: 'pro',
        name: 'Pro Pack',
        description: 'Pacote profissional com 15% de desconto',
        tokens: 1000,
        priceETH: 0.0425,
        priceUSD: 85,
        discount: 15
      },
      {
        id: 'enterprise',
        name: 'Enterprise Pack',
        description: 'Pacote enterprise com 25% de desconto',
        tokens: 5000,
        priceETH: 0.1875,
        priceUSD: 375,
        discount: 25
      }
    ];
    
    this.userTokenBalance = 0;
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
    
    // Carregar saldo de tokens
    this.loadTokenBalance();
    
    // Renderizar componente
    this.render();
    this.attachEventListeners();
    
    // Escutar eventos de autenticação
    window.addEventListener('user:authenticated', () => {
      this.loadTokenBalance();
    });
    
    window.addEventListener('user:logout', () => {
      this.userTokenBalance = 0;
      this.updateBalanceDisplay();
    });
    
    // Escutar eventos de compra de tokens
    window.addEventListener('tokens:purchased', (event) => {
      this.loadTokenBalance();
    });
  }

  /**
   * Renderiza o componente na UI
   */
  render() {
    const template = `
      <div class="token-package-selector">
        <h2 class="section-title">Comprar Tokens IACAI</h2>
        <div class="token-balance-display">
          <span class="balance-label">Seu saldo:</span>
          <span class="balance-amount"><span id="user-token-balance">0</span> IACAI</span>
        </div>
        
        <div class="package-grid">
          ${this.packages.map(pkg => this.renderPackageCard(pkg)).join('')}
        </div>

        <div class="purchase-options" id="purchase-options">
          <h3>Complete sua compra</h3>
          <p id="selected-package-info">Selecione um pacote primeiro</p>
          
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
        
        <div class="token-usage-info">
          <h3>Como usar seus tokens</h3>
          <div class="token-costs">
            <div class="cost-item">
              <span class="operation">Análise Terraform</span>
              <span class="cost">1 Token</span>
            </div>
            <div class="cost-item">
              <span class="operation">Scan Checkov</span>
              <span class="cost">2 Tokens</span>
            </div>
            <div class="cost-item">
              <span class="operation">Análise LLM</span>
              <span class="cost">5 Tokens</span>
            </div>
            <div class="cost-item">
              <span class="operation">Análise de Preview</span>
              <span class="cost">3 Tokens</span>
            </div>
            <div class="cost-item">
              <span class="operation">Auditoria de Segurança</span>
              <span class="cost">10 Tokens</span>
            </div>
            <div class="cost-item">
              <span class="operation">Otimização de Custo</span>
              <span class="cost">5 Tokens</span>
            </div>
          </div>
        </div>
      </div>
    `;
    
    this.container.innerHTML = template;
    this.updateBalanceDisplay();
  }

  /**
   * Renderiza o card para um pacote específico
   */
  renderPackageCard(pkg) {
    return `
      <div class="package-card ${pkg.popular ? 'popular' : ''}" data-package="${pkg.id}">
        ${pkg.popular ? '<div class="popular-badge">Mais Popular</div>' : ''}
        <h3 class="package-name">${pkg.name}</h3>
        <p class="package-description">${pkg.description}</p>
        
        <div class="token-amount">
          <span class="token-count">${pkg.tokens}</span>
          <span class="token-label">IACAI Tokens</span>
        </div>
        
        <div class="package-price">
          <span class="eth-price">${pkg.priceETH} ETH</span>
          <span class="usd-price">~$${pkg.priceUSD}</span>
          ${pkg.discount > 0 ? `<span class="discount-badge">-${pkg.discount}%</span>` : ''}
        </div>
        
        <button class="select-package-btn" data-package="${pkg.id}">Selecionar Pacote</button>
      </div>
    `;
  }

  /**
   * Anexa event listeners aos elementos da UI
   */
  attachEventListeners() {
    // Botões de seleção de pacote
    const packageButtons = this.container.querySelectorAll('.select-package-btn');
    packageButtons.forEach(button => {
      button.addEventListener('click', (e) => {
        const packageId = e.target.getAttribute('data-package');
        this.selectPackage(packageId);
      });
    });
    
    // Botão de compra com ETH
    const ethButton = this.container.querySelector('#purchase-eth');
    if (ethButton) {
      ethButton.addEventListener('click', () => {
        if (this.selectedPackage) {
          this.purchaseWithETH(this.selectedPackage);
        }
      });
    }
    
    // Botão de compra com cartão
    const cardButton = this.container.querySelector('#purchase-card');
    if (cardButton) {
      cardButton.addEventListener('click', () => {
        if (this.selectedPackage) {
          this.purchaseWithCard(this.selectedPackage);
        }
      });
    }
  }

  /**
   * Seleciona um pacote específico
   */
  selectPackage(packageId) {
    // Encontrar o pacote selecionado
    const pkg = this.packages.find(p => p.id === packageId);
    if (!pkg) return;
    
    this.selectedPackage = packageId;
    
    // Atualizar UI
    this.container.querySelectorAll('.package-card').forEach(card => {
      card.classList.remove('selected');
    });
    
    const selectedCard = this.container.querySelector(`.package-card[data-package="${packageId}"]`);
    if (selectedCard) {
      selectedCard.classList.add('selected');
    }
    
    // Atualizar texto de informação
    const packageInfoEl = this.container.querySelector('#selected-package-info');
    if (packageInfoEl) {
      packageInfoEl.textContent = `${pkg.name} selecionado - ${pkg.tokens} tokens por ${pkg.priceETH} ETH (~$${pkg.priceUSD})`;
    }
    
    // Atualizar botões de compra
    const ethButton = this.container.querySelector('#purchase-eth');
    const cardButton = this.container.querySelector('#purchase-card');
    
    if (ethButton) {
      ethButton.classList.remove('disabled');
      ethButton.disabled = false;
      ethButton.querySelector('.price-amount').textContent = `${pkg.priceETH} ETH`;
    }
    
    if (cardButton) {
      cardButton.classList.remove('disabled');
      cardButton.disabled = false;
      cardButton.querySelector('.price-amount').textContent = `$${pkg.priceUSD}`;
    }
    
    // Chamar callback
    this.onPackageSelect(pkg);
  }

  /**
   * Carrega o saldo de tokens do usuário
   */
  async loadTokenBalance() {
    try {
      // Verificar autenticação
      if (!window.Web3Auth || !window.Web3Auth.isAuthenticated()) {
        this.userTokenBalance = 0;
        this.updateBalanceDisplay();
        return;
      }
      
      if (this.mockMode) {
        // Em desenvolvimento, usar localStorage
        this.userTokenBalance = parseInt(localStorage.getItem('tokenBalance') || '0', 10);
      } else {
        // Em produção, chamar a API para obter saldo
        const walletAddress = window.Web3Auth.getWalletAddress();
        if (!walletAddress) {
          this.userTokenBalance = 0;
        } else {
          // Aqui chamaríamos a API ou contrato
          // Simular com localStorage por enquanto
          this.userTokenBalance = parseInt(localStorage.getItem('tokenBalance') || '0', 10);
        }
      }
      
      this.updateBalanceDisplay();
    } catch (error) {
      console.error('Error loading token balance:', error);
    }
  }

  /**
   * Atualiza o display de saldo
   */
  updateBalanceDisplay() {
    const balanceElement = this.container.querySelector('#user-token-balance');
    if (balanceElement) {
      balanceElement.textContent = this.userTokenBalance;
    }
  }

  /**
   * Comprar tokens com ETH
   */
  async purchaseWithETH(packageId) {
    // Verificar autenticação
    if (!window.Web3Auth || !window.Web3Auth.isAuthenticated()) {
      await this.handleLogin();
      return;
    }
    
    // Verificar acesso via NFT
    const userTier = window.IaCNFT ? window.IaCNFT.getCurrentTier() : null;
    if (!userTier) {
      this.showMessage('É necessário ter um NFT de acesso para comprar tokens');
      this.redirectToNFTPurchase();
      return;
    }
    
    const pkg = this.packages.find(p => p.id === packageId);
    if (!pkg) return;
    
    try {
      this.showTransactionStatus('Iniciando transação com sua carteira...');
      
      if (this.mockMode) {
        await this.simulatePurchase(pkg);
      } else {
        await this.executePurchase(pkg);
      }
      
      // Atualizar UI após compra
      this.updateTransactionStatus('Transação confirmada! Tokens adquiridos com sucesso.', 100);
      setTimeout(() => {
        this.hideTransactionStatus();
      }, 2000);
      
      // Atualizar saldo
      this.userTokenBalance += pkg.tokens;
      localStorage.setItem('tokenBalance', this.userTokenBalance.toString());
      this.updateBalanceDisplay();
      
      // Disparar evento
      window.dispatchEvent(new CustomEvent('tokens:purchased', {
        detail: { amount: pkg.tokens, package: pkg.id }
      }));
      
      // Chamar callback
      this.onPurchase({ success: true, package: pkg.id, method: 'eth', amount: pkg.tokens });
      
    } catch (error) {
      console.error('Error during token purchase:', error);
      this.updateTransactionStatus(`Erro na transação: ${error.message}`, 0);
      setTimeout(() => {
        this.hideTransactionStatus();
      }, 3000);
      
      // Chamar callback com erro
      this.onPurchase({ success: false, package: pkg.id, method: 'eth', error: error.message });
    }
  }

  /**
   * Comprar tokens com cartão (via onramp)
   */
  async purchaseWithCard(packageId) {
    // Verificar autenticação
    if (!window.Web3Auth || !window.Web3Auth.isAuthenticated()) {
      await this.handleLogin();
      return;
    }
    
    // Verificar acesso via NFT
    const userTier = window.IaCNFT ? window.IaCNFT.getCurrentTier() : null;
    if (!userTier) {
      this.showMessage('É necessário ter um NFT de acesso para comprar tokens');
      this.redirectToNFTPurchase();
      return;
    }
    
    const pkg = this.packages.find(p => p.id === packageId);
    if (!pkg) return;
    
    try {
      this.showTransactionStatus('Iniciando fluxo de compra com cartão...');
      
      if (this.mockMode) {
        await this.simulateOnramp(pkg);
      } else {
        await this.startOnramp(pkg);
      }
      
      // Atualizar UI após compra
      this.updateTransactionStatus('Tokens adquiridos com sucesso!', 100);
      setTimeout(() => {
        this.hideTransactionStatus();
      }, 2000);
      
      // Atualizar saldo
      this.userTokenBalance += pkg.tokens;
      localStorage.setItem('tokenBalance', this.userTokenBalance.toString());
      this.updateBalanceDisplay();
      
      // Disparar evento
      window.dispatchEvent(new CustomEvent('tokens:purchased', {
        detail: { amount: pkg.tokens, package: pkg.id }
      }));
      
      // Chamar callback
      this.onPurchase({ success: true, package: pkg.id, method: 'card', amount: pkg.tokens });
      
    } catch (error) {
      console.error('Error during token onramp:', error);
      this.updateTransactionStatus(`Erro na transação: ${error.message}`, 0);
      setTimeout(() => {
        this.hideTransactionStatus();
      }, 3000);
      
      // Chamar callback com erro
      this.onPurchase({ success: false, package: pkg.id, method: 'card', error: error.message });
    }
  }

  /**
   * Executa compra real de tokens com web3
   */
  async executePurchase(pkg) {
    const walletAddress = window.Web3Auth.getWalletAddress();
    if (!walletAddress) throw new Error('Carteira não disponível');
    
    // Configurar transação
    const transaction = {
      to: window.TOKEN_CONTRACT_ADDRESS,
      value: ethers.utils.parseEther(pkg.priceETH.toString()),
      data: this.encodeTokenPurchase(pkg.id)
    };
    
    this.updateTransactionStatus('Confirmando transação na carteira...', 20);
    
    // Executar transação
    const txResult = await window.Web3Auth.executeTransaction(transaction);
    if (!txResult) throw new Error('Falha na transação');
    
    this.updateTransactionStatus('Transação enviada. Aguardando confirmação na blockchain...', 50);
    
    // Aguardar confirmação (simplificado)
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    this.updateTransactionStatus('Tokens sendo transferidos para sua carteira...', 80);
    
    // Simular confirmação final
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    return txResult;
  }

  /**
   * Inicia fluxo de onramp para compra com cartão
   */
  async startOnramp(pkg) {
    const walletAddress = window.Web3Auth.getWalletAddress();
    if (!walletAddress) throw new Error('Carteira não disponível');
    
    this.updateTransactionStatus('Redirecionando para pagamento...', 30);
    
    // Em produção, integração com Privy Onramp ou provedor similar
    await new Promise(resolve => setTimeout(resolve, 1500));
    
    this.updateTransactionStatus('Processando pagamento...', 60);
    
    // Simular processamento
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    this.updateTransactionStatus('Pagamento confirmado! Tokens sendo transferidos...', 80);
    
    // Simular transferência
    await new Promise(resolve => setTimeout(resolve, 1500));
  }

  /**
   * Simula compra para desenvolvimento
   */
  async simulatePurchase(pkg) {
    this.updateTransactionStatus(`Simulando compra de ${pkg.tokens} tokens...`, 20);
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    this.updateTransactionStatus('Transação enviada. Aguardando confirmação...', 50);
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    this.updateTransactionStatus('Tokens sendo transferidos para sua carteira...', 80);
    await new Promise(resolve => setTimeout(resolve, 1500));
  }

  /**
   * Simula onramp para desenvolvimento
   */
  async simulateOnramp(pkg) {
    this.updateTransactionStatus(`Simulando onramp para compra de ${pkg.tokens} tokens...`, 20);
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    this.updateTransactionStatus('Processando pagamento...', 50);
    await new Promise(resolve => setTimeout(resolve, 2000));
    
    this.updateTransactionStatus('ETH recebido na wallet. Comprando tokens...', 80);
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
        this.showMessage('É necessário conectar uma carteira para continuar.');
        return false;
      }
    } else {
      this.showMessage('É necessário conectar uma carteira para continuar.');
      return false;
    }
  }

  /**
   * Redireciona para página de compra de NFT
   */
  redirectToNFTPurchase() {
    // Em produção, usar pathname do router
    setTimeout(() => {
      window.location.href = 'pricing.html#nft-access';
    }, 2000);
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
   * Mostra mensagem ao usuário
   */
  showMessage(message, type = 'error') {
    alert(message);
    console.log(`[${type.toUpperCase()}] ${message}`);
  }

  /**
   * Codifica chamada de função para compra de token
   */
  encodeTokenPurchase(packageId) {
    // Em produção, usar biblioteca ethers.js ou web3.js para codificar
    // const iface = new ethers.utils.Interface(['function buyTokens(uint8 packageId)']);
    // return iface.encodeFunctionData('buyTokens', [pkgIdMap[packageId]]);
    
    // Valor simulado
    const pkgIdMap = { starter: 1, power: 2, pro: 3, enterprise: 4 };
    return `0x7d3d6522000000000000000000000000000000000000000000000000000000000000000${pkgIdMap[packageId] || 0}`;
  }
}

// Exportar componente
window.TokenPackageSelector = TokenPackageSelector;
