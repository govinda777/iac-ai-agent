// Privy Authentication Integration for IaC AI Agent
document.addEventListener('DOMContentLoaded', function() {
    // Configuration
    const PRIVY_APP_ID = 'cmgh6un8w007bl10ci0tgitwp'; // ID do aplicativo Privy.io
    const isMockMode = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
    
    // Elements
    const authContainer = document.getElementById('auth-container');
    const userProfileSection = document.getElementById('user-profile');
    const connectWalletButtons = document.querySelectorAll('.connect-wallet-btn');
    const logoutButton = document.getElementById('logout-btn');
    const mockWalletAddress = '0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb';
    
    // State
    let privyReady = false;
    let userWallet = null;
    let userAuthenticated = false;
    
    // Privy Configuration
    const privyConfig = {
        embeddedWallets: {
            createOnLogin: 'all-users',
            noPromptOnSignature: true
        },
        appearance: {
            theme: 'light',
            accentColor: '#3182ce',
            logo: '/img/logo.svg'
        },
        loginMethods: ['wallet', 'email']
    };
    
    // Mock functionality for development
    function setupMockAuth() {
        console.log('Using mock authentication mode');
        
        connectWalletButtons.forEach(button => {
            button.addEventListener('click', function(e) {
                e.preventDefault();
                mockLogin();
            });
        });
        
        if (localStorage.getItem('mockAuthenticated')) {
            mockLogin();
        }
    }
    
    function mockLogin() {
        userAuthenticated = true;
        userWallet = mockWalletAddress;
        localStorage.setItem('mockAuthenticated', 'true');
        localStorage.setItem('mockWallet', mockWalletAddress);
        updateUIAfterLogin();
    }
    
    function mockLogout() {
        userAuthenticated = false;
        userWallet = null;
        localStorage.removeItem('mockAuthenticated');
        localStorage.removeItem('mockWallet');
        updateUIAfterLogout();
    }
    
    // Privy integration
    function initializePrivy() {
        if (typeof privy === 'undefined') {
            console.error('Privy SDK not loaded');
            return;
        }
        
        privy.init({
            appId: PRIVY_APP_ID,
            config: privyConfig,
            onReady: handlePrivyReady,
            onChange: handlePrivyChange
        });
    }
    
    function handlePrivyReady(privyState) {
        privyReady = true;
        console.log('Privy ready', privyState);
        
        if (privyState.authenticated) {
            handleAuthenticated(privyState);
        } else {
            updateUIAfterLogout();
        }
        
        // Setup login buttons
        connectWalletButtons.forEach(button => {
            button.addEventListener('click', function(e) {
                e.preventDefault();
                privy.login();
            });
        });
        
        // Setup logout
        if (logoutButton) {
            logoutButton.addEventListener('click', function(e) {
                e.preventDefault();
                privy.logout();
            });
        }
    }
    
    function handlePrivyChange(privyState) {
        console.log('Privy state changed', privyState);
        
        if (privyState.authenticated && !userAuthenticated) {
            handleAuthenticated(privyState);
        } else if (!privyState.authenticated && userAuthenticated) {
            updateUIAfterLogout();
        }
    }
    
    async function handleAuthenticated(privyState) {
        userAuthenticated = true;
        
        // Get wallet address if available
        const wallets = await privy.getWallets();
        if (wallets.length > 0) {
            userWallet = wallets[0].address;
        } else if (privyState.user.wallet) {
            userWallet = privyState.user.wallet.address;
        }
        
        updateUIAfterLogin();
    }
    
    // UI Updates
    function updateUIAfterLogin() {
        // Update UI elements based on authentication status
        document.body.classList.add('user-authenticated');
        
        // Update any wallet display elements
        const walletDisplays = document.querySelectorAll('.wallet-address');
        const shortWallet = userWallet ? `${userWallet.substring(0, 6)}...${userWallet.substring(userWallet.length - 4)}` : '';
        
        walletDisplays.forEach(display => {
            display.textContent = shortWallet;
        });
        
        // Show user profile section if it exists
        if (userProfileSection) {
            userProfileSection.style.display = 'block';
        }
        
        // Hide connect wallet buttons
        connectWalletButtons.forEach(button => {
            button.style.display = 'none';
        });
        
        // Show logout button
        if (logoutButton) {
            logoutButton.style.display = 'block';
        }
        
        // Update login-required elements
        const loginRequiredElements = document.querySelectorAll('.login-required');
        loginRequiredElements.forEach(el => {
            el.classList.remove('disabled');
        });
        
        // Dispatch custom event for other components
        window.dispatchEvent(new CustomEvent('user:authenticated', {
            detail: { wallet: userWallet }
        }));
    }
    
    function updateUIAfterLogout() {
        userAuthenticated = false;
        userWallet = null;
        
        document.body.classList.remove('user-authenticated');
        
        // Hide user profile section
        if (userProfileSection) {
            userProfileSection.style.display = 'none';
        }
        
        // Show connect wallet buttons
        connectWalletButtons.forEach(button => {
            button.style.display = 'block';
        });
        
        // Hide logout button
        if (logoutButton) {
            logoutButton.style.display = 'none';
        }
        
        // Update login-required elements
        const loginRequiredElements = document.querySelectorAll('.login-required');
        loginRequiredElements.forEach(el => {
            el.classList.add('disabled');
        });
        
        // Dispatch custom event for other components
        window.dispatchEvent(new CustomEvent('user:logout'));
    }
    
    // Initialize based on environment
    if (isMockMode) {
        setupMockAuth();
    } else {
        // Load Privy SDK
        const privyScript = document.createElement('script');
        privyScript.src = 'https://sdk.privy.io/v1.0/privy.js';
        privyScript.async = true;
        privyScript.onload = initializePrivy;
        document.head.appendChild(privyScript);
    }
    
    // Export methods for other scripts
    window.IaCAuth = {
        isAuthenticated: () => userAuthenticated,
        getWalletAddress: () => userWallet,
        login: () => isMockMode ? mockLogin() : (privyReady ? privy.login() : console.error('Privy not ready')),
        logout: () => isMockMode ? mockLogout() : (privyReady ? privy.logout() : console.error('Privy not ready'))
    };
});
