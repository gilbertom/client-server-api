package file

import (
	"fmt"
	"os"
)

func SaveToFile(bid string) error {
	filePath := "cotacao.txt"

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Erro ao criar arquivo:", err)
			return err
		}
		defer file.Close()

		if _, err := file.WriteString(fmt.Sprintf("Dólar: %s\n", bid)); err != nil {
			fmt.Println("Erro ao escrever no arquivo:", err)
			return err
		}
	} else if err == nil {
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Erro ao abrir arquivo:", err)
			return err
		}
		defer file.Close()

		if _, err := file.WriteString(fmt.Sprintf("Dólar: %s\n", bid)); err != nil {
			fmt.Println("Erro ao escrever no arquivo:", err)
			return err
		}
	} else {
		fmt.Println("Erro ao verificar arquivo:", err)
		return err
	}
	return nil
}
