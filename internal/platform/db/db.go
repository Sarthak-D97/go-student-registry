package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Sarthak-D97/go_stuAPI/entity"
	"github.com/Sarthak-D97/go_stuAPI/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// debugLog writes a single NDJSON debug entry for this debug session.
func debugLog(hypothesisID, message string, data map[string]any) {
	// #region agent log
	entry := map[string]any{
		"id":           fmt.Sprintf("log_%d_db", time.Now().UnixNano()),
		"timestamp":    time.Now().UnixMilli(),
		"location":     "internal/platform/db/db.go:NewPostgres",
		"message":      message,
		"data":         data,
		"runId":        "db-init",
		"hypothesisId": hypothesisID,
	}

	b, err := json.Marshal(entry)
	if err != nil {
		return
	}

	f, err := os.OpenFile("/Users/sarthak/Downloads/Dev/VS Projects/Go/student_api/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()

	if _, err := f.Write(append(b, '\n')); err != nil {
		return
	}
	// #endregion agent log
}

// NewPostgres creates a new GORM Postgres connection using the provided config.
func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	debugLog("H1", "attempting Postgres connection", map[string]any{
		"host": cfg.DBHost,
		"port": cfg.DBPort,
		"user": cfg.DBUser,
		"db":   cfg.DBName,
	})

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		debugLog("H1", "Postgres connection failed", map[string]any{
			"error": err.Error(),
		})
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		debugLog("H1", "Postgres sql.DB acquisition failed", map[string]any{
			"error": err.Error(),
		})
		return nil, err
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	if err := db.AutoMigrate(&entity.Student{}); err != nil {
		debugLog("H1", "Postgres automigrate failed", map[string]any{
			"error": err.Error(),
		})
		return nil, err
	}

	log.Println("connected to Postgres and ran migrations")
	debugLog("H1", "Postgres connection and migrations succeeded", nil)

	return db, nil
}
