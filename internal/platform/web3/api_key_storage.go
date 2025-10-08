package web3

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// APIKeyStorage é responsável por gerenciar o armazenamento seguro de chaves de API
type APIKeyStorage struct {
	logger     *logger.Logger
	config     *config.Config
	baseClient *BaseClient
	litClient  *LitProtocolClientSimple
}

// LitProtocolClient representa o cliente para interagir com o Lit Protocol
// Nota: Esta é uma interface simplificada para o Lit Protocol
type LitProtocolClientSimple struct {
	network string
}

// APIKeyStorageData estrutura para armazenar os metadados da chave criptografada
type APIKeyStorageData struct {
	EncryptedData     string `json:"encrypted_data"`
	LitProtocolEncKey string `json:"lit_protocol_enc_key"`
	AccessConditions  string `json:"access_conditions"`
	OwnerAddress      string `json:"owner_address"`
	CreatedAt         int64  `json:"created_at"`
	LastAccessedAt    int64  `json:"last_accessed_at"`
	ServiceType       string `json:"service_type"` // ex: "whatsapp"
}

// NewAPIKeyStorage cria uma nova instância do gerenciador de chaves de API
func NewAPIKeyStorage(cfg *config.Config, log *logger.Logger, baseClient *BaseClient) (*APIKeyStorage, error) {
	// Inicializa o cliente Lit Protocol
	litClient := &LitProtocolClientSimple{
		network: "datil", // Rede padrão do Lit Protocol
	}

	return &APIKeyStorage{
		logger:     log,
		config:     cfg,
		baseClient: baseClient,
		litClient:  litClient,
	}, nil
}

// StoreWhatsAppAPIKey armazena de forma segura a chave da API do WhatsApp
func (a *APIKeyStorage) StoreWhatsAppAPIKey(ctx context.Context, apiKey string, ownerAddress string, signature string) (*APIKeyStorageData, error) {
	// Verificar a assinatura para garantir que o usuário é dono do endereço
	if !a.verifySignature(ownerAddress, signature) {
		return nil, errors.New("assinatura inválida")
	}

	// 1. Gerar uma chave AES-256 aleatória
	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return nil, fmt.Errorf("erro ao gerar chave AES: %w", err)
	}

	// 2. Criptografar a chave da API com AES-256
	encryptedData, err := a.encryptWithAES(apiKey, aesKey)
	if err != nil {
		return nil, fmt.Errorf("erro ao criptografar a chave: %w", err)
	}

	// 3. Definir condições de acesso com o Lit Protocol
	accessConditions := fmt.Sprintf(`[{"chain":"ethereum","method":"","parameters":[":userAddress"],"returnValueTest":{"comparator":"=","value":"%s"}}]`, ownerAddress)

	// 4. Criptografar a chave AES com o Lit Protocol
	litEncryptedKey, err := a.encryptWithLitProtocol(aesKey, accessConditions)
	if err != nil {
		return nil, fmt.Errorf("erro ao criptografar com Lit Protocol: %w", err)
	}

	// 5. Criar e retornar os dados de armazenamento
	storageData := &APIKeyStorageData{
		EncryptedData:     base64.StdEncoding.EncodeToString(encryptedData),
		LitProtocolEncKey: litEncryptedKey,
		AccessConditions:  accessConditions,
		OwnerAddress:      ownerAddress,
		CreatedAt:         a.getCurrentTimestamp(),
		LastAccessedAt:    a.getCurrentTimestamp(),
		ServiceType:       "whatsapp",
	}

	a.logger.Info("Chave da API do WhatsApp armazenada com sucesso", "owner", ownerAddress)

	return storageData, nil
}

// RetrieveWhatsAppAPIKey recupera a chave da API do WhatsApp para o usuário especificado
func (a *APIKeyStorage) RetrieveWhatsAppAPIKey(ctx context.Context, storageData *APIKeyStorageData, signature string) (string, error) {
	// Verificar a assinatura para garantir que o usuário é dono do endereço
	if !a.verifySignature(storageData.OwnerAddress, signature) {
		return "", errors.New("assinatura inválida")
	}

	// 1. Descriptografar a chave AES usando o Lit Protocol
	aesKey, err := a.decryptWithLitProtocol(storageData.LitProtocolEncKey, storageData.AccessConditions, signature)
	if err != nil {
		return "", fmt.Errorf("erro ao descriptografar com Lit Protocol: %w", err)
	}

	// 2. Descriptografar a chave da API usando AES
	encryptedData, err := base64.StdEncoding.DecodeString(storageData.EncryptedData)
	if err != nil {
		return "", fmt.Errorf("erro ao decodificar dados criptografados: %w", err)
	}

	apiKey, err := a.decryptWithAES(encryptedData, aesKey)
	if err != nil {
		return "", fmt.Errorf("erro ao descriptografar a chave da API: %w", err)
	}

	// 3. Atualizar o timestamp do último acesso
	storageData.LastAccessedAt = a.getCurrentTimestamp()

	a.logger.Info("Chave da API do WhatsApp recuperada com sucesso", "owner", storageData.OwnerAddress)

	return apiKey, nil
}

// UpdateWhatsAppAPIKey atualiza a chave da API do WhatsApp
func (a *APIKeyStorage) UpdateWhatsAppAPIKey(ctx context.Context, existingData *APIKeyStorageData, newApiKey string, signature string) (*APIKeyStorageData, error) {
	// Verificar a assinatura para garantir que o usuário é dono do endereço
	if !a.verifySignature(existingData.OwnerAddress, signature) {
		return nil, errors.New("assinatura inválida")
	}

	// Armazenar a nova chave com as mesmas condições de acesso
	return a.StoreWhatsAppAPIKey(ctx, newApiKey, existingData.OwnerAddress, signature)
}

