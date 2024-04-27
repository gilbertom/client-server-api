package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gilbertom/client-server-api/internal/infra/database"
	"github.com/gilbertom/client-server-api/internal/infra/dto"
)

func main() {
	http.HandleFunc("/cotacao", GetCotacao)
	fmt.Println("Servidor escutando na porta 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}

// GetCotacao é um handler que faz uma solicitação HTTP para a API externa
func GetCotacao(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200 * time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Fatalf("Erro ao fazer a solicitação HTTP: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Requisição na API de Cotação cancelada por decurso de timeout de 200ms")
			w.WriteHeader(http.StatusGatewayTimeout)
		} else {
			fmt.Println("Erro ao fazer a requisição:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erro ao ler o corpo da resposta: %v", err)
	}
	defer resp.Body.Close()
	
	var c dto.Cotacao
	if err := json.Unmarshal(body, &c); err != nil {
		log.Fatalf("Erro ao decodificar o JSON: %v", err)
	}

	ctxDb, cancelDb := context.WithTimeout(context.Background(), 10 * time.Millisecond)
	defer cancelDb()

	err = database.CreateCotacao(ctxDb, c)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("Inserção no Banco de Dados cancelada por decurso de timeout de 10ms")
			w.WriteHeader(http.StatusGatewayTimeout)
		} else {
			fmt.Println("Erro ao fazer insert na tabela currency:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	var bid dto.BidOutput
	bid.Bid = c.UsdBrl.Bid

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(bid)
}
