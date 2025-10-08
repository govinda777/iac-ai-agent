/**
 * Web3 Service Integration
 * 
 * Este módulo fornece um serviço para integração entre o frontend e backend web3.
 */

class Web3Service {
  constructor() {
    this.baseUrl = '/api/auth/web3';
    this.mockMode = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
    this.userToken = null;
    this.userInfo = null;
  }

  /**
   * Verifica token de autenticação com o backend
   * @param {string} token - O token de autenticação do Privy
   * @returns {Promise<Object>} Informações do usuário autenticado
   */
  async verifyToken(token) {
    if (this.mockMode) {
      console.log('[MOCK] Verificando token web3');
      // Simula delay de rede
      await new Promise(resolve => setTimeout(resolve, 500));
      
      // Mock da resposta do backend
      const mockResponse = {
        user_id: 'mock-user-123',
        wallet_address: window.Web3Auth ? window.Web3Auth.getWalletAddress() : '0x742d35Cc6634C0532925a3b844Bc9e7595f0bDce',
        has_nft_access: true,
        nft_tier: 2,
        token_balance: '1000 IACAI',
        allowed_operations: [
          'terraform_analysis',
          'checkov_scan',
          'llm_analysis',
          'preview_analysis',
          'security_audit',
        ]
      };
      
      this.userToken = token;
      this.userInfo = mockResponse;
      return mockResponse;
    }
    
    try {
      const response = await fetch(`${this.baseUrl}/verify`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ token }),
      });
      
      if (!response.ok) {
        throw new Error(`Erro na verificação: ${response.status} ${response.statusText}`);
      }
      
      const data = await response.json();
      this.userToken = token;
      this.userInfo = data;
      return data;
    } catch (error) {
      console.error('Error verifying token with backend:', error);
      throw error;
    }
  }

  /**
   * Verifica se um usuário tem acesso a uma operação específica
   * @param {string} walletAddress - Endereço da carteira
   * @param {string} operation - Nome da operação
   * @returns {Promise<Object>} Resultado da verificação
   */
  async checkAccess(walletAddress, operation) {
    if (this.mockMode) {
      console.log(`[MOCK] Verificando acesso: ${operation}`);
      await new Promise(resolve => setTimeout(resolve, 300));
      
      // Determina acesso com base no tier armazenado localmente
      const tier = localStorage.getItem('nftTier');
      let allowed = false;
      
      // Operações básicas que todos podem fazer
      const basicOps = ['view_docs', 'basic_analysis'];
      
      if (basicOps.includes(operation)) {
        allowed = true;
      } else if (tier === 'basic') {
        allowed = ['terraform_analysis', 'checkov_scan'].includes(operation);
      } else if (tier === 'pro') {
        allowed = ['terraform_analysis', 'checkov_scan', 'llm_analysis', 
                   'preview_analysis', 'security_audit'].includes(operation);
      } else if (tier === 'enterprise') {
        allowed = true; // Acesso a tudo
      }
      
      return {
        allowed,
        within_rate_limit: true,
        message: allowed ? undefined : 'Operação não permitida para seu nível de acesso'
      };
    }
    
    try {
      const response = await fetch(`${this.baseUrl}/check-access`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
          wallet_address: walletAddress,
          operation 
        }),
      });
      
      if (!response.ok) {
        throw new Error(`Erro na verificação de acesso: ${response.status} ${response.statusText}`);
      }
      
      return await response.json();
    } catch (error) {
      console.error('Error checking access:', error);
      throw error;
    }
  }

  /**
   * Obtém o custo em tokens de uma operação
   * @param {string} operation - Nome da operação
   * @returns {Promise<Object>} Informações de custo
   */
  async getTokenCost(operation) {
    if (this.mockMode) {
      console.log(`[MOCK] Obtendo custo para operação: ${operation}`);
      await new Promise(resolve => setTimeout(resolve, 200));
      
      // Custos simulados por operação
      const costs = {
        'terraform_analysis': '1 IACAI',
        'checkov_scan': '2 IACAI',
        'llm_analysis': '5 IACAI',
        'preview_analysis': '3 IACAI',
        'security_audit': '10 IACAI',
        'cost_optimization': '5 IACAI',
        'full_review': '15 IACAI'
      };
      
      return {
        operation,
        token_cost: costs[operation] || '1 IACAI'
      };
    }
    
    try {
      const response = await fetch(`${this.baseUrl}/token-cost`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ operation }),
      });
      
      if (!response.ok) {
        throw new Error(`Erro ao obter custo: ${response.status} ${response.statusText}`);
      }
      
      return await response.json();
    } catch (error) {
      console.error('Error getting token cost:', error);
      throw error;
    }
  }

  /**
   * Verifica e retorna informações de usuário autenticado
   * @returns {Promise<Object>} Informações do usuário ou null se não autenticado
   */
  async getUserInfo() {
    // Se já temos informações de usuário, retorná-las
    if (this.userInfo) {
      return this.userInfo;
    }
    
    // Se temos um token salvo, verificá-lo
    const savedToken = localStorage.getItem('authToken');
    if (savedToken) {
      try {
        const userInfo = await this.verifyToken(savedToken);
        return userInfo;
      } catch (error) {
        console.error('Error retrieving user info:', error);
        localStorage.removeItem('authToken');
        return null;
      }
    }
    
    return null;
  }

  /**
   * Registra o gasto de tokens para uma operação
   * @param {string} walletAddress - Endereço da carteira
   * @param {string} operation - Nome da operação
   * @returns {Promise<Object>} Resultado da transação
   */
  async spendTokens(walletAddress, operation) {
    if (this.mockMode) {
      console.log(`[MOCK] Gastando tokens para: ${operation}`);
      await new Promise(resolve => setTimeout(resolve, 500));
      
      // Obter custo simulado
      const costResponse = await this.getTokenCost(operation);
      const costAmount = parseInt(costResponse.token_cost.split(' ')[0], 10);
      
      // Verificar saldo suficiente
      const currentBalance = parseInt(localStorage.getItem('tokenBalance') || '0', 10);
      if (currentBalance < costAmount) {
        throw new Error(`Saldo insuficiente: tem ${currentBalance}, necessário ${costAmount}`);
      }
      
      // Atualizar saldo
      const newBalance = currentBalance - costAmount;
      localStorage.setItem('tokenBalance', newBalance.toString());
      
      // Disparar evento
      window.dispatchEvent(new CustomEvent('tokens:spent', {
        detail: { amount: costAmount, operation }
      }));
      
      return {
        success: true,
        operation,
        cost: costAmount,
        new_balance: newBalance
      };
    }
    
    try {
      const response = await fetch(`${this.baseUrl}/spend-tokens`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ 
          wallet_address: walletAddress,
          operation 
        }),
      });
      
      if (!response.ok) {
        throw new Error(`Erro ao gastar tokens: ${response.status} ${response.statusText}`);
      }
      
      return await response.json();
    } catch (error) {
      console.error('Error spending tokens:', error);
      throw error;
    }
  }
}

// Inicializar e exportar serviço
window.web3Service = new Web3Service();

// Auto-conectar quando Web3Auth estiver pronto
window.addEventListener('user:authenticated', async (e) => {
  try {
    if (window.Web3Auth) {
      const token = await window.Web3Auth.getAccessToken();
      if (token) {
        // Salvar token para uso posterior
        localStorage.setItem('authToken', token);
        
        // Verificar com backend
        await window.web3Service.verifyToken(token);
      }
    }
  } catch (error) {
    console.error('Error during auto-connect:', error);
  }
});
