package abstraction

import (
	"context"
	"time"
)

type Migrator interface {
	RunMigrations() error
}

type VersionedEntity interface {
	GetDomainID() string
	GetDomainIDFieldName() string
	GetCollectionName() string

	CreateNewEmpty() VersionedEntity

	SetCreatedAt(createdAt *time.Time)
	SetUpdatedAt(updatedAt *time.Time)
}

type GenericRepository interface {
	CreateEntity(ctx context.Context, entity VersionedEntity) error
	GetEntity(ctx context.Context, id string) (VersionedEntity, error)
	UpdateEntity(ctx context.Context, entity VersionedEntity) (VersionedEntity, error)
	DeleteEntity(ctx context.Context, id string) error
}
