package integration_tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/jackc/pgx/v5/stdlib"

	"zahne/internal/controller/restapi/v1"
	"zahne/internal/repository/persistent"
	"zahne/internal/usecase/patient"
)

// setupTestDatabase creates a PostgreSQL testcontainer and returns the database connection
// along with a cleanup function to terminate the container
func setupTestDatabase(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	dbName := "testdb"
	dbUser := "testuser"
	dbPassword := "testpass"

	postgresContainer, err := postgres.Run(
		ctx,
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, host, port.Port(), dbName)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	// Create patients table
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS patients (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			cpf VARCHAR(255) NOT NULL,
			phone VARCHAR(255),
			email VARCHAR(255),
			date_of_birth TIMESTAMP NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		t.Fatalf("failed to create patients table: %v", err)
	}

	cleanup := func() {
		db.Close()
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	return db, cleanup
}

// setupPatientRouter creates a Gin router with patient endpoints configured
func setupPatientRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	patientRepo := persistent.NewPatientRepository(db)
	patientUseCase := patient.New(patientRepo)
	patientHandler := v1.NewPatientHandler(patientUseCase)

	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/patient", patientHandler.Create)
		apiV1.GET("/patient/:id", patientHandler.Get)
		apiV1.GET("/patients", patientHandler.GetAll)
		apiV1.PUT("/patient/:id", patientHandler.Update)
		apiV1.DELETE("/patient/:id", patientHandler.Delete)
	}

	return router
}
