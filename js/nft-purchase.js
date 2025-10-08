// NFT Purchase Flow for IaC AI Agent
document.addEventListener('DOMContentLoaded', function() {
    // Configuration
    const isMockMode = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
    let NFT_CONTRACT_ADDRESS;
    
    // Tenta obter o endereço do contrato NFT do ambiente
    try {
        NFT_CONTRACT_ADDRESS = process.env.NFT_CONTRACT_ADDRESS || '';
    } catch (e) {
        // Fallback para desenvolvimento
        NFT_CONTRACT_ADDRESS = '0x123456789abcdef123456789abcdef12345678';
        console.warn('Usando endereço de contrato NFT de desenvolvimento');
    }
    
    // Disponibilizar globalmente para outros componentes
    window.NFT_CONTRACT_ADDRESS = NFT_CONTRACT_ADDRESS;
    
    // Elements
    const nftPurchaseButtons = document.querySelectorAll('.buy-nft-btn');
    const tierButtons = document.querySelectorAll('.select-tier-btn');
    
    // State
    let selectedTier = null;
    let currentETHPrice = 0;
    let isPurchaseInProgress = false;
    
    // Tiers configuration
    const tiers = {
        'basic': {
            id: 1,
            name: 'Basic Access',
            priceETH: 0.01,
            priceUSD: 25
        },
        'pro': {
            id: 2,
            name: 'Pro Access',
            priceETH: 0.05,
            priceUSD: 125
        },
        'enterprise': {
            id: 3,
            name: 'Enterprise Access',
            priceETH: 0.2,
            priceUSD: 500
        }
    };
    
    // Initialize
    function init() {
        // Inicializa o NFTTierSelector para containers específicos
        initNFTTierSelectors();
        
        // Setup tier selection para botões legados
        if (tierButtons.length) {
            tierButtons.forEach(button => {
                button.addEventListener('click', function(e) {
                    e.preventDefault();
                    const tierId = this.getAttribute('data-tier');
                    selectTier(tierId);
                });
            });
        }
        
        // Setup purchase buttons para botões legados
        if (nftPurchaseButtons.length) {
            nftPurchaseButtons.forEach(button => {
                button.addEventListener('click', function(e) {
                    e.preventDefault();
                    
                    const purchaseType = this.getAttribute('data-purchase-type'); // 'eth' or 'card'
                    const tierId = this.getAttribute('data-tier') || selectedTier;
                    
                    if (!tierId) {
                        showMessage('Selecione um tier primeiro', 'error');
                        return;
                    }
                    
                    if (!window.Web3Auth || !window.Web3Auth.isAuthenticated()) {
                        showLoginPrompt();
                        return;
                    }
                    
                    if (purchaseType === 'eth') {
                        purchaseWithETH(tierId);
                    } else if (purchaseType === 'card') {
                        purchaseWithCard(tierId);
                    }
                });
            });
        }
        
        // Listen for authentication events
        window.addEventListener('user:authenticated', function(e) {
            updatePurchaseButtonsState();
        });
        
        window.addEventListener('user:logout', function() {
            updatePurchaseButtonsState();
        });
        
        // Initial state update
        updatePurchaseButtonsState();
    }
    
    // Inicializa os seletores de NFT
    function initNFTTierSelectors() {
        const containers = document.querySelectorAll('.nft-tier-container');
        
        containers.forEach(container => {
            // Cria um novo seletor para cada container
            new NFTTierSelector({
                container: `#${container.id}`,
                onTierSelect: (tier) => {
                    console.log(`Tier selected: ${tier.id}`, tier);
                    // Trigger legacy events para compatibilidade
                    window.IaCNFT.selectTier(tier.id);
                },
                onPurchase: (result) => {
                    console.log(`Purchase result:`, result);
                    if (result.success) {
                        showMessage(`NFT ${result.tier} adquirido com sucesso!`, 'success');
                    } else {
                        showMessage(`Erro na compra: ${result.error}`, 'error');
                    }
                }
            });
        });
    }
    
    function selectTier(tierId) {
        selectedTier = tierId;
        
        // Update UI to highlight selected tier
        document.querySelectorAll('.pricing-card').forEach(card => {
            card.classList.remove('selected');
        });
        
        const selectedCard = document.querySelector(`.pricing-card[data-tier="${tierId}"]`);
        if (selectedCard) {
            selectedCard.classList.add('selected');
        }
        
        // Update purchase buttons
        nftPurchaseButtons.forEach(button => {
            button.setAttribute('data-tier', tierId);
            
            // Update price display if present
            const priceElement = button.querySelector('.price-amount');
            if (priceElement && tiers[tierId]) {
                priceElement.textContent = `${tiers[tierId].priceETH} ETH`;
            }
        });
    }
    
    function updatePurchaseButtonsState() {
        const isAuthenticated = window.Web3Auth && window.Web3Auth.isAuthenticated();
        
        nftPurchaseButtons.forEach(button => {
            if (!isAuthenticated) {
                button.classList.add('requires-login');
                button.querySelector('.button-text').textContent = 'Conectar Wallet';
            } else {
                button.classList.remove('requires-login');
                button.querySelector('.button-text').textContent = button.getAttribute('data-original-text') || 'Comprar NFT';
            }
            
            if (isPurchaseInProgress) {
                button.disabled = true;
                button.classList.add('processing');
            } else {
                button.disabled = false;
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
    
    // Purchase flow
    async function purchaseWithETH(tierId) {
        if (!tiers[tierId]) {
            showMessage('Tier inválido', 'error');
            return;
        }
        
        try {
            isPurchaseInProgress = true;
            updatePurchaseButtonsState();
            showMessage('Iniciando compra com ETH...', 'info');
            
            if (isMockMode) {
                // Mock purchase for development
                await simulatePurchase(tierId);
            } else {
                // Real purchase using Web3
                await performWeb3Purchase(tierId);
            }
        } catch (error) {
            console.error('Error during purchase:', error);
            showMessage(`Erro na compra: ${error.message}`, 'error');
        } finally {
            isPurchaseInProgress = false;
            updatePurchaseButtonsState();
        }
    }
    
    async function purchaseWithCard(tierId) {
        if (!tiers[tierId]) {
            showMessage('Tier inválido', 'error');
            return;
        }
        
        try {
            isPurchaseInProgress = true;
            updatePurchaseButtonsState();
            showMessage('Iniciando compra com cartão...', 'info');
            
            if (isMockMode) {
                // Mock onramp for development
                await simulateOnramp(tierId);
            } else {
                // Real onramp using Privy
                await startPrivyOnramp(tierId);
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
    async function performWeb3Purchase(tierId) {
        if (typeof window.ethereum === 'undefined') {
            throw new Error('MetaMask não encontrado');
        }
        
        const tier = tiers[tierId];
        const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
        const userAddress = accounts[0];
        
        // Call smart contract
        // This is a simplified example - actual implementation would need ethers.js or web3.js
        showMessage(`Enviando transação para comprar NFT ${tier.name}...`, 'info');
        
        // Simulate transaction processing
        await new Promise(resolve => setTimeout(resolve, 2000));
        
        // Handle success (in real implementation, wait for transaction confirmation)
        showMessage(`NFT ${tier.name} adquirido com sucesso!`, 'success');
        
        // Update local state and UI
        localStorage.setItem('nftTier', tierId);
        updateNFTStatus(tierId);
    }
    
    async function startPrivyOnramp(tierId) {
        const tier = tiers[tierId];
        
        // In production, this would use Privy's onramp feature
        showMessage(`Iniciando onramp para comprar ${tier.name}...`, 'info');
        
        // Show onramp modal (simplified mock)
        const onrampModal = document.getElementById('onramp-modal');
        if (onrampModal) {
            document.getElementById('onramp-tier-name').textContent = tier.name;
            document.getElementById('onramp-tier-price').textContent = `${tier.priceETH} ETH (~$${tier.priceUSD})`;
            onrampModal.style.display = 'flex';
            
            // Handle form submission
            document.getElementById('onramp-form').addEventListener('submit', function(e) {
                e.preventDefault();
                simulateOnramp(tierId);
                onrampModal.style.display = 'none';
            });
        } else {
            await simulateOnramp(tierId);
        }
    }
    
    // Mock implementations for development
    async function simulatePurchase(tierId) {
        const tier = tiers[tierId];
        
        showMessage(`Simulando compra do NFT ${tier.name}...`, 'info');
        await new Promise(resolve => setTimeout(resolve, 2000));
        
        showMessage('Transação enviada. Aguardando confirmação...', 'info');
        await new Promise(resolve => setTimeout(resolve, 3000));
        
        showMessage(`NFT ${tier.name} adquirido com sucesso!`, 'success');
        
        // Update local state and UI
        localStorage.setItem('nftTier', tierId);
        updateNFTStatus(tierId);
    }
    
    async function simulateOnramp(tierId) {
        const tier = tiers[tierId];
        
        showMessage(`Simulando onramp para compra do NFT ${tier.name}...`, 'info');
        await new Promise(resolve => setTimeout(resolve, 2000));
        
        showMessage('Processando pagamento...', 'info');
        await new Promise(resolve => setTimeout(resolve, 3000));
        
        showMessage('ETH recebido na wallet. Mintando NFT...', 'info');
        await new Promise(resolve => setTimeout(resolve, 2000));
        
        showMessage(`NFT ${tier.name} adquirido com sucesso!`, 'success');
        
        // Update local state and UI
        localStorage.setItem('nftTier', tierId);
        updateNFTStatus(tierId);
    }
    
    // UI Updates
    function updateNFTStatus(tierId) {
        // Update any NFT status indicators on the page
        const tierName = tiers[tierId]?.name || 'Unknown';
        
        const statusElements = document.querySelectorAll('.nft-status');
        statusElements.forEach(el => {
            el.textContent = tierName;
            el.classList.remove('no-access', 'basic-access', 'pro-access', 'enterprise-access');
            el.classList.add(`${tierId}-access`);
        });
        
        // Update access-controlled features
        document.body.classList.remove('has-basic-access', 'has-pro-access', 'has-enterprise-access');
        document.body.classList.add(`has-${tierId}-access`);
        
        // Dispatch custom event for other components
        window.dispatchEvent(new CustomEvent('nft:purchased', {
            detail: { tier: tierId, tierName }
        }));
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
    window.IaCNFT = {
        selectTier,
        purchaseWithETH,
        purchaseWithCard,
        getTierInfo: (tierId) => tiers[tierId] || null,
        getCurrentTier: () => localStorage.getItem('nftTier') || null,
        getContractAddress: () => NFT_CONTRACT_ADDRESS
    };
    
    // Initialize
    init();
});
