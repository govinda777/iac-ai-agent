package web3

import (
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// WalletTokenGenerator gera tokens para autenticação da wallet com Nation.fun
type WalletTokenGenerator struct {
	config *config.Config
	logger *logger.Logger
}

// NewWalletTokenGenerator cria um novo gerador de tokens de wallet
func NewWalletTokenGenerator(cfg *config.Config, log *logger.Logger) *WalletTokenGenerator {
	return &WalletTokenGenerator{
		config: cfg,
		logger: log,
	}
}

// GenerateToken gera um token de autenticação para Nation.fun
func (wtg *WalletTokenGenerator) GenerateToken() (string, error) {
	// Verificar se já temos um token configurado
	if wtg.config.Web3.WalletToken != "" {
		wtg.logger.Info("Usando WALLET_TOKEN existente")
		return wtg.config.Web3.WalletToken, nil
	}

	// Verificar se temos wallet address
	if wtg.config.Web3.WalletAddress == "" {
		return "", fmt.Errorf("WALLET_ADDRESS não configurado")
	}

	// NOTA: Não usamos mais a chave privada diretamente por razões de segurança
	// Em vez disso, geramos um token temporário baseado no endereço da wallet
	// Em produção, use um serviço externo de assinatura ou um token pré-gerado

	// Gerar um token temporário para desenvolvimento
	token := wtg.generateTemporaryToken()

	// Atualizar na configuração para uso futuro
	wtg.config.Web3.WalletToken = token
	wtg.logger.Info("WALLET_TOKEN temporário gerado para desenvolvimento")
	wtg.logger.Warn("Este token é apenas para desenvolvimento! Em produção, use um token válido pré-gerado")

	return token, nil
}

// createToken cria um token de autenticação assinado
func (wtg *WalletTokenGenerator) createToken(privateKey *ecdsa.PrivateKey) (string, error) {
	// Obtém o endereço público da wallet
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("erro ao converter chave pública")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	// Verifica se o endereço da chave privada corresponde ao endereço configurado
	if !strings.EqualFold(address, wtg.config.Web3.WalletAddress) {
		wtg.logger.Warn("O endereço derivado da chave privada não corresponde ao WALLET_ADDRESS configurado",
			"derived_address", address,
			"configured_address", wtg.config.Web3.WalletAddress)
		return "", fmt.Errorf("chave privada não corresponde ao endereço da wallet")
	}

	// Timestamp atual como nonce
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	// Dados para assinatura (endereço + timestamp)
	data := []byte(address + timestamp)

	// Hash dos dados
	hash := crypto.Keccak256Hash(data)

	// Assinar o hash
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", fmt.Errorf("erro ao assinar mensagem: %w", err)
	}

	// Converter assinatura para hexadecimal
	sigHex := hexutil.Encode(signature)

	// Criar token no formato esperado pela Nation.fun
	// Formato: address:timestamp:signature
	token := fmt.Sprintf("%s:%s:%s", address, timestamp, sigHex)

	// Aplica HMAC-SHA256 para finalizar o token
	// Nota: Na prática, você precisaria de uma chave secreta da Nation.fun para este passo
	// Este é um placeholder que simula esse processo
	h := hmac.New(sha256.New, []byte("nation_fun_secret"))
	h.Write([]byte(token))
	tokenHash := hex.EncodeToString(h.Sum(nil))

	// Formato final: prefixo_versão_tokenHash
	finalToken := fmt.Sprintf("nft_v1_%s", tokenHash)

	wtg.logger.Debug("Token gerado com sucesso",
		"address", address,
		"token_format", "nft_v1_*****")

	return finalToken, nil
}

// generateTemporaryToken gera um token temporário para desenvolvimento
func (wtg *WalletTokenGenerator) generateTemporaryToken() string {
	// Endereço da wallet
	address := wtg.config.Web3.WalletAddress

	// Timestamp atual como nonce
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	// Dados para hash
	data := []byte(address + timestamp + "development_only")

	// Criar hash
	h := sha256.New()
	h.Write(data)
	tokenHash := hex.EncodeToString(h.Sum(nil))

	// Formato: dev_timestamp_hash
	return fmt.Sprintf("dev_%s_%s", timestamp, tokenHash[:32])
}

// VerifyWalletOwnership verifica se o token pertence ao endereço
func (wtg *WalletTokenGenerator) VerifyWalletOwnership() (bool, error) {
	// Não podemos mais verificar a propriedade da wallet sem a chave privada
	// Em vez disso, assumimos que o token fornecido é válido
	wtg.logger.Warn("Verificação de propriedade da wallet não é mais possível sem chave privada")
	wtg.logger.Info("Assumindo que o token fornecido é válido")

	// Verifica se temos um token
	if wtg.config.Web3.WalletToken == "" {
		return false, fmt.Errorf("WALLET_TOKEN não configurado")
	}

	// Para desenvolvimento, sempre retorna true
	return true, nil
}

// ValidateNFTOwnership verifica se a wallet possui o NFT
// Esta é uma versão simplificada para o MVP - na implementação final
// precisaria fazer uma consulta à blockchain
func (wtg *WalletTokenGenerator) ValidateNFTOwnership() (bool, error) {
	if wtg.config.Web3.NFTAccessContractAddress == "" {
		return false, fmt.Errorf("NFT_CONTRACT_ADDRESS não configurado")
	}

	walletAddr := common.HexToAddress(wtg.config.Web3.WalletAddress)
	nftAddr := common.HexToAddress(wtg.config.Web3.NFTAccessContractAddress)

	// TODO: Na versão final, verificar na blockchain se a wallet possui o NFT
	// Simulação para MVP:
	wtg.logger.Info("Validando NFT ownership", "wallet", walletAddr.Hex(), "nft_contract", nftAddr.Hex())
	return true, nil
}
