package service

import (
	"database/sql"

	"github.com/luciluz/psiconexo/internal/db"
)

// Service agrupa toda la lógica de negocio.
// Recibe un *db.Queries para hablar con la base de datos.
type Service struct {
	queries *db.Queries
	db      *sql.DB // Lo guardamos por si necesitamos transacciones complejas a futuro
}

// NewService es el constructor
func NewService(queries *db.Queries, dbConn *sql.DB) *Service {
	return &Service{
		queries: queries,
		db:      dbConn,
	}
}

// Aquí pondremos métodos auxiliares comunes
