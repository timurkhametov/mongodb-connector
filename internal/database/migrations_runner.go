package database

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
	"mongo-db-connector/domain/abstraction"
	"mongo-db-connector/internal/infrastructure"
)

const migrationsSourceURL = "file://internal/database/migrations"

var _ abstraction.Migrator = &mongoMigrator{}

type mongoMigrator struct {
	mongoClient *mongo.Client
	config      *infrastructure.MigrationConfig
}

func NewMigrationRunner(client *mongo.Client, config *infrastructure.MigrationConfig) abstraction.Migrator {
	return &mongoMigrator{
		mongoClient: client,
		config:      config,
	}
}

func (m *mongoMigrator) RunMigrations() error {
	config := &mongodb.Config{
		DatabaseName:         m.config.Database,
		MigrationsCollection: m.config.CollectionName,
		TransactionMode:      m.config.UseTransactionMode,
	}

	if m.config.AdvisorLock != nil {
		config.Locking = mongodb.Locking{
			CollectionName: m.config.AdvisorLock.CollectionName,
			Timeout:        m.config.AdvisorLock.Timeout,
			Enabled:        true,
			Interval:       m.config.AdvisorLock.Interval,
		}
	}

	driver, err := mongodb.WithInstance(m.mongoClient, config)
	if err != nil {
		return err
	}

	migrator, err := migrate.NewWithDatabaseInstance(migrationsSourceURL, m.config.Database, driver)
	if err != nil {
		return err
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
