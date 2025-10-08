// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

/**
 * @title IACaiToken
 * @dev Token ERC-20 para pagamentos no sistema IaC AI Agent
 * @author IaC AI Agent Team
 */
contract IACaiToken is ERC20, Ownable, Pausable, ReentrancyGuard {
    // ============================================
    // 📊 ESTRUTURAS E VARIÁVEIS
    // ============================================
    
    struct TokenPackage {
        uint8 id;
        uint256 tokenAmount;
        uint256 price;
        uint8 discountPercent;
        bool isActive;
        string name;
        string description;
    }
    
    struct TokenStats {
        uint256 totalSold;
        uint256 totalRevenue;
        uint256 totalBurned;
        uint256 lastPurchaseTime;
    }
    
    // Mapeamentos
    mapping(uint8 => TokenPackage) public packages;
    mapping(address => uint256) public userPurchaseHistory;
    mapping(address => uint256) public lastPurchaseTime;
    
    // Variáveis de estado
    uint8 public constant MAX_PACKAGES = 10;
    uint256 public constant MAX_SUPPLY = 1_000_000 * 10**18; // 1M tokens
    uint256 public constant MIN_PURCHASE = 0.001 ether;
    uint256 public constant MAX_PURCHASE = 1 ether;
    
    TokenStats public stats;
    
    // ============================================
    // 🎯 EVENTOS
    // ============================================
    
    event TokensPurchased(
        address indexed buyer,
        uint8 packageId,
        uint256 tokenAmount,
        uint256 price,
        uint256 timestamp
    );
    
    event TokensSpent(
        address indexed spender,
        uint256 amount,
        string reason,
        uint256 timestamp
    );
    
    event TokensBurned(
        address indexed burner,
        uint256 amount,
        string reason,
        uint256 timestamp
    );
    
    event PackageUpdated(
        uint8 packageId,
        uint256 tokenAmount,
        uint256 price,
        bool isActive
    );
    
    event EmergencyWithdraw(
        address indexed owner,
        uint256 amount,
        uint256 timestamp
    );
    
    // ============================================
    // 🏗️ CONSTRUTOR
    // ============================================
    
    constructor() ERC20("IaC AI Token", "IACAI") Ownable(msg.sender) {
        // Mint supply inicial para o contrato
        _mint(address(this), MAX_SUPPLY);
        
        // Configurar pacotes iniciais
        _setupPackages();
        
        // Configurar estatísticas iniciais
        stats = TokenStats({
            totalSold: 0,
            totalRevenue: 0,
            totalBurned: 0,
            lastPurchaseTime: block.timestamp
        });
    }
    
    // ============================================
    // 📦 CONFIGURAÇÃO DE PACOTES
    // ============================================
    
    function _setupPackages() private {
        // Pacote 1: Starter Pack
        packages[1] = TokenPackage({
            id: 1,
            tokenAmount: 100 * 10**decimals(),
            price: 0.005 ether,
            discountPercent: 0,
            isActive: true,
            name: "Starter Pack",
            description: "Perfect for getting started"
        });
        
        // Pacote 2: Power Pack (10% desconto)
        packages[2] = TokenPackage({
            id: 2,
            tokenAmount: 500 * 10**decimals(),
            price: 0.0225 ether,
            discountPercent: 10,
            isActive: true,
            name: "Power Pack",
            description: "Great value for regular users"
        });
        
        // Pacote 3: Pro Pack (15% desconto)
        packages[3] = TokenPackage({
            id: 3,
            tokenAmount: 1000 * 10**decimals(),
            price: 0.0425 ether,
            discountPercent: 15,
            isActive: true,
            name: "Pro Pack",
            description: "Best for professionals"
        });
        
        // Pacote 4: Enterprise Pack (25% desconto)
        packages[4] = TokenPackage({
            id: 4,
            tokenAmount: 5000 * 10**decimals(),
            price: 0.1875 ether,
            discountPercent: 25,
            isActive: true,
            name: "Enterprise Pack",
            description: "Maximum value for teams"
        });
    }
    
    // ============================================
    // 💰 COMPRA DE TOKENS
    // ============================================
    
    function buyTokens(uint8 packageId) external payable nonReentrant whenNotPaused {
        TokenPackage storage package = packages[packageId];
        
        // Validações
        require(package.isActive, "Package is not active");
        require(msg.value >= package.price, "Insufficient payment");
        require(msg.value >= MIN_PURCHASE, "Payment too small");
        require(msg.value <= MAX_PURCHASE, "Payment too large");
        require(balanceOf(address(this)) >= package.tokenAmount, "Insufficient token supply");
        
        // Verificar limite de compra por usuário (máximo 1 compra por hora)
        require(
            block.timestamp >= lastPurchaseTime[msg.sender] + 1 hours,
            "Purchase cooldown active"
        );
        
        // Transferir tokens do contrato para o comprador
        _transfer(address(this), msg.sender, package.tokenAmount);
        
        // Atualizar estatísticas
        stats.totalSold += package.tokenAmount;
        stats.totalRevenue += package.price;
        stats.lastPurchaseTime = block.timestamp;
        
        userPurchaseHistory[msg.sender] += package.tokenAmount;
        lastPurchaseTime[msg.sender] = block.timestamp;
        
        // Emitir evento
        emit TokensPurchased(
            msg.sender,
            packageId,
            package.tokenAmount,
            package.price,
            block.timestamp
        );
        
        // Refund se pagou mais que o necessário
        if (msg.value > package.price) {
            payable(msg.sender).transfer(msg.value - package.price);
        }
    }
    
    // ============================================
    // 💸 GASTAR TOKENS
    // ============================================
    
    function spendTokens(
        address spender,
        uint256 amount,
        string memory reason
    ) external onlyOwner nonReentrant {
        require(balanceOf(spender) >= amount, "Insufficient balance");
        require(amount > 0, "Amount must be positive");
        
        // Transferir tokens do usuário para o owner
        _transfer(spender, owner(), amount);
        
        // Emitir evento
        emit TokensSpent(spender, amount, reason, block.timestamp);
    }
    
    // ============================================
    // 🔥 QUEIMAR TOKENS
    // ============================================
    
    function burnTokens(uint256 amount, string memory reason) external onlyOwner {
        require(amount > 0, "Amount must be positive");
        require(balanceOf(address(this)) >= amount, "Insufficient contract balance");
        
        _burn(address(this), amount);
        
        stats.totalBurned += amount;
        
        emit TokensBurned(msg.sender, amount, reason, block.timestamp);
    }
    
    // ============================================
    // ⚙️ FUNÇÕES ADMINISTRATIVAS
    // ============================================
    
    function updatePackage(
        uint8 packageId,
        uint256 tokenAmount,
        uint256 price,
        uint8 discountPercent,
        bool isActive,
        string memory name,
        string memory description
    ) external onlyOwner {
        require(packageId > 0 && packageId <= MAX_PACKAGES, "Invalid package ID");
        
        packages[packageId] = TokenPackage({
            id: packageId,
            tokenAmount: tokenAmount,
            price: price,
            discountPercent: discountPercent,
            isActive: isActive,
            name: name,
            description: description
        });
        
        emit PackageUpdated(packageId, tokenAmount, price, isActive);
    }
    
    function pause() external onlyOwner {
        _pause();
    }
    
    function unpause() external onlyOwner {
        _unpause();
    }
    
    function emergencyWithdraw() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");
        
        payable(owner()).transfer(balance);
        
        emit EmergencyWithdraw(owner(), balance, block.timestamp);
    }
    
    // ============================================
    // 📊 FUNÇÕES DE CONSULTA
    // ============================================
    
    function getPackageInfo(uint8 packageId) external view returns (TokenPackage memory) {
        require(packageId > 0 && packageId <= MAX_PACKAGES, "Invalid package ID");
        return packages[packageId];
    }
    
    function getStats() external view returns (TokenStats memory) {
        return stats;
    }
    
    function getUserStats(address user) external view returns (
        uint256 purchaseHistory,
        uint256 lastPurchase,
        uint256 currentBalance
    ) {
        return (
            userPurchaseHistory[user],
            lastPurchaseTime[user],
            balanceOf(user)
        );
    }
    
    function getAvailableSupply() external view returns (uint256) {
        return balanceOf(address(this));
    }
    
    function getTotalSupply() external view returns (uint256) {
        return totalSupply();
    }
    
    // ============================================
    // 🔒 FUNÇÕES DE SEGURANÇA
    // ============================================
    
    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) internal override whenNotPaused {
        super._beforeTokenTransfer(from, to, amount);
    }
    
    // Função para receber ETH
    receive() external payable {
        // Permitir recebimento de ETH para compras
    }
    
    // Função de fallback
    fallback() external payable {
        revert("Function not found");
    }
}
