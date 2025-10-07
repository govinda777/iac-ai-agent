package web3

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gosouza/iac-ai-agent/pkg/config"
	"github.com/gosouza/iac-ai-agent/pkg/logger"
)

// BaseClient é o cliente para interagir com Base Network (L2 Ethereum)
type BaseClient struct {
	config     *config.Config
	logger     *logger.Logger
	ethClient  *ethclient.Client
	chainID    *big.Int
	rpcURL     string
}

// NewBaseClient cria um novo cliente Base Network
func NewBaseClient(cfg *config.Config, log *logger.Logger) (*BaseClient, error) {
	// Base Mainnet: https://mainnet.base.org
	// Base Goerli Testnet: https://goerli.base.org
	rpcURL := cfg.Web3.BaseRPCURL
	if rpcURL == "" {
		rpcURL = "https://mainnet.base.org" // Default to mainnet
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao Base Network: %w", err)
	}

	// Get chain ID
	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter chain ID: %w", err)
	}

	log.Info("Conectado ao Base Network", "chain_id", chainID.String(), "rpc", rpcURL)

	return &BaseClient{
		config:    cfg,
		logger:    log,
		ethClient: client,
		chainID:   chainID,
		rpcURL:    rpcURL,
	}, nil
}

// NetworkInfo contém informações sobre a rede
type NetworkInfo struct {
	ChainID         *big.Int `json:"chain_id"`
	NetworkName     string   `json:"network_name"`
	IsTestnet       bool     `json:"is_testnet"`
	LatestBlock     uint64   `json:"latest_block"`
	GasPrice        *big.Int `json:"gas_price"`
	SuggestedGasPrice *big.Int `json:"suggested_gas_price"`
	RPCURL          string   `json:"rpc_url"`
}

// GetNetworkInfo obtém informações sobre a rede Base
func (bc *BaseClient) GetNetworkInfo(ctx context.Context) (*NetworkInfo, error) {
	// Get latest block number
	blockNumber, err := bc.ethClient.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter número do bloco: %w", err)
	}

	// Get gas price
	gasPrice, err := bc.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter preço do gas: %w", err)
	}

	networkName := "Base Mainnet"
	isTestnet := false
	
	// Base Mainnet: Chain ID 8453
	// Base Goerli: Chain ID 84531
	// Base Sepolia: Chain ID 84532
	if bc.chainID.Cmp(big.NewInt(84531)) == 0 {
		networkName = "Base Goerli Testnet"
		isTestnet = true
	} else if bc.chainID.Cmp(big.NewInt(84532)) == 0 {
		networkName = "Base Sepolia Testnet"
		isTestnet = true
	}

	return &NetworkInfo{
		ChainID:           bc.chainID,
		NetworkName:       networkName,
		IsTestnet:         isTestnet,
		LatestBlock:       blockNumber,
		GasPrice:          gasPrice,
		SuggestedGasPrice: gasPrice,
		RPCURL:            bc.rpcURL,
	}, nil
}

// WalletBalance representa o saldo de uma wallet
type WalletBalance struct {
	Address      string   `json:"address"`
	BalanceWei   *big.Int `json:"balance_wei"`
	BalanceETH   string   `json:"balance_eth"`
	Nonce        uint64   `json:"nonce"`
	IsContract   bool     `json:"is_contract"`
}

// GetBalance obtém o saldo de uma wallet
func (bc *BaseClient) GetBalance(ctx context.Context, address string) (*WalletBalance, error) {
	addr := common.HexToAddress(address)

	// Get balance
	balance, err := bc.ethClient.BalanceAt(ctx, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter saldo: %w", err)
	}

	// Get nonce
	nonce, err := bc.ethClient.NonceAt(ctx, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter nonce: %w", err)
	}

	// Check if address is a contract
	code, err := bc.ethClient.CodeAt(ctx, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar código: %w", err)
	}
	isContract := len(code) > 0

	// Convert Wei to ETH (1 ETH = 10^18 Wei)
	ethValue := new(big.Float).Quo(
		new(big.Float).SetInt(balance),
		big.NewFloat(1e18),
	)

	return &WalletBalance{
		Address:    address,
		BalanceWei: balance,
		BalanceETH: ethValue.String(),
		Nonce:      nonce,
		IsContract: isContract,
	}, nil
}

// Transaction representa uma transação
type Transaction struct {
	Hash             string   `json:"hash"`
	From             string   `json:"from"`
	To               string   `json:"to"`
	Value            *big.Int `json:"value"`
	ValueETH         string   `json:"value_eth"`
	Gas              uint64   `json:"gas"`
	GasPrice         *big.Int `json:"gas_price"`
	Nonce            uint64   `json:"nonce"`
	Data             string   `json:"data"`
	BlockNumber      uint64   `json:"block_number,omitempty"`
	BlockHash        string   `json:"block_hash,omitempty"`
	Status           uint64   `json:"status,omitempty"` // 1 = success, 0 = failed
	Confirmations    uint64   `json:"confirmations"`
}

