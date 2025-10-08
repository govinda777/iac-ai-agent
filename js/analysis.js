// Analysis and Suggestions Flow for IaC AI Agent
document.addEventListener('DOMContentLoaded', function() {
    // Configuration
    const isMockMode = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1';
    const API_BASE_URL = isMockMode ? '/api/mock' : 'https://api.iacai.com/v1';
    
    // Elements
    const analysisForm = document.getElementById('analysis-form');
    const codeInput = document.getElementById('terraform-code');
    const analysisTypeSelect = document.getElementById('analysis-type');
    const submitButton = document.getElementById('submit-analysis');
    const resultsContainer = document.getElementById('results-container');
    const analysisStatus = document.getElementById('analysis-status');
    const loadingIndicator = document.getElementById('loading-indicator');
    
    // Analysis types and costs
    const analysisTypes = {
        'basic': {
            id: 'basic',
            name: 'Análise Básica',
            description: 'Análise de código Terraform com sugestões simples',
            tokenCost: 1,
            requiredTier: 'basic'
        },
        'security': {
            id: 'security',
            name: 'Análise de Segurança',
            description: 'Scan completo com Checkov para problemas de segurança',
            tokenCost: 2,
            requiredTier: 'basic'
        },
        'llm': {
            id: 'llm',
            name: 'Análise com LLM',
            description: 'Análise inteligente com GPT-4/Claude e sugestões detalhadas',
            tokenCost: 5,
            requiredTier: 'pro'
        },
        'preview': {
            id: 'preview',
            name: 'Preview Analysis',
            description: 'Análise de terraform plan para avaliar mudanças',
            tokenCost: 3,
            requiredTier: 'pro'
        },
        'cost': {
            id: 'cost',
            name: 'Otimização de Custos',
            description: 'Sugestões para redução de custos e estimativas',
            tokenCost: 5,
            requiredTier: 'pro'
        },
        'full': {
            id: 'full',
            name: 'Full Review',
            description: 'Análise completa incluindo todas as verificações',
            tokenCost: 15,
            requiredTier: 'pro'
        }
    };
    
    // State
    let currentAnalysis = null;
    let isAnalysisInProgress = false;
    
    // Initialize
    function init() {
        updateAnalysisForm();
        
        // Setup analysis type options
        if (analysisTypeSelect) {
            populateAnalysisTypes();
            
            analysisTypeSelect.addEventListener('change', function() {
                updateTokenCost();
            });
        }
        
        // Setup form submission
        if (analysisForm) {
            analysisForm.addEventListener('submit', function(e) {
                e.preventDefault();
                submitAnalysis();
            });
        }
        
        // Listen for authentication events
        window.addEventListener('user:authenticated', function(e) {
            updateAnalysisForm();
        });
        
        window.addEventListener('user:logout', function() {
            updateAnalysisForm();
        });
        
        // Listen for token/NFT events
        window.addEventListener('nft:purchased', function(e) {
            updateAnalysisForm();
            populateAnalysisTypes(); // Refresh based on tier
        });
        
        window.addEventListener('tokens:purchased', function(e) {
            updateAnalysisForm();
        });
    }
    
    function populateAnalysisTypes() {
        if (!analysisTypeSelect) return;
        
        // Clear current options
        analysisTypeSelect.innerHTML = '';
        
        // Get user's tier
        const userTier = getUserTier();
        
        // Add available analysis types based on tier
        Object.values(analysisTypes).forEach(type => {
            const tierValid = isTierValid(userTier, type.requiredTier);
            
            const option = document.createElement('option');
            option.value = type.id;
            option.textContent = `${type.name} (${type.tokenCost} tokens)`;
            
            if (!tierValid) {
                option.disabled = true;
                option.textContent += ` - Requer ${type.requiredTier} ou superior`;
            }
            
            analysisTypeSelect.appendChild(option);
        });
        
        // Set default and trigger change
        if (analysisTypeSelect.options.length > 0) {
            analysisTypeSelect.selectedIndex = 0;
            updateTokenCost();
        }
    }
    
    function updateTokenCost() {
        if (!analysisTypeSelect) return;
        
        const selectedType = analysisTypeSelect.value;
        const cost = analysisTypes[selectedType]?.tokenCost || 0;
        
        const costElement = document.getElementById('token-cost');
        if (costElement) {
            costElement.textContent = cost;
        }
    }
    
    function updateAnalysisForm() {
        // Check authentication
        const isAuthenticated = window.IaCAuth && window.IaCAuth.isAuthenticated();
        
        if (!isAuthenticated) {
            disableAnalysisForm('Conecte sua wallet para analisar código');
            return;
        }
        
        // Check NFT access
        const userTier = getUserTier();
        if (!userTier) {
            disableAnalysisForm('É necessário ter um NFT de acesso para analisar código');
            return;
        }
        
        // Check token balance
        const tokenBalance = getTokenBalance();
        const minCost = getMinimumTokenCost();
        
        if (tokenBalance < minCost) {
            disableAnalysisForm(`Saldo insuficiente. Mínimo necessário: ${minCost} tokens`);
            return;
        }
        
        // Enable form
        enableAnalysisForm();
    }
    
    function disableAnalysisForm(message) {
        if (submitButton) {
            submitButton.disabled = true;
            submitButton.querySelector('.button-text').textContent = message || 'Indisponível';
        }
        
        if (analysisForm) {
            analysisForm.classList.add('disabled');
        }
    }
    
    function enableAnalysisForm() {
        if (submitButton) {
            submitButton.disabled = isAnalysisInProgress;
            submitButton.querySelector('.button-text').textContent = 'Analisar Código';
        }
        
        if (analysisForm) {
            analysisForm.classList.remove('disabled');
        }
    }
    
    function getUserTier() {
        // Get tier from NFT module if available
        return window.IaCNFT ? window.IaCNFT.getCurrentTier() : null;
    }
    
    function getTokenBalance() {
        // Get token balance from token module if available
        return window.IaCToken ? window.IaCToken.getTokenBalance() : 0;
    }
    
    function getMinimumTokenCost() {
        // Get the cost of the cheapest analysis type
        return Math.min(...Object.values(analysisTypes).map(type => type.tokenCost));
    }
    
    function isTierValid(userTier, requiredTier) {
        if (!userTier) return false;
        
        // Simple tier hierarchy
        const tiers = {
            'basic': 1,
            'pro': 2,
            'enterprise': 3
        };
        
        return tiers[userTier] >= tiers[requiredTier];
    }
    
    // Analysis submission
    async function submitAnalysis() {
        if (!codeInput || !codeInput.value.trim()) {
            showMessage('Por favor, insira código Terraform para análise', 'error');
            return;
        }
        
        const code = codeInput.value.trim();
        const analysisType = analysisTypeSelect ? analysisTypeSelect.value : 'basic';
        
        // Verify user tier for analysis type
        const userTier = getUserTier();
        if (!isTierValid(userTier, analysisTypes[analysisType].requiredTier)) {
            showMessage(`Seu tier (${userTier}) não permite este tipo de análise. Necessário: ${analysisTypes[analysisType].requiredTier} ou superior.`, 'error');
            return;
        }
        
        // Verify token balance
        const tokenCost = analysisTypes[analysisType].tokenCost;
        const tokenBalance = getTokenBalance();
        
        if (tokenBalance < tokenCost) {
            showMessage(`Saldo insuficiente de tokens. Necessário: ${tokenCost}, Disponível: ${tokenBalance}`, 'error');
            return;
        }
        
        try {
            isAnalysisInProgress = true;
            updateAnalysisForm();
            showLoading(true);
            
            // Deduct tokens
            if (window.IaCToken) {
                const spent = window.IaCToken.spendTokens(tokenCost, `Análise ${analysisTypes[analysisType].name}`);
                if (!spent) {
                    throw new Error('Falha ao debitar tokens');
                }
            }
            
            // Prepare request
            const request = {
                code: code,
                type: analysisType,
                walletAddress: window.IaCAuth ? window.IaCAuth.getWalletAddress() : null,
                tier: userTier
            };
            
            let results;
            if (isMockMode) {
                // Mock analysis for development
                results = await simulateAnalysis(request);
            } else {
                // Real analysis via API
                results = await performAnalysis(request);
            }
            
            // Store and display results
            currentAnalysis = {
                id: results.id || `analysis_${Date.now()}`,
                timestamp: results.timestamp || Date.now(),
                type: analysisType,
                request: request,
                results: results
            };
            
            // Save to history in localStorage
            saveToAnalysisHistory(currentAnalysis);
            
            // Display results
            displayResults(currentAnalysis);
            
        } catch (error) {
            console.error('Error during analysis:', error);
            showMessage(`Erro na análise: ${error.message}`, 'error');
        } finally {
            isAnalysisInProgress = false;
            updateAnalysisForm();
            showLoading(false);
        }
    }
    
    // API interaction
    async function performAnalysis(request) {
        // Actual API call
        const response = await fetch(`${API_BASE_URL}/analyze`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(request),
        });
        
        if (!response.ok) {
            throw new Error(`API error: ${response.status} ${response.statusText}`);
        }
        
        return await response.json();
    }
    
    // Mock implementation for development
    async function simulateAnalysis(request) {
        const startTime = Date.now();
        
        // Show status updates
        updateStatus('Iniciando análise...');
        await sleep(500);
        
        updateStatus('Verificando código...');
        await sleep(1000);
        
        if (request.type === 'security' || request.type === 'full') {
            updateStatus('Executando checkov...');
            await sleep(1500);
        }
        
        if (request.type === 'llm' || request.type === 'full') {
            updateStatus('Consultando LLM...');
            await sleep(2000);
        }
        
        updateStatus('Gerando sugestões...');
        await sleep(1000);
        
        // Generate mock results
        let score = 70 + Math.floor(Math.random() * 20); // 70-90
        
        // Generate issues based on code
        const issues = [];
        
        // Sample issue detection
        if (request.code.includes('aws_s3_bucket') && !request.code.includes('acl = "private"')) {
            issues.push({
                severity: 'HIGH',
                title: 'S3 Bucket Publicly Accessible',
                message: 'O bucket S3 pode estar configurado para acesso público',
                resource: 'aws_s3_bucket',
                remediation: 'Adicione acl = "private" e configure block_public_access',
                code: 'resource "aws_s3_bucket" "example" {\n  bucket = "my-bucket"\n  acl    = "private"\n}'
            });
        }
        
        if (request.code.includes('aws_instance') && !request.code.includes('monitoring')) {
            issues.push({
                severity: 'MEDIUM',
                title: 'Monitoring Disabled',
                message: 'EC2 instance não tem monitoring habilitado',
                resource: 'aws_instance',
                remediation: 'Adicione monitoring = true',
                code: 'resource "aws_instance" "example" {\n  monitoring = true\n}'
            });
        }
        
        // Always add a best practice issue
        issues.push({
            severity: 'LOW',
            title: 'Tag Missing',
            message: 'Resources should have tags for better organization',
            resource: 'all',
            remediation: 'Add tags to all resources',
            code: 'tags = {\n  Environment = "dev"\n  Project     = "example"\n}'
        });
        
        // If LLM analysis, include more detailed suggestions
        let llmAnalysis = null;
        if (request.type === 'llm' || request.type === 'full') {
            llmAnalysis = {
                summary: "O código Terraform apresenta algumas oportunidades de melhoria em relação à segurança e boas práticas.",
                suggestions: [
                    {
                        title: "Implementar configurações de segurança para S3",
                        description: "Os buckets S3 devem ter configurações de segurança apropriadas para evitar exposição acidental de dados.",
                        code: "resource \"aws_s3_bucket_public_access_block\" \"example\" {\n  bucket = aws_s3_bucket.example.id\n  block_public_acls   = true\n  block_public_policy = true\n  ignore_public_acls  = true\n  restrict_public_buckets = true\n}"
                    },
                    {
                        title: "Usar módulos para padronização",
                        description: "Considere usar módulos para melhorar a reutilização e padronização do código.",
                        code: "module \"s3_bucket\" {\n  source  = \"terraform-aws-modules/s3-bucket/aws\"\n  version = \"~> 3.0\"\n  \n  bucket = \"my-bucket\"\n  acl    = \"private\"\n}"
                    }
                ],
                architecture_review: "A arquitetura atual é simples mas pode ser melhorada com o uso de módulos e padrões consistentes."
            };
        }
        
        // Generate mock response
        return {
            id: `analysis_${Date.now()}`,
            timestamp: Date.now(),
            duration_ms: Date.now() - startTime,
            score: score,
            issues: issues,
            summary: `Encontrados ${issues.length} problemas (${issues.filter(i => i.severity === 'HIGH').length} críticos)`,
            llm_analysis: llmAnalysis,
            sections: [
                "Executive Summary",
                "Security Issues",
                "Best Practices",
                ...(request.type === 'cost' || request.type === 'full' ? ["Cost Optimization"] : []),
                "Detailed Findings"
            ]
        };
    }
    
    // Results display
    function displayResults(analysis) {
        if (!resultsContainer) return;
        
        // Clear existing results
        resultsContainer.innerHTML = '';
        
        // Add results header
        const header = document.createElement('div');
        header.className = 'results-header';
        header.innerHTML = `
            <h2>Análise ${analysisTypes[analysis.type].name}</h2>
            <div class="results-meta">
                <span class="results-time">${formatDateTime(analysis.timestamp)}</span>
                <span class="results-score">Score: ${analysis.results.score}/100</span>
            </div>
        `;
        resultsContainer.appendChild(header);
        
        // Add summary section
        const summary = document.createElement('div');
        summary.className = 'results-summary';
        summary.innerHTML = `
            <h3>Executive Summary</h3>
            <p>${analysis.results.summary}</p>
        `;
        resultsContainer.appendChild(summary);
        
        // Add issues section
        const issuesSection = document.createElement('div');
        issuesSection.className = 'results-issues';
        issuesSection.innerHTML = `
            <h3>Issues Found</h3>
        `;
        
        // Group issues by severity
        const issuesBySeverity = {
            HIGH: [],
            MEDIUM: [],
            LOW: []
        };
        
        analysis.results.issues.forEach(issue => {
            if (issuesBySeverity[issue.severity]) {
                issuesBySeverity[issue.severity].push(issue);
            } else {
                issuesBySeverity.LOW.push(issue);
            }
        });
        
        // Create issue list
        const issuesList = document.createElement('div');
        issuesList.className = 'issues-list';
        
        // Add high severity issues
        if (issuesBySeverity.HIGH.length > 0) {
            const highSeverity = document.createElement('div');
            highSeverity.className = 'severity-group severity-high';
            highSeverity.innerHTML = `
                <h4>Critical Issues (${issuesBySeverity.HIGH.length})</h4>
            `;
            
            issuesBySeverity.HIGH.forEach(issue => {
                highSeverity.appendChild(createIssueCard(issue));
            });
            
            issuesList.appendChild(highSeverity);
        }
        
        // Add medium severity issues
        if (issuesBySeverity.MEDIUM.length > 0) {
            const mediumSeverity = document.createElement('div');
            mediumSeverity.className = 'severity-group severity-medium';
            mediumSeverity.innerHTML = `
                <h4>Medium Issues (${issuesBySeverity.MEDIUM.length})</h4>
            `;
            
            issuesBySeverity.MEDIUM.forEach(issue => {
                mediumSeverity.appendChild(createIssueCard(issue));
            });
            
            issuesList.appendChild(mediumSeverity);
        }
        
        // Add low severity issues
        if (issuesBySeverity.LOW.length > 0) {
            const lowSeverity = document.createElement('div');
            lowSeverity.className = 'severity-group severity-low';
            lowSeverity.innerHTML = `
                <h4>Best Practices (${issuesBySeverity.LOW.length})</h4>
            `;
            
            issuesBySeverity.LOW.forEach(issue => {
                lowSeverity.appendChild(createIssueCard(issue));
            });
            
            issuesList.appendChild(lowSeverity);
        }
        
        issuesSection.appendChild(issuesList);
        resultsContainer.appendChild(issuesSection);
        
        // Add LLM analysis section if available
        if (analysis.results.llm_analysis) {
            const llmSection = document.createElement('div');
            llmSection.className = 'results-llm';
            llmSection.innerHTML = `
                <h3>LLM Analysis</h3>
                <div class="llm-summary">
                    <p>${analysis.results.llm_analysis.summary}</p>
                </div>
                <h4>Advanced Suggestions</h4>
            `;
            
            const suggestionsList = document.createElement('div');
            suggestionsList.className = 'suggestions-list';
            
            analysis.results.llm_analysis.suggestions.forEach(suggestion => {
                const suggestionCard = document.createElement('div');
                suggestionCard.className = 'suggestion-card';
                suggestionCard.innerHTML = `
                    <h5>${suggestion.title}</h5>
                    <p>${suggestion.description}</p>
                    <div class="code-sample">
                        <pre><code>${escapeHtml(suggestion.code)}</code></pre>
                    </div>
                `;
                suggestionsList.appendChild(suggestionCard);
            });
            
            llmSection.appendChild(suggestionsList);
            
            if (analysis.results.llm_analysis.architecture_review) {
                const architectureReview = document.createElement('div');
                architectureReview.className = 'architecture-review';
                architectureReview.innerHTML = `
                    <h4>Architecture Review</h4>
                    <p>${analysis.results.llm_analysis.architecture_review}</p>
                `;
                llmSection.appendChild(architectureReview);
            }
            
            resultsContainer.appendChild(llmSection);
        }
        
        // Add export buttons
        const exportSection = document.createElement('div');
        exportSection.className = 'results-export';
        exportSection.innerHTML = `
            <button class="btn btn-outline export-json">Export JSON</button>
            <button class="btn btn-outline export-pdf">Export PDF</button>
        `;
        resultsContainer.appendChild(exportSection);
        
        // Add event listeners for export buttons
        exportSection.querySelector('.export-json').addEventListener('click', () => {
            exportAnalysisAsJSON(analysis);
        });
        
        exportSection.querySelector('.export-pdf').addEventListener('click', () => {
            exportAnalysisAsPDF(analysis);
        });
        
        // Show results
        resultsContainer.style.display = 'block';
    }
    
    function createIssueCard(issue) {
        const card = document.createElement('div');
        card.className = `issue-card issue-${issue.severity.toLowerCase()}`;
        card.innerHTML = `
            <div class="issue-header">
                <div class="issue-severity">${issue.severity}</div>
                <div class="issue-title">${issue.title}</div>
            </div>
            <div class="issue-content">
                <p>${issue.message}</p>
                <div class="issue-resource">Resource: <code>${issue.resource}</code></div>
                <div class="issue-remediation">
                    <p><strong>Sugestão de correção:</strong></p>
                    <pre><code>${escapeHtml(issue.code)}</code></pre>
                </div>
            </div>
        `;
        return card;
    }
    
    // Export functions
    function exportAnalysisAsJSON(analysis) {
        const dataStr = JSON.stringify(analysis, null, 2);
        const blob = new Blob([dataStr], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = `iacai_analysis_${analysis.id}.json`;
        link.click();
        URL.revokeObjectURL(url);
    }
    
    function exportAnalysisAsPDF(analysis) {
        // In a real implementation, this would generate a PDF
        // For now, just show a message
        showMessage('Export to PDF functionality is coming soon!', 'info');
    }
    
    // History functions
    function saveToAnalysisHistory(analysis) {
        const history = JSON.parse(localStorage.getItem('analysisHistory') || '[]');
        
        // Add to history (limited to last 10)
        history.unshift({
            id: analysis.id,
            timestamp: analysis.timestamp,
            type: analysis.type,
            score: analysis.results.score
        });
        
        if (history.length > 10) {
            history.pop();
        }
        
        localStorage.setItem('analysisHistory', JSON.stringify(history));
    }
    
    // Utility functions
    function updateStatus(message) {
        if (analysisStatus) {
            analysisStatus.textContent = message;
        }
    }
    
    function showLoading(isLoading) {
        if (loadingIndicator) {
            loadingIndicator.style.display = isLoading ? 'block' : 'none';
        }
    }
    
    function showMessage(message, type = 'info') {
        console.log(`[${type.toUpperCase()}] ${message}`);
        
        // Check if message container exists, create if not
        let msgContainer = document.querySelector('.message-container');
        if (!msgContainer) {
            msgContainer = document.createElement('div');
            msgContainer.className = 'message-container';
            document.body.appendChild(msgContainer);
        }
        
        // Create message element
        const msgElement = document.createElement('div');
        msgElement.className = `message message-${type}`;
        msgElement.innerHTML = `
            <div class="message-content">${message}</div>
            <button class="message-close">&times;</button>
        `;
        
        // Add to container
        msgContainer.appendChild(msgElement);
        
        // Auto-remove after delay
        setTimeout(() => {
            msgElement.classList.add('fade-out');
            setTimeout(() => msgElement.remove(), 500);
        }, 5000);
        
        // Close button
        msgElement.querySelector('.message-close').addEventListener('click', () => {
            msgElement.classList.add('fade-out');
            setTimeout(() => msgElement.remove(), 500);
        });
    }
    
    function formatDateTime(timestamp) {
        const date = new Date(timestamp);
        return date.toLocaleString();
    }
    
    function escapeHtml(text) {
        return text
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }
    
    function sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
    
    // Export methods for other scripts
    window.IaCAnalysis = {
        submitAnalysis,
        getAnalysisTypes: () => ({ ...analysisTypes }),
        getCurrentAnalysis: () => currentAnalysis,
        getAnalysisHistory: () => JSON.parse(localStorage.getItem('analysisHistory') || '[]')
    };
    
    // Initialize
    init();
});
