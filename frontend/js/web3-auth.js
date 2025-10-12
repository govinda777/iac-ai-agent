/**
 * Web3 Authentication Integration Module
 * 
 * Este módulo gerencia a autenticação web3 usando Privy.io e
 * fornece integração com a API backend para verificação de tokens.
 */

// Configuração
const PRIVY_APP_ID = 'cmgh6un8w007bl10ci0tgitwp';
const BACKEND_AUTH_URL = '/api/auth/web3/verify';
const isMockMode = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';

// Estado
let privyReady = false;
let privyClient = null;
let userAuthenticated = false;
let userProfile = null;
let walletAddress = null;
let embeddedWallets = [];

/**
 * Inicializa o cliente Privy e configura autenticação
 */
export async function initializeWeb3Auth() {
  // Se estiver em modo de desenvolvimento local, use autenticação simulada
  if (isMockMode) {
    console.log('Using mock Web3 authentication mode');
    setupMockAuth();
    return;
  }

  try {
    // Importar SDK do Privy (garantindo que só carrega uma vez)
    if (!window.privy) {
      await loadPrivySDK();
    }

    // Configurar Privy
    const privyConfig = {
      embeddedWallets: {
        createOnLogin: 'all-users',
        noPromptOnSignature: true,
      },
      appearance: {
        theme: 'light',
        accentColor: '#3182ce',
        logo: '/img/logo.svg'
      },
      loginMethods: ['wallet', 'email'],
      defaultChain: 'base',
    };

    // Inicializar Privy
    await window.privy.init({
      appId: PRIVY_APP_ID,
      config: privyConfig,
      onReady: handlePrivyReady,
      onChange: handlePrivyChange
    });

    console.log('Privy SDK initialized');
  } catch (error) {
    console.error('Error initializing Privy SDK:', error);
  }
}

/**
 * Carrega o SDK do Privy
 */
async function loadPrivySDK() {
  return new Promise((resolve, reject) => {
    const script = document.createElement('script');
    script.src = 'https://sdk.privy.io/v1.0/privy.js';
    script.async = true;
    script.onload = () => resolve(window.privy);
    script.onerror = reject;
    document.head.appendChild(script);
  });
}

/**
 * Função chamada quando o Privy estiver pronto
 */
async function handlePrivyReady(privyState) {
  privyReady = true;
  privyClient = window.privy;
  console.log('Privy ready:', privyState);

  // Se o usuário já estiver autenticado
  if (privyState.authenticated) {
    await handleAuthenticated(privyState);
  } else {
    updateUIAfterLogout();
  }

  // Adicionar listeners a botões de login
  setupLoginButtons();
}

/**
 * Função chamada quando o estado do Privy mudar
 */
async function handlePrivyChange(privyState) {
  console.log('Privy state changed:', privyState);

  if (privyState.authenticated && !userAuthenticated) {
    await handleAuthenticated(privyState);
  } else if (!privyState.authenticated && userAuthenticated) {
    updateUIAfterLogout();
  }
}

/**
 * Função para lidar com autenticação bem-sucedida
 */
async function handleAuthenticated(privyState) {
  userAuthenticated = true;
  userProfile = privyState.user || {};

  try {
    // Obter carteiras conectadas
    const wallets = await privyClient.getWallets();
    
    // Usar a primeira carteira como principal (pode ser embedded ou externa)
    if (wallets.length > 0) {
      walletAddress = wallets[0].address;
      
      // Verificar se é uma embedded wallet
      if (wallets[0].walletClientType === 'privy') {
        embeddedWallets = [wallets[0]];
      }
    }
    
    // Se não encontrou wallet nas conectadas, tentar obter da sessão
    if (!walletAddress && privyState.user && privyState.user.wallet) {
      walletAddress = privyState.user.wallet.address;
    }

    // Verificar o token com o backend
    const token = await privyClient.getAccessToken();
    await verifyWithBackend(token);

    // Atualizar UI
    updateUIAfterLogin();
  } catch (error) {
    console.error('Error handling authentication:', error);
    // Em caso de erro de verificação, fazer logout
    privyClient.logout();
  }
}

/**
 * Verifica o token com o backend
 */
async function verifyWithBackend(token) {
  if (isMockMode) {
    console.log('Mock mode: Skipping backend verification');
    return true;
  }

  try {
    const response = await fetch(BACKEND_AUTH_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ token }),
    });

    if (!response.ok) {
      throw new Error(`Backend verification failed: ${response.status}`);
    }

    const data = await response.json();
    console.log('Backend verification successful:', data);
    return data;
  } catch (error) {
    console.error('Error verifying with backend:', error);
    throw error;
  }
}

/**
 * Configura botões de login
 */
function setupLoginButtons() {
  const loginButtons = document.querySelectorAll('.connect-wallet-btn');
  
  loginButtons.forEach(button => {
    button.addEventListener('click', async (e) => {
      e.preventDefault();
      await login();
    });
  });

  const logoutButton = document.getElementById('logout-btn');
  if (logoutButton) {
    logoutButton.addEventListener('click', async (e) => {
      e.preventDefault();
      await logout();
    });
  }
}

/**
 * Inicia o fluxo de login
 */
export async function login() {
  if (isMockMode) {
    mockLogin();
    return;
  }

  if (!privyReady) {
    console.error('Privy not ready');
    return;
  }

  try {
    await privyClient.login();
  } catch (error) {
    console.error('Error during login:', error);
  }
}

/**
 * Faz logout
 */
