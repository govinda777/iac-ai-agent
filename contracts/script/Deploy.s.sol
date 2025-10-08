// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/IACaiToken.sol";
import "../src/NationPassNFT.sol";
import "../src/AgentContract.sol";

/**
 * @title Deploy
 * @dev Script principal de deploy dos contratos
 * @author IaC AI Agent Team
 */
contract Deploy is Script {
    // ============================================
    // 📊 VARIÁVEIS DE CONFIGURAÇÃO
    // ============================================
    
    struct DeploymentConfig {
        address deployer;
        address tokenOwner;
        address nftOwner;
        address agentOwner;
        bool verifyContracts;
        bool pauseAfterDeploy;
    }
    
    struct DeploymentAddresses {
        address tokenContract;
        address nftContract;
        address agentContract;
    }
    
    // ============================================
    // 🚀 FUNÇÃO PRINCIPAL DE DEPLOY
    // ============================================
    
    function run() external {
        // Configuração do deploy
        DeploymentConfig memory config = _getDeploymentConfig();
        
        // Log de início
        console.log("🚀 Starting IaC AI Agent deployment...");
        console.log("📍 Deployer:", config.deployer);
        console.log("🔐 Verify contracts:", config.verifyContracts);
        
        // Deploy dos contratos
        DeploymentAddresses memory addresses = _deployContracts(config);
        
        // Configuração pós-deploy
        _postDeploySetup(addresses, config);
        
        // Log de conclusão
        console.log("✅ Deployment completed successfully!");
        console.log("📋 Contract addresses:");
        console.log("   Token Contract:", addresses.tokenContract);
        console.log("   NFT Contract:", addresses.nftContract);
        console.log("   Agent Contract:", addresses.agentContract);
        
        // Salvar endereços em arquivo
        _saveAddresses(addresses);
    }
    
    // ============================================
    // 🏗️ DEPLOY DOS CONTRATOS
    // ============================================
    
    function _deployContracts(
        DeploymentConfig memory config
    ) private returns (DeploymentAddresses memory) {
        DeploymentAddresses memory addresses;
        
        // Deploy do Token Contract
        console.log("📦 Deploying IACaiToken...");
        IACaiToken token = new IACaiToken();
        addresses.tokenContract = address(token);
        console.log("✅ IACaiToken deployed at:", addresses.tokenContract);
        
        // Deploy do NFT Contract
        console.log("🎫 Deploying NationPassNFT...");
        NationPassNFT nft = new NationPassNFT();
        addresses.nftContract = address(nft);
        console.log("✅ NationPassNFT deployed at:", addresses.nftContract);
        
        // Deploy do Agent Contract
        console.log("🤖 Deploying AgentContract...");
        AgentContract agent = new AgentContract(
            addresses.tokenContract,
            addresses.nftContract
        );
        addresses.agentContract = address(agent);
        console.log("✅ AgentContract deployed at:", addresses.agentContract);
        
        return addresses;
    }
    
    // ============================================
    // ⚙️ CONFIGURAÇÃO PÓS-DEPLOY
    // ============================================
    
    function _postDeploySetup(
        DeploymentAddresses memory addresses,
        DeploymentConfig memory config
    ) private {
        console.log("⚙️ Setting up post-deployment configuration...");
        
        // Configurar Agent Contract como caller autorizado
        IACaiToken token = IACaiToken(addresses.tokenContract);
        AgentContract agent = AgentContract(addresses.agentContract);
        
        // Transferir ownership se necessário
        if (config.tokenOwner != config.deployer) {
            console.log("🔄 Transferring token ownership...");
            token.transferOwnership(config.tokenOwner);
        }
        
        if (config.nftOwner != config.deployer) {
            console.log("🔄 Transferring NFT ownership...");
            NationPassNFT nft = NationPassNFT(addresses.nftContract);
            nft.transferOwnership(config.nftOwner);
        }
        
        if (config.agentOwner != config.deployer) {
            console.log("🔄 Transferring agent ownership...");
            agent.transferOwnership(config.agentOwner);
        }
        
        // Pausar contratos se configurado
        if (config.pauseAfterDeploy) {
            console.log("⏸️ Pausing contracts...");
            token.pause();
            nft.pause();
            agent.pause();
        }
        
        console.log("✅ Post-deployment setup completed");
    }
    
    // ============================================
    // 🔧 CONFIGURAÇÃO DO DEPLOY
    // ============================================
    
    function _getDeploymentConfig() private view returns (DeploymentConfig memory) {
        return DeploymentConfig({
            deployer: msg.sender,
            tokenOwner: vm.envOr("TOKEN_OWNER", msg.sender),
            nftOwner: vm.envOr("NFT_OWNER", msg.sender),
            agentOwner: vm.envOr("AGENT_OWNER", msg.sender),
            verifyContracts: vm.envOr("VERIFY_CONTRACTS", true),
            pauseAfterDeploy: vm.envOr("PAUSE_AFTER_DEPLOY", false)
        });
    }
    
    // ============================================
    // 💾 SALVAR ENDEREÇOS
    // ============================================
    
    function _saveAddresses(DeploymentAddresses memory addresses) private {
        string memory network = vm.envOr("NETWORK", "local");
        string memory filename = string(abi.encodePacked("deployments/", network, ".json"));
        
        // Criar JSON com endereços
        string memory json = string(abi.encodePacked(
            '{\n',
            '  "network": "', network, '",\n',
            '  "deploymentTime": "', vm.toString(block.timestamp), '",\n',
            '  "deployer": "', vm.toString(msg.sender), '",\n',
            '  "contracts": {\n',
            '    "IACaiToken": "', vm.toString(addresses.tokenContract), '",\n',
            '    "NationPassNFT": "', vm.toString(addresses.nftContract), '",\n',
            '    "AgentContract": "', vm.toString(addresses.agentContract), '"\n',
            '  }\n',
            '}'
        ));
        
        // Salvar arquivo
        vm.writeFile(filename, json);
        console.log("💾 Addresses saved to:", filename);
    }
    
    // ============================================
    // 🔍 FUNÇÕES DE VERIFICAÇÃO
    // ============================================
    
    function verifyContracts(DeploymentAddresses memory addresses) external {
        console.log("🔍 Verifying contracts...");
        
        // Verificar Token Contract
        try this.verifyTokenContract(addresses.tokenContract) {
            console.log("✅ Token contract verified");
        } catch {
            console.log("❌ Token contract verification failed");
        }
        
        // Verificar NFT Contract
        try this.verifyNFTContract(addresses.nftContract) {
            console.log("✅ NFT contract verified");
        } catch {
            console.log("❌ NFT contract verification failed");
        }
        
        // Verificar Agent Contract
        try this.verifyAgentContract(addresses.agentContract) {
            console.log("✅ Agent contract verified");
        } catch {
            console.log("❌ Agent contract verification failed");
        }
    }
    
    function verifyTokenContract(address tokenAddress) external {
        IACaiToken token = IACaiToken(tokenAddress);
        
        // Verificações básicas
        require(token.name() == "IaC AI Token", "Invalid token name");
        require(token.symbol() == "IACAI", "Invalid token symbol");
        require(token.totalSupply() == 1_000_000 * 10**18, "Invalid total supply");
        
        console.log("✅ Token contract verification passed");
    }
    
    function verifyNFTContract(address nftAddress) external {
        NationPassNFT nft = NationPassNFT(nftAddress);
        
        // Verificações básicas
        require(keccak256(bytes(nft.name())) == keccak256(bytes("Nation Pass NFT")), "Invalid NFT name");
        require(keccak256(bytes(nft.symbol())) == keccak256(bytes("NATION")), "Invalid NFT symbol");
        
        console.log("✅ NFT contract verification passed");
    }
    
    function verifyAgentContract(address agentAddress) external {
        AgentContract agent = AgentContract(agentAddress);
        
        // Verificações básicas
        require(agent.config().baseTokenCost > 0, "Invalid base token cost");
        require(agent.config().maxAnalysisLength > 0, "Invalid max analysis length");
        
        console.log("✅ Agent contract verification passed");
    }
}
