package utils

import (
	"encoding/json"
	"fmt"
)

// ToJSON converte um objeto para JSON
func ToJSON(v interface{}) (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao converter para JSON: %w", err)
	}
	return string(data), nil
}

// FromJSON converte JSON para um objeto
func FromJSON(data string, v interface{}) error {
	if err := json.Unmarshal([]byte(data), v); err != nil {
		return fmt.Errorf("erro ao fazer parse de JSON: %w", err)
	}
	return nil
}

// PrettyPrint imprime JSON formatado
func PrettyPrint(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
