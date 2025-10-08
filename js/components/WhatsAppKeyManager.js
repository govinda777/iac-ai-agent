/**
 * WhatsApp Key Manager Component
 * 
 * Componente React para gerenciar a chave da API do WhatsApp de forma segura
 * utilizando Web3 e Lit Protocol para criptografia.
 */

import React, { useState, useEffect } from 'react';
import { whatsappKeyManager } from '../whatsapp-api-key.js';

const WhatsAppKeyManager = () => {
  const [isInitialized, setIsInitialized] = useState(false);
  const [walletAddress, setWalletAddress] = useState('');
  const [apiKey, setApiKey] = useState('');
  const [hasStoredKey, setHasStoredKey] = useState(false);
  const [retrievedKey, setRetrievedKey] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  // Inicializar o gerenciador ao carregar o componente
  useEffect(() => {
    const init = async () => {
      try {
        await whatsappKeyManager.initialize();
        setIsInitialized(true);
        setHasStoredKey(whatsappKeyManager.hasStoredAPIKey());
      } catch (err) {
        setError('Erro ao inicializar: ' + err.message);
      }
    };
    
    init();
  }, []);

  // Função para conectar a carteira
  const handleConnectWallet = async () => {
    setIsLoading(true);
    setError('');
    setSuccess('');
    
    try {
      const address = await whatsappKeyManager.connectWallet();
      setWalletAddress(address);
      setHasStoredKey(whatsappKeyManager.hasStoredAPIKey());
      setSuccess('Carteira conectada com sucesso');
    } catch (err) {
      setError('Erro ao conectar carteira: ' + err.message);
    } finally {
      setIsLoading(false);
    }
  };

  // Função para armazenar a chave
  const handleStoreKey = async () => {
    if (!apiKey.trim()) {
      setError('Por favor, insira uma chave de API');
      return;
    }
    
    setIsLoading(true);
    setError('');
    setSuccess('');
    
    try {
      await whatsappKeyManager.storeAPIKey(apiKey);
      setHasStoredKey(true);
      setSuccess('Chave armazenada com sucesso');
      setApiKey('');
    } catch (err) {
      setError('Erro ao armazenar chave: ' + err.message);
    } finally {
      setIsLoading(false);
    }
  };

  // Função para recuperar a chave
  const handleRetrieveKey = async () => {
    setIsLoading(true);
    setError('');
    setSuccess('');
    setRetrievedKey('');
    
    try {
      const key = await whatsappKeyManager.retrieveAPIKey();
      setRetrievedKey(key);
      setSuccess('Chave recuperada com sucesso');
    } catch (err) {
      setError('Erro ao recuperar chave: ' + err.message);
    } finally {
      setIsLoading(false);
    }
  };

  // Função para excluir a chave
  const handleDeleteKey = async () => {
    setIsLoading(true);
    setError('');
    setSuccess('');
    
    try {
      await whatsappKeyManager.deleteAPIKey();
      setHasStoredKey(false);
      setRetrievedKey('');
      setSuccess('Chave excluída com sucesso');
    } catch (err) {
      setError('Erro ao excluir chave: ' + err.message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="whatsapp-key-manager">
      <h2>Gerenciador de Chave da API do WhatsApp</h2>
      
      {!isInitialized ? (
        <p>Inicializando...</p>
      ) : (
        <div className="key-manager-container">
          {/* Status da carteira */}
          <div className="wallet-status">
            {walletAddress ? (
              <p>Carteira conectada: {walletAddress.substring(0, 6)}...{walletAddress.substring(walletAddress.length - 4)}</p>
            ) : (
              <button 
                className="connect-wallet-btn" 
                onClick={handleConnectWallet}
                disabled={isLoading}
              >
                Conectar Carteira
              </button>
            )}
          </div>
          
          {/* Armazenar chave */}
          {walletAddress && (
            <div className="api-key-form">
              <h3>Armazenar Chave da API</h3>
              <div className="input-group">
                <input
                  type="password"
                  value={apiKey}
                  onChange={(e) => setApiKey(e.target.value)}
                  placeholder="Insira a chave da API do WhatsApp"
                  disabled={isLoading}
                />
                <button 
                  onClick={handleStoreKey}
                  disabled={isLoading || !apiKey.trim()}
                >
                  Armazenar
                </button>
              </div>
            </div>
          )}
          
          {/* Recuperar ou excluir chave */}
          {walletAddress && hasStoredKey && (
            <div className="stored-key-actions">
              <h3>Gerenciar Chave Armazenada</h3>
              <div className="button-group">
                <button 
                  onClick={handleRetrieveKey}
                  disabled={isLoading}
                >
                  Recuperar Chave
                </button>
                <button 
                  className="delete-btn"
                  onClick={handleDeleteKey}
                  disabled={isLoading}
                >
                  Excluir Chave
                </button>
              </div>
              
              {retrievedKey && (
                <div className="retrieved-key">
                  <h4>Chave Recuperada:</h4>
                  <div className="key-display">
                    <code>{retrievedKey}</code>
                  </div>
                </div>
              )}
            </div>
          )}
          
          {/* Mensagens de feedback */}
          {error && <div className="error-message">{error}</div>}
          {success && <div className="success-message">{success}</div>}
          {isLoading && <div className="loading-indicator">Processando...</div>}
        </div>
      )}

      {/* Informações sobre segurança */}
      <div className="security-info">
        <h3>Como funciona?</h3>
        <ul>
          <li>Sua chave é criptografada localmente usando AES-256</li>
          <li>A chave de criptografia é protegida pelo Lit Protocol</li>
          <li>Apenas você, com sua carteira Web3, pode descriptografar a chave</li>
          <li>Nenhuma informação sensível é armazenada em texto simples</li>
        </ul>
      </div>
    </div>
  );
};

export default WhatsAppKeyManager;
