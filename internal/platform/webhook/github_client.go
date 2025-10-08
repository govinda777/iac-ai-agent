package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/govinda777/iac-ai-agent/pkg/config"
	"github.com/govinda777/iac-ai-agent/pkg/logger"
)

// GitHubClient é o cliente para interagir com a API do GitHub
type GitHubClient struct {
	config     *config.Config
	logger     *logger.Logger
	httpClient *http.Client
	token      string
	baseURL    string
}

// PRFile representa um arquivo modificado em um PR
type PRFile struct {
	SHA       string `json:"sha"`
	Filename  string `json:"filename"`
	Status    string `json:"status"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Changes   int    `json:"changes"`
	Patch     string `json:"patch"`
}

// NewGitHubClient cria uma nova instância do cliente GitHub
func NewGitHubClient(cfg *config.Config, log *logger.Logger) *GitHubClient {
	return &GitHubClient{
		config:     cfg,
		logger:     log,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      cfg.GitHub.Token,
		baseURL:    "https://api.github.com",
	}
}

// GetPullRequest busca informações de um PR
func (gc *GitHubClient) GetPullRequest(owner, repo string, prNumber int) (*PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", gc.baseURL, owner, repo, prNumber)

	resp, err := gc.doRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API retornou %d: %s", resp.StatusCode, string(body))
	}

	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &pr, nil
}

// GetPRFiles busca arquivos modificados em um PR
func (gc *GitHubClient) GetPRFiles(owner, repo string, prNumber int) ([]*PRFile, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/files", gc.baseURL, owner, repo, prNumber)

	resp, err := gc.doRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API retornou %d: %s", resp.StatusCode, string(body))
	}

	var files []*PRFile
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return files, nil
}

// GetFileContent busca conteúdo de um arquivo
func (gc *GitHubClient) GetFileContent(owner, repo, path, ref string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s?ref=%s", gc.baseURL, owner, repo, path, ref)

	resp, err := gc.doRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitHub API retornou %d: %s", resp.StatusCode, string(body))
	}

	var fileResp struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&fileResp); err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	// Decodifica base64
	if fileResp.Encoding == "base64" {
		// Implementar decodificação base64 se necessário
		return fileResp.Content, nil
	}

	return fileResp.Content, nil
}

// PostComment posta um comentário em um PR
func (gc *GitHubClient) PostComment(owner, repo string, prNumber int, body string) error {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments", gc.baseURL, owner, repo, prNumber)

	payload := map[string]string{
		"body": body,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar payload: %w", err)
	}

	resp, err := gc.doRequest("POST", url, jsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao postar comentário: %d - %s", resp.StatusCode, string(body))
	}

	gc.logger.Info("Comentário postado com sucesso", "pr", prNumber)
	return nil
}

// PostReview posta um review completo em um PR
func (gc *GitHubClient) PostReview(owner, repo string, prNumber int, event, body string, comments []ReviewComment) error {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/reviews", gc.baseURL, owner, repo, prNumber)

	payload := map[string]interface{}{
		"event": event, // APPROVE, REQUEST_CHANGES, COMMENT
		"body":  body,
	}

	if len(comments) > 0 {
		payload["comments"] = comments
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar payload: %w", err)
	}

	resp, err := gc.doRequest("POST", url, jsonData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao postar review: %d - %s", resp.StatusCode, string(body))
	}

	gc.logger.Info("Review postado com sucesso", "pr", prNumber, "event", event)
	return nil
}

// doRequest executa uma requisição HTTP para a API do GitHub
func (gc *GitHubClient) doRequest(method, url string, body []byte) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	// Headers
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "Bearer "+gc.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	gc.logger.Debug("GitHub API request", "method", method, "url", url)

	return gc.httpClient.Do(req)
}

// PullRequest representa um PR do GitHub
type PullRequest struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	State   string `json:"state"`
	HTMLURL string `json:"html_url"`
	Head    struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
}

// ReviewComment representa um comentário em linha de código
type ReviewComment struct {
	Path     string `json:"path"`
	Position int    `json:"position,omitempty"`
	Line     int    `json:"line"`
	Body     string `json:"body"`
	Side     string `json:"side,omitempty"` // LEFT, RIGHT
}
