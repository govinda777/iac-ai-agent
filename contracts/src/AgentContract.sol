// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "./IACaiToken.sol";
import "./NationPassNFT.sol";

/**
 * @title AgentContract
 * @dev Contrato principal que gerencia o agente IaC AI
 * @author IaC AI Agent Team
 */
contract AgentContract is Ownable, Pausable, ReentrancyGuard {
    // ============================================
    // 📊 ESTRUTURAS E VARIÁVEIS
    // ============================================
    
    struct AnalysisRequest {
        uint256 requestId;
        address requester;
        string analysisType;
        string inputData;
        uint256 tokenCost;
        uint256 timestamp;
        bool isCompleted;
        string result;
        string status;
    }
    
    struct AgentConfig {
        uint256 baseTokenCost;
        uint256 maxAnalysisLength;
        uint256 maxConcurrentRequests;
        bool isActive;
        string version;
    }
    
    struct UserStats {
        uint256 totalAnalyses;
        uint256 totalTokensSpent;
        uint256 lastAnalysisTime;
        bool isActive;
    }
    
    // Mapeamentos
    mapping(uint256 => AnalysisRequest) public analysisRequests;
    mapping(address => UserStats) public userStats;
    mapping(string => uint256) public analysisTypeCosts;
    mapping(address => bool) public authorizedCallers;
    
    // Contratos relacionados
    IACaiToken public tokenContract;
    NationPassNFT public nftContract;
    
    // Variáveis de estado
    uint256 private _requestIdCounter;
    AgentConfig public config;
    uint256 public constant MAX_REQUEST_LENGTH = 10000;
    uint256 public constant MIN_TOKEN_COST = 1;
    uint256 public constant MAX_TOKEN_COST = 1000;
    
    // ============================================
    // 🎯 EVENTOS
    // ============================================
    
    event AnalysisRequested(
        uint256 indexed requestId,
        address indexed requester,
        string analysisType,
        uint256 tokenCost,
        uint256 timestamp
    );
    
    event AnalysisCompleted(
        uint256 indexed requestId,
        address indexed requester,
        string result,
        uint256 timestamp
    );
    
    event AnalysisFailed(
        uint256 indexed requestId,
        address indexed requester,
        string reason,
        uint256 timestamp
    );
    
    event TokensRefunded(
        address indexed user,
        uint256 amount,
        string reason,
        uint256 timestamp
    );
    
    event ConfigUpdated(
        uint256 baseTokenCost,
        uint256 maxAnalysisLength,
        bool isActive,
        string version
    );
    
    event AuthorizedCallerAdded(address indexed caller);
    event AuthorizedCallerRemoved(address indexed caller);
    
    // ============================================
    // 🏗️ CONSTRUTOR
    // ============================================
    
    constructor(
        address _tokenContract,
        address _nftContract
    ) Ownable(msg.sender) {
        tokenContract = IACaiToken(_tokenContract);
        nftContract = NationPassNFT(_nftContract);
        
        // Configuração inicial
        config = AgentConfig({
            baseTokenCost: 10,
            maxAnalysisLength: 5000,
            maxConcurrentRequests: 100,
            isActive: true,
            version: "1.0.0"
        });
        
        // Configurar custos por tipo de análise
        analysisTypeCosts["infrastructure"] = 10;
        analysisTypeCosts["security"] = 15;
        analysisTypeCosts["performance"] = 12;
        analysisTypeCosts["cost"] = 8;
        analysisTypeCosts["compliance"] = 20;
        
        // Autorizar o owner como caller
        authorizedCallers[msg.sender] = true;
    }
    
    // ============================================
    // 🔍 SOLICITAÇÃO DE ANÁLISE
    // ============================================
    
    function requestAnalysis(
        string memory analysisType,
        string memory inputData
    ) external nonReentrant whenNotPaused {
        require(config.isActive, "Agent is not active");
        require(bytes(inputData).length <= config.maxAnalysisLength, "Input too long");
        require(bytes(inputData).length > 0, "Input cannot be empty");
        
        // Verificar se o usuário tem NFT ativo
        require(_hasActiveNFT(msg.sender), "No active NFT required");
        
        // Calcular custo em tokens
        uint256 tokenCost = _calculateTokenCost(analysisType, inputData);
        require(tokenCost >= MIN_TOKEN_COST && tokenCost <= MAX_TOKEN_COST, "Invalid token cost");
        
        // Verificar saldo de tokens
        require(tokenContract.balanceOf(msg.sender) >= tokenCost, "Insufficient token balance");
        
        // Verificar limite de requisições simultâneas
        require(_getActiveRequestCount(msg.sender) < 3, "Too many concurrent requests");
        
        // Gerar novo ID de requisição
        _requestIdCounter++;
        uint256 requestId = _requestIdCounter;
        
        // Criar requisição
        analysisRequests[requestId] = AnalysisRequest({
            requestId: requestId,
            requester: msg.sender,
            analysisType: analysisType,
            inputData: inputData,
            tokenCost: tokenCost,
            timestamp: block.timestamp,
            isCompleted: false,
            result: "",
            status: "pending"
        });
        
        // Cobrar tokens
        tokenContract.spendTokens(msg.sender, tokenCost, "Analysis request");
        
        // Atualizar estatísticas do usuário
        userStats[msg.sender].totalAnalyses++;
        userStats[msg.sender].totalTokensSpent += tokenCost;
        userStats[msg.sender].lastAnalysisTime = block.timestamp;
        userStats[msg.sender].isActive = true;
        
        // Emitir evento
        emit AnalysisRequested(
            requestId,
            msg.sender,
            analysisType,
            tokenCost,
            block.timestamp
        );
    }
    
    // ============================================
    // ✅ COMPLETAR ANÁLISE
    // ============================================
    
    function completeAnalysis(
        uint256 requestId,
        string memory result
    ) external onlyAuthorizedCaller {
        AnalysisRequest storage request = analysisRequests[requestId];
        
        require(request.requestId != 0, "Request does not exist");
        require(!request.isCompleted, "Request already completed");
        require(bytes(result).length > 0, "Result cannot be empty");
        
        // Atualizar requisição
        request.isCompleted = true;
        request.result = result;
        request.status = "completed";
        
        // Emitir evento
        emit AnalysisCompleted(
            requestId,
            request.requester,
            result,
            block.timestamp
        );
    }
    
    // ============================================
    // ❌ FALHAR ANÁLISE
    // ============================================
    
    function failAnalysis(
        uint256 requestId,
        string memory reason
    ) external onlyAuthorizedCaller {
        AnalysisRequest storage request = analysisRequests[requestId];
        
        require(request.requestId != 0, "Request does not exist");
        require(!request.isCompleted, "Request already completed");
        
        // Atualizar requisição
        request.isCompleted = true;
        request.status = "failed";
        request.result = reason;
        
        // Reembolsar tokens
        tokenContract.transfer(request.requester, request.tokenCost);
        
        // Emitir eventos
        emit AnalysisFailed(requestId, request.requester, reason, block.timestamp);
        emit TokensRefunded(request.requester, request.tokenCost, reason, block.timestamp);
    }
    
    // ============================================
    // ⚙️ FUNÇÕES ADMINISTRATIVAS
    // ============================================
    
    function updateConfig(
        uint256 baseTokenCost,
        uint256 maxAnalysisLength,
        uint256 maxConcurrentRequests,
        bool isActive,
        string memory version
    ) external onlyOwner {
        require(baseTokenCost > 0, "Base cost must be positive");
        require(maxAnalysisLength > 0, "Max length must be positive");
        require(maxConcurrentRequests > 0, "Max concurrent must be positive");
        
        config = AgentConfig({
            baseTokenCost: baseTokenCost,
            maxAnalysisLength: maxAnalysisLength,
            maxConcurrentRequests: maxConcurrentRequests,
            isActive: isActive,
            version: version
        });
        
        emit ConfigUpdated(baseTokenCost, maxAnalysisLength, isActive, version);
    }
    
    function updateAnalysisTypeCost(
        string memory analysisType,
        uint256 cost
    ) external onlyOwner {
        require(cost >= MIN_TOKEN_COST && cost <= MAX_TOKEN_COST, "Invalid cost");
        analysisTypeCosts[analysisType] = cost;
    }
    
    function addAuthorizedCaller(address caller) external onlyOwner {
        authorizedCallers[caller] = true;
        emit AuthorizedCallerAdded(caller);
    }
    
    function removeAuthorizedCaller(address caller) external onlyOwner {
        authorizedCallers[caller] = false;
        emit AuthorizedCallerRemoved(caller);
    }
    
    function pause() external onlyOwner {
        _pause();
    }
    
    function unpause() external onlyOwner {
        _unpause();
    }
    
    function emergencyWithdraw() external onlyOwner {
        uint256 balance = address(this).balance;
        if (balance > 0) {
            payable(owner()).transfer(balance);
        }
    }
    
    // ============================================
    // 📊 FUNÇÕES DE CONSULTA
    // ============================================
    
    function getAnalysisRequest(uint256 requestId) external view returns (AnalysisRequest memory) {
        require(analysisRequests[requestId].requestId != 0, "Request does not exist");
        return analysisRequests[requestId];
    }
    
    function getUserStats(address user) external view returns (UserStats memory) {
        return userStats[user];
    }
    
    function getAnalysisTypeCost(string memory analysisType) external view returns (uint256) {
        return analysisTypeCosts[analysisType];
    }
    
    function getActiveRequestCount(address user) external view returns (uint256) {
        return _getActiveRequestCount(user);
    }
    
    function getTotalRequests() external view returns (uint256) {
        return _requestIdCounter;
    }
    
    function getConfig() external view returns (AgentConfig memory) {
        return config;
    }
    
    // ============================================
    // 🔒 FUNÇÕES INTERNAS
    // ============================================
    
    function _hasActiveNFT(address user) private view returns (bool) {
        uint256 activeNFT = nftContract.getUserActiveNFT(user);
        return activeNFT != 0 && nftContract.isNFTActive(activeNFT);
    }
    
    function _calculateTokenCost(
        string memory analysisType,
        string memory inputData
    ) private view returns (uint256) {
        uint256 baseCost = analysisTypeCosts[analysisType];
        if (baseCost == 0) {
            baseCost = config.baseTokenCost;
        }
        
        // Adicionar custo baseado no tamanho da entrada
        uint256 sizeCost = bytes(inputData).length / 100; // 1 token por 100 caracteres
        
        return baseCost + sizeCost;
    }
    
    function _getActiveRequestCount(address user) private view returns (uint256) {
        uint256 count = 0;
        for (uint256 i = 1; i <= _requestIdCounter; i++) {
            AnalysisRequest memory request = analysisRequests[i];
            if (request.requester == user && !request.isCompleted) {
                count++;
            }
        }
        return count;
    }
    
    // ============================================
    // 🔒 MODIFIERS
    // ============================================
    
    modifier onlyAuthorizedCaller() {
        require(authorizedCallers[msg.sender], "Not authorized caller");
        _;
    }
    
    // Função para receber ETH
    receive() external payable {
        // Permitir recebimento de ETH
    }
    
    // Função de fallback
    fallback() external payable {
        revert("Function not found");
    }
}
