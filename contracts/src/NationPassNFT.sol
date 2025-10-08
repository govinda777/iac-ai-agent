// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

/**
 * @title NationPassNFT
 * @dev NFT ERC-721 para controle de acesso ao sistema IaC AI Agent
 * @author IaC AI Agent Team
 */
contract NationPassNFT is ERC721, Ownable, Pausable, ReentrancyGuard {
    using Counters for Counters.Counter;
    
    // ============================================
    // 📊 ESTRUTURAS E VARIÁVEIS
    // ============================================
    
    enum Tier {
        BASIC,    // 0 - Acesso básico
        PRO,      // 1 - Acesso profissional
        ENTERPRISE // 2 - Acesso empresarial
    }
    
    struct TierInfo {
        string name;
        uint256 price;
        uint256 maxSupply;
        uint256 currentSupply;
        bool isActive;
        string description;
        string imageURI;
    }
    
    struct NFTMetadata {
        Tier tier;
        uint256 mintTime;
        uint256 expiresAt;
        bool isActive;
        string customURI;
    }
    
    // Mapeamentos
    mapping(uint256 => NFTMetadata) public nftMetadata;
    mapping(Tier => TierInfo) public tierInfo;
    mapping(address => uint256[]) public userNFTs;
    mapping(address => bool) public hasActiveNFT;
    
    // Contadores
    Counters.Counter private _tokenIdCounter;
    
    // Variáveis de estado
    string private _baseTokenURI;
    uint256 public constant MAX_TOTAL_SUPPLY = 10000;
    uint256 public constant TIER_EXPIRY_DURATION = 365 days; // 1 ano
    
    // ============================================
    // 🎯 EVENTOS
    // ============================================
    
    event NFTMinted(
        address indexed to,
        uint256 tokenId,
        Tier tier,
        uint256 price,
        uint256 expiresAt
    );
    
    event NFTUpgraded(
        uint256 tokenId,
        Tier oldTier,
        Tier newTier,
        uint256 upgradePrice
    );
    
    event NFTExpired(
        uint256 tokenId,
        address indexed owner,
        Tier tier
    );
    
    event TierUpdated(
        Tier tier,
        uint256 price,
        uint256 maxSupply,
        bool isActive
    );
    
    event BaseURIUpdated(string newBaseURI);
    
    // ============================================
    // 🏗️ CONSTRUTOR
    // ============================================
    
    constructor() ERC721("Nation Pass NFT", "NATION") Ownable(msg.sender) {
        _baseTokenURI = "https://api.nation.fun/metadata/";
        
        // Configurar tiers iniciais
        _setupTiers();
    }
    
    // ============================================
    // 📦 CONFIGURAÇÃO DE TIERS
    // ============================================
    
    function _setupTiers() private {
        // Tier Basic
        tierInfo[Tier.BASIC] = TierInfo({
            name: "Basic Access",
            price: 0.01 ether,
            maxSupply: 5000,
            currentSupply: 0,
            isActive: true,
            description: "Basic access to IaC AI Agent",
            imageURI: "https://api.nation.fun/images/basic.png"
        });
        
        // Tier Pro
        tierInfo[Tier.PRO] = TierInfo({
            name: "Pro Access",
            price: 0.05 ether,
            maxSupply: 3000,
            currentSupply: 0,
            isActive: true,
            description: "Professional access with advanced features",
            imageURI: "https://api.nation.fun/images/pro.png"
        });
        
        // Tier Enterprise
        tierInfo[Tier.ENTERPRISE] = TierInfo({
            name: "Enterprise Access",
            price: 0.2 ether,
            maxSupply: 2000,
            currentSupply: 0,
            isActive: true,
            description: "Enterprise access with full features",
            imageURI: "https://api.nation.fun/images/enterprise.png"
        });
    }
    
    // ============================================
    // 🎫 MINT DE NFTS
    // ============================================
    
    function mintNFT(Tier tier) external payable nonReentrant whenNotPaused {
        TierInfo storage tierData = tierInfo[tier];
        
        // Validações
        require(tierData.isActive, "Tier is not active");
        require(msg.value >= tierData.price, "Insufficient payment");
        require(tierData.currentSupply < tierData.maxSupply, "Tier sold out");
        require(totalSupply() < MAX_TOTAL_SUPPLY, "Total supply exceeded");
        require(!hasActiveNFT[msg.sender], "User already has active NFT");
        
        // Gerar novo token ID
        _tokenIdCounter.increment();
        uint256 tokenId = _tokenIdCounter.current();
        
        // Mint do NFT
        _safeMint(msg.sender, tokenId);
        
        // Configurar metadata
        nftMetadata[tokenId] = NFTMetadata({
            tier: tier,
            mintTime: block.timestamp,
            expiresAt: block.timestamp + TIER_EXPIRY_DURATION,
            isActive: true,
            customURI: ""
        });
        
        // Atualizar estatísticas
        tierData.currentSupply++;
        userNFTs[msg.sender].push(tokenId);
        hasActiveNFT[msg.sender] = true;
        
        // Emitir evento
        emit NFTMinted(
            msg.sender,
            tokenId,
            tier,
            tierData.price,
            block.timestamp + TIER_EXPIRY_DURATION
        );
        
        // Refund se pagou mais que o necessário
        if (msg.value > tierData.price) {
            payable(msg.sender).transfer(msg.value - tierData.price);
        }
    }
    
    // ============================================
    // ⬆️ UPGRADE DE TIER
    // ============================================
    
    function upgradeNFT(uint256 tokenId, Tier newTier) external payable nonReentrant {
        require(_exists(tokenId), "Token does not exist");
        require(ownerOf(tokenId) == msg.sender, "Not the owner");
        
        NFTMetadata storage metadata = nftMetadata[tokenId];
        TierInfo storage newTierData = tierInfo[newTier];
        
        // Validações
        require(metadata.isActive, "NFT is not active");
        require(newTierData.isActive, "New tier is not active");
        require(uint256(newTier) > uint256(metadata.tier), "Can only upgrade to higher tier");
        require(newTierData.currentSupply < newTierData.maxSupply, "New tier sold out");
        
        uint256 upgradePrice = newTierData.price - tierInfo[metadata.tier].price;
        require(msg.value >= upgradePrice, "Insufficient payment for upgrade");
        
        // Atualizar tier
        Tier oldTier = metadata.tier;
        metadata.tier = newTier;
        
        // Atualizar estatísticas
        tierInfo[oldTier].currentSupply--;
        newTierData.currentSupply++;
        
        // Emitir evento
        emit NFTUpgraded(tokenId, oldTier, newTier, upgradePrice);
        
        // Refund se pagou mais que o necessário
        if (msg.value > upgradePrice) {
            payable(msg.sender).transfer(msg.value - upgradePrice);
        }
    }
    
    // ============================================
    // ⚙️ FUNÇÕES ADMINISTRATIVAS
    // ============================================
    
    function updateTier(
        Tier tier,
        uint256 price,
        uint256 maxSupply,
        bool isActive,
        string memory name,
        string memory description,
        string memory imageURI
    ) external onlyOwner {
        tierInfo[tier] = TierInfo({
            name: name,
            price: price,
            maxSupply: maxSupply,
            currentSupply: tierInfo[tier].currentSupply,
            isActive: isActive,
            description: description,
            imageURI: imageURI
        });
        
        emit TierUpdated(tier, price, maxSupply, isActive);
    }
    
    function setBaseURI(string memory newBaseURI) external onlyOwner {
        _baseTokenURI = newBaseURI;
        emit BaseURIUpdated(newBaseURI);
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
    }
    
    // ============================================
    // 📊 FUNÇÕES DE CONSULTA
    // ============================================
    
    function getNFTMetadata(uint256 tokenId) external view returns (NFTMetadata memory) {
        require(_exists(tokenId), "Token does not exist");
        return nftMetadata[tokenId];
    }
    
    function getTierInfo(Tier tier) external view returns (TierInfo memory) {
        return tierInfo[tier];
    }
    
    function getUserNFTs(address user) external view returns (uint256[] memory) {
        return userNFTs[user];
    }
    
    function getUserActiveNFT(address user) external view returns (uint256) {
        uint256[] memory userTokens = userNFTs[user];
        for (uint256 i = 0; i < userTokens.length; i++) {
            if (nftMetadata[userTokens[i]].isActive) {
                return userTokens[i];
            }
        }
        return 0; // Nenhum NFT ativo encontrado
    }
    
    function isNFTActive(uint256 tokenId) external view returns (bool) {
        require(_exists(tokenId), "Token does not exist");
        return nftMetadata[tokenId].isActive && 
               block.timestamp < nftMetadata[tokenId].expiresAt;
    }
    
    function getTotalSupply() public view returns (uint256) {
        return _tokenIdCounter.current();
    }
    
    function getTierSupply(Tier tier) external view returns (uint256) {
        return tierInfo[tier].currentSupply;
    }
    
    // ============================================
    // 🔗 FUNÇÕES DE METADATA
    // ============================================
    
    function _baseURI() internal view override returns (string memory) {
        return _baseTokenURI;
    }
    
    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        require(_exists(tokenId), "Token does not exist");
        
        NFTMetadata memory metadata = nftMetadata[tokenId];
        
        if (bytes(metadata.customURI).length > 0) {
            return metadata.customURI;
        }
        
        return string(abi.encodePacked(
            _baseURI(),
            "nft/",
            Strings.toString(tokenId),
            "?tier=",
            Strings.toString(uint256(metadata.tier))
        ));
    }
    
    // ============================================
    // 🔒 FUNÇÕES DE SEGURANÇA
    // ============================================
    
    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 tokenId,
        uint256 batchSize
    ) internal override whenNotPaused {
        super._beforeTokenTransfer(from, to, tokenId, batchSize);
        
        // Atualizar status de NFT ativo do usuário
        if (from != address(0) && to != address(0)) {
            hasActiveNFT[from] = false;
            hasActiveNFT[to] = true;
        }
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
