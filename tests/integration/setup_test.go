package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	goose "github.com/pressly/goose/v3"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/xuri/excelize/v2"

	"leetcode-spaced-repetition/controllers"
	dbmigrations "leetcode-spaced-repetition/db/migrations"
	"leetcode-spaced-repetition/internal"
	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/repositories"
	"leetcode-spaced-repetition/services"

	"go.uber.org/zap"
)

var (
	testDB     *sql.DB
	testRouter *gin.Engine
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(
		ctx,
		"postgres:16",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start postgres container: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = pgContainer.Terminate(ctx) }()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get connection string: %v\n", err)
		os.Exit(1)
	}

	testDB, err = internal.GetDBConnFromDSN(connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open DB: %v\n", err)
		os.Exit(1)
	}
	defer testDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set goose dialect: %v\n", err)
		os.Exit(1)
	}
	goose.SetBaseFS(dbmigrations.MigrationFiles)
	if err := goose.Up(testDB, "."); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run migrations: %v\n", err)
		os.Exit(1)
	}

	gin.SetMode(gin.TestMode)
	logger := zap.NewNop()

	problemsRepo := repositories.NewProblemPostgresRepository(testDB, logger)
	cardStateRepo := repositories.NewProblemCardStatePostgresRepository(testDB, logger)
	problemsService := services.NewProblemsService(problemsRepo, cardStateRepo, logger)

	testRouter = gin.New()
	controllers.RegisterRoutes(testRouter, problemsService, logger)

	os.Exit(m.Run())
}

// insertTestProblem seeds a problem row directly via SQL to satisfy FK constraints.
func insertTestProblem(t *testing.T, p models.Problem) {
	t.Helper()
	_, err := testDB.Exec(
		`INSERT INTO problems (id, title, slug, description, difficulty)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (id) DO NOTHING`,
		p.ID, p.Title, p.Slug, p.Description, p.Difficulty,
	)
	if err != nil {
		t.Fatalf("insertTestProblem: %v", err)
	}
}

// insertTestTopic seeds a topic for a problem.
func insertTestTopic(t *testing.T, problemID int, topic string) {
	t.Helper()
	_, err := testDB.Exec(
		`INSERT INTO problem_topics (problem_id, topic)
		 VALUES ($1, $2)
		 ON CONFLICT (problem_id, topic) DO NOTHING`,
		problemID, topic,
	)
	if err != nil {
		t.Fatalf("insertTestTopic: %v", err)
	}
}

// clearSubmissions deletes all rows from problem_submissions between tests.
func clearSubmissions(t *testing.T) {
	t.Helper()
	if _, err := testDB.Exec(`DELETE FROM problem_submissions`); err != nil {
		t.Fatalf("clearSubmissions: %v", err)
	}
}

// clearCardStates deletes all rows from problem_card_states between tests.
func clearCardStates(t *testing.T) {
	t.Helper()
	if _, err := testDB.Exec(`DELETE FROM problem_card_states`); err != nil {
		t.Fatalf("clearCardStates: %v", err)
	}
}

// makeXLSXFile builds an in-memory .xlsx file with a header row followed by the
// provided data rows and returns a *bytes.Buffer that implements io.Reader.
func makeXLSXFile(t *testing.T, dataRows [][]string) *bytes.Buffer {
	t.Helper()
	f := excelize.NewFile()
	sheet := f.GetSheetName(0)

	header := []string{"Problem #", "Date", "Time Taken", "Confidence"}
	for col, v := range header {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		if err := f.SetCellValue(sheet, cell, v); err != nil {
			t.Fatalf("makeXLSXFile header: %v", err)
		}
	}

	for rowIdx, row := range dataRows {
		for colIdx, v := range row {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, rowIdx+2)
			if err := f.SetCellValue(sheet, cell, v); err != nil {
				t.Fatalf("makeXLSXFile data: %v", err)
			}
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		t.Fatalf("makeXLSXFile write: %v", err)
	}
	return buf
}

// testProblem is a convenience fixture for seeding a known problem.
var testProblem = models.Problem{
	ID:          1,
	Title:       "Two Sum",
	Slug:        "two-sum",
	Description: "Given an array of integers...",
	Difficulty:  models.EasyDifficulty,
}
