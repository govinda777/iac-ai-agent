package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileExists verifica se um arquivo existe
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// IsDirectory verifica se o path é um diretório
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ListTerraformFiles lista todos os arquivos .tf em um diretório
func ListTerraformFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".tf") {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("erro ao listar arquivos: %w", err)
	}

	return files, nil
}

// ReadFile lê conteúdo de um arquivo
func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo: %w", err)
	}
	return string(content), nil
}

// WriteFile escreve conteúdo em um arquivo
func WriteFile(path, content string) error {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo: %w", err)
	}
	return nil
}

// EnsureDir garante que um diretório existe
func EnsureDir(dir string) error {
	if !FileExists(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("erro ao criar diretório: %w", err)
		}
	}
	return nil
}
