package db_test

import (
	"os"
	"log"
	"fmt"
	"testing"
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/luciluz/psiconexo/internal/db"
)

func TestPsychologistNameNotNull(t *testing.T) {
	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	// Crear tablas (si no existen)
	schemaContent, err := os.ReadFile("../schema.sql")
	if err != nil {
		log.Fatal("Error leyendo schema.sql: ", err)
	}
	if _, err := conn.Exec(string(schemaContent)); err != nil {
		log.Fatal("Error ejecutando schema: ", err)
	}
	
	// Inicializar las queries de SQLC
	queries := db.New(conn)
	ctx := context.Background()

	_, err = queries.CreatePsychologist(ctx, db.CreatePsychologistParams{
		Name: "luci",
		Email: "psicologo@mail",
		Phone: sql.NullString{String: "1234", Valid: true},
	})

	psy, err := queries.ListPsychologists(ctx)

	if len(psy) == 0 {
		t.Error("No se cargó psicologo")
	}
	if len(psy) > 1 {
		t.Error("Hay más de un ppsicoloco")
	}

	luci := psy[0]
	
	fmt.Println("El psicologo es", luci.Name)
}