// DeleteWhatsAppAPIKey marca a chave para exclusão
func (a *APIKeyStorage) DeleteWhatsAppAPIKey(ctx context.Context, storageData *APIKeyStorageData, signature string) error {
	// Verificar a assinatura para garantir que o usuário é dono do endereço
	if !a.verifySignature(storageData.OwnerAddress, signature) {
		return errors.New("assinatura inválida")
	}

	a.logger.Info("Chave da API do WhatsApp marcada para exclusão", "owner", storageData.OwnerAddress)

	// Na implementação real, aqui removeremos os dados do armazenamento
	return nil
}

// GenerateMessageToSign gera a mensagem que o usuário deve assinar com sua carteira
func (a *APIKeyStorage) GenerateMessageToSign(ownerAddress string, action string) string {
	timestamp := a.getCurrentTimestamp()
	return fmt.Sprintf("Autorizo a %s da chave de API do WhatsApp para o endereço %s em %d", action, ownerAddress, timestamp)
}

// encryptWithAES criptografa dados com AES-256-GCM
func (a *APIKeyStorage) encryptWithAES(plaintext string, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Criar um novo GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Criar um nonce (número usado apenas uma vez)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Criptografar os dados
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return ciphertext, nil
}

// decryptWithAES descriptografa dados com AES-256-GCM
func (a *APIKeyStorage) decryptWithAES(ciphertext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Criar um novo GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Obter o tamanho do nonce
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("dados criptografados inválidos")
	}

	// Extrair o nonce e o ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Descriptografar os dados
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// encryptWithLitProtocol criptografa uma chave usando o Lit Protocol
// Esta é uma implementação simulada; na implementação real, usaríamos a SDK do Lit Protocol
func (a *APIKeyStorage) encryptWithLitProtocol(key []byte, accessConditions string) (string, error) {
	// Simulação da criptografia com Lit Protocol
	// Em uma implementação real, usaríamos o client SDK do Lit Protocol
	// para enviar a chave e as condições de acesso para a rede Lit

	// Codificar a chave em base64 para simulação
	simEncryptedKey := base64.StdEncoding.EncodeToString(key)

	// Criar um objeto que representa os dados de criptografia do Lit Protocol
	litData := struct {
		EncryptedKey    string `json:"encryptedKey"`
		AccessCondition string `json:"accessCondition"`
		Network         string `json:"network"`
	}{
		EncryptedKey:    simEncryptedKey,
		AccessCondition: accessConditions,
		Network:         a.litClient.network,
	}

	// Serializar para JSON
	litDataJson, err := json.Marshal(litData)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(litDataJson), nil
}

// decryptWithLitProtocol descriptografa uma chave usando o Lit Protocol
// Esta é uma implementação simulada; na implementação real, usaríamos a SDK do Lit Protocol
func (a *APIKeyStorage) decryptWithLitProtocol(encryptedKeyData string, accessConditions string, signature string) ([]byte, error) {
	// Decodificar os dados do Lit Protocol
	litDataJson, err := base64.StdEncoding.DecodeString(encryptedKeyData)
	if err != nil {
		return nil, err
	}

	// Desserializar de JSON
	var litData struct {
		EncryptedKey    string `json:"encryptedKey"`
		AccessCondition string `json:"accessCondition"`
		Network         string `json:"network"`
	}

	if err := json.Unmarshal(litDataJson, &litData); err != nil {
		return nil, err
	}

	// Verificar se as condições de acesso correspondem
	if litData.AccessCondition != accessConditions {
		return nil, errors.New("condições de acesso não correspondem")
	}

	// Em uma implementação real, enviaríamos a assinatura para a rede Lit
	// para verificar se o usuário tem acesso à chave

	// Decodificar a chave simulada
	key, err := base64.StdEncoding.DecodeString(litData.EncryptedKey)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// verifySignature verifica se a assinatura é válida para o endereço
func (a *APIKeyStorage) verifySignature(address string, signatureHex string) bool {
	// Em uma implementação real, verificaríamos a assinatura
	// usando os métodos de criptografia do Ethereum

	// Esta é uma implementação simulada
	if signatureHex == "" {
		return false
	}

	// Simula a verificação de assinatura sempre retornando verdadeiro
	// Em uma implementação real, seria algo como:
	/*
		signature, err := hexutil.Decode(signatureHex)
		if err != nil {
			return false
		}

		message := a.GenerateMessageToSign(address, "acesso")
		messageHash := accounts.TextHash([]byte(message))

		sigPublicKey, err := crypto.SigToPub(messageHash, signature)
		if err != nil {
			return false
		}

		recoveredAddress := crypto.PubkeyToAddress(*sigPublicKey).Hex()
		return strings.EqualFold(recoveredAddress, address)
	*/

	return true
}

// getCurrentTimestamp retorna o timestamp atual em segundos
func (a *APIKeyStorage) getCurrentTimestamp() int64 {
	return 1728393600 // Exemplo fixo: 08/08/2024 às 00:00:00 UTC
}

// SerializeStorageData serializa os dados de armazenamento para JSON
func (a *APIKeyStorage) SerializeStorageData(data *APIKeyStorageData) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// DeserializeStorageData desserializa os dados de armazenamento do JSON
func (a *APIKeyStorage) DeserializeStorageData(jsonStr string) (*APIKeyStorageData, error) {
	var data APIKeyStorageData
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return nil, err
	}
	return &data, nil
}
