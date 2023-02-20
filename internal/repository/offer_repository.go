package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"mongo-db-connector/domain"
	"mongo-db-connector/domain/abstraction"
	"time"
)

const (
	offerDomainIDFieldName string = "offer_id"
	offersCollection       string = "offers"
)

type DBOffer struct {
	ID            primitive.ObjectID `bson:"_id"`
	*domain.Offer `bson:"inline"`
}

func (dbOffer *DBOffer) GetDomainID() string {
	return dbOffer.OfferID
}

func (dbOffer *DBOffer) GetDomainIDFieldName() string {
	return offerDomainIDFieldName
}

func (dbOffer *DBOffer) GetCollectionName() string {
	return offersCollection
}

func (dbOffer *DBOffer) CreateNewEmpty() abstraction.VersionedEntity {
	return new(DBOffer)
}

func (dbOffer *DBOffer) SetCreatedAt(createdAt *time.Time) {
	dbOffer.CreatedAt = createdAt
}

func (dbOffer *DBOffer) SetUpdatedAt(updatedAt *time.Time) {
	dbOffer.UpdatedAt = updatedAt
}

type offerRepository struct {
	abstraction.GenericRepository
}

var _ abstraction.OfferRepository = &offerRepository{}

func NewOfferRepository(dbClient *mongo.Client, database string) abstraction.OfferRepository {
	db := dbClient.Database(database)
	genericRepo := NewGenericRepository(db, &DBOffer{})

	return &offerRepository{genericRepo}
}
