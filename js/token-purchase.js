// Token Purchase Flow for IaC AI Agent
document.addEventListener('DOMContentLoaded', function() {
    // Configuration
    const isMockMode = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
    let TOKEN_CONTRACT_ADDRESS;
    
    // Tenta obter o endereço do contrato de token do ambiente
    try {
        TOKEN_CONTRACT_ADDRESS = process.env.TOKEN_CONTRACT_ADDRESS || '';
    } catch (e) {
        // Fallback para desenvolvimento
        TOKEN_CONTRACT_ADDRESS = '0x987654321abcdef987654321abcdef98765432';
        console.warn('Usando endereço de contrato de token de desenvolvimento');
    }
    
    // Disponibilizar globalmente
    window.TOKEN_CONTRACT_ADDRESS = TOKEN_CONTRACT_ADDRESS;
    
    // Elements
    const tokenPurchaseButtons = document.querySelectorAll('.buy-tokens-btn');
    const packageButtons = document.querySelectorAll('.select-package-btn');
    const tokenBalanceElements = document.querySelectorAll('.token-balance');
    
    // State
    let selectedPackage = null;
    let userTokenBalance = 0;
    let isPurchaseInProgress = false;
    
    // Token packages configuration
    const packages = {
        'starter': {
            id: 1,
            name: 'Starter Pack',
            tokens: 100,
            priceETH: 0.005,
            priceUSD: 10,
            discount: 0
        },
        'power': {
            id: 2,
            name: 'Power Pack',
            tokens: 500,
            priceETH: 0.0225,
            priceUSD: 45,
            discount: 10
        },
        'pro': {
            id: 3,
            name: 'Pro Pack',
            tokens: 1000,
            priceETH: 0.0425,
            priceUSD: 85,
            discount: 15
        },
        'enterprise': {
            id: 4,
            name: 'Enterprise Pack',
            tokens: 5000,
            priceETH: 0.1875,
            priceUSD: 375,
            discount: 25
        }
    };
    
    // Initialize
    function init() {
        // Load token balance from localStorage for development
        if (isMockMode) {
            userTokenBalance = parseInt(localStorage.getItem('tokenBalance') || '0', 10);
            updateTokenBalanceDisplay();
        }
        
        // Setup package selection
        if (packageButtons.length) {
            packageButtons.forEach(button => {
                button.addEventListener('click', function(e) {
                    e.preventDefault();
                    const packageId = this.getAttribute('data-package');
                    selectPackage(packageId);
                });
            });
        }
        
        // Setup purchase buttons
        if (tokenPurchaseButtons.length) {
            tokenPurchaseButtons.forEach(button => {
                button.addEventListener('click', function(e) {
                    e.preventDefault();
                    
                    const purchaseType = this.getAttribute('data-purchase-type'); // 'eth' or 'card'
                    const packageId = this.getAttribute('data-package') || selectedPackage;
                    
                    if (!packageId) {
                        showMessage('Selecione um pacote primeiro', 'error');
                        return;
                    }
                    
                    if (!window.Web3Auth || !window.Web3Auth.isAuthenticated()) {
                        showLoginPrompt();
                        return;
                    }
                    
                    // Check if user has NFT access (required to buy tokens)
                    const userTier = window.IaCNFT ? window.IaCNFT.getCurrentTier() : null;
                    if (!userTier) {
                        showMessage('É necessário ter um NFT de acesso para comprar tokens', 'error');
                        showNFTPurchasePrompt();
                        return;
                    }
                    
                    if (purchaseType === 'eth') {
                        purchaseWithETH(packageId);
                    } else if (purchaseType === 'card') {
                        purchaseWithCard(packageId);
                    }
                });
            });
        }
        
        // Listen for authentication events
        window.addEventListener('user:authenticated', function(e) {
            updatePurchaseButtonsState();
            loadUserTokenBalance();
        });
        
        window.addEventListener('user:logout', function() {
            updatePurchaseButtonsState();
            userTokenBalance = 0;
            updateTokenBalanceDisplay();
        });
        
        // Listen for NFT purchase events
        window.addEventListener('nft:purchased', function(e) {
            updatePurchaseButtonsState();
        });
        
        // Initial state update
        updatePurchaseButtonsState();
        loadUserTokenBalance();
    }
    
    function selectPackage(packageId) {
        selectedPackage = packageId;
        
        // Update UI to highlight selected package
        document.querySelectorAll('.token-card').forEach(card => {
            card.classList.remove('selected');
        });
        
        const selectedCard = document.querySelector(`.token-card[data-package="${packageId}"]`);
        if (selectedCard) {
            selectedCard.classList.add('selected');
        }
        
        // Update purchase buttons
        tokenPurchaseButtons.forEach(button => {
            button.setAttribute('data-package', packageId);
            
            // Update price display if present
            const priceElement = button.querySelector('.price-amount');
            if (priceElement && packages[packageId]) {
                priceElement.textContent = `${packages[packageId].priceETH} ETH`;
            }
        });
    }
    
    function updatePurchaseButtonsState() {
        const isAuthenticated = window.IaCAuth && window.IaCAuth.isAuthenticated();
        const hasNFT = window.IaCNFT && window.IaCNFT.getCurrentTier();
        
        tokenPurchaseButtons.forEach(button => {
            if (!isAuthenticated) {
                button.classList.add('requires-login');
                button.querySelector('.button-text').textContent = 'Conectar Wallet';
                button.disabled = true;
            } else if (!hasNFT) {
                button.classList.add('requires-nft');
                button.querySelector('.button-text').textContent = 'Comprar NFT';
                button.disabled = true;
            } else {
                button.classList.remove('requires-login', 'requires-nft');
                button.querySelector('.button-text').textContent = button.getAttribute('data-original-text') || 'Comprar Tokens';
                button.disabled = isPurchaseInProgress;
            }
            
            if (isPurchaseInProgress) {
                button.classList.add('processing');
            } else {
                button.classList.remove('processing');
            }
        });
    }
    
    function showLoginPrompt() {
        if (window.Web3Auth) {
            window.Web3Auth.login();
        } else {
            showMessage('Sistema de autenticação não disponível', 'error');
        }
    }
    
    function showNFTPurchasePrompt() {
        // Redirect to NFT purchase page
        window.location.href = 'pricing.html#nft-access';
    }
    
    // Load user token balance
    async function loadUserTokenBalance() {
        if (!window.Web3Auth || !window.Web3Auth.isAuthenticated()) {
            userTokenBalance = 0;
            updateTokenBalanceDisplay();
            return;
        }
        
        if (isMockMode) {
            // Load mock balance from localStorage
            userTokenBalance = parseInt(localStorage.getItem('tokenBalance') || '0', 10);
        } else {
            // In production, call token contract
            try {
                const walletAddress = window.Web3Auth.getWalletAddress();
                if (!walletAddress) {
                    userTokenBalance = 0;
                } else {
                    // Simplified example - actual implementation would use ethers.js or web3.js
                    // const tokenContract = new Contract(TOKEN_CONTRACT_ADDRESS, tokenABI, provider);
                    // userTokenBalance = await tokenContract.balanceOf(walletAddress);
                    // For now, use mock
                    userTokenBalance = parseInt(localStorage.getItem('tokenBalance') || '0', 10);
                }
            } catch (error) {
                console.error('Error loading token balance:', error);
                userTokenBalance = 0;
            }
        }
        
        updateTokenBalanceDisplay();
    }
    
    // Purchase flow
    async function purchaseWithETH(packageId) {
        if (!packages[packageId]) {
            showMessage('Pacote inválido', 'error');
            return;
        }
        
        try {
            isPurchaseInProgress = true;
            updatePurchaseButtonsState();
            showMessage('Iniciando compra com ETH...', 'info');
            
            if (isMockMode) {
                // Mock purchase for development
                await simulatePurchase(packageId);
            } else {
                // Real purchase using Web3
                await performWeb3Purchase(packageId);
            }
        } catch (error) {
            console.error('Error during purchase:', error);
            showMessage(`Erro na compra: ${error.message}`, 'error');
        } finally {
            isPurchaseInProgress = false;
            updatePurchaseButtonsState();
        }
    }
    
    async function purchaseWithCard(packageId) {
        if (!packages[packageId]) {
            showMessage('Pacote inválido', 'error');
            return;
        }
        
        try {
            isPurchaseInProgress = true;
            updatePurchaseButtonsState();
            showMessage('Iniciando compra com cartão...', 'info');
            
            if (isMockMode) {
                // Mock onramp for development
                await simulateOnramp(packageId);
            } else {
                // Real onramp using Privy
                await startPrivyOnramp(packageId);
            }
        } catch (error) {
            console.error('Error during onramp:', error);
            showMessage(`Erro na compra: ${error.message}`, 'error');
        } finally {
            isPurchaseInProgress = false;
            updatePurchaseButtonsState();
        }
    }
    
    // Web3 interaction
    async function performWeb3Purchase(packageId) {
        if (!window.Web3Auth || !window.Web3Auth.isAuthenticated()) {
            throw new Error('Autenticação necessária');
        }
        
        const pkg = packages[packageId];
        const walletAddress = window.Web3Auth.getWalletAddress();
        if (!walletAddress) {
            throw new Error('Carteira não disponível');
        }
        
        // Preparar transação para o contrato de tokens
        showMessage(`Enviando transação para comprar ${pkg.tokens} tokens...`, 'info');
        
        try {
            // Configurar transação
            const transaction = {
                to: window.TOKEN_CONTRACT_ADDRESS,
                value: ethers.utils.parseEther(pkg.priceETH.toString()),
                data: encodeTokenPurchase(packageId)
            };
            
            // Executar transação
            const txResult = await window.Web3Auth.executeTransaction(transaction);
            if (!txResult) throw new Error('Falha na transação');
            
            showMessage('Transação enviada. Aguardando confirmação...', 'info');
            
            // Aguardar confirmação (simplificado)
            await new Promise(resolve => setTimeout(resolve, 2000));
            
            showMessage(`${pkg.tokens} IACAI tokens adquiridos com sucesso!`, 'success');
            
            // Update balance
            userTokenBalance += pkg.tokens;
            localStorage.setItem('tokenBalance', userTokenBalance.toString());
            updateTokenBalanceDisplay();
            
            // Dispatch token purchase event
            window.dispatchEvent(new CustomEvent('tokens:purchased', {
                detail: { amount: pkg.tokens, package: packageId }
            }));
            
            return txResult;
        } catch (error) {
            console.error('Error in token purchase transaction:', error);
            throw new Error(`Erro na transação: ${error.message}`);
        }
    }
    
    // Encode token purchase function call
    function encodeTokenPurchase(packageId) {
        // Em produção, usar ethers.js para codificar a chamada
        // const iface = new ethers.utils.Interface(['function buyTokens(uint8 packageId)']);
        // return iface.encodeFunctionData('buyTokens', [parseInt(packageId)]);
        
        // Valor simulado
        const pkgIdMap = { starter: 1, power: 2, pro: 3, enterprise: 4 };
        return `0x7d3d6522000000000000000000000000000000000000000000000000000000000000000${pkgIdMap[packageId] || 0}`;
    }
    
    async function startPrivyOnramp(packageId) {
        const pkg = packages[packageId];
        
        // In production, this would use Privy's onramp feature
        showMessage(`Iniciando onramp para comprar ${pkg.tokens} tokens...`, 'info');
        
        // Show onramp modal (simplified mock)
        const onrampModal = document.getElementById('onramp-modal');
        if (onrampModal) {
            document.getElementById('onramp-package-name').textContent = pkg.name;
            document.getElementById('onramp-token-amount').textContent = `${pkg.tokens} IACAI`;
            document.getElementById('onramp-package-price').textContent = `${pkg.priceETH} ETH (~$${pkg.priceUSD})`;
            onrampModal.style.display = 'flex';
            
            // Handle form submission
            document.getElementById('onramp-form').addEventListener('submit', function(e) {
                e.preventDefault();
                simulateOnramp(packageId);
                onrampModal.style.display = 'none';
            });
        } else {
            await simulateOnramp(packageId);
        }
    }
    
    // Mock implementations for development
    async function simulatePurchase(packageId) {
        const pkg = packages[packageId];
        
        showMessage(`Simulando compra de ${pkg.tokens} IACAI tokens...`, 'info');
        await new Promise(resolve => setTimeout(resolve, 2000));
        
        showMessage('Transação enviada. Aguardando confirmação...', 'info');
        await new Promise(resolve => setTimeout(resolve, 3000));
        
        showMessage(`${pkg.tokens} IACAI tokens adquiridos com sucesso!`, 'success');
        
        // Update balance
        userTokenBalance += pkg.tokens;
        localStorage.setItem('tokenBalance', userTokenBalance.toString());
        updateTokenBalanceDisplay();
        
        // Dispatch token purchase event
        window.dispatchEvent(new CustomEvent('tokens:purchased', {
            detail: { amount: pkg.tokens, package: packageId }
        }));
    }
    
    async function simulateOnramp(packageId) {
        const pkg = packages[packageId];
        
        showMessage(`Simulando onramp para compra de ${pkg.tokens} IACAI tokens...`, 'info');
        await new Promise(resolve => setTimeout(resolve, 2000));
        
        showMessage('Processando pagamento...', 'info');
        await new Promise(resolve => setTimeout(resolve, 3000));
        
        showMessage('ETH recebido na wallet. Comprando tokens...', 'info');
        await new Promise(resolve => setTimeout(resolve, 2000));
        
        showMessage(`${pkg.tokens} IACAI tokens adquiridos com sucesso!`, 'success');
        
        // Update balance
        userTokenBalance += pkg.tokens;
        localStorage.setItem('tokenBalance', userTokenBalance.toString());
        updateTokenBalanceDisplay();
        
        // Dispatch token purchase event
        window.dispatchEvent(new CustomEvent('tokens:purchased', {
            detail: { amount: pkg.tokens, package: packageId }
        }));
    }
    
    // Token spending
    function spendTokens(amount, purpose) {
        if (userTokenBalance < amount) {
            showMessage(`Saldo insuficiente. Necessário: ${amount} IACAI, Disponível: ${userTokenBalance} IACAI`, 'error');
            return false;
        }
        
        userTokenBalance -= amount;
        localStorage.setItem('tokenBalance', userTokenBalance.toString());
        updateTokenBalanceDisplay();
        
        showMessage(`${amount} IACAI tokens utilizados para: ${purpose}`, 'info');
        
        // Dispatch token spent event
        window.dispatchEvent(new CustomEvent('tokens:spent', {
            detail: { amount, purpose }
        }));
        
        return true;
    }
    
    // UI Updates
    function updateTokenBalanceDisplay() {
        tokenBalanceElements.forEach(el => {
            el.textContent = userTokenBalance;
        });
    }
    
    function showMessage(message, type = 'info') {
        console.log(`[${type.toUpperCase()}] ${message}`);
        
        // Check if message container exists, create if not
        let msgContainer = document.querySelector('.message-container');
        if (!msgContainer) {
            msgContainer = document.createElement('div');
            msgContainer.className = 'message-container';
            document.body.appendChild(msgContainer);
        }
        
        // Create message element
        const msgElement = document.createElement('div');
        msgElement.className = `message message-${type}`;
        msgElement.innerHTML = `
            <div class="message-content">${message}</div>
            <button class="message-close">&times;</button>
        `;
        
        // Add to container
        msgContainer.appendChild(msgElement);
        
        // Auto-remove after delay
        setTimeout(() => {
            msgElement.classList.add('fade-out');
            setTimeout(() => msgElement.remove(), 500);
        }, 5000);
        
        // Close button
        msgElement.querySelector('.message-close').addEventListener('click', () => {
            msgElement.classList.add('fade-out');
            setTimeout(() => msgElement.remove(), 500);
        });
    }
    
    // Export methods for other scripts
    window.IaCToken = {
        selectPackage,
        purchaseWithETH,
        purchaseWithCard,
        spendTokens,
        getPackageInfo: (packageId) => packages[packageId] || null,
        getTokenBalance: () => userTokenBalance,
        refreshBalance: loadUserTokenBalance
    };
    
    // Initialize
    init();
});
