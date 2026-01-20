package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Helper para crear fechas rápidas (Año, Mes, Día)
func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func TestCreateAppointment(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		date        time.Time
		startTime   string
		expectError bool
		setup       func(q *Queries, psyID int64, patID int64)
	}{
		{
			name:        "Success - Cita normal",
			date:        date(2026, 1, 20),
			startTime:   "10:00",
			expectError: false,
		},
		{
			name:        "Fail - Duplicado exacto (Mismo doctor, dia y hora)",
			date:        date(2026, 1, 20),
			startTime:   "10:00",
			expectError: true,
			setup: func(q *Queries, psyID int64, patID int64) {
				// Creamos la "cita original" que va a ocupar el lugar
				_, err := q.CreateAppointment(ctx, CreateAppointmentParams{
					PsychologistID:  psyID,
					PatientID:       patID,
					Date:            date(2026, 1, 20),
					StartTime:       "10:00",
					DurationMinutes: 60,
					// RescheduledFromID es NULL para citas nuevas
					RescheduledFromID: sql.NullInt64{Valid: false},
				})
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := setupTestDB(t)

			// 1. Setup Básico (Crear Doctor y Paciente)
			psy, _ := q.CreatePsychologist(ctx, CreatePsychologistParams{
				Name:  "Dr. Strange",
				Email: "strange@time.com",
			})
			pat, _ := q.CreatePatient(ctx, CreatePatientParams{
				Name:           "Wong",
				Email:          "wong@kamartaj.com",
				PsychologistID: psy.ID,
			})

			// 2. Setup específico del test (ej: crear cita previa)
			if tt.setup != nil {
				tt.setup(q, psy.ID, pat.ID)
			}

			// 3. Ejecución
			_, err := q.CreateAppointment(ctx, CreateAppointmentParams{
				PsychologistID:    psy.ID,
				PatientID:         pat.ID,
				Date:              tt.date,
				StartTime:         tt.startTime,
				DurationMinutes:   60,
				RescheduledFromID: sql.NullInt64{Valid: false},
			})

			// 4. Validación
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListAppointmentsInDateRange(t *testing.T) {
	ctx := context.Background()
	q := setupTestDB(t)

	// 1. Setup: 1 Doctor, 1 Paciente
	psy, _ := q.CreatePsychologist(ctx, CreatePsychologistParams{Name: "Doc", Email: "d@d.com"})
	pat, _ := q.CreatePatient(ctx, CreatePatientParams{Name: "Pat", Email: "p@p.com", PsychologistID: psy.ID})

	// 2. Insertamos 3 citas:
	// - Una la semana pasada (FUERA de rango)
	// - Una hoy (DENTRO de rango)
	// - Una mañana (DENTRO de rango)

	inputs := []struct {
		d    time.Time
		time string
	}{
		{date(2026, 1, 10), "09:00"}, // Pasada
		{date(2026, 1, 20), "10:00"}, // Hoy (Target)
		{date(2026, 1, 21), "11:00"}, // Mañana (Target)
	}

	for _, in := range inputs {
		_, err := q.CreateAppointment(ctx, CreateAppointmentParams{
			PsychologistID:    psy.ID,
			PatientID:         pat.ID,
			Date:              in.d,
			StartTime:         in.time,
			DurationMinutes:   60,
			RescheduledFromID: sql.NullInt64{Valid: false},
		})
		assert.NoError(t, err)
	}

	// 3. Consultamos el rango: Del 19 al 25 de Enero
	startFilter := date(2026, 1, 19)
	endFilter := date(2026, 1, 25)

	list, err := q.ListAppointmentsInDateRange(ctx, ListAppointmentsInDateRangeParams{
		PsychologistID: psy.ID,
		Date:           startFilter,
		Date_2:         endFilter,
	})

	assert.NoError(t, err)

	// 4. Validamos
	// Debería traer solo 2 citas (la del 20 y la del 21). La del 10 queda fuera.
	assert.Len(t, list, 2)
	assert.Equal(t, "10:00", list[0].StartTime)
	assert.Equal(t, "11:00", list[1].StartTime)
}
