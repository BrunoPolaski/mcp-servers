package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/infra/thirdparty/database"
	"github.com/BrunoPolaski/internal-registry-mcp-server/internal/infra/thirdparty/logger"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type fixtureTarget struct {
	FileName  string
	TableName string
}

var loadOrder = []fixtureTarget{
	{FileName: "files.json", TableName: "files"},
	{FileName: "addresses.json", TableName: "addresses"},
	{FileName: "personal_informations.json", TableName: "personal_informations"},
	{FileName: "admins.json", TableName: "admins"},
	{FileName: "analysts.json", TableName: "analysts"},
	{FileName: "customer_relationships.json", TableName: "customer_relationships"},
	{FileName: "persons.json", TableName: "persons"},
	{FileName: "person_addresses.json", TableName: "person_addresses"},
	{FileName: "person_documents.json", TableName: "person_documents"},
	{FileName: "contracted_products.json", TableName: "contracted_products"},
	{FileName: "internal_payment_records.json", TableName: "internal_payment_records"},
	{FileName: "pre_approved_limits.json", TableName: "pre_approved_limits"},
	{FileName: "income_declarations.json", TableName: "income_declarations"},
	{FileName: "api_keys.json", TableName: "api_keys"},
	{FileName: "tokens.json", TableName: "tokens"},
	{FileName: "users.json", TableName: "users"},
	{FileName: "sessions.json", TableName: "sessions"},
}

func main() {
	fixturesDir := flag.String("dir", "./fixtures", "Path to fixtures directory")
	truncate := flag.Bool("truncate", false, "Truncate tables before loading")
	flag.Parse()

	logger.Init()

	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		logger.Error("Error loading location", err)
	}
	time.Local = location

	err = godotenv.Load("../../.env")
	if err != nil {
		logger.Error("Error loading .env file", err)
	}

	if _, err := os.Stat("/.dockerenv"); err == nil {
		logger.Info("Running inside Docker container")
	} else {
		logger.Info("Running outside Docker container")
		if os.Getenv("ENV") == "production" {
			logger.Info("Using production database configuration")
		} else if os.Getenv("ENV") == "dev" {
			logger.Info("Using development database configuration")
		} else {
			os.Setenv("DATABASE_URL", "postgres://postgres:123456@127.0.0.1:5433/internal-registry-mcp?sslmode=disable")
		}
	}

	db := database.NewGormAdapter().Connect()

	if *truncate {
		err := truncateTables(db)
		if err != nil {
			fmt.Printf("failed to truncate tables: %v\n", err)
			os.Exit(1)
		}
	}

	for _, target := range loadOrder {
		path := filepath.Join(*fixturesDir, target.FileName)
		rows, err := loadFixture(path)
		if err != nil {
			fmt.Printf("failed to load %s: %v\n", target.FileName, err)
			os.Exit(1)
		}

		if len(rows) == 0 {
			fmt.Printf("skipping %s (no rows)\n", target.FileName)
			continue
		}

		if err := db.Table(target.TableName).Create(&rows).Error; err != nil {
			fmt.Printf("failed to insert %s: %v\n", target.TableName, err)
			os.Exit(1)
		}
		fmt.Printf("inserted %d rows into %s\n", len(rows), target.TableName)
	}
}

func loadFixture(path string) ([]map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rows []map[string]any
	if err := json.Unmarshal(data, &rows); err != nil {
		return nil, err
	}

	for i := range rows {
		rows[i] = normalizeNumbers(rows[i])
	}

	return rows, nil
}

func normalizeNumbers(row map[string]any) map[string]any {
	for key, value := range row {
		switch v := value.(type) {
		case float64:
			if v == float64(int64(v)) {
				row[key] = int64(v)
			}
		case map[string]any:
			row[key] = normalizeNumbers(v)
		}
	}
	return row
}

func truncateTables(db *gorm.DB) error {
	var tableNames []string
	for _, target := range loadOrder {
		tableNames = append(tableNames, target.TableName)
	}

	stmt := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(tableNames, ", "))
	return db.Exec(stmt).Error
}