export async function logout() {
  if (isMockMode) {
    mockLogout();
    return;
  }

  if (!privyReady) {
    console.error('Privy not ready');
    return;
  }

  try {
    await privyClient.logout();
  } catch (error) {
    console.error('Error during logout:', error);
  }
}

/**
 * Cria uma embedded wallet para o usuário atual
 */
export async function createEmbeddedWallet() {
  if (!userAuthenticated || !privyReady) {
    console.error('User not authenticated or Privy not ready');
    return null;
  }

  try {
    const wallet = await privyClient.createWallet({
      walletClientType: 'privy',
    });
    
    embeddedWallets = [...embeddedWallets, wallet];
    console.log('Embedded wallet created:', wallet);
    return wallet;
  } catch (error) {
    console.error('Error creating embedded wallet:', error);
    return null;
  }
}

/**
 * Executa uma transação usando a carteira
 */
export async function executeTransaction(tx) {
  if (!userAuthenticated || !walletAddress) {
    console.error('User not authenticated or no wallet available');
    return null;
  }

  try {
    // Privy suporta transações diretamente 
    const result = await privyClient.sendTransaction(tx);
    return result;
  } catch (error) {
    console.error('Error executing transaction:', error);
    return null;
  }
}

// UI Updates
function updateUIAfterLogin() {
  // Atualizar o estado do DOM
  document.body.classList.add('user-authenticated');
  
  // Atualizar elementos de exibição da carteira
  const walletDisplays = document.querySelectorAll('.wallet-address');
  const shortWallet = walletAddress ? `${walletAddress.substring(0, 6)}...${walletAddress.substring(walletAddress.length - 4)}` : '';
  
  walletDisplays.forEach(display => {
    display.textContent = shortWallet;
  });
  
  // Mostrar seção de perfil do usuário
  const userProfileSection = document.getElementById('user-profile');
  if (userProfileSection) {
    userProfileSection.style.display = 'block';
  }
  
  // Ocultar botões de conexão de carteira
  const connectWalletButtons = document.querySelectorAll('.connect-wallet-btn');
  connectWalletButtons.forEach(button => {
    button.style.display = 'none';
  });
  
  // Mostrar botão de logout
  const logoutButton = document.getElementById('logout-btn');
  if (logoutButton) {
    logoutButton.style.display = 'block';
  }
  
  // Atualizar elementos que exigem login
  const loginRequiredElements = document.querySelectorAll('.login-required');
  loginRequiredElements.forEach(el => {
    el.classList.remove('disabled');
  });
  
  // Disparar evento personalizado para outros componentes
  window.dispatchEvent(new CustomEvent('user:authenticated', {
    detail: { 
      wallet: walletAddress,
      profile: userProfile
    }
  }));
}

function updateUIAfterLogout() {
  userAuthenticated = false;
  walletAddress = null;
  userProfile = null;
  embeddedWallets = [];
  
  document.body.classList.remove('user-authenticated');
  
  // Ocultar seção de perfil do usuário
  const userProfileSection = document.getElementById('user-profile');
  if (userProfileSection) {
    userProfileSection.style.display = 'none';
  }
  
  // Mostrar botões de conexão de carteira
  const connectWalletButtons = document.querySelectorAll('.connect-wallet-btn');
  connectWalletButtons.forEach(button => {
    button.style.display = 'block';
  });
  
  // Ocultar botão de logout
  const logoutButton = document.getElementById('logout-btn');
  if (logoutButton) {
    logoutButton.style.display = 'none';
  }
  
  // Atualizar elementos que exigem login
  const loginRequiredElements = document.querySelectorAll('.login-required');
  loginRequiredElements.forEach(el => {
    el.classList.add('disabled');
  });
  
  // Disparar evento personalizado para outros componentes
  window.dispatchEvent(new CustomEvent('user:logout'));
}

// Implementação mock para desenvolvimento local
function setupMockAuth() {
  const mockWalletAddress = '0x742d35Cc6634C0532925a3b844Bc9e7595f0bDce';
  
  const connectWalletButtons = document.querySelectorAll('.connect-wallet-btn');
  connectWalletButtons.forEach(button => {
    button.addEventListener('click', function(e) {
      e.preventDefault();
      mockLogin();
    });
  });

  const logoutButton = document.getElementById('logout-btn');
  if (logoutButton) {
    logoutButton.addEventListener('click', function(e) {
      e.preventDefault();
      mockLogout();
    });
  }
  
  if (localStorage.getItem('mockAuthenticated')) {
    mockLogin();
  }
}

function mockLogin() {
  const mockWalletAddress = '0x742d35Cc6634C0532925a3b844Bc9e7595f0bDce';
  userAuthenticated = true;
  walletAddress = mockWalletAddress;
  userProfile = {
    id: 'mock-user-id',
    email: 'user@example.com'
  };
  localStorage.setItem('mockAuthenticated', 'true');
  localStorage.setItem('mockWallet', mockWalletAddress);
  updateUIAfterLogin();
}

function mockLogout() {
  userAuthenticated = false;
  walletAddress = null;
  userProfile = null;
  localStorage.removeItem('mockAuthenticated');
  localStorage.removeItem('mockWallet');
  updateUIAfterLogout();
}

// API pública
export const Web3Auth = {
  initialize: initializeWeb3Auth,
  login,
  logout,
  isAuthenticated: () => userAuthenticated,
  getWalletAddress: () => walletAddress,
  getUserProfile: () => userProfile,
  createEmbeddedWallet,
  executeTransaction
};

// Auto-inicializar quando o DOM estiver pronto
document.addEventListener('DOMContentLoaded', initializeWeb3Auth);
