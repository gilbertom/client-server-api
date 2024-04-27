package database

import (
	"context"
	"database/sql"
	"log"

	"github.com/gilbertom/client-server-api/internal/infra/dto"

	// Blank import for database driver
	_ "github.com/mattn/go-sqlite3"
)

// CreateCotacao cria uma nova entrada no banco de dados
func CreateCotacao(ctx context.Context, c dto.Cotacao) error {
	db, err := sql.Open("sqlite3", "currency.db")
	if err != nil {
		log.Fatalf("Erro ao abrir o banco de dados: %v", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS currency (
						code TEXT,
						codein TEXT,
						name TEXT,
						high TEXT,
						low TEXT,
						varBid TEXT,
						pctChange TEXT,
						bid TEXT,
						ask TEXT,
						timestamp TEXT,
						create_date TEXT)`)
	if err != nil {
		log.Fatalf("Erro ao criar a tabela: %v", err)
		return err
	}

	stmt, err := db.Prepare(`INSERT INTO currency (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatalf("Erro ao preparar a declaração de inserção: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, c.UsdBrl.Code, c.UsdBrl.Codein, c.UsdBrl.Name, c.UsdBrl.High, c.UsdBrl.Low, c.UsdBrl.VarBid, c.UsdBrl.PctChange, c.UsdBrl.Bid, c.UsdBrl.Ask, c.UsdBrl.Timestamp, c.UsdBrl.CreateDate)
	if err != nil {
		log.Fatalf("Erro ao inserir dados na tabela: %v", err)
		return err
	}
	return nil
}
