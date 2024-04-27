package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gilbertom/client-server-api/internal/infra/dto"
	"github.com/gilbertom/client-server-api/internal/infra/file"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300 * time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("Erro ao fazer a solicitação HTTP: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Requisição cancelada por decurso de timeout de 300ms")
		} else {
			fmt.Println("Erro ao fazer a requisição:", err)
		}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erro ao ler o corpo da resposta: %v", err)
	}

	if (resp.StatusCode == http.StatusGatewayTimeout) || (resp.StatusCode == http.StatusInternalServerError) {
		log.Fatalf("Erro ao fazer a requisição. HTTP Code: %d", resp.StatusCode)
	}

	var bid dto.BidOutput
	if err := json.Unmarshal(body, &bid); err != nil {
		log.Fatalf("Erro ao decodificar o JSON: %v", err)
	}

	err = file.SaveToFile(bid.Bid)
	if err != nil {
		log.Fatalf("Erro ao salvar o arquivo: %v", err)
	}

	fmt.Printf("O valor do dólar é R$ %s\n", bid.Bid)
	fmt.Printf("Arquivo cotacao.txt salvo com sucesso\n")
	fmt.Println("Fim do programa client.")
}
