/**
 * WhatsApp API Key Management
 * 
 * Este módulo gerencia o armazenamento e recuperação seguros da chave de API do WhatsApp
 * usando a abordagem de criptografia baseada em Web3 e Lit Protocol.
 */

// Importações necessárias
import { ethers } from 'ethers';
import LitJsSdk from '@lit-protocol/lit-node-client';

// Classe para gerenciamento da chave de API do WhatsApp
class WhatsAppAPIKeyManager {
  constructor() {
    this.litClient = new LitJsSdk.LitNodeClient({
      litNetwork: 'datil', // Rede do Lit Protocol (mainnet)
      debug: false,
    });
    this.initialized = false;
    this.storageKey = 'whatsapp_api_key_storage';
    this.provider = null;
    this.wallet = null;
    this.storageData = null;
  }

  /**
   * Inicializa o cliente Lit Protocol e o provedor Ethereum
   */
  async initialize() {
    if (this.initialized) return;

    // Conectar ao Lit Protocol
    await this.litClient.connect();
    
    // Inicializar provedor Ethereum (usando o provedor injetado do navegador, como MetaMask)
    if (window.ethereum) {
      this.provider = new ethers.providers.Web3Provider(window.ethereum);
    }
    
    this.initialized = true;
    console.log('WhatsApp API Key Manager inicializado');
    
    // Carregar dados do armazenamento local, se existirem
    this.loadStorageData();
  }

  /**
   * Conecta a carteira do usuário
   */
  async connectWallet() {
    if (!this.initialized) await this.initialize();
    
    if (!this.provider) {
      throw new Error('Provedor Web3 não encontrado. Por favor, instale o MetaMask.');
    }
    
    // Solicitar conexão da carteira
    await this.provider.send('eth_requestAccounts', []);
    this.wallet = this.provider.getSigner();
    const address = await this.wallet.getAddress();
    
    console.log('Carteira conectada:', address);
    return address;
  }

  /**
   * Armazena a chave da API do WhatsApp de forma segura
   * @param {string} apiKey - A chave da API do WhatsApp a ser armazenada
   */
  async storeAPIKey(apiKey) {
    if (!this.wallet) {
      await this.connectWallet();
    }
    
    const address = await this.wallet.getAddress();
    
    try {
      // 1. Gerar mensagem para assinatura
      const timestamp = Math.floor(Date.now() / 1000);
      const messageToSign = `Autorizo o armazenamento da chave de API do WhatsApp para o endereço ${address} em ${timestamp}`;
      
      // 2. Solicitar assinatura do usuário
      const signature = await this.wallet.signMessage(messageToSign);
      
      // 3. Gerar uma chave AES-256 aleatória
      const symmetricKey = await window.crypto.subtle.generateKey(
        {
          name: "AES-GCM",
          length: 256
        },
        true,
        ["encrypt", "decrypt"]
      );
      
      // 4. Criptografar a chave da API com AES-256
      const encoder = new TextEncoder();
      const apiKeyEncoded = encoder.encode(apiKey);
      const iv = window.crypto.getRandomValues(new Uint8Array(12));
      
      const encryptedApiKey = await window.crypto.subtle.encrypt(
        {
          name: "AES-GCM",
          iv: iv
        },
        symmetricKey,
        apiKeyEncoded
      );
      
      // 5. Exportar a chave simétrica para criptografar com Lit Protocol
      const exportedSymKey = await window.crypto.subtle.exportKey("raw", symmetricKey);
      
      // 6. Definir condições de acesso para o Lit Protocol
      const accessControlConditions = [
        {
          chain: 'ethereum',
          method: '',
          parameters: [':userAddress'],
          returnValueTest: {
            comparator: '=',
            value: address
          }
        }
      ];
      
      // 7. Criptografar a chave simétrica com o Lit Protocol
      const encryptedSymmetricKey = await this.litClient.saveEncryptionKey({
        accessControlConditions,
        symmetricKey: new Uint8Array(exportedSymKey),
        authSig: this.getLitAuthSig(signature, address),
        chain: 'ethereum',
      });
      
      // 8. Criar objeto para armazenar
      this.storageData = {
        encryptedData: this.arrayBufferToBase64(encryptedApiKey),
        iv: this.arrayBufferToBase64(iv),
        litProtocolEncKey: encryptedSymmetricKey,
        accessControlConditions: JSON.stringify(accessControlConditions),
        ownerAddress: address,
        createdAt: timestamp,
        lastAccessedAt: timestamp,
        serviceType: 'whatsapp'
      };
      
      // 9. Salvar no armazenamento local
      this.saveStorageData();
      
      console.log('Chave da API do WhatsApp armazenada com sucesso');
      return true;
      
    } catch (error) {
      console.error('Erro ao armazenar chave da API:', error);
      throw error;
    }
  }

