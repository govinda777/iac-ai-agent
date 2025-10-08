# <div align="center">🤖 IaC AI Agent</div>

<div align="center">

<h3>AI-powered agent for Infrastructure as Code analysis, security scanning, and optimization.</h3>
<h4>Featuring Web3-native authentication and on-chain payments.</h4>

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Privy](https://img.shields.io/badge/Auth-Privy.io-6366F1?style=flat&logo=ethereum)](https://privy.io)
[![Base Network](https://img.shields.io/badge/L2-Base-0052FF?style=flat&logo=coinbase)](https://base.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

</div>

<br>

<div align="center">
  <img src="https://img.shields.io/badge/Terraform-7B42BC?style=for-the-badge&logo=terraform&logoColor=white" alt="Terraform">
  <img src="https://img.shields.io/badge/AWS-FF9900?style=for-the-badge&logo=amazonaws&logoColor=white" alt="AWS">
  <img src="https://img.shields.io/badge/Azure-0078D4?style=for-the-badge&logo=microsoftazure&logoColor=white" alt="Azure">
  <img src="https://img.shields.io/badge/GCP-4285F4?style=for-the-badge&logo=googlecloud&logoColor=white" alt="GCP">
  <img src="https://img.shields.io/badge/OpenAI-412991?style=for-the-badge&logo=openai&logoColor=white" alt="OpenAI">
</div>

<br>

## 📊 Overview

The **IaC AI Agent** is an intelligent bot that analyzes Terraform code and provides:

- **✅ Security Analysis**: Integrates with Checkov to detect vulnerabilities and misconfigurations.
- **✅ LLM-Powered Insights**: Delivers contextual suggestions using GPT-4 and other large language models.
- **✅ Drift Detection**: Identifies discrepancies between your code and deployed infrastructure.
- **✅ Cost Optimization**: Recommends cost-saving opportunities with detailed estimates.
- **✅ Best Practices**: Enforces architectural patterns and industry-standard best practices.
- **✅ IAM Analysis**: Specializes in analyzing IAM policies for overly permissive or insecure configurations.

## 🚀 Getting Started

Follow these steps to set up and run the IaC AI Agent locally.

### Prerequisites

- **Go**: Version 1.21 or higher
- **Docker**: For running the application in a containerized environment
- **Node.js**: For managing smart contract dependencies
- **Privy.io Account**: For Web3-native authentication
- **OpenAI API Key**: For leveraging LLM-powered analysis

### 1. Clone the Repository

```bash
git clone https://github.com/govinda777/iac-ai-agent
cd iac-ai-agent
```

### 2. Configure Environment Variables

Copy the example environment file and add your credentials:

```bash
cp .env.example .env
```

You will need to add the following required variables to your `.env` file:

- `PRIVY_APP_ID` and `PRIVY_APP_SECRET`: Get these from your [Privy.io dashboard](https://privy.io/).
- `LLM_API_KEY`: Your [OpenAI API key](https://platform.openai.com/api-keys).
- `BASE_RPC_URL`: The RPC URL for the Base network (e.g., `https://goerli.base.org` for testnet).

### 3. Run the Application

You can run the agent using either Go or Docker.

**Using Go:**

```bash
# Install dependencies
go mod download

# Run the application
go run cmd/agent/main.go
```

**Using Docker:**

```bash
# Build and run the Docker container
docker-compose up
```

### 4. Run Tests

The project uses `godog` for BDD testing. To run the tests:

```bash
# Install Godog
go install github.com/cucumber/godog/cmd/godog@latest

# Run all tests
godog test/bdd/features/
```

## ✨ Features

The IaC AI Agent offers a robust set of features designed to improve the quality, security, and cost-effectiveness of your infrastructure.

- **Intelligent Agent System**: The agent is built with a modular and extensible architecture, allowing for different analysis "personalities" (e.g., Security Specialist, Cost Optimizer).
- **Web3-Native Integration**: Authentication is handled via [Privy.io](https://privy.io), allowing users to log in with their wallet or email. The system also supports on-chain payments and access control using NFTs on the Base network.
- **Comprehensive Security Analysis**: In addition to Checkov, the agent includes specialized analyzers for IAM policies, helping to identify and fix potential security risks.
- **Cost Optimization**: The agent provides actionable recommendations for reducing cloud costs, complete with estimated savings.
- **Drift Detection**: Keep your infrastructure in sync with your code by automatically detecting and reporting any drift.

## 🏗️ Architecture

The IaC AI Agent is built with a clean, layered architecture that separates concerns and promotes extensibility.

```mermaid
flowchart TB
    subgraph Frontend
        A[Privy SDK] --- B[Wagmi]
        B --- C[Next.js]
    end

    subgraph "Backend (Go)"
        D[API REST] --- E[Web3 Platform]
        E --- F[LLM]
        D --- G[Analyzers]
        G --- H[Knowledge Base]
    end

    subgraph "Base Network (L2)"
        I[NFT Access] --- J[IACAI Token]
    end

    Frontend --> Backend
    Backend --> "Base Network (L2)"
```

- **Backend (Go)**: The core of the application is a Go backend that exposes a REST API for handling analysis requests. It includes a suite of analyzers, an LLM client for intelligent suggestions, and a knowledge base of best practices.
- **Web3 Integration**: The agent integrates with the Base network for on-chain operations, including NFT-based access control and token payments.
- **Frontend**: A Next.js frontend provides a user-friendly interface for interacting with the agent, featuring seamless authentication via Privy.io.

## 🗺️ Roadmap

The following features are planned for future releases:

- **Preview Analysis**: Analyze `terraform plan` outputs to catch issues before they are applied.
- **Expanded IaC Support**: Add support for other IaC tools like CloudFormation and Pulumi.
- **Web Dashboard**: A comprehensive web interface for visualizing analysis results and managing agent settings.
- **CI/CD Integration**: Native integrations with popular CI/CD platforms like GitHub Actions, GitLab, and Bitbucket.