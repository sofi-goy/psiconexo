package service

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/luciluz/psiconexo/internal/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// Helper para configurar DB en tests de servicio
func setupServiceTest(t *testing.T) *Service {
	// 1. Conexión en memoria
	conn, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	// 2. Activar FK
	_, err = conn.Exec("PRAGMA foreign_keys = ON;")
	assert.NoError(t, err)

	// 3. Cargar Schema (asumiendo que corres el test desde la raíz o ajustando path)
	// NOTA: Go busca el archivo relativo a donde está el test.
	// Si internal/service está a 2 niveles de la raíz, subimos ../../
	schema, err := os.ReadFile("../../internal/db/schema.sql")
	if err != nil {
		// Fallback por si lo corres desde otro lado
		schema, err = os.ReadFile("schema.sql")
	}
	assert.NoError(t, err, "No se encontró schema.sql")

	_, err = conn.Exec(string(schema))
	assert.NoError(t, err)

	queries := db.New(conn)
	return NewService(queries, conn)
}

func TestAppointmentCollision(t *testing.T) {
	svc := setupServiceTest(t)
	ctx := context.Background()

	// 1. Datos Base (Psicólogo y Paciente)
	psy, _ := svc.queries.CreatePsychologist(ctx, db.CreatePsychologistParams{
		Name: "Dr. House", Email: "house@ppth.com",
	})
	pat, _ := svc.queries.CreatePatient(ctx, db.CreatePatientParams{
		Name: "Paciente 1", Email: "p1@test.com", PsychologistID: psy.ID,
	})

	targetDate := time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC)

	// 2. Agendar PRIMER turno (10:00 a 11:00) -> DEBE FUNCIONAR
	_, err := svc.CreateAppointment(ctx, CreateAppointmentRequest{
		PsychologistID: psy.ID,
		PatientID:      pat.ID,
		Date:           targetDate,
		StartTime:      "10:00",
		Duration:       60, // Termina 11:00
	})
	assert.NoError(t, err, "El primer turno debería guardarse ok")

	// 3. Probar colisiones

	cases := []struct {
		name       string
		startTime  string
		duration   int
		shouldFail bool
	}{
		{
			name:       "Colisión Total (Misma hora 10:00)",
			startTime:  "10:00",
			duration:   60,
			shouldFail: true,
		},
		{
			name:       "Colisión Parcial (Empieza antes 09:30, termina 10:30)",
			startTime:  "09:30",
			duration:   60, // Termina 10:30 (Choca con el inicio del otro)
			shouldFail: true,
		},
		{
			name:       "Colisión Parcial (Empieza durante 10:30, termina 11:30)",
			startTime:  "10:30",
			duration:   60,
			shouldFail: true,
		},
		{
			name:       "Sin Colisión (Empieza justo cuando termina el otro 11:00)",
			startTime:  "11:00",
			duration:   60,
			shouldFail: false, // Debería pasar
		},
		{
			name:       "Sin Colisión (Termina justo cuando empieza el otro 10:00)",
			startTime:  "09:00",
			duration:   60,    // Termina 10:00
			shouldFail: false, // Debería pasar
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.CheckAvailability(ctx, psy.ID, targetDate, tc.startTime, tc.duration)
			if tc.shouldFail {
				assert.Error(t, err, "Debería haber detectado colisión")
			} else {
				assert.NoError(t, err, "Debería estar libre")
			}
		})
	}
}