  /**
   * Recupera a chave da API do WhatsApp
   * @returns {Promise<string>} A chave da API do WhatsApp descriptografada
   */
  async retrieveAPIKey() {
    if (!this.storageData) {
      throw new Error('Nenhuma chave de API armazenada encontrada');
    }
    
    if (!this.wallet) {
      await this.connectWallet();
    }
    
    const address = await this.wallet.getAddress();
    
    // Verificar se o endereço atual corresponde ao proprietário
    if (address.toLowerCase() !== this.storageData.ownerAddress.toLowerCase()) {
      throw new Error('Esta chave pertence a outra carteira');
    }
    
    try {
      // 1. Gerar mensagem para assinatura
      const messageToSign = `Autorizo o acesso à chave de API do WhatsApp para o endereço ${address} em ${Math.floor(Date.now() / 1000)}`;
      
      // 2. Solicitar assinatura do usuário
      const signature = await this.wallet.signMessage(messageToSign);
      
      // 3. Recuperar a chave simétrica usando o Lit Protocol
      const accessControlConditions = JSON.parse(this.storageData.accessControlConditions);
      
      const symmetricKey = await this.litClient.getEncryptionKey({
        accessControlConditions,
        toDecrypt: this.storageData.litProtocolEncKey,
        chain: 'ethereum',
        authSig: this.getLitAuthSig(signature, address),
      });
      
      // 4. Importar a chave simétrica para usar na descriptografia
      const importedSymKey = await window.crypto.subtle.importKey(
        "raw",
        symmetricKey,
        {
          name: "AES-GCM",
          length: 256
        },
        false,
        ["decrypt"]
      );
      
      // 5. Descriptografar a chave da API
      const encryptedApiKey = this.base64ToArrayBuffer(this.storageData.encryptedData);
      const iv = this.base64ToArrayBuffer(this.storageData.iv);
      
      const decryptedApiKey = await window.crypto.subtle.decrypt(
        {
          name: "AES-GCM",
          iv: iv
        },
        importedSymKey,
        encryptedApiKey
      );
      
      // 6. Converter o resultado para string
      const decoder = new TextDecoder();
      const apiKey = decoder.decode(decryptedApiKey);
      
      // 7. Atualizar timestamp de último acesso
      this.storageData.lastAccessedAt = Math.floor(Date.now() / 1000);
      this.saveStorageData();
      
      console.log('Chave da API do WhatsApp recuperada com sucesso');
      return apiKey;
      
    } catch (error) {
      console.error('Erro ao recuperar chave da API:', error);
      throw error;
    }
  }

  /**
   * Deleta a chave da API do WhatsApp armazenada
   */
  async deleteAPIKey() {
    if (!this.storageData) {
      throw new Error('Nenhuma chave de API armazenada encontrada');
    }
    
    if (!this.wallet) {
      await this.connectWallet();
    }
    
    const address = await this.wallet.getAddress();
    
    // Verificar se o endereço atual corresponde ao proprietário
    if (address.toLowerCase() !== this.storageData.ownerAddress.toLowerCase()) {
      throw new Error('Esta chave pertence a outra carteira');
    }
    
    try {
      // 1. Gerar mensagem para assinatura
      const messageToSign = `Autorizo a exclusão da chave de API do WhatsApp para o endereço ${address} em ${Math.floor(Date.now() / 1000)}`;
      
      // 2. Solicitar assinatura do usuário
      const signature = await this.wallet.signMessage(messageToSign);
      
      // 3. Remover do armazenamento local
      localStorage.removeItem(this.storageKey);
      this.storageData = null;
      
      console.log('Chave da API do WhatsApp excluída com sucesso');
      return true;
      
    } catch (error) {
      console.error('Erro ao excluir chave da API:', error);
      throw error;
    }
  }

  /**
   * Verifica se existe uma chave armazenada
   * @returns {boolean} Verdadeiro se houver uma chave armazenada
   */
  hasStoredAPIKey() {
    return !!this.storageData;
  }

  /**
   * Converte um ArrayBuffer para string Base64
   * @param {ArrayBuffer} buffer - O buffer a ser convertido
   * @returns {string} String em formato Base64
   */
  arrayBufferToBase64(buffer) {
    const bytes = new Uint8Array(buffer);
    let binary = '';
    for (let i = 0; i < bytes.byteLength; i++) {
      binary += String.fromCharCode(bytes[i]);
    }
    return btoa(binary);
  }

  /**
   * Converte uma string Base64 para ArrayBuffer
   * @param {string} base64 - A string Base64 a ser convertida
   * @returns {ArrayBuffer} O buffer resultante
   */
  base64ToArrayBuffer(base64) {
    const binaryString = atob(base64);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes.buffer;
  }

  /**
   * Cria o objeto de assinatura para o Lit Protocol
   * @param {string} signature - A assinatura obtida do usuário
   * @param {string} address - O endereço da carteira
   * @returns {Object} Objeto de assinatura para o Lit Protocol
   */
  getLitAuthSig(signature, address) {
    return {
      sig: signature,
      derivedVia: "web3.eth.personal.sign",
      signedMessage: "Acesse sua chave de API do WhatsApp",
      address: address,
    };
  }

  /**
   * Salva os dados de armazenamento no localStorage
   */
  saveStorageData() {
    if (this.storageData) {
      localStorage.setItem(this.storageKey, JSON.stringify(this.storageData));
    }
  }

  /**
   * Carrega os dados de armazenamento do localStorage
   */
  loadStorageData() {
    const storedData = localStorage.getItem(this.storageKey);
    if (storedData) {
      try {
        this.storageData = JSON.parse(storedData);
      } catch (error) {
        console.error('Erro ao carregar dados armazenados:', error);
      }
    }
  }
}

// Exportar uma instância singleton do gerenciador
export const whatsappKeyManager = new WhatsAppAPIKeyManager();

// Exemplo de como usar:
/*
// Inicializar
await whatsappKeyManager.initialize();

// Conectar carteira
const address = await whatsappKeyManager.connectWallet();

// Armazenar chave da API
await whatsappKeyManager.storeAPIKey('chave-da-api-123456');

// Recuperar chave da API
const apiKey = await whatsappKeyManager.retrieveAPIKey();

// Verificar se existe chave armazenada
const hasKey = whatsappKeyManager.hasStoredAPIKey();
*/
