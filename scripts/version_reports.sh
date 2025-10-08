#!/bin/bash

# Script para versionar relatórios de qualidade
# Este script deve ser executado após a geração dos relatórios

set -e

REPORTS_DIR="reports/html"
GIT_REPORTS_DIR="reports/versioned"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BRANCH=$(git rev-parse --abbrev-ref HEAD)
COMMIT=$(git rev-parse --short HEAD)

echo "📊 Versionando relatórios de qualidade..."
echo "Branch: $BRANCH"
echo "Commit: $COMMIT"
echo "Timestamp: $TIMESTAMP"

# Criar diretório de relatórios versionados se não existir
mkdir -p "$GIT_REPORTS_DIR"

# Copiar relatórios mais recentes para diretório versionado
if [ -d "$REPORTS_DIR" ]; then
    echo "📁 Copiando relatórios para versionamento..."
    
    # Criar diretório específico para este commit
    VERSION_DIR="$GIT_REPORTS_DIR/${BRANCH}_${COMMIT}_${TIMESTAMP}"
    mkdir -p "$VERSION_DIR"
    
    # Copiar todos os relatórios HTML
    cp "$REPORTS_DIR"/*.html "$VERSION_DIR/" 2>/dev/null || echo "Nenhum relatório HTML encontrado"
    cp "$REPORTS_DIR"/*.json "$VERSION_DIR/" 2>/dev/null || echo "Nenhum relatório JSON encontrado"
    
    # Criar arquivo de metadados
    cat > "$VERSION_DIR/metadata.json" << EOF
{
  "branch": "$BRANCH",
  "commit": "$COMMIT",
  "timestamp": "$TIMESTAMP",
  "date": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "author": "$(git log -1 --pretty=format:%an)",
  "message": "$(git log -1 --pretty=format:%s)",
  "files": [
    $(ls "$VERSION_DIR"/*.html "$VERSION_DIR"/*.json 2>/dev/null | sed 's/^/    "/' | sed 's/$/",/' | sed '$s/,$//')
  ]
}
EOF
    
    echo "✅ Relatórios versionados em: $VERSION_DIR"
    
    # Criar link simbólico para o mais recente
    LATEST_LINK="$GIT_REPORTS_DIR/latest"
    rm -f "$LATEST_LINK"
    ln -sf "$(basename "$VERSION_DIR")" "$LATEST_LINK"
    
    echo "🔗 Link 'latest' atualizado para: $(basename "$VERSION_DIR")"
    
    # Criar índice HTML para navegação
    generate_index_html
    
else
    echo "❌ Diretório de relatórios não encontrado: $REPORTS_DIR"
    echo "💡 Execute os testes primeiro para gerar relatórios"
    exit 1
fi

echo ""
echo "📊 Relatórios versionados com sucesso!"
echo "📁 Localização: $VERSION_DIR"
echo "🔗 Acesso rápido: $GIT_REPORTS_DIR/latest"

# Função para gerar índice HTML
generate_index_html() {
    cat > "$GIT_REPORTS_DIR/index.html" << 'EOF'
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Relatórios de Qualidade - Índice</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: rgba(255, 255, 255, 0.95);
            border-radius: 15px;
            padding: 30px;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
        }
        h1 {
            color: #2c3e50;
            text-align: center;
            margin-bottom: 30px;
        }
        .report-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
            gap: 20px;
        }
        .report-card {
            border: 1px solid #ecf0f1;
            border-radius: 10px;
            padding: 20px;
            background: #f8f9fa;
            transition: transform 0.3s ease;
        }
        .report-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
        }
        .report-title {
            font-weight: bold;
            color: #2c3e50;
            margin-bottom: 10px;
        }
        .report-meta {
            color: #7f8c8d;
            font-size: 0.9em;
            margin-bottom: 15px;
        }
        .report-links {
            display: flex;
            gap: 10px;
        }
        .report-link {
            padding: 8px 16px;
            background: #667eea;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            font-size: 0.9em;
            transition: background 0.3s ease;
        }
        .report-link:hover {
            background: #764ba2;
        }
        .latest-badge {
            background: #27ae60;
            color: white;
            padding: 4px 8px;
            border-radius: 15px;
            font-size: 0.8em;
            font-weight: bold;
            margin-left: 10px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>📊 Relatórios de Qualidade - Índice</h1>
        <div class="report-grid" id="reports-grid">
            <!-- Relatórios serão inseridos aqui via JavaScript -->
        </div>
    </div>
    
    <script>
        // Carregar lista de relatórios
        fetch('reports.json')
            .then(response => response.json())
            .then(data => {
                const grid = document.getElementById('reports-grid');
                data.reports.forEach(report => {
                    const card = document.createElement('div');
                    card.className = 'report-card';
                    card.innerHTML = `
                        <div class="report-title">
                            ${report.branch} - ${report.commit}
                            ${report.isLatest ? '<span class="latest-badge">LATEST</span>' : ''}
                        </div>
                        <div class="report-meta">
                            📅 ${new Date(report.date).toLocaleString('pt-BR')}<br>
                            👤 ${report.author}<br>
                            💬 ${report.message}
                        </div>
                        <div class="report-links">
                            <a href="${report.path}/" class="report-link">📁 Ver Relatórios</a>
                            <a href="${report.path}/metadata.json" class="report-link">📋 Metadados</a>
                        </div>
                    `;
                    grid.appendChild(card);
                });
            })
            .catch(error => {
                console.error('Erro ao carregar relatórios:', error);
                document.getElementById('reports-grid').innerHTML = 
                    '<p>Erro ao carregar lista de relatórios.</p>';
            });
    </script>
</body>
</html>
EOF

    # Gerar arquivo JSON com lista de relatórios
    cat > "$GIT_REPORTS_DIR/reports.json" << EOF
{
  "generated_at": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "reports": [
EOF

    # Adicionar cada relatório ao JSON
    for dir in "$GIT_REPORTS_DIR"/*/; do
        if [ -d "$dir" ] && [ "$(basename "$dir")" != "latest" ]; then
            if [ -f "$dir/metadata.json" ]; then
                # Ler metadados e adicionar ao JSON
                METADATA=$(cat "$dir/metadata.json")
                echo "    $METADATA," >> "$GIT_REPORTS_DIR/reports.json"
            fi
        fi
    done

    # Remover última vírgula e fechar JSON
    sed -i '' '$s/,$//' "$GIT_REPORTS_DIR/reports.json"
    echo "  ]" >> "$GIT_REPORTS_DIR/reports.json"
    echo "}" >> "$GIT_REPORTS_DIR/reports.json"
    
    echo "📄 Índice HTML gerado: $GIT_REPORTS_DIR/index.html"
}
