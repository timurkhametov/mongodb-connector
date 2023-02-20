package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"mongo-db-connector/domain/abstraction"
	"reflect"
	"time"
)

var _ abstraction.GenericRepository = &genericRepository{}

type genericRepository struct {
	collection *mongo.Collection
	baseEntity abstraction.VersionedEntity
}

func NewGenericRepository(database *mongo.Database, baseEntity abstraction.VersionedEntity) abstraction.GenericRepository {
	collectionName := baseEntity.GetCollectionName()

	return &genericRepository{
		collection: database.Collection(collectionName),
		baseEntity: baseEntity,
	}
}

func (repo *genericRepository) CreateEntity(ctx context.Context, entity abstraction.VersionedEntity) error {
	repo.mustInitialized()
	err := repo.checkType(entity)
	if err != nil {
		return err
	}

	if entity == nil {
		return fmt.Errorf("no entities to insert")
	}

	now := time.Now()
	entity.SetCreatedAt(&now)
	entity.SetUpdatedAt(nil)

	_, err = repo.collection.InsertOne(ctx, entity)
	return err
}

func (repo *genericRepository) GetEntity(ctx context.Context, id string) (abstraction.VersionedEntity, error) {
	repo.mustInitialized()

	output := repo.baseEntity.CreateNewEmpty()
	err := repo.checkType(output)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{repo.baseEntity.GetDomainIDFieldName(), id}}
	err = repo.collection.FindOne(ctx, filter).Decode(output)
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (repo *genericRepository) UpdateEntity(ctx context.Context, entity abstraction.VersionedEntity) (abstraction.VersionedEntity, error) {
	repo.mustInitialized()
	err := repo.checkType(entity)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	entity.SetUpdatedAt(&now)

	filter := bson.D{{entity.GetDomainIDFieldName(), entity.GetDomainID()}}
	update := bson.D{{"$set", entity}}

	result, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("no entities were updated")
	}

	return entity, nil
}

func (repo *genericRepository) DeleteEntity(ctx context.Context, id string) error {
	repo.mustInitialized()

	filter := bson.D{{repo.baseEntity.GetDomainIDFieldName(), id}}

	_, err := repo.collection.DeleteOne(ctx, filter)
	return err
}

func (repo *genericRepository) checkType(inputEntity abstraction.VersionedEntity) error {
	if reflect.TypeOf(repo.baseEntity) != reflect.TypeOf(inputEntity) {
		return fmt.Errorf("type mismatch")
	}

	return nil
}

func (repo *genericRepository) mustInitialized() {
	if repo == nil || repo.baseEntity == nil || repo.collection == nil {
		panic("repository is not initialized")
	}
}
