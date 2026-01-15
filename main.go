package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/luciluz/psiconexo/internal/db"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx := context.Background()

	// 1. Conectar a la DB
	conn, err := sql.Open("sqlite3", "./psiconexo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	schemaContent, err := os.ReadFile("internal/db/schema.sql")
	if err != nil {
		log.Fatal("Error leyendo el schema.sql: ", err)
	}

	// Ejecutamos ese SQL en la base de datos
	_, err = conn.Exec(string(schemaContent))
	if err != nil {
		log.Fatal("Error creando tablas: ", err)
	}

	queries := db.New(conn)

	// 3. Crear un Psicólogo (Simulacro)
	// Usamos timestamp en el mail para que no de error de "Unique" si lo corres 2 veces
	psy, err := queries.CreatePsychologist(ctx, db.CreatePsychologistParams{
		Name:  "Lic. Ana Freud",
		Email: fmt.Sprintf("ana_%d@example.com", 1), // Truco barato para pruebas
		Phone: sql.NullString{String: "555-0001", Valid: true},
	})
	if err != nil {
		log.Fatal("Error creando psicólogo: ", err)
	}
	fmt.Printf("Psicólogo creado: %s (ID: %d)\n", psy.Name, psy.ID)

	// 4. Crear un Paciente para Ana
	pat, err := queries.CreatePatient(ctx, db.CreatePatientParams{
		Name:           "Paciente Ejemplo",
		Email:          fmt.Sprintf("paciente_%d@test.com", 1),
		Phone:          sql.NullString{String: "555-1234", Valid: true},
		PsychologistID: psy.ID, // Aquí vinculamos!
	})
	if err != nil {
		log.Fatal("Error creando paciente: ", err)
	}
	fmt.Printf("Paciente creado: %s asignado a Psicólogo ID %d\n", pat.Name, pat.PsychologistID)

	// 5. Listar pacientes de Ana
	patients, err := queries.ListPatients(ctx, psy.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("--- Lista de Pacientes ---")
	for _, p := range patients {
		fmt.Printf("- %s (%s)\n", p.Name, p.Email)
	}
}
