// Package db initializes and works on the database
package db

import (
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var DB *gorm.DB

// UpsertMedia performs an upsert operation on any media item
func UpsertMedia(item any, updateColumns []string) error {
	return DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "username"}, {Name: "external_id"}},
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(item).Error
}

// InitDatabase initializes the database at the given path
func InitDatabase(dbPath string) {
	// Extract directory from path and create if needed
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic("failed to create database directory: " + err.Error())
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}


// MigrateModels migrates the provided models into the database
func MigrateModels(models ...any) error {
	err := DB.AutoMigrate(models...)
	if err != nil {
		return err
	}

	// Create composite unique indexes for each table
	if err := DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_animes_user_external ON animes(username, external_id)").Error; err != nil {
		return err
	}
	if err := DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_mangas_user_external ON mangas(username, external_id)").Error; err != nil {
		return err
	}

	return nil
}
