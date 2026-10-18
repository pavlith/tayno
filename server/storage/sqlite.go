package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"myst/server/core"
	"myst/server/models"
)

// Repository implements the core persistence interfaces with GORM and SQLite.
// One pooled connection serializes writes and keeps in-memory databases stable.
type Repository struct{ db *gorm.DB }

var (
	_ core.UserRepository       = (*Repository)(nil)
	_ core.KeystoreRepository   = (*Repository)(nil)
	_ core.InvitationRepository = (*Repository)(nil)
)

// Open opens a SQLite file (or :memory:) and applies additive schema migrations.
func Open(path string) (*Repository, error) {
	if path == "" {
		return nil, fmt.Errorf("SQLite database path is empty")
	}
	if path != ":memory:" && !strings.HasPrefix(path, "file:") {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
	}
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		TranslateError: true, Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("open SQLite: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	fail := func(err error) (*Repository, error) { _ = sqlDB.Close(); return nil, err }
	for _, pragma := range []string{"PRAGMA busy_timeout = 5000", "PRAGMA foreign_keys = ON", "PRAGMA journal_mode = WAL"} {
		if err := db.Exec(pragma).Error; err != nil {
			return fail(fmt.Errorf("configure SQLite: %w", err))
		}
	}
	if err := db.AutoMigrate(&models.User{}, &models.Keystore{}, &models.Invitation{}); err != nil {
		return fail(fmt.Errorf("migrate SQLite schema: %w", err))
	}
	if err := db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_active_invitation
 ON invitations (keystore_id, inviter_id, invitee_id)
 WHERE status IN ('pending', 'accepted', 'finalized')`).Error; err != nil {
		return fail(fmt.Errorf("create invitation index: %w", err))
	}
	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
