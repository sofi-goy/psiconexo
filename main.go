package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/luciluz/psiconexo/internal/db"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	ctx := context.Background()

	// Conectar a SQLite
	conn, err := sql.Open("sqlite3", "./psiconexo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Crear tablas (si no existen)
	schemaContent, err := os.ReadFile("internal/db/schema.sql")
	if err != nil {
		log.Fatal("Error leyendo schema.sql: ", err)
	}
	if _, err := conn.Exec(string(schemaContent)); err != nil {
		log.Fatal("Error ejecutando schema: ", err)
	}

	// Inicializar las queries de SQLC
	queries := db.New(conn)

	_ = queries // método backer
	_ = ctx

	// Crear el enrutador
	r := chi.NewRouter()

	// Middlewares básicos
	r.Use(middleware.Logger)    // Muestra en consola cada vez que alguien entra
	r.Use(middleware.Recoverer) // Si algo explota, que el servidor no se caiga

	// Definir una ruta básica: GET /
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hola, Psiconexo te ama x_x"))
	})

	// Definir ruta para probar la DB: GET /ping
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		if err := conn.Ping(); err != nil {
			http.Error(w, "Base de datos muerta x_x", 500)
			return
		}
		w.Write([]byte("Yeey! se conectó la base de datos correctamente."))
	})

	// Arrancar el servidor en el puerto 3000
	fmt.Println("Servidor escuchando en http://localhost:3000")
	http.ListenAndServe(":3000", r)
}