// GetTransaction obtém informações de uma transação
func (bc *BaseClient) GetTransaction(ctx context.Context, txHash string) (*Transaction, error) {
	hash := common.HexToHash(txHash)

	// Get transaction
	tx, isPending, err := bc.ethClient.TransactionByHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter transação: %w", err)
	}

	// Get transaction receipt for status and block info
	var receipt *types.Receipt
	var status uint64
	var blockNumber uint64
	var blockHash string
	var confirmations uint64

	if !isPending {
		receipt, err = bc.ethClient.TransactionReceipt(ctx, hash)
		if err == nil {
			status = receipt.Status
			blockNumber = receipt.BlockNumber.Uint64()
			blockHash = receipt.BlockHash.Hex()

			// Calculate confirmations
			currentBlock, err := bc.ethClient.BlockNumber(ctx)
			if err == nil && currentBlock >= blockNumber {
				confirmations = currentBlock - blockNumber + 1
			}
		}
	}

	// Get sender address
	msg := ethereum.CallMsg{
		From:     common.Address{},
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}
	from, err := types.Sender(types.LatestSignerForChainID(bc.chainID), tx)
	if err != nil {
		from = msg.From
	}

	// Convert value to ETH
	ethValue := new(big.Float).Quo(
		new(big.Float).SetInt(tx.Value()),
		big.NewFloat(1e18),
	)

	var toAddress string
	if tx.To() != nil {
		toAddress = tx.To().Hex()
	}

	return &Transaction{
		Hash:          tx.Hash().Hex(),
		From:          from.Hex(),
		To:            toAddress,
		Value:         tx.Value(),
		ValueETH:      ethValue.String(),
		Gas:           tx.Gas(),
		GasPrice:      tx.GasPrice(),
		Nonce:         tx.Nonce(),
		Data:          common.Bytes2Hex(tx.Data()),
		BlockNumber:   blockNumber,
		BlockHash:     blockHash,
		Status:        status,
		Confirmations: confirmations,
	}, nil
}

// WaitForTransaction aguarda uma transação ser minerada
func (bc *BaseClient) WaitForTransaction(ctx context.Context, txHash string) (*Transaction, error) {
	hash := common.HexToHash(txHash)

	bc.logger.Info("Aguardando confirmação da transação", "hash", txHash)

	// Wait for transaction to be mined
	receipt, err := bc.waitForReceipt(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("erro ao aguardar transação: %w", err)
	}

	bc.logger.Info("Transação confirmada", "hash", txHash, "block", receipt.BlockNumber.String(), "status", receipt.Status)

	// Get full transaction details
	return bc.GetTransaction(ctx, txHash)
}

// waitForReceipt aguarda o receipt de uma transação
func (bc *BaseClient) waitForReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	// Try to get receipt immediately
	receipt, err := bc.ethClient.TransactionReceipt(ctx, txHash)
	if err == nil {
		return receipt, nil
	}

	// If not found, wait for it
	// Base Network has ~2 second block time
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timeout := time.After(5 * time.Minute)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timeout:
			return nil, fmt.Errorf("timeout aguardando transação")
		case <-ticker.C:
			receipt, err := bc.ethClient.TransactionReceipt(ctx, txHash)
			if err == nil {
				return receipt, nil
			}
		}
	}
}

// EstimateGas estima o gas necessário para uma transação
func (bc *BaseClient) EstimateGas(ctx context.Context, from, to string, value *big.Int, data []byte) (uint64, error) {
	fromAddr := common.HexToAddress(from)
	toAddr := common.HexToAddress(to)

	msg := ethereum.CallMsg{
		From:  fromAddr,
		To:    &toAddr,
		Value: value,
		Data:  data,
	}

	gasLimit, err := bc.ethClient.EstimateGas(ctx, msg)
	if err != nil {
		return 0, fmt.Errorf("erro ao estimar gas: %w", err)
	}

	bc.logger.Info("Gas estimado", "gas_limit", gasLimit, "from", from, "to", to)
	return gasLimit, nil
}

// GetBlock obtém informações de um bloco
func (bc *BaseClient) GetBlock(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	block, err := bc.ethClient.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, fmt.Errorf("erro ao obter bloco: %w", err)
	}

	return block, nil
}

// GetLatestBlock obtém o bloco mais recente
func (bc *BaseClient) GetLatestBlock(ctx context.Context) (*types.Block, error) {
	block, err := bc.ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter último bloco: %w", err)
	}

	return block, nil
}

// ValidateAddress valida se um endereço é válido
func (bc *BaseClient) ValidateAddress(address string) bool {
	return common.IsHexAddress(address)
}

// Close fecha a conexão com o cliente
func (bc *BaseClient) Close() {
	if bc.ethClient != nil {
		bc.ethClient.Close()
	}
}
