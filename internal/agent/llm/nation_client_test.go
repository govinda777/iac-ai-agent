package llm

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/govinda777/iac-ai-agent/internal/models"
	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
	. "github.com/onsi/gomega"
)

func TestNewNationClient(t *testing.T) {
	g := NewGomegaWithT(t)
	log := logger.New("debug", "text")

	t.Run("should return error if api key is missing", func(t *testing.T) {
		cfg := &config.Config{
			LLM:  config.LLMConfig{APIKey: ""},
			Web3: config.Web3Config{NFTAccessContractAddress: "0x123", WalletToken: "abc"},
		}
		client, err := NewNationClient(cfg, log)
		g.Expect(err).To(HaveOccurred())
		g.Expect(client).To(BeNil())
		g.Expect(err.Error()).To(Equal("Nation.fun API key não configurada"))
	})

	t.Run("should return error if nft address is missing", func(t *testing.T) {
		cfg := &config.Config{
			LLM:  config.LLMConfig{APIKey: "test-key"},
			Web3: config.Web3Config{NFTAccessContractAddress: "", WalletToken: "abc"},
		}
		client, err := NewNationClient(cfg, log)
		g.Expect(err).To(HaveOccurred())
		g.Expect(client).To(BeNil())
		g.Expect(err.Error()).To(Equal("Nation.fun NFT contract address não configurado"))
	})

	t.Run("should return error if wallet token is missing", func(t *testing.T) {
		cfg := &config.Config{
			LLM:  config.LLMConfig{APIKey: "test-key"},
			Web3: config.Web3Config{NFTAccessContractAddress: "0x123", WalletToken: ""},
		}
		client, err := NewNationClient(cfg, log)
		g.Expect(err).To(HaveOccurred())
		g.Expect(client).To(BeNil())
		g.Expect(err.Error()).To(Equal("Nation.fun wallet token não configurado"))
	})

	t.Run("should return client on valid config", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		cfg := &config.Config{
			LLM:  config.LLMConfig{APIKey: "test-key", BaseURL: server.URL},
			Web3: config.Web3Config{NFTAccessContractAddress: "0x123", WalletToken: "abc"},
		}

		client, err := NewNationClient(cfg, log)
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(client).ToNot(BeNil())
	})
}

func TestNationClient_ValidateConnection(t *testing.T) {
	g := NewGomegaWithT(t)
	log := logger.New("debug", "text")

	t.Run("should return nil on successful connection", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			g.Expect(r.URL.Path).To(Equal("/validate"))
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := &NationClient{
			apiKey:      "test-key",
			nftAddress:  "0x123",
			walletToken: "abc",
			baseURL:     server.URL,
			httpClient:  server.Client(),
			logger:      log,
		}

		err := client.ValidateConnection()
		g.Expect(err).ToNot(HaveOccurred())
	})

	t.Run("should return error on failed connection", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := &NationClient{
			apiKey:      "test-key",
			nftAddress:  "0x123",
			walletToken: "abc",
			baseURL:     server.URL,
			httpClient:  server.Client(),
			logger:      log,
		}

		err := client.ValidateConnection()
		g.Expect(err).To(HaveOccurred())
	})
}

func TestNationClient_Generate(t *testing.T) {
	g := NewGomegaWithT(t)
	log := logger.New("debug", "text")

	t.Run("should generate content successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			g.Expect(r.URL.Path).To(Equal("/completions"))

			var reqPayload map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&reqPayload)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(reqPayload["prompt"]).To(Equal("test prompt"))

			respPayload := map[string]interface{}{
				"content":     "test response",
				"model":       "test-model",
				"tokens_used": 10,
				"nft_used":    "0x123",
			}
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(respPayload)
			g.Expect(err).ToNot(HaveOccurred())
		}))
		defer server.Close()

		cfg := &config.Config{Web3: config.Web3Config{DefaultAgentAddress: "0xagent"}}

		client := &NationClient{
			config:      cfg,
			apiKey:      "test-key",
			nftAddress:  "0x123",
			walletToken: "abc",
			baseURL:     server.URL,
			httpClient:  server.Client(),
			logger:      log,
			modelName:   "test-model",
		}

		req := &models.LLMRequest{Prompt: "test prompt"}
		resp, err := client.Generate(req)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(resp.Content).To(Equal("test response"))
		g.Expect(resp.TokensUsed).To(Equal(10))
	})
}

func TestNationClient_GenerateStructured(t *testing.T) {
	g := NewGomegaWithT(t)
	log := logger.New("debug", "text")

	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	t.Run("should generate structured content successfully", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			g.Expect(r.URL.Path).To(Equal("/completions"))

			var reqPayload map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&reqPayload)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(reqPayload["response_format"]).To(Equal("json"))

			structuredResponse := TestStruct{Name: "test", Value: 123}
			responseBytes, _ := json.Marshal(structuredResponse)
			respPayload := map[string]interface{}{
				"content": string(responseBytes),
			}
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(respPayload)
			g.Expect(err).ToNot(HaveOccurred())
		}))
		defer server.Close()

		cfg := &config.Config{Web3: config.Web3Config{DefaultAgentAddress: "0xagent"}}

		client := &NationClient{
			config:      cfg,
			apiKey:      "test-key",
			nftAddress:  "0x123",
			walletToken: "abc",
			baseURL:     server.URL,
			httpClient:  server.Client(),
			logger:      log,
			modelName:   "test-model",
		}

		req := &models.LLMRequest{Prompt: "generate json"}
		var result TestStruct
		err := client.GenerateStructured(req, &result)

		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(result.Name).To(Equal("test"))
		g.Expect(result.Value).To(Equal(123))
	})
}
